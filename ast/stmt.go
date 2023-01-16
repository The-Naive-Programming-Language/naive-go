package ast

import (
	"fmt"
	"strings"
)

type Stmt interface {
	Node
	String() string
	stmtNode()
}

var (
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*PrintStmt)(nil)
)

type ExprStmt struct {
	Expr Expr
}

func (es *ExprStmt) Accept(v Visitor) any {
	return v.VisitExprStmt(es)
}

func (es *ExprStmt) String() string {
	return es.Expr.String() + ";"
}

func (ExprStmt) stmtNode() {}

type ExprList []Expr

func (el ExprList) String() string {
	ss := make([]string, 0, len(el))
	for _, e := range el {
		ss = append(ss, e.String())
	}
	return strings.Join(ss, ", ")
}

func (el ExprList) Empty() bool {
	return len(el) == 0
}

type PrintStmt struct {
	Format string
	Args   ExprList
}

func (ps *PrintStmt) Accept(v Visitor) any {
	return v.VisitPrintStmt(ps)
}

func (ps *PrintStmt) String() string {
	if ps.Args.Empty() {
		return fmt.Sprintf("print(%s)", ps.Format)
	}
	return fmt.Sprintf("PRINT format=%q args=(%s)", ps.Format, ps.Args.String())
}

func (PrintStmt) stmtNode() {}
