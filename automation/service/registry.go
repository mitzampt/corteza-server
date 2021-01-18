package service

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"sort"
)

var (
	// workflow function registry
	functionRegistry = make(map[string]*types.Function)
	typeRegistry     = make(map[string]expr.Var)
)

func FunctionRegistrator(fn *types.Function) {
	functionRegistry[fn.Ref] = fn
}

func RegisteredFunctions() []*types.Function {
	var (
		rr = make([]string, 0, len(functionRegistry))
		ff = make([]*types.Function, 0, len(functionRegistry))
	)

	for ref := range functionRegistry {
		rr = append(rr, ref)
	}

	sort.Strings(rr)

	for _, ref := range rr {
		ff = append(ff, functionRegistry[ref])
	}

	return ff
}

func TypeRegistrator(t expr.Var) {
	typeRegistry[t.Type()] = t
}

func RegisteredTypes() []string {
	var (
		rr = make([]string, 0, len(typeRegistry))
	)

	for ref := range typeRegistry {
		rr = append(rr, ref)
	}

	sort.Strings(rr)

	return rr
}

func registerCoreTypes() {
	TypeRegistrator(&expr.Any{})
	TypeRegistrator(&expr.Boolean{})
	TypeRegistrator(&expr.ID{})
	TypeRegistrator(&expr.Integer{})
	TypeRegistrator(&expr.Integer64{})
	TypeRegistrator(&expr.Unsigned{})
	TypeRegistrator(&expr.Float64{})
	TypeRegistrator(&expr.String{})
	TypeRegistrator(&expr.Datetime{})
	TypeRegistrator(&expr.Duration{})
	TypeRegistrator(&expr.KV{})
	TypeRegistrator(&expr.KVV{})
	TypeRegistrator(&expr.Reader{})
}
