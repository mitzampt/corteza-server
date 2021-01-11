package wfexec

import (
	"time"
)

type (
	ExecResponse interface{}
	partial      struct{}
)

func NewPartial() *partial {
	return &partial{}
}

func DelayExecution(until time.Time) *suspended {
	return &suspended{resumeAt: &until}
}

func WaitForInput() *suspended {
	return &suspended{input: true}
}
