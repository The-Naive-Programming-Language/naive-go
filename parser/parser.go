package parser

import (
	"fmt"

	"naive/ast"
	"naive/scanner"
	"naive/token"
)

type elem struct {
	kind token.Kind
	text string
}

type lookAheadStack []elem

func (s *lookAheadStack) push(kind token.Kind, text string) {
	*s = append(*s, elem{kind, text})
}

func (s *lookAheadStack) pop() (kind token.Kind, text string) {
	if s.empty() {
		panic("no staged tokens")
	}
	last := len(*s) - 1
	var e elem
	e, *s = (*s)[last], (*s)[:last]
	return e.kind, e.text
}

func (s lookAheadStack) empty() bool {
	return len(s) == 0
}

type Parser struct {
	s scanner.Scanner

	kind token.Kind
	text string

	prevKind token.Kind
	prevText string

	lookAhead lookAheadStack

	Statements []ast.Stmt
}

func New(file *token.File, src []byte) *Parser {
	p := &Parser{
		s: *scanner.New(file, src),
	}
	p.nextToken()
	return p
}

func (p *Parser) Parse() {
	for p.kind != token.KindEOF {
		if p.kind == token.KindComment {
			p.discard()
			continue
		}
		p.addStmt(p.parseStatement())
	}
}

func (p *Parser) addStmt(stmt ast.Stmt) {
	p.Statements = append(p.Statements, stmt)
}

func (p *Parser) parseStatement() ast.Stmt {
	if p.kind == token.KindSemicolon {
		p.discard()
		return &ast.EmptyStmt{}
	} else if p.kind == token.KindLet {
		return p.parseDeclStmt()
	} else if p.kind == token.KindPrint {
		return p.parsePrintStmt()
	} else if p.kind == token.KindIdent {
		// look ahead
		return p.branchAssignOrExpr()
	}
	return p.parseExprStmt()
}

func (p *Parser) parseDeclStmt() ast.Stmt {
	p.discard()
	if !p.match(token.KindIdent) {
		panic(fmt.Sprintf("incomplete let-statement, want an identifier but got %s", p.kind.String()))
	}
	ident := p.text
	p.discard()
	var init ast.Expr = ast.Nil{}
	if p.match(token.KindAssign) {
		p.discard()
		init = p.parseExpr()
	}
	p.consume(token.KindSemicolon)
	return &ast.DeclStmt{
		Ident: ident,
		Init:  init,
	}
}

func (p *Parser) branchAssignOrExpr() ast.Stmt {
	ident := p.text
	p.advance()
	if p.match(token.KindAssign) {
		return p.parseAssignStmt(ident)
	}
	p.goBack()
	return p.parseExprStmt()
}

func (p *Parser) parseAssignStmt(ident string) ast.Stmt {
	// skip '='
	p.discard()
	v := p.parseExpr()
	p.consume(token.KindSemicolon)
	return &ast.AssignStmt{
		Ident: ident,
		Expr:  v,
	}
}

func (p *Parser) parseExprStmt() ast.Stmt {
	s := &ast.ExprStmt{
		Expr: p.parseExpr(),
	}
	p.consume(token.KindSemicolon)
	return s
}

func (p *Parser) parseExpr() (ans ast.Expr) {
	return p.parseLogical()
}

func (p *Parser) parseLogical() ast.Expr {
	return p.parseOrClause()
}

func (p *Parser) parseOrClause() (ans ast.Expr) {
	ans = p.parseAndClause()
	for p.matchAny(token.KindOr) {
		p.discard()
		rhs := p.parseAndClause()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  token.KindOr,
		}
	}
	return
}

func (p *Parser) parseAndClause() (ans ast.Expr) {
	ans = p.parseRelational()
	for p.matchAny(token.KindAnd) {
		p.discard()
		rhs := p.parseRelational()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  token.KindAnd,
		}
	}
	return
}

