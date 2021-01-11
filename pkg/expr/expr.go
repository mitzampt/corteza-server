package expr

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
)

type (
	Parser__ interface {
		Parse(string) (Evaluable, error)
		ParseEvaluators(ee ...Evaluator) error
	}

	Evaluable interface {
		Eval(context.Context, Variables) (interface{}, error)
		Test(context.Context, Variables) (bool, error)
	}

	Evaluator interface {
		GetExpr() string
		SetEval(eval Evaluable)
	}

	gvalParser struct {
		lang gval.Language
	}

	gvalEval struct {
		evaluable gval.Evaluable
	}
)

func NewParser() Parser__ {
	return NewGvalParser()
}

func NewGvalParser() *gvalParser {
	return &gvalParser{lang: gval.Full(AllFunctions()...)}
}

func (p *gvalParser) Parse(expr string) (Evaluable, error) {
	var (
		ge  = &gvalEval{}
		err error
	)

	if ge.evaluable, err = p.lang.NewEvaluable(expr); err != nil {
		return nil, err
	}

	return ge, err
}

func (p *gvalParser) ParseEvaluators(ee ...Evaluator) error {
	for _, e := range ee {
		evaluable, err := p.Parse(e.GetExpr())
		if err != nil {
			return err
		}

		e.SetEval(evaluable)
	}

	return nil
}

func (e *gvalEval) Eval(ctx context.Context, scope Variables) (interface{}, error) {
	return e.evaluable(ctx, scope)
}

func (e *gvalEval) Test(ctx context.Context, scope Variables) (bool, error) {
	return e.evaluable.EvalBool(ctx, scope)
}

func Parser(ll ...gval.Language) gval.Language {
	return gval.Full(append(AllFunctions(), ll...)...)
}

func AllFunctions() []gval.Language {
	ff := make([]gval.Language, 0, 100)

	//ff = append(ff, GenericFunctions()...)
	ff = append(ff, StringFunctions()...)
	ff = append(ff, NumericFunctions()...)
	ff = append(ff, TimeFunctions()...)

	return ff
}

// utility function for examples
func eval(e string, p interface{}) {
	result, err := Parser().Evaluate(e, p)
	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		fmt.Printf("%v", result)
	}
}
