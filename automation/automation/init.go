package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
)

const (
	baseRef = "base"
)

func List() []*types.Function {
	return httpSenders()
}
