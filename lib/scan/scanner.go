package scan

type Scanner struct {
	source               []rune
	tokens               []Token
	start, current, line int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
	}
}

func (sc *Scanner) advance() rune {
	sc.current += 1
	return sc.source[sc.current-1]
}

func (sc *Scanner) matchNext(expected rune) bool {
	if sc.IsAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}
	sc.current += 1
	return true
}

func (sc *Scanner) peek() rune {
	if sc.IsAtEnd() {
		return '\x00'
	}
	return sc.source[sc.current]
}

func (sc *Scanner) peekNext() rune {
	if sc.current+1 >= len(sc.source) {
		return '\x00'
	}
	return sc.source[sc.current+1]
}

func (sc *Scanner) IsAtEnd() bool {
	return sc.current >= len(sc.source)
}
