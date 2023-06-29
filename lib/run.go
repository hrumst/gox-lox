package lib

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hrumst/gox-lox/lib/scan"
)

func Run(source string) ([]scan.Token, error) {
	return scan.NewScanner(source).ScanTokens()
}

func runPrintTokens(source string) error {
	tokens, err := Run(source)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter script (type 'q' or 'quit' to exit):")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "q" || line == "quit" {
			break
		}
		runPrintTokens(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error on reading input:", err)
	}
}

func RunFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error on reading input:", err)
	}
	runPrintTokens(string(source))
}
