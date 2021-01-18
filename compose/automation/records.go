package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	recordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)

		Create(ctx context.Context, record *types.Record) (*types.Record, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordID ...uint64) error
	}

	moduleService interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*types.Module, error)
	}

	namespaceService interface {
		FindByID(ctx context.Context, namespaceID uint64) (*types.Namespace, error)
	}

	recordsHandler struct {
		ns  namespaceService
		mod moduleService
		rec recordService
	}
)

func RecordsHandler(ns namespaceService, mod moduleService, rec recordService) *recordsHandler {
	return &recordsHandler{
		ns:  ns,
		mod: mod,
		rec: rec,
	}
}

func (h recordsHandler) lookupByID(ctx context.Context, args *recordsLookupByIDArgs) (results *recordsLookupByIDResults, err error) {
	results = &recordsLookupByIDResults{}
	results.Record, err = h.rec.FindByID(ctx, args.namespaceID, args.moduleID, args.RecordID)
	return
}

func (h recordsHandler) create(ctx context.Context, args *recordsCreateArgs) (results *recordsCreateResults, err error) {
	return
}

func (h recordsHandler) save(ctx context.Context, args *recordsSaveArgs) (results *recordsSaveResults, err error) {
	return
}

func (h recordsHandler) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
	return
}

func (h recordsHandler) delete(ctx context.Context, args *recordsDeleteArgs) (results *recordsDeleteResults, err error) {
	return
}

func (h recordsHandler) restore(ctx context.Context, args *recordsRestoreArgs) (results *recordsRestoreResults, err error) {
	return
}
