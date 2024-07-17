package interpret

import (
	"fmt"
	"lox/ast"
	"lox/env"
	"reflect"
)

var (
	_ ast.ExprVisitor = (*Interpreter)(nil)
	_ ast.StmtVisitor = (*Interpreter)(nil)
)

type Interpreter struct {
	env *env.Env
}

func New() *Interpreter {
	return &Interpreter{
		env: env.New(),
	}
}

func (i *Interpreter) Interpret(expr ast.Expr) {
	obj := i.evaluate(expr)
	fmt.Printf("Interpret obj: %v\n", obj)
}

func (i *Interpreter) InterpretStmt(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		i.executeStmt(stmt)
	}
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) executeStmt(stmt ast.Stmt) {
	stmt.Accept(i)
}

// Stmt visitors
func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) {
	val := i.evaluate(stmt.Expression)
	fmt.Println(val)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) {
	var v any
	if stmt.Expr != nil {
		v = i.evaluate(stmt.Expr)
	}

	i.env.Define(stmt.Name.Lexeme(), v)
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) {
	newEnv := env.New()
	newEnv.SetEnclosing(i.env)
	i.executeBlock(stmt.Statements, newEnv)
}

// Expr visitors
func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return expr.Val
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	right := i.evaluate(expr.Expr)

	switch expr.Op.Type() {
	case ast.BANG:
		return !i.isTruthy(right)
	case ast.MINUS:
		v := right.(float64)
		return -v
	}

	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	// TODO: check before cast value
	switch expr.Op.Type() {
	case ast.MINUS:
		return left.(float64) - right.(float64)
	case ast.SLASH:
		return left.(float64) / right.(float64)
	case ast.STAR:
		return left.(float64) * right.(float64)
	case ast.PLUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 &&
			reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}

		if reflect.TypeOf(left).Kind() == reflect.String &&
			reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string)
		}
	case ast.GREATER:
		return left.(float64) > right.(float64)
	case ast.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case ast.LESS:
		return left.(float64) < right.(float64)
	case ast.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case ast.BANG_EQUAL:
		return !i.isEqual(left, right)
	case ast.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}
func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) any {
	val, err := i.env.Get(expr.Name)
	if err != nil {
		panic(err)
	}

	return val
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) any {
	val := i.evaluate(expr.Value)
	err := i.env.Assign(expr.Name, val)
	if err != nil {
		panic(err)
	}

	return val
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, env *env.Env) {
	prevEnv := i.env
	i.env = env

	for _, stmt := range stmts {
		i.executeStmt(stmt)
	}

	i.env = prevEnv
}

func (i *Interpreter) isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}

	return reflect.DeepEqual(left, right)
}

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if reflect.TypeOf(obj).Kind() == reflect.Bool {
		return reflect.ValueOf(obj).Bool()
	}

	return false
}
