package ast

import "lox/token"

type Expr interface {
	Accept(v ExprVisitor) any
}

// BinaryExpr ...
type BinaryExpr struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(e)
}

// UnaryExpr ...
type UnaryExpr struct {
	Op    token.Token
	Right Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(e)
}

// LiteralExpr ...
type LiteralExpr struct {
	Val any
}

func (e *LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(e)
}

// GroupingExpr ...
type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(e)
}

// VariableExpr ...
type VariableExpr struct {
	Name *token.Token
}

func (e *VariableExpr) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(e)
}

// AssignExpr ...
type AssignExpr struct {
	Name  *token.Token
	Value Expr
}

func (e *AssignExpr) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(e)
}

// LogicalExpr ...
type LogicalExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (s *LogicalExpr) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(s)
}

// CallExpr ...
type CallExpr struct {
	Callee    Expr
	Paren     *token.Token
	Arguments []Expr
}

func (s *CallExpr) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(s)
}

// GetExpr ...
type GetExpr struct {
	Object Expr
	Name   *token.Token
}

func (s *GetExpr) Accept(v ExprVisitor) any {
	return v.VisitGetExpr(s)
}

// SetExpr ...
type SetExpr struct {
	Object Expr
	Name   *token.Token
	Value  Expr
}

func (s *SetExpr) Accept(v ExprVisitor) any {
	return v.VisitSetExpr(s)
}

// ThisExpr ...
type ThisExpr struct {
	Keyword *token.Token
}

func (s *ThisExpr) Accept(v ExprVisitor) any {
	return v.VisitThisExpr(s)
}

// SuperExpr
type SuperExpr struct {
	Keyword *token.Token
	Method  *token.Token
}

func (s *SuperExpr) Accept(v ExprVisitor) any {
	return v.VisitSuperExpr(s)
}
