package wfexec

import (
	"context"
	"time"
)

type (
	ExecResponse interface{}
	partial      struct{}

	message struct {
		ID   uint64
		Kind string
		Body string
		msg  *Expression
	}

	messageEmitter struct {
		stepIdentifier

		Kind string
		Body string
		msg  *Expression
	}

	prompt struct {
		stepIdentifier

		Kind string

		// list of input variables that need to be set
		// see prompt's Exec fn for more details
		Required []string
	}
)

func NewPartial() *partial {
	return &partial{}
}

func DelayExecution(until time.Time) *suspended {
	return &suspended{resumeAt: &until}
}

func WaitForInput() *suspended {
	return &suspended{input: true}
}

func NewMessageEmitter(kind string, msg *Expression) *messageEmitter {
	return &messageEmitter{
		Kind: kind,
		msg:  msg,
	}
}

func (e *messageEmitter) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	var (
		m   = &message{}
		err error
	)

	if m.Body, err = e.msg.eval.EvalString(ctx, r.Scope); err != nil {
		return nil, err
	}

	return m, nil
}

func NewPrompt(kind string, rr ...string) *prompt {
	return &prompt{
		Kind:     kind,
		Required: rr,
	}
}

// Executes prompt step
//
// @todo This is as basic as it gets; we need a more advanced approach
//       either by defining (required) variable types or by validation of values from the scope
//       .
//       Should this be solved in the implementation (automation/service pkg) or
//       here by providing new struct(s) for testing scope?
//
// Idea:
//  - automation/types.Argument could be extended to support slice of validation structs
//    that provide text expression and error message for failed tests
//
func (p *prompt) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if len(p.Required) > 0 && !r.Input.Has(p.Required[0], p.Required[1:]...) {
		return WaitForInput(), nil
	}

	return r.Input, nil
}
