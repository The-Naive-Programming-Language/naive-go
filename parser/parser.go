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
	} else if p.kind == token.KindIdent {
		// look ahead
		return p.branchAssignOrExpr()
	} else if p.kind == token.KindLBrace {
		return p.parseBlock()
	} else if p.kind == token.KindIf {
		return p.parseIfElse()
	} else if p.kind == token.KindElse {
		panic("dangling else")
	} else if p.kind == token.KindWhile {
		return p.parseWhile()
	} else if p.kind == token.KindFn {
		return p.branchNamedFuncOrLambda()
	} else if p.kind == token.KindReturn {
		return p.parseReturn()
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
	return &ast.LetStmt{
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

func (p *Parser) parseBlock() ast.Stmt {
	p.skipComments()
	p.consume(token.KindLBrace)
	p.skipComments()
	blk := &ast.Block{}
	for !p.match(token.KindRBrace) {
		blk.Statements = append(blk.Statements, p.parseStatement())
		p.skipComments()
	}
	p.consume(token.KindRBrace)
	return blk
}

func (p *Parser) parseIfElse() ast.Stmt {
	p.discard()
	cond := p.parseExpr()
	thenArm := p.parseBlock()
	var elseArm ast.Stmt = ast.EmptyStmt{}
	if p.match(token.KindElse) {
		p.discard()
		if p.match(token.KindIf) {
			elseArm = p.parseIfElse()
		} else if p.match(token.KindLBrace) {
			elseArm = p.parseBlock()
		} else {
			panic("incomplete else")
		}
	}
	return &ast.IfElseStmt{
		Cond: cond,
		Then: thenArm,
		Else: elseArm,
	}
}

func (p *Parser) parseWhile() ast.Stmt {
	p.discard()
	cond := p.parseExpr()
	body := p.parseBlock()
	return &ast.WhileStmt{
		Cond: cond,
		Body: body,
	}
}

func (p *Parser) branchNamedFuncOrLambda() ast.Stmt {
	p.advance()
	if p.match(token.KindIdent) {
		p.goBack()
		return p.parseFunction()
	}
	p.goBack()
	return p.parseExprStmt()
}

func (p *Parser) parseFunction() ast.Stmt {
	p.discard()
	if !p.match(token.KindIdent) {
		panic(fmt.Sprintf("when parsing function definition: want %s, got %s", token.KindIdent, p.kind))
	}
	name := p.text
	p.discard()
	p.consume(token.KindLParen)
	params := p.parseIdentList()
	p.consume(token.KindRParen)
	body := p.parseBlock()
	return &ast.FnStmt{
		Ident:  name,
		Params: params,
		Body:   body,
	}
}

func (p *Parser) parseIdentList() (ans []string) {
	if p.match(token.KindRParen) {
		return
	}
	if !p.match(token.KindIdent) {
		panic(fmt.Sprintf("when parsing identifier list: want %s, got %s", token.KindIdent, p.kind))
	}
	ans = append(ans, p.text)
	p.discard()
	for !p.match(token.KindRParen) {
		p.consume(token.KindComma)
		if !p.match(token.KindIdent) {
			panic(fmt.Sprintf("when parsing identifier list: want %s, got %s", token.KindIdent, p.kind))
		}
		ans = append(ans, p.text)
		p.discard()
	}
	return
}

func (p *Parser) parseReturn() ast.Stmt {
	p.discard()
	ret := p.parseExpr()
	p.consume(token.KindSemicolon)
	return &ast.ReturnStmt{
		RetVal: ret,
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
	return p.parseLambda()
}

func (p *Parser) parseLambda() (ans ast.Expr) {
	if !p.match(token.KindFn) {
		return p.parseLogical()
	}
	p.consume(token.KindFn)
	p.consume(token.KindLParen)
	params := p.parseIdentList()
	p.consume(token.KindRParen)
	var body ast.Stmt
	if p.match(token.KindLtRArrow) {
		p.discard()
		e := p.parseExpr()
		body = &ast.Block{
			Statements: []ast.Stmt{
				&ast.ReturnStmt{
					RetVal: e,
				},
			},
		}
	} else {
		body = p.parseBlock()
	}
	return &ast.Lambda{
		Params: params,
		Body:   body,
	}
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
		return p.parseCall()
	}
	op := p.kind
	p.discard()
	return &ast.UnaryExpr{
		X:  p.parseUnary(),
		Op: op,
	}
}

func (p *Parser) parseCall() ast.Expr {
	if !p.match(token.KindIdent) {
		return p.parsePrimary()
	}
	ident := p.text
	p.advance()
	if p.match(token.KindLParen) {
		p.discard()
		args := p.parseExprList()
		p.consume(token.KindRParen)
		return &ast.CallExpr{
			Callee: ident,
			Args:   args,
		}
	}
	// If current token is not '(', leave it unchanged and return a Variable.
	// NOTE: although assignment statements also start with IDENT's, call sites
	// of parseCall can eliminate the possibility.
	return &ast.Variable{
		Ident: ident,
	}
}

func (p *Parser) parseExprList() (ans []ast.Expr) {
	if p.match(token.KindRParen) {
		return
	}
	ans = append(ans, p.parseExpr())
	for !p.match(token.KindRParen) {
		p.consume(token.KindComma)
		ans = append(ans, p.parseExpr())
	}
	return
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

func (p *Parser) skipComments() {
	for p.match(token.KindComment) {
		p.discard()
	}
}
