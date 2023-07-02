package main

import (
	"github.com/hrumst/gox-lox/lib/interpret"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"os"
)

// todo if (hadRuntimeError) System.exit(70);
func main() {
	//source := `fun makeCounter() { //+++
	//   var i = 0;
	//   fun count() {
	//i = i + 1;
	//print i; }
	//   return count;
	// }
	// var counter = makeCounter();
	// counter(); // "1".
	// counter(); // "2".`

	//source := `var i = i;`

	//source := `
	//class Circle {
	//	init(radius, scale) {
	//		this.radius = radius;
	//		this.scale = scale;
	//	}
	//
	//	area() {
	//		return 3.141592653 * this.radius * this.radius * this.scale;
	//	}
	//}
	//var circle = Circle(4, 2);
	//print circle.area();
	//`

	source := `
	class A {
		method() {
			return "Method A";
		}
	}

    class B < A {}
	print B().method();

  	class C < A {
		method() {
			return "Method C";
		}
		test() {
			return super.method() + " from C";
		}
  	}
    print C().test();`

	tokens, err := scan.NewScanner(source).ScanTokens()
	if err != nil {
		panic(err)
	}
	stmts, err := parse.NewParser(tokens).Parse()
	if err != nil {
		panic(err)
	}

	interpreter := interpret.NewInterpreter(os.Stdout)

	resolver := interpret.NewResolver(interpreter)
	if err := resolver.Resolve(stmts); err != nil {
		panic(err)
	}

	if err := interpreter.Interpret(stmts); err != nil {
		panic(err)
	}
}
