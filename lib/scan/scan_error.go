package scan

import "fmt"

type ScanError struct {
	line  int
	where string
	err   error
}

func NewScanError(line int, where string, err error) *ScanError {
	return &ScanError{
		line:  line,
		where: where,
		err:   err,
	}
}

func (se *ScanError) Error() string {
	// ctx for HAD_ERROR ???
	return fmt.Sprintf("[Line %d] Error %s: %s", se.line, se.where, se.err.Error())
}
