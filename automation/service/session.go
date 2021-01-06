package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"sync"
)

type (
	session struct {
		store      store.Storer
		actionlog  actionlog.Recorder
		ac         sessionAccessController
		log        *zap.Logger
		mux        *sync.RWMutex
		pool       map[uint64]*types.Session
		spawnQueue chan *spawn
	}

	spawn struct {
		session chan *wfexec.Session
		graph   *wfexec.Graph
		trace   bool
	}

	sessionAccessController interface {
		CanSearchSessions(context.Context) bool
		CanManageWorkflowSessions(context.Context, *types.Workflow) bool
	}
)

func Session(log *zap.Logger) *session {
	return &session{
		log:        log,
		actionlog:  DefaultActionlog,
		store:      DefaultStore,
		ac:         DefaultAccessControl,
		mux:        &sync.RWMutex{},
		pool:       make(map[uint64]*types.Session),
		spawnQueue: make(chan *spawn),
	}
}

func (svc *session) Find(ctx context.Context, filter types.SessionFilter) (rr types.SessionSet, f types.SessionFilter, err error) {
	var (
		wap = &sessionActionProps{filter: &filter}
	)

	err = func() (err error) {
		if !svc.ac.CanSearchSessions(ctx) {
			return SessionErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchAutomationSessions(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, filter, svc.recordAction(ctx, wap, SessionActionSearch, err)
}

func (svc *session) FindByID(ctx context.Context, sessionID uint64) (res *types.Session, err error) {
	var (
		wap = &sessionActionProps{session: &types.Session{ID: sessionID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if res, err = loadSession(ctx, s, sessionID); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, wap, SessionActionLookup, err)
}

func (svc *session) resumeAll(ctx context.Context) error {
	return nil
}

func (svc *session) suspendAll(ctx context.Context) error {
	return nil
}

// Start new workflow session on a specific step with a given identity and scope
//func (svc *session) Start(g *wfexec.Graph, stepID uint64, i auth.Identifiable, input types.Variables, keepFor int, trace bool) error {
func (svc *session) Start(g *wfexec.Graph, i auth.Identifiable, ssp types.SessionStartParams) error {
	var (
		ctx   = auth.SetIdentityToContext(context.Background(), i)
		ses   = svc.spawn(g, ssp.Trace)
		start wfexec.Step
	)

	if ssp.StepID == 0 {
		// starting step is not explicitly workflows on trigger, find orphan step
		switch oo := g.Orphans(); len(oo) {
		case 1:
			start = oo[0]
		case 0:
			return fmt.Errorf("could not find step without parents")
		default:
			return fmt.Errorf("multiple steps without parents")
		}
	} else if start = g.GetStepByIdentifier(ssp.StepID); start == nil {
		return fmt.Errorf("trigger staring step references nonexisting step")
	}

	ses.CreatedAt = *now()
	ses.CreatedBy = i.Identity()
	ses.Apply(ssp)

	if err := store.CreateAutomationSession(context.TODO(), svc.store, ses); err != nil {
		return err
	}

	return ses.Exec(ctx, start, ssp.Input)
}

// Resume resumes suspended session/state
//
// Session can only be resumed by knowing session and state ID. Resume is an asynchronous operation
func (svc *session) Resume(sessionID, stateID uint64, i auth.Identifiable, input types.Variables) error {
	var (
		ctx = auth.SetIdentityToContext(context.Background(), i)
	)

	defer svc.mux.RUnlock()
	svc.mux.RLock()
	ses := svc.pool[sessionID]
	if ses == nil {
		return errors.NotFound("session not found")
	}

	return ses.Resume(ctx, stateID, input)
}

// spawns a new session
//
// We need initial context for the session because we want to catch all cancellations or timeouts from there
// and not from any potential HTTP requests or similar temporary context that can prematurely destroy a workflow session
func (svc *session) spawn(g *wfexec.Graph, trace bool) (ses *types.Session) {
	s := &spawn{make(chan *wfexec.Session, 1), g, trace}

	// Send new-session request
	svc.spawnQueue <- s

	// blocks until session is set
	ses = types.NewSession(<-s.session)

	svc.mux.Lock()
	svc.pool[ses.ID] = ses
	svc.mux.Unlock()
	return ses
}

func (svc *session) Watch(ctx context.Context) {
	go func() {
		defer sentry.Recover()
		defer svc.log.Info("stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case s := <-svc.spawnQueue:
				//
				s.session <- wfexec.NewSession(ctx, s.graph, wfexec.SetHandler(svc.stateChangeHandler(ctx)))
				// case time for a pool cleanup
				// @todo cleanup pool when sessions are complete
			}
		}

		// @todo serialize sessions & suspended states
		//svc.suspendAll(ctx)
	}()

	svc.log.Debug("watcher initialized")
}

func (svc *session) stateChangeHandler(ctx context.Context) wfexec.StateChangeHandler {
	return func(i int, state *wfexec.State, s *wfexec.Session) {
		log := svc.log.With(zap.Uint64("sessionID", s.ID()))

		defer svc.mux.RUnlock()
		svc.mux.RLock()
		ses := svc.pool[s.ID()]
		if ses == nil {
			log.Warn("could not find session to update")
			return
		}

		var err error

		rq := state.MakeRequest()

		// @todo collect all info and finalize the step
		ses.Trace = append(ses.Trace, &types.SessionTraceStep{
			ID:         nextID(),
			CallerStep: 0,
			StateID:    rq.StateID,
			CallerID:   0,
			StepID:     0,
			Depth:      0,
			Scope:      nil,
			Duration:   0,
		})

		switch i {
		case wfexec.SessionStepSuspended:
			// @todo handle step suspension!
		case wfexec.SessionSuspended:
			ses.SuspendedAt = now()
			ses.Status = types.SessionSuspended

		case wfexec.SessionCompleted:
			ses.SuspendedAt = nil
			ses.CompletedAt = now()
			ses.Status = types.SessionCompleted

		case wfexec.SessionFailed:
			ses.SuspendedAt = nil
			ses.CompletedAt = now()
			ses.Error = state.Error()
			ses.Status = types.SessionFailed

		default:
			return
		}

		if err = svc.store.UpdateAutomationSession(ctx, ses); err != nil {
			log.Error("failed to update session", zap.Error(err))
		} else {
			log.Debug("session updated", zap.Stringer("status", ses.Status))
		}
	}
}

func loadSession(ctx context.Context, s store.Storer, sessionID uint64) (res *types.Session, err error) {
	if sessionID == 0 {
		return nil, SessionErrInvalidID()
	}

	if res, err = store.LookupAutomationSessionByID(ctx, s, sessionID); errors.IsNotFound(err) {
		return nil, SessionErrNotFound()
	}

	return
}