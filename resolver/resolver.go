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

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      *dst.Stack[map[string]bool]

	currentFunc  FunctionType
	currentClass ClassType
}

func NewResolver(i *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:  i,
		scopes:       dst.NewStack[map[string]bool](),
		currentFunc:  FT_NONE,
		currentClass: CT_NONE,
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

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) any {
	r.beginScope()
	r.resolveListStmt(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FT_FUNCTION)

	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)
	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	if r.currentFunc == FT_NONE {
		panic(fmt.Errorf("%s Can't return from top-level code.", stmt.KeyWord.String()))
	}

	if stmt.Value != nil {
		if r.currentFunc == FT_INITIALIZER {
			panic(fmt.Errorf("%v Can't return a value from an initializer.", stmt.KeyWord))
		}

		r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = CT_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.SuperClass != nil {
		if stmt.Name.Lexeme() == stmt.SuperClass.Name.Lexeme() {
			panic(fmt.Errorf("A class '%s' can't inherit from itself.", stmt.Name.Lexeme()))
		}

		r.currentClass = CT_SUBCLASS
		r.resolveExpr(stmt.SuperClass)

		r.beginScope()
		peek := r.scopes.Peek().Val
		peek["super"] = true
	}

	r.beginScope()
	peek := r.scopes.Peek().Val
	peek["this"] = true
	for _, method := range stmt.Methods {
		declaration := FT_METHOD
		if method.Name.Lexeme() == "init" {
			declaration = FT_INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()

	if stmt.SuperClass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass

	return nil
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

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) any {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *ast.ThisExpr) any {
	if r.currentClass == CT_NONE {
		panic(fmt.Errorf("%v Can't use 'this' outside of a class.", expr.Keyword))
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *ast.SuperExpr) any {
	if r.currentClass == CT_NONE {
		panic("Can't use 'super' outside of a class.")
	} else if r.currentClass != CT_SUBCLASS {
		panic("Can't use 'super' in a class with no superclass.")
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}
