package interpreter

import (
	"math/big"

	"naive/ast"
	"naive/parser"
	"naive/token"
)

var _ ast.Visitor = (*Interpreter)(nil)

type Env struct {
	bindings  map[string]any
	enclosing *Env
}

func newGlobalEnv() *Env {
	return newLocalEnv(nil)
}

func newLocalEnv(enclosing *Env) *Env {
	return &Env{
		bindings:  make(map[string]any),
		enclosing: enclosing,
	}
}

func (e *Env) Define(k string, v any) (shadow bool) {
	_, shadow = e.bindings[k]
	e.bindings[k] = v
	return
}

func (e *Env) Lookup(k string) (v any, present bool) {
	curr := e
	for curr != nil {
		v, present = curr.bindings[k]
		if present {
			return v, present
		}
		curr = curr.enclosing
	}
	return nil, false
}

func (e *Env) Assign(k string, v any) (done bool) {
	curr := e
	for curr != nil {
		_, present := curr.bindings[k]
		if present {
			curr.bindings[k] = v
			return true
		}
		curr = curr.enclosing
	}
	return false
}

type Interpreter struct {
	P *parser.Parser

	env *Env
}

func New(filename string, src []byte) *Interpreter {
	i := &Interpreter{
		P: parser.New(
			token.NewFile(filename),
			src,
		),
		env: newGlobalEnv(),
	}
	i.setupBuiltins()
	return i
}

func (i *Interpreter) setupBuiltins() {
	i.env.Define("print", BuiltinPrint{})
	i.env.Define("println", BuiltinPrintLn{})
	i.env.Define("format", BuiltinFormat{})
	i.env.Define("getline", BuiltinGetLine{})
}

func Default() *Interpreter {
	return New("", nil)
}

func (i *Interpreter) Interpret() {
	i.P.Parse()
	// fmt.Println(i.P.Statements)
	for _, stmt := range i.P.Statements {
		stmt.Accept(i)
	}
}

func (Interpreter) VisitIntegerValue(expr ast.IntegerValue) any {
	return expr.Value
}

func (Interpreter) VisitFloatValue(expr ast.FloatValue) any {
	return expr.Value
}

func (Interpreter) VisitCharValue(expr ast.CharValue) any {
	return expr.Value
}

func (Interpreter) VisitStringValue(expr ast.StringValue) any {
	return expr.Value
}

func (Interpreter) VisitTrue(expr ast.True) any {
	return true
}

func (Interpreter) VisitFalse(expr ast.False) any {
	return false
}

func (*Interpreter) VisitNil(_ ast.Nil) any {
	return nil
}

func (i *Interpreter) VisitVariable(expr *ast.Variable) any {
	v, ok := i.env.Lookup(expr.Ident)
	if !ok {
		panic("undefined variable: " + expr.Ident)
	}
	return v
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	lhs, rhs := expr.Lhs.Accept(i), expr.Rhs.Accept(i)
	switch expr.Op {
	case token.KindAdd:
		return doAdd(lhs, rhs)
	case token.KindSub:
		return doSub(lhs, rhs)
	case token.KindMul:
		return doMul(lhs, rhs)
	case token.KindDiv:
		return doDiv(lhs, rhs)
	case token.KindMod:
		return doMod(lhs, rhs)

	case token.KindEq:
		return doEq(lhs, rhs)
	case token.KindNe:
		return doNe(lhs, rhs)
	case token.KindGt:
		return doGt(lhs, rhs)
	case token.KindGe:
		return doGe(lhs, rhs)
	case token.KindLt:
		return doLt(lhs, rhs)
	case token.KindLe:
		return doLe(lhs, rhs)

	case token.KindAnd:
		return doLogicalAnd(lhs, rhs)
	case token.KindOr:
		return doLogicalOr(lhs, rhs)

	default:
		panic("unreachable")
	}
}

func doAdd(lhs, rhs any) any {
	if x, y := isFloat(lhs), isFloat(rhs); x || y {
		return doFloatAdd(toFloat(lhs), toFloat(rhs))
	}
	return doIntegerAdd(toInt(lhs), toInt(rhs))
}

func doFloatAdd(lhs *big.Float, rhs *big.Float) *big.Float {
	ans := big.NewFloat(0)
	return ans.Add(lhs, rhs)
}

func doIntegerAdd(lhs *big.Int, rhs *big.Int) *big.Int {
	ans := big.NewInt(0)
	return ans.Add(lhs, rhs)
}

func doSub(lhs, rhs any) any {
	if x, y := isFloat(lhs), isFloat(rhs); x || y {
		return doFloatSub(toFloat(lhs), toFloat(rhs))
	}
	return doIntegerSub(toInt(lhs), toInt(rhs))
}

