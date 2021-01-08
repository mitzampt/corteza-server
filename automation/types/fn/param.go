package fn

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	ParamSet []*Param
	Param    struct {
		Name     string      `json:"name,omitempty"`
		Types    []paramType `json:"type,omitempty"`
		Required bool        `json:"required,omitempty"`
		Array    bool        `json:"array,omitempty"`
		Meta     *ParamMeta  `json:"meta,omitempty"`
	}

	ParamMeta struct {
		Label       string                 `json:"label,omitempty"`
		Description string                 `json:"description,omitempty"`
		Visual      map[string]interface{} `json:"visual,omitempty"`
	}

	paramType string
	paramOpt  func(p *Param)
)

const (
	TypeAny      paramType = "any"
	TypeString   paramType = "string"
	TypeInt      paramType = "int"
	TypeUint64   paramType = "uint64"
	TypeFloat64  paramType = "float64"
	TypeBool     paramType = "bool"
	TypeDuration paramType = "duration"
	TypeTime     paramType = "time"
	TypeReader   paramType = "reader"
	TypeKV       paramType = "kv"
)

func NewParam(name string, opts ...paramOpt) *Param {
	p := &Param{Name: name}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func Required(p *Param) { p.Required = !p.Required }
func Array(p *Param)    { p.Array = !p.Array }

func Any(p *Param)      { p.Types = append(p.Types, TypeAny) }
func String(p *Param)   { p.Types = append(p.Types, TypeString) }
func Int(p *Param)      { p.Types = append(p.Types, TypeInt) }
func Uint64(p *Param)   { p.Types = append(p.Types, TypeUint64) }
func Float64(p *Param)  { p.Types = append(p.Types, TypeFloat64) }
func Bool(p *Param)     { p.Types = append(p.Types, TypeBool) }
func Duration(p *Param) { p.Types = append(p.Types, TypeDuration) }
func Time(p *Param)     { p.Types = append(p.Types, TypeTime) }
func Reader(p *Param)   { p.Types = append(p.Types, TypeReader) }
func KV(p *Param)       { p.Types = append(p.Types, TypeKV) }

// CheckInput validates (at compile-time) input data (arguments)
func (set ParamSet) CheckInput(ee *wfexec.Expressions) error {
	for _, p := range set {
		if p.Required && !ee.Has(p.Name) {
			return fmt.Errorf("%q is required", p.Name)
		}
	}

	return nil
}

// CheckOutput validates (at compile-time) output data (arguments)
func (set ParamSet) CheckOutput(ee *wfexec.Expressions) error {
	var ind = make(map[string]bool)
	for _, p := range set {
		ind[p.Name] = true
	}

	for _, name := range ee.Names() {
		if !ind[name] {
			return fmt.Errorf("unknown result %q used", name)
		}
	}

	return nil
}
