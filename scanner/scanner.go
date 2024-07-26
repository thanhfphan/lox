package scanner

import (
	"fmt"
	"lox/token"
	"os"
	"strconv"
)

type Scanner struct {
	source []rune
	tokens []*token.Token

	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source: source,
		tokens: []*token.Token{},
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

func (s *Scanner) addToken(t token.Type) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t token.Type, literal any) {
	text := s.source[s.start:s.current]
	token := token.New(t, string(text), literal, s.line)
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		t := token.BANG
		if s.match('=') {
			t = token.BANG_EQUAL
		}
		s.addToken(t)
	case '=':
		t := token.EQUAL
		if s.match('=') {
			t = token.EQUAL_EQUAL
		}
		s.addToken(t)
	case '<':
		t := token.LESS
		if s.match('=') {
			t = token.LESS_EQUAL
		}
		s.addToken(t)
	case '>':
		t := token.GREATER
		if s.match('=') {
			t = token.GREATER_EQUAL
		}
		s.addToken(t)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
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
	tokenType := token.ToToken(string(text))
	if tokenType == token.UNKNOWN {
		tokenType = token.IDENTIFIER
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
	s.addTokenLiteral(token.NUMBER, num)
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
	s.addTokenLiteral(token.STRING, string(value))
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	eofToken := token.New(token.EOF, "", nil, s.line)
	s.tokens = append(s.tokens, eofToken)

	return s.tokens
}

func (s *Scanner) error(line int, msg string) {
	fmt.Printf("[line %d] Error: %s\n", line, msg)
	os.Exit(1)
}