func doFloatSub(lhs *big.Float, rhs *big.Float) *big.Float {
	ans := big.NewFloat(0)
	return ans.Sub(lhs, rhs)
}

func doIntegerSub(lhs *big.Int, rhs *big.Int) *big.Int {
	ans := big.NewInt(0)
	return ans.Sub(lhs, rhs)
}

func doMul(lhs, rhs any) any {
	if x, y := isFloat(lhs), isFloat(rhs); x || y {
		return doFloatMul(toFloat(lhs), toFloat(rhs))
	}
	return doIntegerMul(toInt(lhs), toInt(rhs))
}

func doFloatMul(lhs *big.Float, rhs *big.Float) *big.Float {
	ans := big.NewFloat(0)
	return ans.Mul(lhs, rhs)
}

func doIntegerMul(lhs *big.Int, rhs *big.Int) *big.Int {
	ans := big.NewInt(0)
	return ans.Mul(lhs, rhs)
}

func doDiv(lhs, rhs any) any {
	if x, y := isFloat(lhs), isFloat(rhs); x || y {
		return doFloatDiv(toFloat(lhs), toFloat(rhs))
	}
	return doIntegerDiv(toInt(lhs), toInt(rhs))
}

func doFloatDiv(lhs *big.Float, rhs *big.Float) *big.Float {
	ans := big.NewFloat(0)
	return ans.Quo(lhs, rhs)
}

func doIntegerDiv(lhs *big.Int, rhs *big.Int) *big.Int {
	ans := big.NewInt(0)
	return ans.Div(lhs, rhs)
}

func doMod(lhs, rhs any) any {
	if x, y := isFloat(lhs), isFloat(rhs); x || y {
		return doFloatMod(toFloat(lhs), toFloat(rhs))
	}
	return doIntegerMod(toInt(lhs), toInt(rhs))
}

func doFloatMod(lhs *big.Float, rhs *big.Float) *big.Float {
	panic("unimplemented")
}

func doIntegerMod(lhs *big.Int, rhs *big.Int) *big.Int {
	ans := big.NewInt(0)
	return ans.Mod(lhs, rhs)
}

func isFloat(x any) bool {
	_, ok := x.(*big.Float)
	return ok
}

func toFloat(x0 any) *big.Float {
	switch x := x0.(type) {
	case *big.Float:
		return x
	case *big.Int:
		ans, _ := new(big.Float).SetString(x.String())
		return ans
	case bool:
		if x {
			return big.NewFloat(1)
		} else {
			return big.NewFloat(0)
		}
	default:
		panic("unreachable")
	}
}

func toInt(x0 any) *big.Int {
	switch x := x0.(type) {
	case *big.Int:
		return x
	case bool:
		if x {
			return big.NewInt(1)
		} else {
			return big.NewInt(0)
		}
	default:
		panic("unreachable")
	}
}

func doEq(lhs, rhs any) bool {
	v, _ := branchByNumberType(lhs, rhs, doFloatEq, doIntegerEq).(bool)
	return v
}

func doFloatEq(lhs, rhs *big.Float) any {
	return lhs.Cmp(rhs) == 0
}

func doIntegerEq(lhs, rhs *big.Int) any {
	return lhs.Cmp(rhs) == 0
}

func doNe(lhs, rhs any) any {
	return branchByNumberType(lhs, rhs, doFloatNe, doIntegerNe)
}

func doFloatNe(lhs, rhs *big.Float) any {
	return lhs.Cmp(rhs) != 0
}

func doIntegerNe(lhs, rhs *big.Int) any {
	return lhs.Cmp(rhs) != 0
}

func doGt(lhs, rhs any) any {
	return branchByNumberType(lhs, rhs, func(lhs, rhs *big.Float) any {
		return lhs.Cmp(rhs) > 0
	}, func(lhs, rhs *big.Int) any {
		return lhs.Cmp(rhs) > 0
	})
}

func doGe(lhs, rhs any) any {
	return branchByNumberType(lhs, rhs, func(lhs, rhs *big.Float) any {
		return lhs.Cmp(rhs) >= 0
	}, func(lhs, rhs *big.Int) any {
		return lhs.Cmp(rhs) >= 0
	})
}

func doLt(lhs, rhs any) any {
	return branchByNumberType(lhs, rhs, func(lhs, rhs *big.Float) any {
		return lhs.Cmp(rhs) < 0
	}, func(lhs, rhs *big.Int) any {
		return lhs.Cmp(rhs) < 0
	})
}

