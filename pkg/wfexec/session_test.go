package wfexec

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"testing"
	"time"
)

type (
	sesTestStep struct {
		stepIdentifier
		name string
		exec func(context.Context, *ExecRequest) (ExecResponse, error)
	}

	sesTestTemporal struct {
		stepIdentifier
		delay time.Duration
		until time.Time
	}
)

var (
	// used for testing to produce lower numbers that are easier to inspect and compare
	testID = atomic.NewUint64(0)
)

func (s *sesTestStep) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if s.exec != nil {
		return s.exec(ctx, r)
	}

	return expr.Variables{
		"counter": r.Scope.Int("counter", 0) + 1,
		"path":    r.Scope.String("path", "") + "/" + s.name,
		s.name:    "executed",
	}, nil
}

func (s *sesTestTemporal) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if s.until.IsZero() {
		s.until = now().Add(s.delay)
	}

	if now().Before(s.until) {
		return DelayExecution(s.until), nil
	}

	return expr.Variables{
		"waitForMoment": "executed",
	}, nil
}

func TestSession_TwoStepWorkflow(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf)

		s1 = &sesTestStep{name: "s1"}
		s2 = &sesTestStep{name: "s2"}
	)

	wf.AddStep(s1, s2) // 1st execute s1 then s2
	ses.Exec(ctx, s1, expr.Variables{"two": 1, "three": 1})
	ses.Wait(ctx)
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	req.Equal("/s1/s2", ses.Result().String("path", ""))
}

func TestSession_SplitAndMerge(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf)

		start  = &sesTestStep{name: "start"}
		split1 = &sesTestStep{name: "split1"}
		split2 = &sesTestStep{name: "split2"}
		split3 = &sesTestStep{name: "split3"}

		end = JoinGateway(split1, split2, split3)
	)

	wf.AddStep(start, split1, split2, split3)
	wf.AddStep(split1, end)
	wf.AddStep(split2, end)
	wf.AddStep(split3, end)
	ses.Exec(ctx, start, nil)
	ses.Wait(ctx)
	req.True(ses.Idle())
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	// split3 only!
	req.Equal("/start/split3", ses.Result().String("path", ""))
	req.Contains(ses.Result(), "split1")
	req.Contains(ses.Result(), "split2")
	req.Contains(ses.Result(), "split3")
}

func TestSession_Delays(t *testing.T) {
	var (
		// how fast we want to go (lower = faster)
		//
		unit  = time.Millisecond
		delay = unit * 3

		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf,
			// for testing we need much shorter worker intervals
			SetWorkerInterval(unit),
		)

		start = &sesTestStep{name: "start"}

		waitForMoment = &sesTestTemporal{delay: delay}

		waitForInputStateId atomic.Uint64
		waitForInput        = &sesTestStep{name: "waitForInput", exec: func(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
			if !r.Input.Has("input") {
				waitForInputStateId.Store(r.StateID)
				return WaitForInput(), nil
			}

			return expr.Variables{
				"input":        r.Input["input"],
				"waitForInput": "executed",
			}, nil
		}}
	)

	//ses.log = logger.MakeDebugLogger().Sugar()

	ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()

	wf.AddStep(start, waitForMoment)
	wf.AddStep(waitForMoment, waitForInput)

	req.NoError(ses.Exec(ctx, start, nil))

	// wait-for-moment step needs to be executed before we can resume wait-for-input
	ses.Wait(ctx)
	time.Sleep(delay + unit)
	req.NotZero(waitForInputStateId.Load())

	// should not be completed yet...
	req.True(ses.Idle())
	req.True(ses.Suspended())

	// push in the input
	req.NoError(ses.Resume(ctx, waitForInputStateId.Load(), expr.Variables{"input": "foo"}))

	req.False(ses.Suspended())
	ses.Wait(ctx)
	time.Sleep(2 * unit)

	// should not be completed yet...
	req.True(ses.Idle())
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	req.Contains(ses.Result(), "waitForMoment")
	req.Contains(ses.Result(), "waitForInput")
	req.Equal("foo", ses.Result().String("input", ""))
}

func bmSessionSimpleStepSequence(c uint64, b *testing.B) {
	var (
		ctx = context.Background()
		g   = NewGraph()
		err error
	)

	for i := uint64(1); i <= c; i++ {
		s := &sesTestStep{name: "start"}
		s.SetID(i)
		g.AddStep(s)
		if i > 1 {
			g.AddParent(s, g.StepByID(i-1))
		}
	}

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ses := NewSession(ctx, g)
		if err = ses.Exec(ctx, g.StepByID(1), nil); err != nil {
			b.Fatal(err.Error())
		}

		ses.Wait(ctx)
	}
	b.StopTimer()
}

func BenchmarkSessionSimple1StepSequence(b *testing.B)       { bmSessionSimpleStepSequence(1, b) }
func BenchmarkSessionSimple10StepSequence(b *testing.B)      { bmSessionSimpleStepSequence(10, b) }
func BenchmarkSessionSimple100StepSequence(b *testing.B)     { bmSessionSimpleStepSequence(100, b) }
func BenchmarkSessionSimple1000StepSequence(b *testing.B)    { bmSessionSimpleStepSequence(1000, b) }
func BenchmarkSessionSimple10000StepSequence(b *testing.B)   { bmSessionSimpleStepSequence(10000, b) }
func BenchmarkSessionSimple100000StepSequence(b *testing.B)  { bmSessionSimpleStepSequence(100000, b) }
func BenchmarkSessionSimple1000000StepSequence(b *testing.B) { bmSessionSimpleStepSequence(1000000, b) }