func (p *Parser) parseRelational() (ans ast.Expr) {
	ans = p.parseTerm()
	for p.matchAny(token.KindEq, token.KindNe, token.KindLt, token.KindGt, token.KindLe, token.KindGe) {
		op := p.kind
		p.discard()
		rhs := p.parseTerm()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  op,
		}
	}
	return
}

func (p *Parser) parseTerm() (ans ast.Expr) {
	ans = p.parseFactor()
	for p.matchAny(token.KindAdd, token.KindSub) {
		op := p.kind
		p.discard()
		rhs := p.parseFactor()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  op,
		}
	}
	return
}

func (p *Parser) parseFactor() (ans ast.Expr) {
	ans = p.parseUnary()
	for p.matchAny(token.KindMul, token.KindDiv, token.KindMod) {
		op := p.kind
		p.discard()
		rhs := p.parseUnary()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  op,
		}
	}
	return
}

func (p *Parser) parseUnary() (ans ast.Expr) {
	if !p.matchAny(token.KindNot, token.KindSub) {
		return p.parsePrimary()
	}
	op := p.kind
	p.discard()
	return &ast.UnaryExpr{
		X:  p.parseUnary(),
		Op: op,
	}
}

func (p *Parser) parsePrimary() (ans ast.Expr) {
	if p.match(token.KindLParen) {
		return p.parseGroupingExpr()
	}

	if p.match(token.KindInt) {
		ans = ast.NewIntegerValue(p.text)
	} else if p.match(token.KindFloat) {
		ans = ast.NewFloatValue(p.text)
	} else if p.match(token.KindChar) {
		ans = ast.NewCharValue(p.text)
	} else if p.match(token.KindString) {
		ans = ast.NewStringValue(p.text)
	} else if p.match(token.KindTrue) {
		ans = ast.True{}
	} else if p.match(token.KindFalse) {
		ans = ast.False{}
	} else if p.match(token.KindIdent) {
		ans = &ast.Variable{
			Ident: p.text,
		}
	} else {
		panic("unreachable")
	}
	p.discard()
	return
}

func (p *Parser) parseGroupingExpr() (ans *ast.GroupingExpr) {
	p.consume(token.KindLParen)
	e := p.parseExpr()
	p.consume(token.KindRParen)
	return &ast.GroupingExpr{
		Expr: e,
	}
}

func (p *Parser) parsePrintStmt() ast.Stmt {
	p.discard()
	p.consume(token.KindLParen)
	if p.kind != token.KindString {
		panic("expect token STRING, actual: " + p.kind.String())
	}
	s := &ast.PrintStmt{
		Format: p.text[1 : len(p.text)-1],
	}
	p.discard()
	for p.match(token.KindComma) {
		p.discard()
		s.Args = append(s.Args, p.parseExpr())
	}
	p.consumeMany(token.KindRParen, token.KindSemicolon)
	return s
}

func (p *Parser) nextToken() {
	if !p.lookAhead.empty() {
		p.kind, p.text = p.lookAhead.pop()
		return
	}
	_, p.kind, p.text = p.s.Scan()
}

func (p *Parser) match(kind token.Kind) bool {
	return p.kind == kind
}

func (p *Parser) matchAny(kinds ...token.Kind) bool {
	for _, k := range kinds {
		if p.kind == k {
			return true
		}
	}
	return false
}

func (p *Parser) discard() {
	p.nextToken()
}

func (p *Parser) consume(kind token.Kind) {
	if p.kind != kind {
		panic(fmt.Sprintf("expect token: %s, actual: %s", kind, p.kind))
	}
	p.discard()
}

func (p *Parser) consumeMany(kinds ...token.Kind) {
	for _, k := range kinds {
		p.consume(k)
	}
}

func (p *Parser) advance() {
	p.prevKind, p.prevText = p.kind, p.text
	p.nextToken()
}

func (p *Parser) goBack() {
	p.lookAhead.push(p.kind, p.text)
	p.kind, p.text = p.prevKind, p.prevText
}
