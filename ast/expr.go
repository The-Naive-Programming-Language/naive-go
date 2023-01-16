package ast

import (
	"math/big"
	"strings"
	"unicode/utf8"

	"naive/token"
)

// Expr represents an expression. It is a kind of Node.
type Expr interface {
	Node
	String() string
	// exprNode is a marker interface.
	exprNode()
}

var (
	_ Expr = (*IntegerValue)(nil)
	_ Expr = (*FloatValue)(nil)
	_ Expr = (*CharValue)(nil)
	_ Expr = (*StringValue)(nil)

	_ Expr = (*True)(nil)
	_ Expr = (*False)(nil)

	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*GroupingExpr)(nil)

	_ Expr = (*Block)(nil)
)

type IntegerValue struct {
	Value *big.Int
}

func NewIntegerValue(text string) *IntegerValue {
	ans := &IntegerValue{}
	if strings.HasPrefix(text, "0b") {
		ans.Value, _ = new(big.Int).SetString(text[2:], 2)
	} else if strings.HasPrefix(text, "0o") {
		ans.Value, _ = new(big.Int).SetString(text[2:], 8)
	} else if strings.HasPrefix(text, "0x") {
		ans.Value, _ = new(big.Int).SetString(text[2:], 16)
	} else {
		ans.Value, _ = new(big.Int).SetString(text, 10)
	}
	return ans
}

func (iv IntegerValue) Accept(v Visitor) any {
	return v.VisitIntegerValue(iv)
}

func (iv IntegerValue) String() string {
	return iv.Value.String()
}

func (IntegerValue) exprNode() {}

type FloatValue struct {
	Value *big.Float
}

func NewFloatValue(text string) *FloatValue {
	ans := &FloatValue{}
	ans.Value, _ = new(big.Float).SetString(text)
	return ans
}

func (fv FloatValue) Accept(v Visitor) any {
	return v.VisitFloatValue(fv)
}

func (fv FloatValue) String() string {
	return fv.Value.String()
}

func (FloatValue) exprNode() {}

type CharValue struct {
	Value rune
}

func NewCharValue(text string) CharValue {
	r, _ := utf8.DecodeRuneInString(text[1:])
	return CharValue{
		Value: r,
	}
}

func (cv CharValue) Accept(v Visitor) any {
	return v.VisitCharValue(cv)
}

func (cv CharValue) String() string {
	return "'" + string(cv.Value) + "'"
}

func (CharValue) exprNode() {}

type StringValue struct {
	Value string
}

func NewStringValue(text string) StringValue {
	return StringValue{
		Value: text[1 : len(text)-1],
	}
}

func (sv StringValue) Accept(v Visitor) any {
	return v.VisitStringValue(sv)
}

func (sv StringValue) String() string {
	return "\"" + sv.Value + "\""
}

func (StringValue) exprNode() {}

// True represents a node of boolean value true.
type True struct{}

func (t True) Accept(v Visitor) any {
	return v.VisitTrue(t)
}

func (True) String() string {
	return "true"
}

func (True) exprNode() {}

// False represents a node of boolean value false.
type False struct{}

func (f False) Accept(v Visitor) any {
	return v.VisitFalse(f)
}

func (False) String() string {
	return "false"
}

func (False) exprNode() {}

type Nil struct{}

func (n Nil) Accept(v Visitor) any {
	return v.VisitNil(n)
}

func (Nil) String() string {
	return "nil"
}

func (Nil) exprNode() {}

// UnaryExpr represents a node of unary expression.
type UnaryExpr struct {
	X  Expr
	Op token.Kind
}

func (ue *UnaryExpr) Accept(v Visitor) any {
	return v.VisitUnaryExpr(ue)
}

func (ue *UnaryExpr) String() string {
	return ue.Op.String() + " " + ue.X.String()
}

func (*UnaryExpr) exprNode() {}

// BinaryExpr represents a node of binary expression.
type BinaryExpr struct {
	Lhs, Rhs Expr
	Op       token.Kind
}

func (be *BinaryExpr) Accept(v Visitor) any {
	return v.VisitBinaryExpr(be)
}

func (be *BinaryExpr) String() string {
	return be.Lhs.String() + " " + be.Op.String() + " " + be.Rhs.String()
}

func (*BinaryExpr) exprNode() {}

type GroupingExpr struct {
	Expr Expr
}

func (ge *GroupingExpr) Accept(v Visitor) any {
	return v.VisitGroupingExpr(ge)
}

func (ge *GroupingExpr) String() string {
	return "(" + ge.Expr.String() + ")"
}

func (*GroupingExpr) exprNode() {}

type Variable struct {
	Ident string
}

func (ve *Variable) Accept(v Visitor) any {
	return v.VisitVariable(ve)
}

func (ve *Variable) String() string {
	return "VAR " + ve.Ident
}

func (*Variable) exprNode() {}

type Block struct {
	Statements []Stmt
}

func (blk *Block) Accept(v Visitor) any {
	return v.VisitBlock(blk)
}

func (blk *Block) String() string {
	ss := make([]string, 0, len(blk.Statements))
	for _, stmt := range blk.Statements {
		ss = append(ss, stmt.String())
	}
	return "{" + strings.Join(ss, "; ") + "}"
}

func (*Block) exprNode() {}

func (*Block) stmtNode() {}
