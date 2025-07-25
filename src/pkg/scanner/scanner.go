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
}

func NewScanner(source string) *Scanner {
	tokens := make([]token.Token, 0)
	return &Scanner{source: source, tokens: tokens, start: 0, current: 0, line: 1}
}

func (scan *Scanner) ScanTokens() ([]token.Token, error) {
	for !scan.isEOF() {
		scan.start = scan.current
		err := scan.scanToken()
		if err != nil {
			return nil, err
		}
	}

	scan.tokens = append(scan.tokens, *token.NewToken(token.EOF, "", nil, scan.line))
	return scan.tokens, nil
}

func (scan *Scanner) scanToken() error {
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
			for !scan.isEOF() && scan.peek() != '\n' {
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
			break
		}
		fallthrough
	default:
		if isDigit(c) {
			scan.addNumber()
		} else if isAlpha(c) {
			scan.addIdentifier()
		} else {
			return lox_error.NewError(*token.NewToken(token.ERROR, string(c), nil, scan.line), "Unexpected character.")
		}
	}
	return nil
}

func (scan *Scanner) addString() error {
	if scan.isEOF() {
		return lox_error.NewError(*token.NewToken(token.ERROR, "", nil, scan.line), "Unterminated string.")
	}

	for !scan.isEOF() && scan.peek() != '"' {
		if scan.peek() == '\n' {
			scan.line++
		}

		scan.advance()
	}

	// Last '"'
	scan.advance()

    str := scan.source[scan.start + 1:scan.current - 1]
	scan.addTokenLiteral(token.STRING, str)
	return nil
}

func (scan *Scanner) addNumber() error {
    // TODO: Support negative numbers
    // Scan until number terminates
	hasDec := false
	for !scan.isEOF() {
		if scan.peek() == '.' && isDigit(scan.peekTwice()) && !hasDec {
			hasDec = true
			scan.advance()
		} else if isDigit(scan.peek()) {
            scan.advance()
		} else {
			break
		}
	}

    // Get string-formatted number from source
    numStr := scan.source[scan.start:scan.current]

    // Validate number
	if numStr[len(numStr) - 1] == '.' {
		return lox_error.NewError(*token.NewToken(token.ERROR, "", nil, scan.line), "Invalid number (trailing decimal point).")
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return lox_error.NewError(*token.NewToken(token.ERROR, "", nil, scan.line), "Invalid number.")
	}

    // Add token
	scan.addTokenLiteral(token.NUMBER, num)
	return nil
}

func (scan *Scanner) addIdentifier() {
	for !scan.isEOF() && isAlphaNumeric(scan.peek()) {
		scan.advance()
	}

    identifier := scan.source[scan.start:scan.current]
	if typ, ok := token.Keywords[identifier]; ok {
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
	text := scan.source[scan.start:scan.current]
	scan.tokens = append(scan.tokens, *token.NewToken(typ, text, literal, scan.line))
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
	return isDigit(c) || isAlpha(c)
}
