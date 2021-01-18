package automation

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	. "github.com/cortezaproject/corteza-server/automation/types"
//	"github.com/cortezaproject/corteza-server/pkg/expr"
//	"github.com/cortezaproject/corteza-server/pkg/version"
//	"io"
//	"net/http"
//	"net/url"
//	"strings"
//	"time"
//)
//
//var (
//	stdHttpSendParameters = []*Param{
//		NewParam("url", Types(&expr.String{}), Required),
//		NewParam("header", Types(&expr.KV{})),
//		NewParam("headerAuthBearer", Types(&expr.String{})),
//		NewParam("headerAuthUsername", Types(&expr.String{})),
//		NewParam("headerAuthPassword", Types(&expr.String{})),
//		NewParam("headerContentType", Types(&expr.String{})),
//		NewParam("timeout", Types(&expr.Duration{})),
//	}
//
//	stdHttpPayloadParameters = []*Param{
//		NewParam("form", Types(&expr.KVV{})),
//		NewParam("body", Types(&expr.String{}, &expr.Reader{}, &expr.Any{})),
//	}
//
//	stdHttpSendResults = []*Param{
//		NewParam("status", Types(&expr.String{})),
//		NewParam("statusCode", Types(&expr.Integer{})),
//		NewParam("header", Types(&expr.KV{})),
//		NewParam("contentLength", Types(&expr.Integer{})),
//		NewParam("body", Types(&expr.String{})),
//	}
//)
//
//func httpSend(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//	var (
//		req *http.Request
//		rsp *http.Response
//	)
//
//	req, err = makeHttpRequest(ctx, in)
//	if err != nil {
//		return nil, err
//	}
//
//	rsp, err = http.DefaultClient.Do(req)
//	if err != nil {
//		return
//	}
//
//	out = expr.Vars{
//		"status":        rsp.Status,
//		"statusCode":    rsp.StatusCode,
//		"header":        rsp.Header,
//		"contentLength": rsp.ContentLength,
//		"body":          rsp.Body,
//	}
//
//	return
//}
//
//func httpSenders() []*Function {
//	return []*Function{
//		httpSendRequest(),
//		httpSendRequestGet(),
//		httpSendRequestPost(),
//		httpSendRequestPut(),
//		httpSendRequestPatch(),
//		httpSendRequestDelete(),
//	}
//}
//
//func httpSendRequest() *Function {
//	return &Function{
//		Ref: "httpSendRequest",
//		//Meta: &FunctionMeta{},
//		Parameters: append(append(
//			[]*Param{NewParam("method", Types(&expr.String{}), Required)},
//			stdHttpSendParameters...),
//			stdHttpPayloadParameters...,
//		),
//		Results: stdHttpSendResults,
//		Handler: httpSend,
//	}
//}
//
//func httpSendRequestGet() *Function {
//	return &Function{
//		Ref: "httpSendRequestGet",
//		//Meta:       &FunctionMeta{},
//		Parameters: stdHttpSendParameters,
//		Results:    stdHttpSendResults,
//		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodGet}))
//		},
//	}
//}
//
//func httpSendRequestPost() *Function {
//	return &Function{
//		Ref: "httpSendRequestPost",
//		//Meta:       &FunctionMeta{},
//		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
//		Results:    stdHttpSendResults,
//		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPost}))
//		},
//	}
//}
//
//func httpSendRequestPut() *Function {
//	return &Function{
//		Ref: "httpSendRequestPut",
//		//Meta:       &FunctionMeta{},
//		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
//		Results:    stdHttpSendResults,
//		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPut}))
//		},
//	}
//}
//
//func httpSendRequestPatch() *Function {
//	return &Function{
//		Ref: "httpSendRequestPatch",
//		//Meta:       &FunctionMeta{},
//		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
//		Results:    stdHttpSendResults,
//		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodPatch}))
//		},
//	}
//}
//
//func httpSendRequestDelete() *Function {
//	return &Function{
//		Ref: "httpSendRequestDelete",
//		//Meta:       &FunctionMeta{},
//		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
//		Results:    stdHttpSendResults,
//		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
//			return httpSend(ctx, in.Merge(expr.Vars{"method": http.MethodDelete}))
//		},
//	}
//}
