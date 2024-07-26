package resolver

import (
	"fmt"
	"lox/ast"
	"lox/dst"
	"lox/interpreter"
	"lox/token"
)

var (
	_ ast.ExprVisitor = (*Resolver)(nil)
	_ ast.StmtVisitor = (*Resolver)(nil)
)

type FunctionType int

const (
	NONE     FunctionType = 1
	FUNCTION FunctionType = 2
)

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      *dst.Stack[map[string]bool]
	currentFunc FunctionType
}

func NewResolver(i *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter: i,
		scopes:      dst.NewStack[map[string]bool](),
		currentFunc: NONE,
	}
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) resolveListStmt(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveFunction(f *ast.FunctionStmt, funcType FunctionType) {
	r.beginScope()
	enclosingFunc := r.currentFunc
	r.currentFunc = funcType
	for _, param := range f.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveListStmt(f.Body)
	r.endScope()
	r.currentFunc = enclosingFunc
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *token.Token) {
	pointer := r.scopes.Peek()
	dept := 0
	for pointer != nil {
		if pointer.Val[name.Lexeme()] {
			r.interpreter.Resolve(expr, dept)
			return
		}
		pointer = pointer.Next
		dept++
	}
}

func (r *Resolver) declare(name *token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek().Val
	if _, has := scope[name.Lexeme()]; has {
		panic(fmt.Errorf("%s Already a variable with this name in this scope.", name.String()))
	}
	scope[name.Lexeme()] = false
}

func (r *Resolver) define(name *token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek().Val
	scope[name.Lexeme()] = true
}

func (r *Resolver) Resolve(stmts []ast.Stmt) {
	r.resolveListStmt(stmts)
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) {
	r.beginScope()
	r.resolveListStmt(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)

	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)
	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	if r.currentFunc == NONE {
		panic(fmt.Errorf("%s Can't return from top-level code.", stmt.KeyWord.String()))
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) any {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) any {
	if !r.scopes.IsEmpty() {
		scope := r.scopes.Peek().Val
		if val, has := scope[expr.Name.Lexeme()]; has && !val {
			panic(expr.Name.String() + "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil
}
