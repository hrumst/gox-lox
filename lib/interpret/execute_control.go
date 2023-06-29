package interpret

import "github.com/hrumst/gox-lox/lib/scan"

const (
	executeBreakControlType = iota
	executeContinueControlType
	executeReturnControlType
)

type executeControl struct {
	controlType int
	value       *scan.LoxValue
}

func newReturnExecuteControl(value *scan.LoxValue) executeControl {
	return executeControl{
		controlType: executeReturnControlType,
		value:       value,
	}
}

func newContinueExecuteControl() executeControl {
	return executeControl{
		controlType: executeContinueControlType,
	}
}

func newBreakExecuteControl() executeControl {
	return executeControl{
		controlType: executeBreakControlType,
	}
}

func (r executeControl) isBreak() bool {
	return r.controlType == executeBreakControlType
}

func (r executeControl) isContinue() bool {
	return r.controlType == executeContinueControlType
}

func (r executeControl) isReturn() bool {
	return r.controlType == executeReturnControlType
}