func doLe(lhs, rhs any) any {
	return branchByNumberType(lhs, rhs, func(lhs, rhs *big.Float) any {
		return lhs.Cmp(rhs) <= 0
	}, func(lhs, rhs *big.Int) any {
		return lhs.Cmp(rhs) <= 0
	})
}

type (
	FloatOp func(*big.Float, *big.Float) any
	IntOp   func(*big.Int, *big.Int) any
)

func branchByNumberType(lhs, rhs any, fo FloatOp, io IntOp) any {
	if isFloat(lhs) || isFloat(rhs) {
		return fo(toFloat(lhs), toFloat(rhs))
	}
	return io(toInt(lhs), toInt(rhs))
}

func doLogicalAnd(lhs, rhs any) any {
	if isTruthy(lhs) {
		return rhs
	}
	return lhs
}

func doLogicalOr(lhs, rhs any) any {
	if isTruthy(lhs) {
		return lhs
	}
	return rhs
}

func isTruthy(x any) bool {
	return !isFalsy(x)
}

func isFalsy(x0 any) bool {
	switch x := x0.(type) {
	case bool:
		return !x
	case nil:
		return true
	}
	return false
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	x := expr.X.Accept(i)
	switch expr.Op {
	case token.KindSub:
		return doNeg(x)
	case token.KindNot:
		return doLogicalNot(x)
	default:
		panic("unreachable")
	}
}

func doNeg(x any) any {
	if isFloat(x) {
		return doFloatNeg(x)
	}
	return doIntegerNeg(x)
}

func doFloatNeg(x0 any) *big.Float {
	switch x := x0.(type) {
	case *big.Float:
		return big.NewFloat(0).Neg(x)
	default:
		panic("unreachable")
	}
}

func doIntegerNeg(x0 any) *big.Int {
	switch x := x0.(type) {
	case *big.Int:
		return big.NewInt(0).Neg(x)
	case bool:
		if x {
			return big.NewInt(-1)
		} else {
			return big.NewInt(0)
		}
	default:
		panic("unreachable")
	}
}

func doLogicalNot(x any) any {
	return isFalsy(x)
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return expr.Expr.Accept(i)
}

func (i *Interpreter) VisitLetStmt(stmt *ast.LetStmt) any {
	init := stmt.Init.Accept(i)
	i.env.Define(stmt.Ident, init)
	return nil
}

func (i *Interpreter) VisitAssignStmt(stmt *ast.AssignStmt) any {
	v := stmt.Expr.Accept(i)
	if !i.env.Assign(stmt.Ident, v) {
		panic("assignment to undefined variable " + stmt.Ident)
	}
	return nil
}

func (i *Interpreter) VisitIfElseStmt(stmt *ast.IfElseStmt) any {
	if isTruthy(stmt.Cond.Accept(i)) {
		return stmt.Then.Accept(i)
	}
	return stmt.Else.Accept(i)
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) any {
	for isTruthy(stmt.Cond.Accept(i)) {
		stmt.Body.Accept(i)
	}
	return nil
}

func (i *Interpreter) VisitFnStmt(stmt *ast.FnStmt) any {
	i.env.Define(stmt.Ident, &Func{
		Name:   stmt.Ident,
		Params: stmt.Params,
		Body:   stmt.Body,
		Env:    i.env,
	})
	i.env = newLocalEnv(i.env)
	return nil
}

type Return struct {
	RetVal any
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	ret := stmt.RetVal.Accept(i)
	// TODO: better solutions?
	panic(&Return{RetVal: ret})
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) any {
	args := make([]any, 0, len(expr.Args))
	for _, a := range expr.Args {
		args = append(args, a.Accept(i))
	}
	v, ok := i.env.Lookup(expr.Callee)
	if !ok {
		panic("undefined callable object '" + expr.Callee + "'")
	}
	f, ok := v.(Callable)
	if !ok {
		panic("calling non-callable object")
	}
	return f.Call(args, i)
}

func (i *Interpreter) VisitLambda(expr *ast.Lambda) any {
	ans := &Func{
		Name:   "<anonymous>",
		Params: expr.Params,
		Body:   expr.Body,
		Env:    i.env,
	}
	i.env = newLocalEnv(i.env)
	return ans
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) any {
	stmt.Expr.Accept(i)
	return nil
}

func (i *Interpreter) VisitEmptyStmt(stmt *ast.EmptyStmt) any {
	return nil
}

func (i *Interpreter) VisitBlock(blk *ast.Block) any {
	outer := i.env
	i.env = newLocalEnv(outer)
	defer func() {
		i.env = outer
	}()
	for _, stmt := range blk.Statements {
		stmt.Accept(i)
	}
	return nil
}
