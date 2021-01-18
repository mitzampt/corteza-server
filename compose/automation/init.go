package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
)

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
		//hNamespace.Restore(),

		//hModule.LookupByID(),
		//hModule.Save(),
		//hModule.Create(),
		//hModule.Update(),
		//hModule.Delete(),
		//hModule.Restore(),

		hRecords.LookupByID(),
		hRecords.Save(),
		hRecords.Create(),
		hRecords.Update(),
		hRecords.Delete(),
		hRecords.Restore(),

		//hPage.LookupByID(),
		//hPage.Save(),
		//hPage.Create(),
		//hPage.Update(),
		//hPage.Delete(),
		//hPage.Restore(),

		//hChart.LookupByID(),
		//hChart.Save(),
		//hChart.Create(),
		//hChart.Update(),
		//hChart.Delete(),
		//hChart.Restore(),
	}
}
