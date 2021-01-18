package expr

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"reflect"
	"strings"
	"time"
)

type (
	Var interface {
		Setter
		Type() string
		Is(interface{}) bool
		Get() interface{}
	}

	Setter interface {
		Set(interface{}, ...string) error
	}

	Decoder interface {
		Decode(reflect.Value) error
	}

	Dict interface {
		Dict() map[string]interface{}
	}
)

func ReqNoPath(t string, pp []string) error {
	if len(pp) > 0 {
		return fmt.Errorf("setting values with path on %s is not supported", t)
	}

	return nil
}

type Vars map[string]interface{}

func (Vars) Type() string          { return "variables" }
func (Vars) Is(v interface{}) bool { _, is := v.(Vars); return is }
func (vv Vars) Get() interface{}   { return vv }
func (vv *Vars) Set(new interface{}, pp ...string) error {
	if len(pp) > 0 {
		p := pp[0]
		if _, has := (*vv)[p]; !has {
			if len(pp) > 1 {
				return fmt.Errorf("%q does not exist, can not set value with path %q", p, strings.Join(pp[1:], "."))
			}
			(*vv)[p] = new
		}

		if sub, is := (*vv)[p].(Vars); is {
			return sub.Set(new, pp[1:]...)
		} else if sub, is := (*vv)[p].(Setter); is {
			return sub.Set(new, pp[1:]...)
		} else {
			if len(pp) > 1 {
				return fmt.Errorf("can not set value with path %q on %T", strings.Join(pp[1:], "."), (*vv)[p])
			} else {
				(*vv)[p] = new
			}
		}

		return nil
	}

	switch casted := new.(type) {
	case Vars:
		*vv = casted
		return nil
	default:
		return fmt.Errorf("unable to cast type %T to Vars", vv)
	}
}

// Assign takes base variables and assigns all new variables
func (vv Vars) Merge(nn ...Vars) Vars {
	var (
		out = Vars{}
	)

	nn = append([]Vars{vv}, nn...)
	for i := range nn {
		for k, v := range nn[i] {
			out[k] = v
		}
	}

	return out
}

// Returns true if all keys are present
func (vv Vars) Has(key string, kk ...string) bool {
	for _, key = range append([]string{key}, kk...) {
		if _, has := vv[key]; !has {
			return false
		}
	}

	return true
}

// Returns true if all keys are present
func (vv Vars) Any(key string, kk ...string) bool {
	for _, key = range append([]string{key}, kk...) {
		if _, has := vv[key]; has {
			return true
		}
	}

	return false
}

func (vv Vars) Dict() map[string]interface{} {
	dict := make(map[string]interface{})
	for k, v := range vv {
		switch v := v.(type) {
		case Dict:
			dict[k] = v.Dict()

		case Var:
			dict[k] = v.Get()

		default:
			dict[k] = v
		}

	}

	return dict
}

func (vv Vars) Decode(dst interface{}) (err error) {
	dstRef := reflect.ValueOf(dst)

	if dstRef.Kind() != reflect.Ptr {
		return fmt.Errorf("expecting a pointer, not a value")
	}

	if dstRef.IsNil() {
		return fmt.Errorf("nil pointer passed")
	}

	dstRef = dstRef.Elem()

	length := dstRef.NumField()
	for i := 0; i < length; i++ {
		t := dstRef.Type().Field(i)

		keyName := t.Tag.Get("var")
		if keyName == "" {
			keyName = strings.ToLower(t.Name[:1]) + t.Name[1:]
		}

		value, has := vv[keyName]
		if !has {
			continue
		} else if um, is := value.(Decoder); is {
			if err = um.Decode(dstRef.Field(i)); err != nil {
				return
			}
		} else if get, is := value.(Var); is {
			dstRef.Field(i).Set(reflect.ValueOf(get.Get()))
		} else {
			dstRef.Field(i).Set(reflect.ValueOf(value))
		}
	}

	return
}

func (vv *Vars) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = Vars{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into Variables: %w", string(b), err)
		}
	}

	return nil
}

func (vv Vars) Value() (driver.Value, error) {
	return json.Marshal(vv)
}

type Any struct{ v interface{} }

func NewAny(new interface{}) Var    { return &Any{v: new} }
func (Any) New(new interface{}) Var { return &Any{v: new} }
func (Any) Type() string            { return "any" }
func (Any) Is(v interface{}) bool   { return true }
func (v Any) Get() interface{}      { return v.v }
func (v *Any) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}
	v.v = new
	return nil
}

type Boolean struct{ v bool }

func NewBoolean(new bool) Var         { return &Boolean{v: new} }
func (Boolean) New(new bool) Var      { return &Boolean{v: new} }
func (Boolean) Type() string          { return "boolean" }
func (Boolean) Is(v interface{}) bool { _, is := v.(bool); return is }
func (v Boolean) Get() interface{}    { return v.v }
func (v *Boolean) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToBoolE(new)
	return err
}
func (v Boolean) Decode(x reflect.Value) error { x.SetBool(v.v); return nil }

type ID struct{ v uint64 }

func NewID(new uint64) Var       { return &ID{v: new} }
func (ID) New(new uint64) Var    { return &ID{v: new} }
func (ID) Type() string          { return "id" }
func (ID) Is(v interface{}) bool { _, is := v.(uint64); return is }
func (v ID) Get() interface{}    { return v.v }
func (v *ID) Set(new interface{}, pp ...string) (err error) {
	return SetIDWithPath(&v.v, new, pp...)
}

func (v ID) Decode(x reflect.Value) error { x.SetUint(v.v); return nil }

type Integer struct{ v int }

