package ast

import "lox/token"

type Stmt interface {
	Accept(v StmtVisitor)
}

// PrintStmt ...
type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(v StmtVisitor) {
	v.VisitPrintStmt(s)
}

// ExpressionStmt ...
type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) {
	v.VisitExpressionStmt(s)
}

// VarStmt ...
type VarStmt struct {
	Name        *token.Token
	Initializer Expr
}

func (s *VarStmt) Accept(v StmtVisitor) {
	v.VisitVarStmt(s)
}

// BlockStmt ...
type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) {
	v.VisitBlockStmt(s)
}

// IfStmt ...
type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

func (s *IfStmt) Accept(v StmtVisitor) {
	v.VisitIfStmt(s)
}

// WhileStmt ...
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) Accept(v StmtVisitor) {
	v.VisitWhileStmt(s)
}

// FunctionStmt
type FunctionStmt struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Stmt
}

func (s *FunctionStmt) Accept(v StmtVisitor) {
	v.VisitFunctionStmt(s)
}

// ReturnStmt
type ReturnStmt struct {
	KeyWord *token.Token
	Value   Expr
}

func (s *ReturnStmt) Accept(v StmtVisitor) {
	v.VisitReturnStmt(s)
}

// ClassStmt
type ClassStmt struct {
	Name    *token.Token
	Methods []Stmt
}

func (s *ClassStmt) Accept(v StmtVisitor) {
	v.VisitClassStmt(s)
}
