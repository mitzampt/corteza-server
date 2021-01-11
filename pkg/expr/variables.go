package expr

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

type (
	// Variables uses same structure as expr.Variables
	// but implements scan/value to support serialization into store
	Variables map[string]interface{}
)

func ParseVariables(ss []string) (p Variables, err error) {
	p = Variables{}
	if len(ss) == 0 {
		return
	}

	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (vv *Variables) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = Variables{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into Variables: %w", string(b), err)
		}
	}

	return nil
}

func (vv Variables) Value() (driver.Value, error) {
	return json.Marshal(vv)
}

// Assign takes base variables and assigns all new variables
func (vv Variables) Merge(nn ...Variables) Variables {
	var (
		out = Variables{}
	)

	nn = append([]Variables{vv}, nn...)
	for i := range nn {
		for k, v := range nn[i] {
			out[k] = v
		}
	}

	return out
}

// Returns true if all keys are present
func (vv Variables) Has(key string, kk ...string) bool {
	for _, key = range append([]string{key}, kk...) {
		if _, has := vv[key]; !has {
			return false
		}
	}

	return true
}

// Returns true if all keys are present
func (vv Variables) Any(key string, kk ...string) bool {
	for _, key = range append([]string{key}, kk...) {
		if _, has := vv[key]; has {
			return true
		}
	}

	return false
}

func (vv Variables) String(key string, def ...string) string {
	if v, has := vv[key]; has {
		if o, err := cast.ToStringE(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return ""
}

func (vv Variables) Bool(key string, def ...bool) bool {
	if v, has := vv[key]; has {
		if o, err := cast.ToBoolE(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return false
}

func (vv Variables) Int(key string, def ...int) int {
	if v, has := vv[key]; has {
		if o, err := cast.ToIntE(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (vv Variables) Int64(key string, def ...int64) int64 {
	if v, has := vv[key]; has {
		if o, err := cast.ToInt64E(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (vv Variables) Uint64(key string, def ...uint64) uint64 {
	if v, has := vv[key]; has {
		if o, err := cast.ToUint64E(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (vv Variables) Float64(key string, def ...float64) float64 {
	if v, has := vv[key]; has {
		if o, err := cast.ToFloat64E(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (vv Variables) Duration(key string, def ...time.Duration) time.Duration {
	if v, has := vv[key]; has {
		if o, err := cast.ToDurationE(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (vv Variables) Time(key string, def ...time.Time) time.Time {
	if v, has := vv[key]; has {
		if o, err := cast.ToTimeE(v); err == nil {
			return o
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return time.Time{}
}
