package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
)

type (
	Type struct {
	}

	typeSetPayload struct {
		Set []string `json:"set"`
	}
)

func (Type) New() *Type {
	ctrl := &Type{}
	return ctrl
}

func (ctrl Type) List(ctx context.Context, r *request.TypeList) (interface{}, error) {
	return typeSetPayload{Set: service.RegisteredTypes()}, nil
}
