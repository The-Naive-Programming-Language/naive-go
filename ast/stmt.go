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
	_ Stmt = (*LetStmt)(nil)
	_ Stmt = (*IfElseStmt)(nil)
	_ Stmt = (*WhileStmt)(nil)
	_ Stmt = (*FnStmt)(nil)
	_ Stmt = (*ReturnStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*EmptyStmt)(nil)

	_ Stmt = (*Block)(nil)
)

type LetStmt struct {
	Ident string
	Init  Expr
}

func (ds *LetStmt) Accept(v Visitor) any {
	return v.VisitLetStmt(ds)
}

func (ds *LetStmt) String() string {
	return fmt.Sprintf("LET ident=%s init=%s", ds.Ident, ds.Init.String())
}

func (*LetStmt) stmtNode() {}

type FnStmt struct {
	Ident  string
	Params []string
	Body   Stmt
}

func (fs *FnStmt) Accept(v Visitor) any {
	return v.VisitFnStmt(fs)
}

func (fs *FnStmt) String() string {
	return fmt.Sprintf("fn %s (%s) %s", fs.Ident, strings.Join(fs.Params, ", "), fs.Body)
}

func (*FnStmt) stmtNode() {}

type ReturnStmt struct {
	RetVal Expr
}

func (rs *ReturnStmt) Accept(v Visitor) any {
	return v.VisitReturnStmt(rs)
}

func (rs *ReturnStmt) String() string {
	return fmt.Sprintf("return %s;", rs.RetVal)
}

func (*ReturnStmt) stmtNode() {}

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

type IfElseStmt struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (ie *IfElseStmt) Accept(v Visitor) any {
	return v.VisitIfElseStmt(ie)
}

func (ie *IfElseStmt) String() string {
	return fmt.Sprintf("if %s %s %s", ie.Cond, ie.Then, ie.Else)
}

func (ie *IfElseStmt) stmtNode() {}
func (ie *IfElseStmt) exprNode() {}

type WhileStmt struct {
	Cond Expr
	Body Stmt
}

func (ws *WhileStmt) Accept(v Visitor) any {
	return v.VisitWhileStmt(ws)
}

func (ws *WhileStmt) String() string {
	return fmt.Sprintf("while %s %s", ws.Cond, ws.Body)
}

func (ws *WhileStmt) stmtNode() {}

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

type EmptyStmt struct{}

func (EmptyStmt) Accept(v Visitor) any {
	return v.VisitEmptyStmt(&EmptyStmt{})
}

func (EmptyStmt) String() string {
	return ";"
}

func (EmptyStmt) stmtNode() {}
