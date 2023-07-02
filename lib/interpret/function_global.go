package interpret

import (
	"github.com/hrumst/gox-lox/lib/scan"
	"time"
)

type ClockFunction struct {
}

func (c ClockFunction) String() string {
	return "[function] clock"
}

func NewClockFunction() *ClockFunction {
	return &ClockFunction{}
}

func (c ClockFunction) Arity() int {
	return 0
}

func (c ClockFunction) Call(args []*scan.LoxValue) (*scan.LoxValue, error) {
	timeMs := time.Now().Unix()
	return scan.NewFloatLoxValue(float64(timeMs)), nil
}
