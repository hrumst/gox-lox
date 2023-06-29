package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/scan"
)

type RuntimeError struct {
	message string
	token   *scan.Token
}

func NewRuntimeError(message string, token *scan.Token) *RuntimeError {
	return &RuntimeError{
		message: message,
		token:   token,
	}
}

func ConvertToRuntimeError(message string, err error, token *scan.Token) *RuntimeError {
	return NewRuntimeError(
		fmt.Sprintf("%s: %s", message, err.Error()),
		token,
	)
}

func (re *RuntimeError) Error() string {
	errStr := re.message
	if re.token != nil {
		errStr = errStr + fmt.Sprintf("\nat line: %d, token: %s", re.token.Line, re.token.Lexeme)
	}
	return errStr
}
