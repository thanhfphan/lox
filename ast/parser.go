package ast

import "fmt"

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parser() Expr {
	return p.expression()
}

func (p *Parser) ParserStmt() []Stmt {
	stmts := []Stmt{}
	for !p.isAtEnd() {
		stmts = append(stmts, p.stmt())
	}
	return stmts
}

func (p *Parser) stmt() Stmt {
	if p.match(PRINT) {
		return p.printStmt()
	}

	return p.expressionStmt()
}

func (p *Parser) printStmt() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")

	return &PrintStmt{expr}
}

func (p *Parser) expressionStmt() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")

	return &ExpressionStmt{expr}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type() == SEMICOLON {
			return
		}

		switch p.peek().Type() {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

/*
expression     → literal

	| unary
	| binary
	| grouping ;
*/
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparision()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparision()

		expr = &BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr

}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparision() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()

		expr = &BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.factor()

		expr = &BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()

		expr = &BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}
	return expr
}

// unary          → ( "!" | "-" ) unary
//
//	| primary ;
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right := p.unary()

		return &UnaryExpr{
			Op:   *op,
			Expr: right,
		}
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//
//	| "(" expression ")" ;
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{false}
	}
	if p.match(TRUE) {
		return &LiteralExpr{true}
	}
	if p.match(NIL) {
		return &LiteralExpr{nil}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{p.previous().Literal()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{
			Expr: expr,
		}
	}

	panic("Expect expression.")
}

func (p *Parser) consume(t TokenType, msg string) *Token {
	if p.check(t) {
		return p.advance()
	}

	err := fmt.Errorf("consume TokenType: %s error, msg: %s", t, msg)
	panic(err)
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type() == t
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}
