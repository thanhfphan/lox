package scan

import (
	"fmt"
	"lox/ast"
	"os"
	"strconv"
)

type Scanner struct {
	source []rune
	tokens []*ast.Token

	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source: source,
		tokens: []*ast.Token{},
		line:   1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(t ast.TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t ast.TokenType, literal any) {
	text := s.source[s.start:s.current]
	token := ast.NewToken(t, string(text), literal, s.line)
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(ast.LEFT_PAREN)
	case ')':
		s.addToken(ast.RIGHT_PAREN)
	case '{':
		s.addToken(ast.LEFT_BRACE)
	case '}':
		s.addToken(ast.RIGHT_BRACE)
	case ',':
		s.addToken(ast.COMMA)
	case '.':
		s.addToken(ast.DOT)
	case '-':
		s.addToken(ast.MINUS)
	case '+':
		s.addToken(ast.PLUS)
	case ';':
		s.addToken(ast.SEMICOLON)
	case '*':
		s.addToken(ast.STAR)
	case '!':
		t := ast.BANG
		if s.match('=') {
			t = ast.BANG_EQUAL
		}
		s.addToken(t)
	case '=':
		t := ast.EQUAL
		if s.match('=') {
			t = ast.EQUAL_EQUAL
		}
		s.addToken(t)
	case '<':
		t := ast.LESS
		if s.match('=') {
			t = ast.LESS_EQUAL
		}
		s.addToken(t)
	case '>':
		t := ast.GREATER
		if s.match('=') {
			t = ast.GREATER_EQUAL
		}
		s.addToken(t)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if (s.current + 1) >= len(s.source) {
		return rune(0)
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType := ast.ToToken(string(text))
	if tokenType == ast.UNKNOWN {
		tokenType = ast.IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peek()) {
		// Consume the "."
		s.advance()
	}

	num, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		panic(err)
	}
	s.addTokenLiteral(ast.NUMBER, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error(s.line, "Unterminated string.")
		return
	}

	// Chop the closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(ast.STRING, string(value))
}

func (s *Scanner) ScanTokens() []*ast.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eofToken := ast.NewToken(ast.EOF, "", nil, s.line)
	s.tokens = append(s.tokens, eofToken)

	return s.tokens
}

func (s *Scanner) error(line int, msg string) {
	fmt.Printf("[line %d] Error: %s\n", line, msg)
	os.Exit(1)
}
