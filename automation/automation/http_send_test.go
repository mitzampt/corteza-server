package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpRequestMaker(t *testing.T) {
	validateBody := func(r *require.Assertions, req *http.Request, expected string) {
		reader, err := req.GetBody()
		r.NoError(err)
		body, err := ioutil.ReadAll(reader)
		r.NoError(err)

		r.Equal(expected, string(body))

	}

	t.Run("basic get", func(t *testing.T) {
		var (
			r  = require.New(t)
			in = expr.Vars{
				"url": expr.NewString("http://localhost/test"),
			}
			req, err = makeHttpRequest(context.Background(), in)
		)

		r.NoError(err)
		r.Equal("GET", req.Method)
		r.Equal("http://localhost/test", req.URL.String())
	})
	//
	//j := func(in expr.Vars) expr.Vars {
	//	j, err := json.Marshal(in)
	//	if err != nil {
	//		panic(err)
	//	}
	//	out := expr.Vars{}
	//	err = json.Unmarshal(j, &out)
	//	if err != nil {
	//		panic(err)
	//	}
	//	return out
	//}

	t.Run("post form", func(t *testing.T) {
		var (
			r  = require.New(t)
			in = expr.Vars{
				"form": expr.NewKVV(map[string][]string{
					"a": []string{"a"},
					"b": []string{"b", "b"},
					"i": []string{"42"},
				}),
			}
			req, err = makeHttpRequest(context.Background(), in)
		)

		r.NoError(err)
		r.Equal("POST", req.Method)
		validateBody(r, req, "a=a&b=b&b=b&i=42")
	})
}
