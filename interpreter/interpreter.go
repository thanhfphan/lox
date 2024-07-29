package interpreter

import (
	"fmt"
	"lox/ast"
	"lox/env"
	"lox/token"
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
func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) any {
	val := i.evaluate(stmt.Expression)
	fmt.Println(val)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) any {
	var v any
	if stmt.Initializer != nil {
		v = i.evaluate(stmt.Initializer)
	}

	i.env.Define(stmt.Name.Lexeme(), v)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) any {
	newEnv := env.New(i.env)
	i.executeBlock(stmt.Statements, newEnv)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
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

func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) any {
	i.env.Define(stmt.Name.Lexeme(), nil)

	methods := map[string]*Function{}
	for _, method := range stmt.Methods {
		f := NewFunction(method, i.env)
		methods[method.Name.Lexeme()] = f
	}

	c := NewClass(stmt.Name.Lexeme(), methods)
	i.env.Assign(stmt.Name, c)
	return nil
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
	case token.BANG:
		return !i.isTruthy(right)
	case token.MINUS:
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
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	case token.PLUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 &&
			reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}

		if reflect.TypeOf(left).Kind() == reflect.String &&
			reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string)
		}
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
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
	if expr.Operator.Type() == token.OR {
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

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) any {
	obj := i.evaluate(expr.Object)
	ins, ok := obj.(*Instance)
	if !ok {
		panic(fmt.Errorf("%s Only instances have properties.", expr.Name))
	}

	return ins.Get(expr.Name)
}

func (i *Interpreter) VisitSetExpr(expr *ast.SetExpr) any {
	obj := i.evaluate(expr.Object)
	ins, ok := obj.(*Instance)
	if !ok {
		panic(fmt.Errorf("%s Only instances have fields.", expr.Name))
	}

	val := i.evaluate(expr.Value)
	ins.Set(expr.Name, val)

	return val
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

func (i *Interpreter) lookUpVariable(name *token.Token, expr ast.Expr) any {
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
