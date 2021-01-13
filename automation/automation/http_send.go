package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	. "github.com/cortezaproject/corteza-server/automation/types/fn"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	stdHttpSendParameters = []*Param{
		NewParam("url", String, Required),
		NewParam("header", KV),
		NewParam("headerAuthBearer", String),
		NewParam("headerAuthUsername", String),
		NewParam("headerAuthPassword", String),
		NewParam("headerContentType", String),
		NewParam("timeout", Duration),
	}

	stdHttpPayloadParameters = []*Param{
		NewParam("form", KV),
		NewParam("body", String, Reader, Any),
	}

	stdHttpSendResults = []*Param{
		NewParam("status", String),
		NewParam("statusCode", Int),
		NewParam("header", KV),
		NewParam("contentLength", Int),
		NewParam("body", String),
	}
)

func makeHttpRequest(ctx context.Context, in expr.Variables) (req *http.Request, err error) {
	var (
		body   io.Reader
		header = make(http.Header)
		method = strings.ToUpper(in.String("method"))
	)

	if method == "" && in.Any("form", "body") {
		// when no method is set and form or body are passed
		method = http.MethodPost
	}

	err = func() error {
		switch method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
		default:
			return nil
		}

		// @todo handle (multiple) file upload as well

		if in.Has("form") {
			if in.Has("body") {
				return fmt.Errorf("can not not use form and body parameters at the same time")
			}

			var form url.Values
			form, err = cast.ToStringMapStringSliceE(in["form"])
			if err != nil {
				return fmt.Errorf("failed to resolve form values: %w", err)
			}

			header.Add("Content-Type", "application/x-www-form-urlencoded")
			payload := &bytes.Buffer{}
			if _, err = payload.WriteString(form.Encode()); err != nil {
				return err
			}

			body = payload
			return nil
		}

		if !in.Has("body") {
			return nil
		}

		switch payload := in["body"].(type) {
		case string:
			body = strings.NewReader(payload)
		case []byte:
			body = bytes.NewReader(payload)
		case io.Reader:
			body = payload
		default:
			aux := &bytes.Buffer{}
			payload = aux
			return json.NewEncoder(aux).Encode(in["body"])
		}

		return nil
	}()
	if err != nil {
		return nil, err
	}

	if in.Has("timeout") {
		var tfn context.CancelFunc
		ctx, tfn = context.WithTimeout(ctx, in.Duration("timeout"))
		defer tfn()
	}

	req, err = http.NewRequestWithContext(ctx, method, in.String("url"), body)
	if err != nil {
		return nil, err
	}

	header.Set("User-Agent", in.String("headerUserAgent", "Corteza-Automation-Client/"+version.Version))

	if in.Has("header") {
		for k, v := range in["header"].(map[string]string) {
			header.Add(k, v)
		}
	}

	switch {
	case in.Has("headerAuthBearer"):
		header.Add("Authorization", "Bearer "+in.String("headerAuthBearer"))
	case in.Any("headerAuthUsername", "headerAuthPassword"):
		req.SetBasicAuth(
			in.String("headerAuthUsername"),
			in.String("headerAuthPassword"),
		)
	}

	if in.Has("headerContentType") {
		header.Add("Content-Type", in.String("headerContentType"))
	}

	req.Header = header

	return
}

func httpSend(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
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

	out = expr.Variables{
		"status":        rsp.Status,
		"statusCode":    rsp.StatusCode,
		"header":        rsp.Header,
		"contentLength": rsp.ContentLength,
		"body":          rsp.Body,
	}

	return
}

func httpSenders() []*types.Function {
	return []*types.Function{
		httpSendRequest(),
		httpSendGetRequest(),
		httpSendPostRequest(),
		httpSendPutRequest(),
		httpSendPatchRequest(),
		httpSendDeleteRequest(),
	}
}

func httpSendRequest() *types.Function {
	return &types.Function{
		Ref: "httpRequest",
		//Meta: &FunctionMeta{},
		Parameters: append(append(
			[]*Param{NewParam("method", String, Required)},
			stdHttpSendParameters...),
			stdHttpPayloadParameters...,
		),
		Results: stdHttpSendResults,
		Handler: httpSend,
	}
}

func httpSendGetRequest() *types.Function {
	return &types.Function{
		Ref: "httpGetRequest",
		//Meta:       &FunctionMeta{},
		Parameters: stdHttpSendParameters,
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
			return httpSend(ctx, in.Merge(expr.Variables{"method": http.MethodGet}))
		},
	}
}

func httpSendPostRequest() *types.Function {
	return &types.Function{
		Ref: "httpPostRequest",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
			return httpSend(ctx, in.Merge(expr.Variables{"method": http.MethodPost}))
		},
	}
}

func httpSendPutRequest() *types.Function {
	return &types.Function{
		Ref: "httpPutRequest",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
			return httpSend(ctx, in.Merge(expr.Variables{"method": http.MethodPut}))
		},
	}
}

func httpSendPatchRequest() *types.Function {
	return &types.Function{
		Ref: "httpPatchRequest",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
			return httpSend(ctx, in.Merge(expr.Variables{"method": http.MethodPatch}))
		},
	}
}

func httpSendDeleteRequest() *types.Function {
	return &types.Function{
		Ref: "httpDeleteRequest",
		//Meta:       &FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in expr.Variables) (out expr.Variables, err error) {
			return httpSend(ctx, in.Merge(expr.Variables{"method": http.MethodDelete}))
		},
	}
}
