package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	// Used for expression steps, arguments/results mapping and for input validation
	Expr struct {
		// Expression to evaluate over the input variables; results will be set to scope under variable Name
		Expr string `json:"expr,omitempty"`

		eval expr.Evaluable

		// Variable name to copy results of the expression to
		Name string `json:"name"`

		// Expected type of the input value
		Type string `json:"type,omitempty"`

		// Set of tests that can be run before input is evaluated and result copied to scope
		Tests TestSet `json:"tests,omitempty"`
	}

	ExprSet []*Expr

	// WorkflowStepExpression is created from WorkflowStep with kind=expressions
	expressionsStep struct {
		identifiableStep
		Set ExprSet
	}
)

func NewExpr(name, typ, expr string) (e *Expr, err error) {
	return &Expr{Expr: expr, Name: name, Type: typ}, nil
}

func (e Expr) GetExpr() string              { return e.Expr }
func (e *Expr) SetEval(eval expr.Evaluable) { e.eval = eval }
func (e Expr) Eval(ctx context.Context, scope expr.Vars) (interface{}, error) {
	return e.eval.Eval(ctx, scope)
}
func (e Expr) Test(ctx context.Context, scope expr.Vars) (bool, error) {
	return e.eval.Test(ctx, scope)
}

func (set ExprSet) Validate(ctx context.Context, in expr.Vars) (TestSet, error) {
	var (
		out TestSet
		vv  TestSet
		err error

		// Copy/create scope
		scope = expr.Vars.Merge(in)
	)

	for _, e := range set {
		vv, err = e.Tests.Validate(ctx, expr.Vars(scope))
		if err != nil {
			return nil, err
		}

		out = append(out, vv...)
	}

	return out, nil
}

func (set ExprSet) Eval(ctx context.Context, in expr.Vars) (expr.Vars, error) {
	var (
		err error

		// Copy/create scope
		scope = expr.Vars.Merge(in)
		out   = expr.Vars{}
		tmp   interface{}
	)

	for _, e := range set {
		tmp, err = e.eval.Eval(ctx, scope)
		if err != nil {
			return nil, err
		}

		// @todo figure out what type do we need to use from e.Type
		new := expr.NewAny(tmp)

		if err = scope.Set(new, expr.Path(e.Name)...); err != nil {
			return nil, err
		}

		if err = out.Set(new, expr.Path(e.Name)...); err != nil {
			return nil, err
		}

	}

	return out, nil
}

func ExpressionsStep(ee ...*Expr) *expressionsStep {
	return &expressionsStep{Set: ee}
}

func (s *expressionsStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	result, err := s.Set.Eval(ctx, r.Scope.Merge(r.Input))
	if err != nil {
		return nil, err
	}

	return result, nil
}
