package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	Function struct {
	}

	functionSetPayload struct {
		Set []*types.Function `json:"set"`
	}
)

func (Function) New() *Function {
	ctrl := &Function{}
	return ctrl
}

func (ctrl Function) List(ctx context.Context, r *request.FunctionList) (interface{}, error) {
	return functionSetPayload{Set: service.RegisteredFunctions()}, nil
}
