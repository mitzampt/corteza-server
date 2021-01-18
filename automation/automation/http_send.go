package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	stdHttpSendParameters = []*Param{
		NewParam("url", Types(&expr.String{}), Required),
		NewParam("header", Types(&expr.KV{})),
		NewParam("headerAuthBearer", Types(&expr.String{})),
		NewParam("headerAuthUsername", Types(&expr.String{})),
		NewParam("headerAuthPassword", Types(&expr.String{})),
		NewParam("headerContentType", Types(&expr.String{})),
		NewParam("timeout", Types(&expr.Duration{})),
	}

	stdHttpPayloadParameters = []*Param{
		NewParam("form", Types(&expr.KVV{})),
		NewParam("body", Types(&expr.String{}, &expr.Reader{}, &expr.Any{})),
	}

	stdHttpSendResults = []*Param{
		NewParam("status", Types(&expr.String{})),
		NewParam("statusCode", Types(&expr.Integer{})),
		NewParam("header", Types(&expr.KV{})),
		NewParam("contentLength", Types(&expr.Integer{})),
		NewParam("body", Types(&expr.String{})),
	}
)

func makeHttpRequest(ctx context.Context, in expr.Vars) (req *http.Request, err error) {
	var (
		body   io.Reader
		header = make(http.Header)
		args   = struct {
			Method             string
			Form               url.Values
			Body               interface{}
			Timeout            time.Duration
			Url                string
			Header             http.Header
			HeaderUserAgent    string
			HeaderAuthBearer   string
			HeaderAuthUsername string
			HeaderAuthPassword string
			HeaderContentType  string
		}{
			HeaderUserAgent: "Corteza-Automation-Client/" + version.Version,
		}
	)

	if err = in.Decode(&args); err != nil {
		return
	}

	args.Method = strings.ToUpper(args.Method)

	if args.Method == "" && (len(args.Form) > 0 || args.Body != nil) {
		// when no method is set and form or body are passed
		args.Method = http.MethodPost
	}

	err = func() error {
		switch args.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
		default:
			return nil
		}

		// @todo handle (multiple) file upload as well

		if len(args.Form) > 0 {
			if args.Body != nil {
				return fmt.Errorf("can not not use form and body parameters at the same time")
			}

			//var form url.Values
			//form, err = cast.ToStringMapStringSliceE(in["form"])
			//if err != nil {
			//	return fmt.Errorf("failed to resolve form values: %w", err)
			//}

			header.Add("Content-Type", "application/x-www-form-urlencoded")
			payload := &bytes.Buffer{}
			if _, err = payload.WriteString(args.Form.Encode()); err != nil {
				return err
			}

			body = payload
			return nil
		}

		if args.Body != nil {
			return nil
		}

		switch payload := args.Body.(type) {
		case string:
			body = strings.NewReader(payload)
		case []byte:
			body = bytes.NewReader(payload)
		case io.Reader:
			body = payload
		default:
			aux := &bytes.Buffer{}
			payload = aux
			return json.NewEncoder(aux).Encode(args.Body)
		}

		return nil
	}()
	if err != nil {
		return nil, err
	}

	if args.Timeout > 0 {
		var tfn context.CancelFunc
		ctx, tfn = context.WithTimeout(ctx, args.Timeout)
		defer tfn()
	}

	req, err = http.NewRequestWithContext(ctx, args.Method, args.Url, body)
	if err != nil {
		return nil, err
	}

	header.Set("User-Agent", args.HeaderUserAgent)

	if len(args.Header) > 0 {
		for k, vv := range args.Header {
			for _, v := range vv {
				header.Add(k, v)
			}
		}
	}

	switch {
	case len(args.HeaderAuthBearer) > 0:
		header.Add("Authorization", "Bearer "+args.HeaderAuthBearer)
	case len(args.HeaderAuthPassword+args.HeaderAuthPassword) > 0:
		req.SetBasicAuth(
			args.HeaderAuthPassword,
			args.HeaderAuthPassword,
		)
	}

	if len(args.HeaderContentType) > 0 {
		header.Add("Content-Type", args.HeaderContentType)
	}

	req.Header = header

	return
}

func httpSend(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	req, err = makeHttpRequest(ctx, in)
	if err != nil {
		return nil, err
	}

	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	out = expr.Vars{
		"status":        rsp.Status,
		"statusCode":    rsp.StatusCode,
		"header":        rsp.Header,
		"contentLength": rsp.ContentLength,
		"body":          rsp.Body,
	}

	return
}

func httpSenders() []*Function {
	return []*Function{
		httpSendRequest(),
		httpSendRequestGet(),
		httpSendRequestPost(),
		httpSendRequestPut(),
		httpSendRequestPatch(),
		httpSendRequestDelete(),
	}
}

func httpSendRequest() *Function {
	return &Function{
		Ref: "httpSendRequest",
		//Meta: &FunctionMeta{},
		Parameters: append(append(
			[]*Param{NewParam("method", Types(&expr.String{}), Required)},
			stdHttpSendParameters...),
			stdHttpPayloadParameters...,
		),
		Results: stdHttpSendResults,
		Handler: httpSend,
	}
}

func httpSendRequestGet() *Function {
	return &Function{
		Ref: "httpSendRequestGet",
		//Meta:       &FunctionMeta{},
		Parameters: stdHttpSendParameters,
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodGet}))
		},
	}
}

func httpSendRequestPost() *Function {
	return &Function{
		Ref: "httpSendRequestPost",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPost}))
		},
	}
}

func httpSendRequestPut() *Function {
	return &Function{
		Ref: "httpSendRequestPut",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPut}))
		},
	}
}

func httpSendRequestPatch() *Function {
	return &Function{
		Ref: "httpSendRequestPatch",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPatch}))
		},
	}
}

func httpSendRequestDelete() *Function {
	return &Function{
		Ref: "httpSendRequestDelete",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodDelete}))
		},
	}
}
