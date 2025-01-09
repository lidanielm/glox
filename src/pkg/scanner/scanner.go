package scanner

import (
	"strconv"

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

// TODO: Rewrite with custom error
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
	case '"':
		scan.addString()
	case 'o':
		if scan.matchNext('r') {
			scan.addToken(token.OR)
		}
	default:
		if isDigit(c) {
			scan.addNumber()
		} else if isAlpha(c) {
			scan.addIdentifier()
		} else {
			return scan.err.New(scan.line, "Unexpected character.")
		}
	}
	return nil
}

func (scan *Scanner) addString() {
	if scan.isEOF() {
		scan.err.New(scan.line, "Unterminated string.")
		return
	}

	str := make([]byte, 0)
	for !scan.isEOF() && scan.peek() != '"' {
		if scan.peek() == '\n' {
			scan.line++
		}

		str = append(str, scan.advance())
	}

	// Last '"'
	scan.advance()

	scan.addTokenLiteral(token.STRING, string(str))
}

func (scan *Scanner) addNumber() {
	numStr := make([]byte, 0)
	hasDec := false
	for !scan.isEOF() {
		if scan.peek() == '.' && isDigit(scan.peekTwice()) && !hasDec {
			hasDec = true
			scan.advance()
		} else if isDigit(scan.peek()) {
			numStr = append(numStr, scan.advance())
		} else {
			break
		}
	}

	if numStr[len(numStr) - 1] == '.' {
		scan.err.New(scan.line, "Invalid number (trailing decimal point).")
		return
	}

	num, err := strconv.ParseFloat(string(numStr), 64)
	if err != nil {
		scan.err.New(scan.line, "Invalid number.")
		return
	}
	scan.addTokenLiteral(token.NUMBER, num)
}

func (scan *Scanner) addIdentifier() {
	idStr := make([]byte, 0)

	for isAlphaNumeric(scan.peek()) {
		idStr = append(idStr, scan.advance())
	}

	if typ, ok := token.Keywords[string(idStr)]; ok {
		scan.addToken(typ)
	} else {
		scan.addToken(token.IDENTIFIER)
	}
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

func (scan *Scanner) peekTwice() byte {
	if scan.current + 1 >= len(scan.source) {
		return '\u0000'
	}

	return scan.source[scan.current + 1]	
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

func isDigit(c byte) bool {
	return int(c) >= 48 && int(c) <= 57
}

func isAlpha(c byte) bool {
	return (int(c) >= 65 && int(c) <= 90) || (int(c) >= 97 && int(c) <= 122) || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) && isAlpha(c)
}