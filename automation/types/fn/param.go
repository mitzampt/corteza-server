package fn

type (
	Param struct {
		Name     string      `json:"name"`
		Types    []paramType `json:"type"`
		Required bool        `json:"required"`
		Array    bool        `json:"array"`
		Meta     ParamMeta   `json:"meta"`
	}

	ParamMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
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
func Array(p *Param)    { p.Required = !p.Required }

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
