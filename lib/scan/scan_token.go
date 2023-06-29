package scan

import (
	"fmt"
	"strconv"
	"unicode"
)

func (sc *Scanner) ScanTokens() ([]Token, error) {
	for !sc.IsAtEnd() {
		sc.start = sc.current
		if err := sc.scanToken(); err != nil {
			return nil, err
		}
	}
	sc.tokens = append(sc.tokens, NewToken(EOF, "", nil, sc.line))
	return sc.tokens, nil
}

func (sc *Scanner) scanToken() error {
	sourceChar := sc.advance()
	switch sourceChar {
	case '(':
		sc.addToken(LEFT_PAREN)
	case ')':
		sc.addToken(RIGHT_PAREN)
	case '{':
		sc.addToken(LEFT_BRACE)
	case '}':
		sc.addToken(RIGHT_BRACE)
	case ',':
		sc.addToken(COMMA)
	case '.':
		sc.addToken(DOT)
	case '-':
		sc.addToken(MINUS)
	case '+':
		sc.addToken(PLUS)
	case ';':
		sc.addToken(SEMICOLON)
	case '*':
		sc.addToken(STAR)
	case '!':
		if sc.matchNext('=') {
			sc.addToken(BANG_EQUAL)
		} else {
			sc.addToken(BANG)
		}
	case '=':
		if sc.matchNext('=') {
			sc.addToken(EQUAL_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '<':
		if sc.matchNext('=') {
			sc.addToken(LESS_EQUAL)
		} else {
			sc.addToken(LESS)
		}
	case '>':
		if sc.matchNext('=') {
			sc.addToken(GREATER_EQUAL)
		} else {
			sc.addToken(GREATER)
		}
	case '/':
		if sc.matchNext('/') {
			for sc.peek() != '\n' && !sc.IsAtEnd() {
				sc.advance()
			}
			// omit all comment tokens
		} else if sc.matchNext('*') {
			var terminated bool
			for !sc.IsAtEnd() {
				if sc.matchNext('\n') {
					sc.line += 1
				}
				if sc.matchNext('*') && sc.matchNext('/') {
					terminated = true
					break
				}
				sc.advance()
			}
			if !terminated {
				return NewScanError(sc.line, strconv.Itoa(sc.current), fmt.Errorf("comment not terminated"))
			}
		} else {
			sc.addToken(SLASH)
		}
	case '"':
		if err := sc.string(); err != nil {
			return err
		}

	case '\n':
		sc.line += 1
	case ' ', '\r', '\t':
		break

	default:
		if sc.isDigit(sourceChar) {
			if err := sc.number(); err != nil {
				return err
			}
		} else if sc.isAlpha(sourceChar) {
			sc.identifier()
		} else {
			// todo typing for fmt.Errorf("unexpected character.")
			return NewScanError(sc.line, strconv.Itoa(sc.current), fmt.Errorf("unexpected character"))
		}
	}

	return nil
}

func (sc *Scanner) addToken(tokenType TokenType) {
	sc.addTokenWithLiteral(tokenType, nil)
}

func (sc *Scanner) addTokenWithLiteral(tokenType TokenType, literal *Literal) {
	lexeme := string(sc.source[sc.start:sc.current])
	sc.tokens = append(sc.tokens, NewToken(tokenType, lexeme, literal, sc.line))
}

func (sc *Scanner) isDigit(char rune) bool {
	return unicode.IsDigit(char)
}

func (sc *Scanner) isAlpha(char rune) bool {
	return unicode.IsLetter(char) || char == '_'
}

func (sc *Scanner) isAlphaNumeric(char rune) bool {
	return sc.isAlpha(char) || sc.isDigit(char)
}

func (sc *Scanner) identifier() {
	for sc.isAlphaNumeric(sc.peek()) {
		sc.advance()
	}
	identifierStr := string(sc.source[sc.start:sc.current])
	if keywordType, ok := keywordsToToken[identifierStr]; ok {
		sc.addToken(keywordType)
	} else {
		sc.addToken(IDENTIFIER)
	}
}

func (sc *Scanner) number() error {
	for sc.isDigit(sc.peek()) {
		sc.advance()
	}

	if sc.peek() == '.' && sc.isDigit(sc.peekNext()) {
		sc.advance()
		for sc.isDigit(sc.peek()) {
			sc.advance()
		}
	}

	floatVal, err := strconv.ParseFloat(string(sc.source[sc.start:sc.current]), 64)
	if err != nil {
		return NewScanError(sc.line, strconv.Itoa(sc.current), err)
	}
	sc.addTokenWithLiteral(NUMBER, NewLiteral(NewFloatLoxValue(floatVal)))
	return nil
}

func (sc *Scanner) string() error {
	for sc.peek() != '"' && !sc.IsAtEnd() {
		if sc.peek() == '\n' {
			sc.line += 1
		}
		sc.advance()
	}

	if sc.IsAtEnd() {
		return NewScanError(sc.line, strconv.Itoa(sc.current), fmt.Errorf("unterminated string"))
	}
	sc.advance()

	val := string(sc.source[sc.start+1 : sc.current-1])
	sc.addTokenWithLiteral(STRING, NewLiteral(NewStringLoxValue(val)))
	return nil
}
