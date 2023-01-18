package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"naive/ast"
	"naive/token"
)

func TestParser_Parse(t *testing.T) {
}

func TestParser_parseStmt(t *testing.T) {
	Convey("arith", t, func() {
		p := New(nil, []byte("1 + 2 * 3 - 4;"))

		s := p.parseStatement()

		So(s, ShouldNotBeNil)
		es, ok := s.(*ast.ExprStmt)
		So(ok, ShouldBeTrue)
		be, ok := es.Expr.(*ast.BinaryExpr)
		So(ok, ShouldBeTrue)
		So(be.Op, ShouldEqual, token.KindSub)

		L, ok := be.Lhs.(*ast.BinaryExpr)
		So(ok, ShouldBeTrue)
		So(L.Op, ShouldEqual, token.KindAdd)
		LR, ok := L.Rhs.(*ast.BinaryExpr)
		So(ok, ShouldBeTrue)
		So(LR.Op, ShouldEqual, token.KindMul)
	})
}

func TestParser_parseExpr(t *testing.T) {
	Convey("arith", t, func() {
		Convey("int atom", func() {
			p := New(nil, []byte("1 + 2 + 3"))
			_ = p.parseExpr()

		})
	})
}

func TestParser_parseStatement(t *testing.T) {
	Convey("consecutive semicolons", t, func() {
		p := New(nil, []byte(";;;;"))
		s := p.parseStatement()
		es, ok := s.(*ast.EmptyStmt)
		So(ok, ShouldBeTrue)
		So(es, ShouldNotBeNil)
	})
}

func TestParser_parseDeclStmt(t *testing.T) {
	Convey("valid", t, func() {
		Convey("w/o init", func() {
			p := New(nil, []byte("let a;"))
			s := p.parseDeclStmt()
			ds, ok := s.(*ast.LetStmt)
			So(ok, ShouldBeTrue)
			So(ds.Ident, ShouldEqual, "a")
			So(ds.Init, ShouldHaveSameTypeAs, ast.Nil{})
		})
	})
}
