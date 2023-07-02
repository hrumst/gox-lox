package main

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/interpret"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"os"
)

func main() {
	source := `fun makeCounter() { //+++
	  var i = 0;
	  fun count() {
	i = i + 1;
	print i; }
	  return count;
	}
	var counter = makeCounter();
	counter(); // "1".
	counter(); // "2".`

	tokens, err := scan.NewScanner(source).ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		os.Exit(70)
	}
	stmts, err := parse.NewParser(tokens).Parse()
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		os.Exit(70)
	}

	interpreter := interpret.NewInterpreter(os.Stdout)

	resolver := interpret.NewResolver(interpreter)
	if err := resolver.Resolve(stmts); err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		os.Exit(70)
	}

	if err := interpreter.Interpret(stmts); err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		os.Exit(70)
	}
}