func NewInteger(new int) Var          { return &Integer{v: new} }
func (Integer) New(new int) Var       { return &Integer{v: new} }
func (Integer) Type() string          { return "integer" }
func (Integer) Is(v interface{}) bool { _, is := v.(int); return is }
func (v Integer) Get() interface{}    { return v.v }
func (v *Integer) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToIntE(new)
	return err
}
func (v Integer) Decode(x reflect.Value) error { x.SetInt(int64(v.v)); return nil }

type Integer64 struct{ v int64 }

func NewInteger64(new int64) Var        { return &Integer64{v: new} }
func (Integer64) New(new int64) Var     { return &Integer64{v: new} }
func (Integer64) Type() string          { return "integer64" }
func (Integer64) Is(v interface{}) bool { _, is := v.(int64); return is }
func (v Integer64) Get() interface{}    { return v.v }
func (v *Integer64) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToInt64E(new)
	return err
}
func (v Integer64) Decode(x reflect.Value) error { x.SetInt(v.v); return nil }

type Unsigned struct{ v uint }

func NewUnsigned(new uint) Var         { return &Unsigned{v: new} }
func (Unsigned) New(new uint) Var      { return &Unsigned{v: new} }
func (Unsigned) Type() string          { return "unsigned" }
func (Unsigned) Is(v interface{}) bool { _, is := v.(uint); return is }
func (v Unsigned) Get() interface{}    { return v.v }
func (v *Unsigned) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToUintE(v)
	return err
}
func (v Unsigned) Decode(x reflect.Value) error { x.SetUint(uint64(v.v)); return nil }

type Float struct{ v float64 }

func NewFloat(new float64) Var      { return &Float{v: new} }
func (Float) New(new float64) Var   { return &Float{v: new} }
func (Float) Type() string          { return "float" }
func (Float) Is(v interface{}) bool { _, is := v.(float64); return is }
func (v Float) Get() interface{}    { return v.v }
func (v *Float) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToFloat64E(new)
	return err
}
func (v Float) Decode(x reflect.Value) error { x.SetFloat(v.v); return nil }

type String struct{ v string }

func NewString(new string) Var       { return &String{v: new} }
func (String) New(new string) Var    { return &String{v: new} }
func (String) Type() string          { return "string" }
func (String) Is(v interface{}) bool { _, is := v.(string); return is }
func (v String) Get() interface{}    { return v.v }
func (v *String) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToStringE(new)
	return err
}
func (v String) Decode(x reflect.Value) error { x.SetString(v.v); return nil }

type Datetime struct{ v time.Time }

func NewDatetime(new time.Time) Var    { return &Datetime{v: new} }
func (Datetime) New(new time.Time) Var { return &Datetime{v: new} }
func (Datetime) Type() string          { return "datetime" }
func (Datetime) Is(v interface{}) bool { _, is := v.(time.Time); return is }
func (v Datetime) Get() interface{}    { return v.v }
func (v *Datetime) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToTimeE(new)
	return err
}

type Duration struct{ v time.Duration }

func NewDuration(new time.Duration) Var    { return &Duration{v: new} }
func (Duration) New(new time.Duration) Var { return &Duration{v: new} }
func (Duration) Type() string              { return "duration" }
func (Duration) Is(v interface{}) bool     { _, is := v.(time.Duration); return is }
func (v Duration) Get() interface{}        { return v.v }
func (v *Duration) Set(new interface{}, pp ...string) (err error) {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	v.v, err = cast.ToDurationE(new)
	return err
}

type KV struct{ v map[string]string }

func NewKV(new map[string]string) Var    { return &KV{v: new} }
func (KV) New(new map[string]string) Var { return &KV{v: new} }
func (KV) Type() string                  { return "keyValue" }
func (KV) Is(v interface{}) bool         { _, err := cast.ToStringMapStringE(v); return err != nil }
func (v KV) Get() interface{}            { return v.v }
func (v *KV) Set(new interface{}, pp ...string) error {
	return SetKVWithPath(&v.v, new, pp...)
}

type KVV struct{ v map[string][]string }

func NewKVV(new map[string][]string) Var    { return &KVV{v: new} }
func (KVV) New(new map[string][]string) Var { return &KVV{v: new} }
func (KVV) Type() string                    { return "keyValues" }
func (KVV) Is(v interface{}) bool           { _, err := cast.ToStringMapStringSliceE(v); return err != nil }
func (v KVV) Get() interface{}              { return v.v }
func (v *KVV) Set(new interface{}, pp ...string) error {
	switch len(pp) {
	case 0:
		tmp, err := cast.ToStringMapStringSliceE(new)
		if err != nil {
			return err
		}

		v.v = tmp
	case 1:
		tmp, err := cast.ToStringSliceE(new)
		if err != nil {
			return err
		}

		v.v[pp[0]] = tmp
	default:
		return fmt.Errorf("can not set values to KVV with path deeper than 1 level")
	}

	return nil
}

type Reader struct{ v io.Reader }

func NewReader(new io.Reader) Var    { return &Reader{v: new} }
func (Reader) New(new io.Reader) Var { return &Reader{v: new} }
func (Reader) Type() string          { return "reader" }
func (Reader) Is(v interface{}) bool { _, is := v.(io.Reader); return is }
func (v Reader) Get() interface{}    { return v.v }
func (v *Reader) Set(new interface{}, pp ...string) error {
	if err := ReqNoPath(v.Type(), pp); err != nil {
		return err
	}

	var ok bool
	v.v, ok = new.(io.Reader)
	if !ok {
		return fmt.Errorf("unable to cast type %T to io.Reader", v)
	}

	return nil
}
