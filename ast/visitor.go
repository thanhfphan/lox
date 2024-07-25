package ast

type ExprVisitor interface {
	VisitLiteralExpr(*LiteralExpr) any
	VisitGroupingExpr(*GroupingExpr) any
	VisitUnaryExpr(*UnaryExpr) any
	VisitBinaryExpr(*BinaryExpr) any
	VisitVariableExpr(*VariableExpr) any
	VisitAssignExpr(*AssignExpr) any
	VisitLogicalExpr(*LogicalExpr) any
	VisitCallExpr(*CallExpr) any
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitVarStmt(stmt *VarStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitWhileStmt(stmt *WhileStmt)
	VisitFunctionStmt(*FunctionStmt) any
	VisitReturnStmt(*ReturnStmt) any
}
