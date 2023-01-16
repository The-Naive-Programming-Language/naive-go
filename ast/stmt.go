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
	_ Stmt = (*DeclStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*PrintStmt)(nil)
	_ Stmt = (*EmptyStmt)(nil)

	_ Stmt = (*Block)(nil)
)

type DeclStmt struct {
	Ident string
	Init  Expr
}

func (ds *DeclStmt) Accept(v Visitor) any {
	return v.VisitDeclStmt(ds)
}

func (ds *DeclStmt) String() string {
	return fmt.Sprintf("LET ident=%s init=%s", ds.Ident, ds.Init.String())
}

func (*DeclStmt) stmtNode() {}

type AssignStmt struct {
	Ident string
	Expr  Expr
}

func (as *AssignStmt) Accept(v Visitor) any {
	return v.VisitAssignStmt(as)
}

func (as *AssignStmt) String() string {
	return fmt.Sprintf("ASSIGN ident=%s expr=%s", as.Ident, as.Expr.String())
}

func (*AssignStmt) stmtNode() {}

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

type EmptyStmt struct{}

func (EmptyStmt) Accept(v Visitor) any {
	return v.VisitEmptyStmt(&EmptyStmt{})
}

func (EmptyStmt) String() string {
	return ";"
}

func (EmptyStmt) stmtNode() {}
