package ast

type Visitor interface {
	VisitIntegerValue(expr IntegerValue) any
	VisitFloatValue(expr FloatValue) any
	VisitStringValue(expr StringValue) any
	VisitCharValue(expr CharValue) any

	VisitTrue(expr True) any
	VisitFalse(expr False) any

	VisitBinaryExpr(expr *BinaryExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any

	VisitExprStmt(stmt *ExprStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
}
