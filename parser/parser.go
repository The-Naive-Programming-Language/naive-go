package parser

import (
	"fmt"
	"strings"

	"naive/ast"
	"naive/scanner"
	"naive/token"
)

type Parser struct {
	s scanner.Scanner

	kind token.Kind
	text string

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
	}
	if p.kind == token.KindPrint {
		return p.parsePrintStmt()
	}
	return p.parseExprStmt()
}

func (p *Parser) parseExprStmt() ast.Stmt {
	s := &ast.ExprStmt{
		Expr: p.parseExpr(),
	}
	p.consume(token.KindSemicolon)
	return s
}

func (p *Parser) parseExpr() (ans ast.Expr) {
	if p.matchAny(token.KindInt, token.KindFloat) {
		return p.parseArithExpr()
	}
	if p.matchAny(token.KindTrue, token.KindFalse) {
		return p.parseBoolExpr()
	}
	if p.match(token.KindLParen) {
		return p.parseGroupingExpr()
	}

	if p.match(token.KindString) {
		begin, end := 0, len(p.text)
		if strings.HasPrefix(p.text, "\"") {
			begin++
		}
		if strings.HasSuffix(p.text, "\"") && begin < end {
			end--
		}
		ans = ast.NewStringValue(p.text[begin:end])
	} else if p.match(token.KindChar) {
		ans = ast.NewCharValue(p.text)
	}
	p.discard()
	return ans
}

func (p *Parser) parseArithExpr() ast.Expr {
	return p.parseTerm()
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
	ans = p.parseArithAtom()
	for p.matchAny(token.KindMul, token.KindDiv, token.KindMod) {
		op := p.kind
		p.discard()
		rhs := p.parseArithAtom()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  op,
		}
	}
	return
}

func (p *Parser) parseArithAtom() (ans ast.Expr) {
	if p.match(token.KindLParen) {
		return p.parseGroupingExpr()
	}
	if p.match(token.KindInt) {
		ans = ast.NewIntegerValue(p.text)
	} else if p.match(token.KindFloat) {
		ans = ast.NewFloatValue(p.text)
	} else {
		panic("unreachable")
	}
	p.discard()
	return
}

func (p *Parser) parseBoolExpr() ast.Expr {
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
	ans = p.parseNotClause()
	for p.matchAny(token.KindAnd) {
		p.discard()
		rhs := p.parseNotClause()
		ans = &ast.BinaryExpr{
			Lhs: ans,
			Rhs: rhs,
			Op:  token.KindAnd,
		}
	}
	return
}

func (p *Parser) parseNotClause() (ans ast.Expr) {
	if p.match(token.KindLParen) {
		return p.parseGroupingExpr()
	}
	if p.match(token.KindNot) {
		return &ast.UnaryExpr{
			X:  p.parseNotClause(),
			Op: token.KindNot,
		}
	}
	if p.match(token.KindTrue) {
		ans = ast.True{}
	} else if p.match(token.KindFalse) {
		ans = ast.False{}
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
		Format: p.text,
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
