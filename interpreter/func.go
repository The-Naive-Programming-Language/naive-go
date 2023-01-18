package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"naive/ast"
)

type Callable interface {
	Call(args []any, i *Interpreter) any
}

// Func is the runtime representation of named function
type Func struct {
	Name   string
	Params []string
	Body   ast.Stmt
	Env    *Env
}

func (f *Func) Call(args []any, i *Interpreter) (ans any) {
	if len(args) != len(f.Params) {
		panic(fmt.Sprintf(
			"function %s takes %d positional arguments but %d are provided",
			f.Name, len(f.Params), len(args)))
	}

	old := i.env
	i.env = newLocalEnv(f.Env)
	for j := range args {
		k, v := f.Params[j], args[j]
		i.env.Define(k, v)
	}
	defer func() {
		i.env = old
	}()

	defer func() {
		r0 := recover()
		if r0 == nil {
			return
		}
		switch r := r0.(type) {
		case *Return:
			ans = r.RetVal
		default:
			panic(r0)
		}
	}()
	f.Body.Accept(i)

	return
}

type BuiltinPrint struct{}

func (BuiltinPrint) Call(args []any, i *Interpreter) any {
	fmt.Print(args...)
	return nil
}

type BuiltinPrintLn struct{}

func (BuiltinPrintLn) Call(args []any, i *Interpreter) any {
	fmt.Println(args...)
	return nil
}

type BuiltinFormat struct{}

func (BuiltinFormat) Call(args []any, i *Interpreter) any {
	if len(args) < 1 {
		panic("function format takes at least 1 argument, but none provided")
	}
	f, ok := args[0].(string)
	if !ok {
		panic("type mismatch: 1st argument of function format shall be of type 'String'")
	}
	f = strings.ReplaceAll(f, "{}", "%v")
	return fmt.Sprintf(f, args...)
}

type BuiltinGetLine struct{}

func (BuiltinGetLine) Call(args []any, i *Interpreter) any {
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanLines)
	if s.Scan() {
		return s.Text()
	}
	return ""
}
