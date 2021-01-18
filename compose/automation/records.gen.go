package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/records.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"time"
)

type (
	recordsLookupByIDArgs struct {
		hasRecordID bool

		RecordID uint64

		hasModule bool

		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleFull   *types.Module

		hasNamespace bool

		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceFull   *types.Namespace
	}

	recordsLookupByIDResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of lookupByID function:
// func (h records) lookupByID(ctx context.Context, args *recordsLookupByIDArgs) (results *recordsLookupByIDResults, err error) {
//    return
// }
func (h recordsHandler) LookupByID() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsLookupByID",
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recordID",
				Types: []string{(expr.ID{}).Type()}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Module{}).Type()}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Namespace{}).Type()}, Required: true,
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsLookupByIDArgs{
					hasRecordID:  in.Has("recordID"),
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
				}

				results *recordsLookupByIDResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Module to go type
			switch casted := args.Module.(type) {
			case uint64:
				args.moduleID = casted
			case string:
				args.moduleHandle = casted
			case *types.Module:
				args.moduleFull = casted
			}

			// Converting Namespace to go type
			switch casted := args.Namespace.(type) {
			case uint64:
				args.namespaceID = casted
			case string:
				args.namespaceHandle = casted
			case *types.Namespace:
				args.namespaceFull = casted
			}

			if results, err = h.lookupByID(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}

