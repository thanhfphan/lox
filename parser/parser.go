package parser

import (
	"fmt"
	"lox/ast"

	"lox/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parser() ast.Expr {
	return p.expression()
}

func (p *Parser) ParserStmt() []ast.Stmt {
	stmts := []ast.Stmt{}
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(token.CLASS) {
		return p.classDeclaration()
	}
	if p.match(token.FUN) {
		return p.function("function")
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.stmt()
}

func (p *Parser) stmt() ast.Stmt {
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.FOR) {
		return p.forStmt()
	}
	if p.match(token.IF) {
		return p.ifStmt()
	}
	if p.match(token.RETURN) {
		return p.returnStmt()
	}
	if p.match(token.WHILE) {
		return p.whileStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return &ast.BlockStmt{
			Statements: p.block(),
		}
	}

	return p.expressionStmt()
}

func (p *Parser) returnStmt() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")

	return &ast.ReturnStmt{
		KeyWord: keyword,
		Value:   value,
	}
}

func (p *Parser) function(kind string) *ast.FunctionStmt {
	funcName := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")

	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")

	parameters := []*token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name."))
	}
	for p.match(token.COMMA) {
		parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name."))
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body := p.block()

	return &ast.FunctionStmt{
		Name:   funcName,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) forStmt() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initStmt ast.Stmt
	if p.match(token.SEMICOLON) {
		initStmt = nil
	} else if p.match(token.VAR) {
		initStmt = p.varDeclaration()
	} else {
		initStmt = p.expressionStmt()
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.stmt()

	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				body,
				&ast.ExpressionStmt{
					Expression: increment,
				}},
		}
	}

	if condition == nil {
		condition = &ast.LiteralExpr{
			Val: true,
		}
	}

	body = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initStmt != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				initStmt,
				body,
			},
		}
	}

	return body
}

func (p *Parser) whileStmt() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.stmt()

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) ifStmt() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.stmt()
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch = p.stmt()
	}

	return &ast.IfStmt{
		Condition: condition,
		Then:      thenBranch,
		Else:      elseBranch,
	}
}

func (p *Parser) classDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect class name.")
	p.consume(token.LEFT_BRACE, "Expect '{' before class body.")

	methods := []*ast.FunctionStmt{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(token.RIGHT_BRACE, "Epect '}' after class body.")

	return &ast.ClassStmt{
		Name:    name,
		Methods: methods,
	}
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name")
	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")

	return &ast.VarStmt{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block")

	return statements
}

func (p *Parser) printStmt() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")

	return &ast.PrintStmt{
		Expression: expr,
	}
}

func (p *Parser) expressionStmt() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")

	return &ast.ExpressionStmt{
		Expression: expr,
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type() == token.SEMICOLON {
			return
		}

		switch p.peek().Type() {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()
	if p.match(token.EQUAL) {
		equals := p.previous()
		val := p.assignment()

		switch v := expr.(type) {
		case *ast.VariableExpr:
			return &ast.AssignExpr{
				Name:  v.Name,
				Value: val,
			}
		case *ast.GetExpr:
			return &ast.SetExpr{
				Object: v.Object,
				Name:   v.Name,
				Value:  val,
			}
		}

		panic(fmt.Errorf("Invalid assignment target: %s", equals.Literal()))
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		op := p.previous()
		right := p.and()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		op := p.previous()
		right := p.equality()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

/*
expression     → literal

	| unary
	| binary
	| grouping ;
*/
func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() ast.Expr {
	expr := p.comparision()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparision()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr

}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparision() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		op := p.previous()
		right := p.factor()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    *op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(token.SLASH, token.STAR) {
		op := p.previous()
		right := p.unary()
		expr = &ast.BinaryExpr{
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
func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		op := p.previous()
		right := p.unary()
		return &ast.UnaryExpr{
			Op:    *op,
			Right: right,
		}
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(token.DOT) {
			name := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			expr = &ast.GetExpr{
				Name:   name,
				Object: expr,
			}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	args := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(token.COMMA) {
			args = append(args, p.expression())
		}
	}

	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments")

	return &ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: args,
	}
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//
//	| "(" expression ")" ;
func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.LiteralExpr{Val: false}
	}
	if p.match(token.TRUE) {
		return &ast.LiteralExpr{Val: true}
	}
	if p.match(token.NIL) {
		return &ast.LiteralExpr{Val: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return &ast.LiteralExpr{Val: p.previous().Literal()}
	}
	if p.match(token.THIS) {
		return &ast.ThisExpr{
			Keyword: p.previous(),
		}
	}
	if p.match(token.IDENTIFIER) {
		return &ast.VariableExpr{
			Name: p.previous(),
		}
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.GroupingExpr{
			Expression: expr,
		}
	}

	panic("Expect expression.")
}

func (p *Parser) consume(t token.Type, msg string) *token.Token {
	if p.check(t) {
		return p.advance()
	}

	err := fmt.Errorf("consume TokenType: %s error, msg: %s", t, msg)
	panic(err)
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type() == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type() == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}
