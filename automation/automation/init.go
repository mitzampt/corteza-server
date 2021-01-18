package automation

import atypes "github.com/cortezaproject/corteza-server/automation/types"

func RegisterFunctions(reg func(*atypes.Function)) {
	httpRequest.register(reg)
}