type (
	recordsSaveArgs struct {
		hasRecord bool
	}

	recordsSaveResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of save function:
// func (h records) save(ctx context.Context, args *recordsSaveArgs) (results *recordsSaveResults, err error) {
//    return
// }
func (h recordsHandler) Save() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsSave",
		Meta: &atypes.FunctionMeta{
			Short: "Save record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "record",
				Types: []string{},
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsSaveArgs{
					hasRecord: in.Has("record"),
				}

				results *recordsSaveResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			if results, err = h.save(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}

type (
	recordsCreateArgs struct {
		hasModule bool

		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleFull   *types.Module

		hasNamespace bool

		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceFull   *types.Namespace

		hasValues bool

		Values types.RecordValueSet

		hasLabels bool

		Labels label.Labels

		hasOwnedBy bool

		OwnedBy uint64

		hasCreatedBy bool

		CreatedBy uint64

		hasCreatedAt bool

		CreatedAt *time.Time

		hasUpdatedBy bool

		UpdatedBy uint64

		hasUpdatedAt bool

		UpdatedAt *time.Time

		hasDeletedBy bool

		DeletedBy uint64

		hasDeletedAt bool

		DeletedAt *time.Time
	}

	recordsCreateResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of create function:
// func (h records) create(ctx context.Context, args *recordsCreateArgs) (results *recordsCreateResults, err error) {
//    return
// }
func (h recordsHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsCreate",
		Meta: &atypes.FunctionMeta{
			Short: "Creates and stores a new record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Module{}).Type()}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Namespace{}).Type()}, Required: true,
			},
			{
				Name:  "values",
				Types: []string{(expr.KV{}).Type()},
			},
			{
				Name:  "labels",
				Types: []string{(expr.KV{}).Type()},
			},
			{
				Name:  "ownedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record owner",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "createdBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record creator",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "createdAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
			{
				Name:  "updatedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record updater",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "updatedAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
			{
				Name:  "deletedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record updater",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "deletedAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsCreateArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasValues:    in.Has("values"),
					hasLabels:    in.Has("labels"),
					hasOwnedBy:   in.Has("ownedBy"),
					hasCreatedBy: in.Has("createdBy"),
					hasCreatedAt: in.Has("createdAt"),
					hasUpdatedBy: in.Has("updatedBy"),
					hasUpdatedAt: in.Has("updatedAt"),
					hasDeletedBy: in.Has("deletedBy"),
					hasDeletedAt: in.Has("deletedAt"),
				}

				results *recordsCreateResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Module to go type
			switch casted := args.Module.(type) {
			case uint64:
				args.moduleID = casted
			case string:
				args.moduleHandle = casted
			case *types.Module:
				args.moduleFull = casted
			}

			// Converting Namespace to go type
			switch casted := args.Namespace.(type) {
			case uint64:
				args.namespaceID = casted
			case string:
				args.namespaceHandle = casted
			case *types.Namespace:
				args.namespaceFull = casted
			}

			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}

type (
	recordsUpdateArgs struct {
		hasModule bool

		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleFull   *types.Module

		hasNamespace bool

		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceFull   *types.Namespace

		hasValues bool

		Values types.RecordValueSet

		hasLabels bool

		Labels label.Labels

		hasOwnedBy bool

		OwnedBy uint64

		hasCreatedBy bool

		CreatedBy uint64

		hasCreatedAt bool

		CreatedAt *time.Time

		hasUpdatedBy bool

		UpdatedBy uint64

		hasUpdatedAt bool

		UpdatedAt *time.Time

		hasDeletedBy bool

		DeletedBy uint64

		hasDeletedAt bool

		DeletedAt *time.Time
	}

	recordsUpdateResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of update function:
// func (h records) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
//    return
// }
func (h recordsHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsUpdate",
		Meta: &atypes.FunctionMeta{
			Short: "Updates an existing record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Module{}).Type()}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Namespace{}).Type()}, Required: true,
			},
			{
				Name:  "values",
				Types: []string{(expr.KV{}).Type()},
			},
			{
				Name:  "labels",
				Types: []string{(expr.KV{}).Type()},
			},
			{
				Name:  "ownedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record owner",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "createdBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record creator",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "createdAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
			{
				Name:  "updatedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record updater",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "updatedAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
			{
				Name:  "deletedBy",
				Types: []string{(expr.ID{}).Type()},
				Meta: &atypes.ParamMeta{
					Label:  "Record updater",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
			{
				Name:  "deletedAt",
				Types: []string{(expr.Datetime{}).Type()},
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsUpdateArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasValues:    in.Has("values"),
					hasLabels:    in.Has("labels"),
					hasOwnedBy:   in.Has("ownedBy"),
					hasCreatedBy: in.Has("createdBy"),
					hasCreatedAt: in.Has("createdAt"),
					hasUpdatedBy: in.Has("updatedBy"),
					hasUpdatedAt: in.Has("updatedAt"),
					hasDeletedBy: in.Has("deletedBy"),
					hasDeletedAt: in.Has("deletedAt"),
				}

				results *recordsUpdateResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Module to go type
			switch casted := args.Module.(type) {
			case uint64:
				args.moduleID = casted
			case string:
				args.moduleHandle = casted
			case *types.Module:
				args.moduleFull = casted
			}

			// Converting Namespace to go type
			switch casted := args.Namespace.(type) {
			case uint64:
				args.namespaceID = casted
			case string:
				args.namespaceHandle = casted
			case *types.Namespace:
				args.namespaceFull = casted
			}

			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}

type (
	recordsDeleteArgs struct {
		hasRecordID bool

		RecordID uint64

		hasModule bool

		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleFull   *types.Module

		hasNamespace bool

		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceFull   *types.Namespace
	}

	recordsDeleteResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of delete function:
// func (h records) delete(ctx context.Context, args *recordsDeleteArgs) (results *recordsDeleteResults, err error) {
//    return
// }
func (h recordsHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsDelete",
		Meta: &atypes.FunctionMeta{
			Short: "Soft deletes compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recordID",
				Types: []string{(expr.ID{}).Type()}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Module{}).Type()}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Namespace{}).Type()}, Required: true,
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsDeleteArgs{
					hasRecordID:  in.Has("recordID"),
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
				}

				results *recordsDeleteResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Module to go type
			switch casted := args.Module.(type) {
			case uint64:
				args.moduleID = casted
			case string:
				args.moduleHandle = casted
			case *types.Module:
				args.moduleFull = casted
			}

			// Converting Namespace to go type
			switch casted := args.Namespace.(type) {
			case uint64:
				args.namespaceID = casted
			case string:
				args.namespaceHandle = casted
			case *types.Namespace:
				args.namespaceFull = casted
			}

			if results, err = h.delete(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}

type (
	recordsRestoreArgs struct {
		hasRecordID bool

		RecordID uint64

		hasModule bool

		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleFull   *types.Module

		hasNamespace bool

		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceFull   *types.Namespace
	}

	recordsRestoreResults struct {
		Record *types.Record
	}
)

//
//
// expects implementation of restore function:
// func (h records) restore(ctx context.Context, args *recordsRestoreArgs) (results *recordsRestoreResults, err error) {
//    return
// }
func (h recordsHandler) Restore() *atypes.Function {
	return &atypes.Function{
		Ref: "composeRecordsRestore",
		Meta: &atypes.FunctionMeta{
			Short: "Soft deletes compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recordID",
				Types: []string{(expr.ID{}).Type()}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Module{}).Type()}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{(expr.ID{}).Type(), (expr.String{}).Type(), (Namespace{}).Type()}, Required: true,
			},
		},

		Results: []*atypes.Param{

			atypes.NewParam("record",
				atypes.Types(&Record{}),
			),
		},

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &recordsRestoreArgs{
					hasRecordID:  in.Has("recordID"),
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
				}

				results *recordsRestoreResults
			)

			if err = in.Decode(&args); err != nil {
				return
			}

			// Converting Module to go type
			switch casted := args.Module.(type) {
			case uint64:
				args.moduleID = casted
			case string:
				args.moduleHandle = casted
			case *types.Module:
				args.moduleFull = casted
			}

			// Converting Namespace to go type
			switch casted := args.Namespace.(type) {
			case uint64:
				args.namespaceID = casted
			case string:
				args.namespaceHandle = casted
			case *types.Namespace:
				args.namespaceFull = casted
			}

			if results, err = h.restore(ctx, args); err != nil {
				return
			}

			out = expr.Vars{
				"record": (Record{}).New(results.Record),
			}

			return
		},
	}
}
