package ast

// Node represents a node in AST's
type Node interface {
	Accept(v Visitor) any
}
