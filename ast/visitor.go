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
	VisitPrintStmt(stmt *PrintStmt) any
	VisitExpressionStmt(stmt *ExpressionStmt) any
	VisitVarStmt(stmt *VarStmt) any
	VisitBlockStmt(stmt *BlockStmt) any
	VisitIfStmt(stmt *IfStmt) any
	VisitWhileStmt(stmt *WhileStmt) any
	VisitFunctionStmt(*FunctionStmt) any
	VisitReturnStmt(*ReturnStmt) any
	VisitClassStmt(*ClassStmt) any
}
