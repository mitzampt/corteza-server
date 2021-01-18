package expr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVars_UnmarshalVar(t *testing.T) {
	var (
		req = require.New(t)

		dst = &struct {
			Int        int
			String     string `var:"STRING"`
			RawString  string `var:"rawString"`
			Bool       bool
			Unexisting byte
		}{}

		vars = Vars{
			"int":       NewInteger(42),
			"STRING":    NewString("foo"),
			"rawString": "raw",
			"bool":      NewBoolean(true),
			"missing":   NewBoolean(true),
		}
	)

	req.NoError(vars.Decode(dst))
	req.Equal(42, dst.Int)
	req.Equal("foo", dst.String)
	req.Equal("raw", dst.RawString)
	req.Equal(true, dst.Bool)
	req.Empty(dst.Unexisting)
}

func TestVars_Set(t *testing.T) {
	var (
		req = require.New(t)

		vars = Vars{
			"rawInt": 42,
			"int":    NewInteger(42),
			"sub": Vars{
				"foo": "foo",
			},
		}
	)

	req.NoError(vars.Set(123, "rawInt"))
	req.Equal(123, vars["rawInt"])

	req.NoError(vars.Set(123, "int"))
	req.Equal(123, vars["int"].(Var).Get().(int))

	req.NoError(vars.Set("bar", "sub", "foo"))
	req.Equal("bar", vars["sub"].(Vars)["foo"].(string))
}

func TestKV_Set(t *testing.T) {
	var (
		req = require.New(t)

		vars = &KV{v: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	req.NoError(vars.Set("v11", "k1"))
	req.Equal("v11", vars.v["k1"])

}
