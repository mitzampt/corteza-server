package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types/fn"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	FunctionHandler func(ctx context.Context, in expr.Variables) (expr.Variables, error)

	// workflow functions are defined in the core code and through plugins
	Function struct {
		Ref        string        `json:"ref,omitempty"`
		Meta       *FunctionMeta `json:"meta,omitempty"`
		Parameters fn.ParamSet   `json:"parameters,omitempty"`
		Results    fn.ParamSet   `json:"results,omitempty"`

		Handler FunctionHandler `json:"-"`
	}

	FunctionMeta struct {
		Name        string                 `json:"name,omitempty"`
		Description string                 `json:"description,omitempty"`
		Visual      map[string]interface{} `json:"visual,omitempty"`
	}

	functionStep struct {
		identifiableStep
		def       *Function
		arguments ExprSet
		results   ExprSet
	}
)

func FunctionStep(def *Function, arguments, results ExprSet) (f *functionStep, err error) {
	f = &functionStep{def: def, arguments: arguments, results: results}
	return
}

func (f *functionStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	var (
		args, results expr.Variables
		err           error
	)

	if len(f.arguments) > 0 {
		// Arguments defined, get values from scope and use them when calling
		// function/handler
		args, err = f.arguments.Eval(ctx, expr.Variables(r.Scope.Merge(r.Input)))
		if err != nil {
			return nil, err
		}
	}

	results, err = f.def.Handler(ctx, args)
	if err != nil {
		return nil, err
	}

	if len(f.results) == 0 {
		// No results defined, nothing to return
		return expr.Variables{}, nil
	}

	results, err = f.results.Eval(ctx, results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
