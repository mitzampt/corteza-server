package automation

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/spf13/cast"
	"strings"
)

type Record struct{ v *types.Record }

func (Record) New(new *types.Record) expr.Var { return &Record{v: new} }
func (Record) Type() string                   { return "ComposeRecord" }
func (Record) Is(v interface{}) bool          { _, is := v.(types.Record); return is }
func (v Record) Get() interface{}             { return v.v }
func (v *Record) Set(new interface{}, pp ...string) (err error) {
	if len(pp) == 0 {
		var ok bool
		v.v, ok = new.(*types.Record)
		if !ok {
			return fmt.Errorf("unable to cast type %T to types.Record", v)
		}
	} else {
		return setRecordWithPath(v.v, new, pp...)
	}
	return nil
}

func setRecordValueWithPath(rec *types.Record, val interface{}, pp ...string) (err error) {
	if len(pp) == 0 {
		if rvs, is := val.(types.RecordValueSet); !is {
			return fmt.Errorf("expecting type RecordValueSet")
		} else {
			rec.Values = rvs
			return nil
		}
	}

	if len(pp) > 2 {
		return fmt.Errorf("invalid path for record value: %q", strings.Join(pp, "."))
	}

	rv := &types.RecordValue{Name: pp[0]}
	if rv.Value, err = cast.ToStringE(val); err != nil {
		return
	}

	if len(pp) == 2 {
		if rv.Place, err = cast.ToUintE(val); err != nil {
			return
		}
	}

	rec.Values = rec.Values.Set(rv)
	return nil
}

func setRecordWithPath(rec *types.Record, val interface{}, pp ...string) (err error) {
	switch pp[0] {
	case "ID":
		return expr.SetIDWithPath(&rec.ID, val, pp[1:]...)

	case "moduleID":
		return expr.SetIDWithPath(&rec.ModuleID, val, pp[1:]...)

	case "namespaceID":
		return expr.SetIDWithPath(&rec.NamespaceID, val, pp[1:]...)

	case "values":
		return setRecordValueWithPath(rec, val, pp[1:]...)

	case "labels":
		return expr.SetKVWithPath(&rec.Labels, val, pp[1:]...)

	case "ownedBy":
		return expr.SetIDWithPath(&rec.OwnedBy, val, pp[1:]...)

	case "createdAt":
		return expr.SetTimeWithPath(&rec.CreatedAt, val, pp[1:]...)

	case "createdBy":
		return expr.SetIDWithPath(&rec.CreatedBy, val, pp[1:]...)

	case "updatedAt":
		return expr.SetTimeWithPath(rec.UpdatedAt, val, pp[1:]...)

	case "updatedBy":
		return expr.SetIDWithPath(&rec.UpdatedBy, val, pp[1:]...)

	case "deletedAt":
		return expr.SetTimeWithPath(rec.DeletedAt, val, pp[1:]...)

	case "deletedBy":
		return expr.SetIDWithPath(&rec.DeletedBy, val, pp[1:]...)

	default:
		return fmt.Errorf("unknown record field %q", pp[0])

	}

}

type Module struct{ v *types.Module }

func (Module) New(new *types.Module) expr.Var { return &Module{v: new} }
func (Module) Type() string                   { return "ComposeModule" }
func (Module) Is(v interface{}) bool          { _, is := v.(types.Module); return is }
func (v Module) Get() interface{}             { return v.v }
func (v *Module) Set(new interface{}, pp ...string) (err error) {
	// @todo implement setting via path
	var ok bool
	v.v, ok = new.(*types.Module)
	if !ok {
		return fmt.Errorf("unable to cast type %T to types.Module", v)
	}
	return nil
}

type Namespace struct{ v *types.Namespace }

func (Namespace) New(new *types.Namespace) expr.Var { return &Namespace{v: new} }
func (Namespace) Type() string                      { return "ComposeNamespace" }
func (Namespace) Is(v interface{}) bool             { _, is := v.(types.Namespace); return is }
func (v Namespace) Get() interface{}                { return v.v }
func (v *Namespace) Set(new interface{}, pp ...string) (err error) {
	// @todo implement setting via path
	var ok bool
	v.v, ok = new.(*types.Namespace)
	if !ok {
		return fmt.Errorf("unable to cast type %T to types.Namespace", v)
	}
	return nil
}
