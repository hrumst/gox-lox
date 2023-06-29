package parse

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/scan"
)

type ParseError struct {
	token scan.Token
	err   error
}

func NewParseError(token scan.Token, err error) *ParseError {
	return &ParseError{
		token: token,
		err:   err,
	}
}

func (pe *ParseError) Error() string {
	if pe.token.Type == scan.EOF {
		return fmt.Sprintf(
			"%d at end %s",
			pe.token.Line,
			pe.err.Error(),
		)
	}
	return fmt.Sprintf(
		"%d at '%s'(%s). %s",
		pe.token.Line,
		pe.token.Lexeme,
		pe.token.Type,
		pe.err.Error(),
	)
}
