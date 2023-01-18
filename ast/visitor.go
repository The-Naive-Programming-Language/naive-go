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
	VisitCallExpr(expr *CallExpr) any
	VisitLambda(expr *Lambda) any

	VisitLetStmt(stmt *LetStmt) any
	VisitAssignStmt(stmt *AssignStmt) any
	VisitIfElseStmt(stmt *IfElseStmt) any
	VisitWhileStmt(stmt *WhileStmt) any
	VisitFnStmt(stmt *FnStmt) any
	VisitReturnStmt(stmt *ReturnStmt) any
	VisitExprStmt(stmt *ExprStmt) any
	VisitEmptyStmt(stmt *EmptyStmt) any

	VisitBlock(blk *Block) any
}
