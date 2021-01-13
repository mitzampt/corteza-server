package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
)

func List() []*types.Function {
	return httpSenders()
}
