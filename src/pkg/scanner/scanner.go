package scanner

import (
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Scanner struct {
	source string
	tokens []token.Token
	start int
	current int
	line int
	err lox_error.Error
}

func NewScanner(source string) *Scanner {
	tokens := make([]token.Token, 0)
	err := lox_error.NewError()
	return &Scanner{source: source, tokens: tokens, start: 0, current: 0, line: 1, err: *err}
}

func (scan *Scanner) ScanTokens() (tokens []token.Token, err error) {
	for !scan.isEOF() {
		scan.start = scan.current
		scan.scanToken()
	}

	scan.tokens = append(scan.tokens, *token.NewToken(token.EOF, "", nil, scan.line))
	return scan.tokens, nil
}

func (scan *Scanner) scanToken() *lox_error.Error {
	switch c := scan.advance(); c {
	case '(':
		scan.addToken(token.LEFT_PAREN)
	case ')':
		scan.addToken(token.RIGHT_PAREN)
	case '{':
		scan.addToken(token.LEFT_BRACE)
	case '}':
		scan.addToken(token.RIGHT_BRACE)
	case ',':
		scan.addToken(token.COMMA)
	case '.':
		scan.addToken(token.DOT)
	case '-':
		scan.addToken(token.MINUS)
	case '+':
		scan.addToken(token.PLUS)
	case ';':
		scan.addToken(token.SEMICOLON)
	case '*':
		scan.addToken(token.STAR)
	case '!':
		if scan.matchNext('=') {
			scan.addToken(token.BANG_EQUAL)
		} else {
			scan.addToken(token.BANG)
		}
	case '=':
		if scan.matchNext('=') {
			scan.addToken(token.EQUAL_EQUAL)
		} else {
			scan.addToken(token.EQUAL)
		}
	case '>':
		if scan.matchNext('=') {
			scan.addToken(token.GREATER_EQUAL)
		} else {
			scan.addToken(token.GREATER)
		}
	case '<':
		if scan.matchNext('=') {
			scan.addToken(token.LESS_EQUAL)
		} else {
			scan.addToken(token.LESS)
		}
	case '/':
		if scan.matchNext('/') {
			for scan.peek() != '\n' && !scan.isEOF() {
				scan.advance()
			}
		} else {
			scan.addToken(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		scan.line++
	default:
		return scan.err.New(scan.line, "Unexpected character.")
	}
	return nil
}

func (scan *Scanner) matchNext(expected byte) bool {
	if scan.isEOF() {
		return false
	}
	if scan.source[scan.current] != expected {
		return false
	}
	scan.current++
	return true
}

func (scan *Scanner) peek() byte {
	// Return current byte
	return scan.source[scan.current]
}

func (scan *Scanner) advance() byte {
	// Return current byte and advance pointer
	c := scan.source[scan.current]
	scan.current++
	return c
}

func (scan *Scanner) addToken(typ token.TokenType) {
	scan.addTokenLiteral(typ, nil)
}

func (scan *Scanner) addTokenLiteral(typ token.TokenType, literal interface{}) {
	text := scan.source[scan.start:scan.current + 1]
	scan.tokens = append(scan.tokens, *token.NewToken(token.EOF, text, literal, scan.line))
}

func (scan *Scanner) isEOF() bool {
	return scan.current >= len(scan.source)
}

