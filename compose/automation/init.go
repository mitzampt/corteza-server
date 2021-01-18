package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
)

func RegisterFunctions(reg func(*types.Function)) {
	//namespaces.Register(reg)
	//modules.Register(reg)
	records.register(reg)
}

func List(records recordService) []*types.Function {
	var (
		// hNamespace
		// hModule

		hRecords = &recordsHandler{
			ns:  nil, // @todo
			mod: nil, // @todo
			rec: records,
		}

		// hPage
		// hChart
	)

	return []*types.Function{
		//hNamespace.LookupByID(),
		//hNamespace.Save(),
		//hNamespace.Create(),
		//hNamespace.Update(),
		//hNamespace.Delete(),

		//hModule.LookupByID(),
		//hModule.Save(),
		//hModule.Create(),
		//hModule.Update(),
		//hModule.Delete(),

		hRecords.LookupByID(),
		hRecords.Save(),
		hRecords.Create(),
		hRecords.Update(),
		hRecords.Delete(),

		//hPage.LookupByID(),
		//hPage.Save(),
		//hPage.Create(),
		//hPage.Update(),
		//hPage.Delete(),

		//hChart.LookupByID(),
		//hChart.Save(),
		//hChart.Create(),
		//hChart.Update(),
		//hChart.Delete(),
	}
}
