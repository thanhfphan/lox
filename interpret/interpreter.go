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

var (
	GLOBAL_ENV = env.New(nil)
)

func init() {
	GLOBAL_ENV.Define("clock", NewClock())
}

type Interpreter struct {
	env    *env.Env
	locals map[ast.Expr]int
}

func New() *Interpreter {
	return &Interpreter{
		env:    GLOBAL_ENV,
		locals: make(map[ast.Expr]int),
	}
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
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
	if stmt.Initializer != nil {
		v = i.evaluate(stmt.Initializer)
	}

	i.env.Define(stmt.Name.Lexeme(), v)
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) {
	newEnv := env.New(i.env)
	i.executeBlock(stmt.Statements, newEnv)
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) any {
	fun := NewFunction(stmt, i.env)
	i.env.Define(stmt.Name.Lexeme(), fun)
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	var value any
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	r := &Return{
		Value: value,
	}
	// hack to back top of the Stack
	panic(r)
}

// Expr visitors
func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return expr.Val
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	right := i.evaluate(expr.Right)

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
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) any {
	val := i.evaluate(expr.Value)

	distance, has := i.locals[expr]
	if has {
		i.env.AssignAt(distance, expr.Name, val)
	} else {
		GLOBAL_ENV.Assign(expr.Name, val)
	}

	return val
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) any {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type() == ast.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) any {
	callee := i.evaluate(expr.Callee)

	args := []any{}
	for _, item := range expr.Arguments {
		args = append(args, i.evaluate(item))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic("Can't parse callee to LoxCallable")
	}
	if len(args) != function.Arity() {
		panic(fmt.Errorf("Expect %d arguments but got %d.", function.Arity(), len(args)))
	}

	return function.Call(i, args)
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, env *env.Env) {
	prevEnv := i.env
	defer func() {
		i.env = prevEnv
	}()
	i.env = env

	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *Interpreter) lookUpVariable(name *ast.Token, expr ast.Expr) any {
	distance, has := i.locals[expr]
	if has {
		return i.env.GetAt(distance, name.Lexeme())
	}

	val, err := GLOBAL_ENV.Get(name)
	if err != nil {
		panic(fmt.Errorf("GlobalEnv.Get error: %w", err))
	}

	return val
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
