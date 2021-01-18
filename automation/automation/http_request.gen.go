package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/http_request.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	httpRequest = &httpRequestHandler{}
)

func (h httpRequestHandler) register(reg func(*atypes.Function)) {
	reg(h.Send())
}

type (
	httpRequestSendArgs struct {
		hasUrl bool

		Url string

		hasMethod bool

		Method string

		hasParams bool

		Params url.Values

		hasHeaders bool

		Headers http.Header

		hasHeaderAuthBearer bool

		HeaderAuthBearer string

		hasHeaderAuthUsername bool

		HeaderAuthUsername string

		hasHeaderAuthPassword bool

		HeaderAuthPassword string

		hasHeaderUserAgent bool

		HeaderUserAgent string

		hasHeaderContentType bool

		HeaderContentType string

		hasTimeout bool

		Timeout time.Duration

		hasForm bool

		Form url.Values

		hasBody bool

		Body       interface{}
		bodyString string
		bodyStream io.Reader
		bodyRaw    interface{}
	}

	httpRequestSendResults struct {
		Status        string
		StatusCode    int
		Headers       http.Header
		ContentLength int64
		ContentType   string
		Body          io.Reader
	}
)

//
//
// expects implementation of send function:
// func (h httpRequest) send(ctx context.Context, args *httpRequestSendArgs) (results *httpRequestSendResults, err error) {
//    return
// }
func (h httpRequestHandler) Send() *atypes.Function {
	return &atypes.Function{
		Ref: "httpRequestSend",
		Meta: &atypes.FunctionMeta{
			Short: "Sends HTTP request",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "url",
				Types: []string{(expr.String{}).Type()}, Required: true,
			},
			{
				Name:  "method",
				Types: []string{(expr.String{}).Type()}, Required: true,
			},
			{
				Name:  "params",
				Types: []string{(expr.KVV{}).Type()},
			},
			{
				Name:  "headers",
				Types: []string{(expr.KVV{}).Type()},
			},
			{
				Name:  "headerAuthBearer",
				Types: []string{(expr.String{}).Type()},
			},
			{
				Name:  "headerAuthUsername",
				Types: []string{(expr.String{}).Type()},
			},
			{
				Name:  "headerAuthPassword",
				Types: []string{(expr.String{}).Type()},
			},
			{
				Name:  "headerUserAgent",
				Types: []string{(expr.String{}).Type()},
			},
			{
				Name:  "headerContentType",
				Types: []string{(expr.String{}).Type()},
			},
			{
				Name:  "timeout",
				Types: []string{(expr.Duration{}).Type()},
			},
			{
				Name:  "form",
				Types: []string{(expr.KVV{}).Type()},
			},
			{
				Name:  "body",
				Types: []string{(expr.String{}).Type(), (expr.Reader{}).Type(), (expr.Any{}).Type()},
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("status",
				atypes.Types(&expr.String{}),
			),

			atypes.NewParam("statusCode",
				atypes.Types(&expr.Integer{}),
			),

			atypes.NewParam("headers",
				atypes.Types(&expr.KVV{}),
			),

			atypes.NewParam("contentLength",
				atypes.Types(&expr.Integer64{}),
			),

			atypes.NewParam("contentType",
				atypes.Types(&expr.String{}),
			),

			atypes.NewParam("body",
				atypes.Types(&expr.Reader{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &httpRequestSendArgs{
					hasUrl:                in.Has("url"),
					hasMethod:             in.Has("method"),
					hasParams:             in.Has("params"),
					hasHeaders:            in.Has("headers"),
					hasHeaderAuthBearer:   in.Has("headerAuthBearer"),
					hasHeaderAuthUsername: in.Has("headerAuthUsername"),
					hasHeaderAuthPassword: in.Has("headerAuthPassword"),
					hasHeaderUserAgent:    in.Has("headerUserAgent"),
					hasHeaderContentType:  in.Has("headerContentType"),
					hasTimeout:            in.Has("timeout"),
					hasForm:               in.Has("form"),
					hasBody:               in.Has("body"),
				}

				results *httpRequestSendResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Body to go type
			switch casted := args.Body.(type) {
			case string:
				args.bodyString = casted
			case io.Reader:
				args.bodyStream = casted
			case interface{}:
				args.bodyRaw = casted
			}

			if results, err = h.send(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"status": (expr.String{}).New(results.Status),

				"statusCode": (expr.Integer{}).New(results.StatusCode),

				"headers": (expr.KVV{}).New(results.Headers),

				"contentLength": (expr.Integer64{}).New(results.ContentLength),

				"contentType": (expr.String{}).New(results.ContentType),

				"body": (expr.Reader{}).New(results.Body),
			}

			return
		},
	}
}
