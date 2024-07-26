package token

import "fmt"

type Type string

const (
	// Single-character tokens.
	LEFT_PAREN  Type = "("
	RIGHT_PAREN Type = ")"
	LEFT_BRACE  Type = "{"
	RIGHT_BRACE Type = "}"
	COMMA       Type = ","
	DOT         Type = "."
	MINUS       Type = "-"
	PLUS        Type = "+"
	SEMICOLON   Type = ";"
	SLASH       Type = "\\"
	STAR        Type = "*"

	// One or two character tokens.
	BANG          Type = "!"
	BANG_EQUAL    Type = "!="
	EQUAL         Type = "="
	EQUAL_EQUAL   Type = "=="
	GREATER       Type = ">"
	GREATER_EQUAL Type = ">="
	LESS          Type = "<"
	LESS_EQUAL    Type = "<="

	// Literals.
	IDENTIFIER Type = "INDENTIFIER"
	STRING     Type = "STRING"
	NUMBER     Type = "NUMBER"

	// Keywords.
	AND    Type = "AND"
	CLASS  Type = "CLASS"
	ELSE   Type = "ESLE"
	FALSE  Type = "FALSE"
	FUN    Type = "FUN"
	FOR    Type = "FOR"
	IF     Type = "IF"
	NIL    Type = "NIL"
	OR     Type = "OR"
	PRINT  Type = "PRINT"
	RETURN Type = "RETURN"
	SUPER  Type = "SUPER"
	THIS   Type = "THIS"
	TRUE   Type = "TRUE"
	VAR    Type = "VAR"
	WHILE  Type = "WHILE"

	EOF Type = "EOF"

	UNKNOWN Type = "unknown"
)

type Token struct {
	tokenType Type
	lexeme    string
	literal   any
	line      int
}

func New(t Type, lexeme string, literal any, line int) *Token {
	return &Token{
		tokenType: t,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t *Token) Type() Type {
	return t.tokenType
}

func (t *Token) Literal() any {
	return t.literal
}

func (t *Token) Lexeme() string {
	return t.lexeme
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %v", t.tokenType, t.lexeme, t.literal)
}

func ToToken(text string) Type {
	switch text {
	case "and":
		return AND
	case "class":
		return CLASS
	case "else":
		return ELSE
	case "false":
		return FALSE
	case "for":
		return FOR
	case "fun":
		return FUN
	case "if":
		return IF
	case "nil":
		return NIL
	case "or":
		return OR
	case "print":
		return PRINT
	case "return":
		return RETURN
	case "super":
		return SUPER
	case "this":
		return THIS
	case "true":
		return TRUE
	case "var":
		return VAR
	case "while":
		return WHILE
	}

	return UNKNOWN
}
