package ast

type Visitor interface {
	VisitIntegerValue(expr IntegerValue) any
	VisitFloatValue(expr FloatValue) any
	VisitStringValue(expr StringValue) any
	VisitCharValue(expr CharValue) any

	VisitTrue(expr True) any
	VisitFalse(expr False) any
	VisitNil(expr Nil) any

	VisitVariable(expr *Variable) any

	VisitBinaryExpr(expr *BinaryExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any

	VisitDeclStmt(stmt *DeclStmt) any
	VisitAssignStmt(stmt *AssignStmt) any
	VisitExprStmt(stmt *ExprStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
	VisitEmptyStmt(stmt *EmptyStmt) any
}
