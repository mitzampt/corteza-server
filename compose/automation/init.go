package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

func RegisterFunctions(reg func(*types.Function)) {
	//namespaces.Register(reg)
	//modules.Register(reg)
	records.register(reg)
}

func RegisterTypes(reg func(p expr.Var)) {
	reg(&Namespace{})
	reg(&Module{})
	reg(&Record{})
}
