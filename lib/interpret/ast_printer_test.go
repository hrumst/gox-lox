package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	type testCase struct {
		expr     parse.Expression
		expected string
	}

	tcs := []testCase{
		{
			expr: parse.NewBinaryExpression(
				parse.NewUnaryExpression(
					scan.NewToken(scan.MINUS, "-", nil, 1),
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewFloatLoxValue(123.),
						),
					),
				),
				scan.NewToken(scan.STAR, "*", nil, 1),
				parse.NewGroupingExpression(
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewFloatLoxValue(45.67),
						),
					),
				),
			),
			expected: "((123 -) (45.67 group) *)",
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("ast_printer_test_case_%d", i),
			func(t *testing.T) {
				result, err := NewAstPrinter(true).Print(tc.expr)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			},
		)
	}
}
