package scanner

import (
	"pkg/tokentype",
	"pkg/error",
	"errors"
)

type Scanner struct {
	source string
	tokens []Token
	start int
	current int
	line int
	err Error
}

func NewScanner(source string, err) *Scanner {
	tokens := make([]Token)
	err := NewError()
	return &Scanner{source: source, tokens: tokens, start: 0, current: 0, line: 1, err: err}
}

func (scan *Scanner) ScanTokens() (tokens []Token, err error) {
	while (!isEOF()) {
		scan.start = scan.current
		_ := scanToken()
	}

	scan.tokens = append(scan.tokens, NewToken(EOF, "", null, line))
	return scan.tokens, nil
}

func (scan *Scanner) scanToken() *Error {
	switch c := scan.advance(); c {
	case '(':
		scan.addToken(LEFT_PAREN)
	case ')':
		scan.addToken(RIGHT_PAREN)
	case '{':
		scan.addToken(LEFT_BRACE)
	case '}':
		scan.addToken(RIGHT_BRACE)
	case ',':
		scan.addToken(COMMA)
	case '.':
		scan.addToken(DOT)
	case '-':
		scan.addToken(MINUS)
	case '+':
		scan.addToken(PLUS)
	case ';':
		scan.addToken(SEMICOLON)
	case '/':
		scan.addToken(SLASH)
	case '*':
		scan.addToken(STAR)
	case '!':
		if scan.matchNext('=') {
			scan.addToken(BANG_EQUAL)
		} else {
			scan.addToken(BANG)
		}
	case '=':
		if scan.matchNext('=') {
			scan.addToken(EQUAL_EQUAL)
		} else {
			scan.addToken(EQUAL)
		}
	case '>':
		if scan.matchNext('=') {
			scan.addToken(GREATER_EQUAL)
		} else {
			scan.addToken(GREATER)
		}
	case '<':
		if scan.matchNext('=') {
			scan.addToken(LESS_EQUAL)
		} else {
			scan.addToken(LESS)
		}
	case '/':
		if scan.matchNext('/') {
			while peek() != '\n' && !isEOF() {
				advance()
			}
		} else {
			scan.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		line++
	default:
		return err.New(line, "Unexpected character.")
	}
	return nil
}

func (scan *Scanner) matchNext(byte expected) bool {
	if scan.isEOF() {
		return false
	}
	if scan.source[current] != expected {
		return false
	}
	current++
	return true
}

func (scan *Scanner) peek() byte {
	// Return current byte
	return scan.source[current]
}

func (scan *Scanner) advance() byte {
	// Return current byte and advance pointer
	c := scan.source[current]
	current++
	return c
}

func (scan *Scanner) addToken(typ TokenType) {
	addToken(typ, nil)
}

func (scan *Scanner) addToken(typ TokenType, literal interface{}) {
	string text = scan.source[start:current + 1]
	scan.tokens = append(scan.tokens, NewToken(EOF, text, literal, line))
}

func (scan *Scanner) isEOF() bool {
	return scan.current >= scan.source.length()
}

