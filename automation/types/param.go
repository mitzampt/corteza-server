package types

import (
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	ParamSet []*Param
	Param    struct {
		Name     string     `json:"name,omitempty"`
		Types    []string   `json:"types,omitempty"`
		Required bool       `json:"required,omitempty"`
		SetOf    bool       `json:"setOf,omitempty"`
		Meta     *ParamMeta `json:"meta,omitempty"`
	}

	ParamMeta struct {
		Label       string                 `json:"label,omitempty"`
		Description string                 `json:"description,omitempty"`
		Visual      map[string]interface{} `json:"visual,omitempty"`
	}

	paramOpt func(p *Param)
)

//const
func NewParam(name string, opts ...paramOpt) *Param {
	p := &Param{Name: name}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func Required(p *Param) { p.Required = !p.Required }
func SetOf(p *Param)    { p.SetOf = !p.SetOf }

func Types(tt ...expr.Var) paramOpt {
	return func(p *Param) {
		for _, t := range tt {
			p.Types = append(p.Types, t.Type())
		}
	}
}

//
//// CheckInput validates (at compile-time) input data (arguments)
//func (set ParamSet) CheckInput(ee *wfexec.Expressions) error {
//	for _, p := range set {
//		if p.Required && !ee.Has(p.Name) {
//			return fmt.Errorf("%q is required", p.Name)
//		}
//	}
//
//	return nil
//}
//
//// CheckOutput validates (at compile-time) output data (arguments)
//func (set ParamSet) CheckOutput(ee *wfexec.Expressions) error {
//	var ind = make(map[string]bool)
//	for _, p := range set {
//		ind[p.Name] = true
//	}
//
//	for _, name := range ee.Names() {
//		if !ind[name] {
//			return fmt.Errorf("unknown result %q used", name)
//		}
//	}
//
//	return nil
//}
