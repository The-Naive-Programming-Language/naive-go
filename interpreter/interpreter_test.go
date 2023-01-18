package interpreter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInterpreter_VisitBlock(t *testing.T) {
	Convey("nested", t, func() {
		src := `let a = 42;

{
    let a = 3.14;
    println(a);

    {
        println(a);
        let a = 2.718;
        println(a);
    }

    println(a);
}

println(a);`
		interp := New("", []byte(src))
		interp.Interpret()
	})
}
