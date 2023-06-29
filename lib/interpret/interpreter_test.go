package interpret

import (
	"bytes"
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterpreter_Evaluate(t *testing.T) {
	type testCase struct {
		stmts    []parse.Statement
		expected string
	}

	tcs := []testCase{
		{
			// print "one";
			stmts: []parse.Statement{
				parse.NewStmtPrint(
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewStringLoxValue("one"),
						),
					),
				),
			},
			expected: "one\n",
		}, {
			// print "string" + 4 == "string4";
			stmts: []parse.Statement{
				parse.NewStmtPrint(
					parse.NewBinaryExpression(
						parse.NewBinaryExpression(
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewStringLoxValue("string"),
								),
							),
							scan.NewToken(scan.PLUS, "+", nil, 0),
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewFloatLoxValue(4.),
								),
							),
						),
						scan.NewToken(scan.EQUAL_EQUAL, "==", nil, 0),
						parse.NewLiteralExpression(
							scan.NewLiteral(
								scan.NewStringLoxValue("string4"),
							),
						),
					),
				),
			},
			expected: "true\n",
		}, {
			// print 13.4 >= (5 + -3) - --1 == (7 != 6);
			[]parse.Statement{
				parse.NewStmtPrint(
					parse.NewBinaryExpression(
						parse.NewBinaryExpression(
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewFloatLoxValue(13.4),
								),
							),
							scan.NewToken(scan.GREATER_EQUAL, ">=", nil, 0),
							parse.NewBinaryExpression(
								parse.NewGroupingExpression(
									parse.NewBinaryExpression(
										parse.NewLiteralExpression(
											scan.NewLiteral(
												scan.NewFloatLoxValue(5.),
											),
										),
										scan.NewToken(scan.PLUS, "+", nil, 0),
										parse.NewUnaryExpression(
											scan.NewToken(scan.MINUS, "-", nil, 0),
											parse.NewLiteralExpression(
												scan.NewLiteral(
													scan.NewFloatLoxValue(3.),
												),
											),
										),
									),
								),
								scan.NewToken(scan.MINUS, "-", nil, 0),
								parse.NewUnaryExpression(
									scan.NewToken(scan.MINUS, "-", nil, 0),
									parse.NewUnaryExpression(
										scan.NewToken(scan.MINUS, "-", nil, 0),
										parse.NewLiteralExpression(
											scan.NewLiteral(
												scan.NewFloatLoxValue(1.),
											),
										),
									),
								),
							),
						),
						scan.NewToken(scan.EQUAL_EQUAL, "==", nil, 0),
						parse.NewGroupingExpression(
							parse.NewBinaryExpression(
								parse.NewLiteralExpression(
									scan.NewLiteral(
										scan.NewFloatLoxValue(7),
									),
								),
								scan.NewToken(scan.BANG_EQUAL, "!=", nil, 0),
								parse.NewLiteralExpression(
									scan.NewLiteral(
										scan.NewFloatLoxValue(6.),
									),
								),
							),
						),
					),
				),
			},
			"true\n",
		}, {
			/*  var a = "global A";
			var b;
			var c;
			{
				print a;
				print b;
				var a = "inner A";
				print a;
				b = "inner from global B";
			}
			print a;
			print b;
			print c; */
			stmts: []parse.Statement{
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "a", nil, 0),
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewStringLoxValue("global A"),
						),
					),
				),
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "b", nil, 0),
					nil,
				),
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "c", nil, 0),
					nil,
				),
				parse.NewStmtBlock(
					[]parse.Statement{
						parse.NewStmtPrint(
							parse.NewVariableExpression(
								scan.NewToken(scan.IDENTIFIER, "a", nil, 0),
							),
						),
						parse.NewStmtPrint(
							parse.NewVariableExpression(
								scan.NewToken(scan.IDENTIFIER, "b", nil, 0),
							),
						),
						parse.NewStmtVar(
							scan.NewToken(scan.IDENTIFIER, "a", nil, 0),
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewStringLoxValue("inner A"),
								),
							),
						),
						parse.NewStmtPrint(
							parse.NewVariableExpression(
								scan.NewToken(scan.IDENTIFIER, "a", nil, 0),
							),
						),
						parse.NewStmtExpression(
							parse.NewAssignExpression(
								scan.NewToken(scan.IDENTIFIER, "b", nil, 0),
								parse.NewLiteralExpression(
									scan.NewLiteral(
										scan.NewStringLoxValue("inner from global B"),
									),
								),
							),
						),
					},
				),
				parse.NewStmtPrint(
					parse.NewVariableExpression(
						scan.NewToken(scan.IDENTIFIER, "a", nil, 0),
					),
				),
				parse.NewStmtPrint(
					parse.NewVariableExpression(
						scan.NewToken(scan.IDENTIFIER, "b", nil, 0),
					),
				),
				parse.NewStmtPrint(
					parse.NewVariableExpression(
						scan.NewToken(scan.IDENTIFIER, "c", nil, 0),
					),
				),
			},
			expected: "global A\nnil\ninner A\nglobal A\ninner from global B\nnil\n",
		}, {
			// for (var i = 0; i < 10; i = i + 1) print i;
			stmts: []parse.Statement{
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewFloatLoxValue(1.),
						),
					),
				),
				parse.NewStmtWhile(
					parse.NewBinaryExpression(
						parse.NewVariableExpression(
							scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
						),
						scan.NewToken(scan.LESS, "<", nil, 0),
						parse.NewLiteralExpression(
							scan.NewLiteral(
								scan.NewFloatLoxValue(10.),
							),
						),
					),
					parse.NewStmtBlock(
						[]parse.Statement{
							parse.NewStmtPrint(
								parse.NewVariableExpression(
									scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
								),
							),
							parse.NewStmtExpression(
								parse.NewAssignExpression(
									scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
									parse.NewBinaryExpression(
										parse.NewVariableExpression(
											scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
										),
										scan.NewToken(scan.PLUS, "+", nil, 0),
										parse.NewLiteralExpression(
											scan.NewLiteral(
												scan.NewFloatLoxValue(1.),
											),
										),
									),
								),
							),
						},
					),
				),
			},
			expected: "1\n2\n3\n4\n5\n6\n7\n8\n9\n",
		}, {
			/*
				var i = 0;
				while (true) {
					i = i + 1;
					if (i < 3) {
						continue;
					}
					print i;
					if (i > 5) {
						break;
					}
				}
			*/
			stmts: []parse.Statement{
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewFloatLoxValue(1.),
						),
					),
				),
				parse.NewStmtWhile(
					parse.NewLiteralExpression(
						scan.NewLiteral(
							scan.NewBooleanLoxValue(true),
						),
					),
					parse.NewStmtBlock(
						[]parse.Statement{
							parse.NewStmtExpression(
								parse.NewAssignExpression(
									scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
									parse.NewBinaryExpression(
										parse.NewVariableExpression(
											scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
										),
										scan.NewToken(scan.PLUS, "+", nil, 0),
										parse.NewLiteralExpression(
											scan.NewLiteral(
												scan.NewFloatLoxValue(1.),
											),
										),
									),
								),
							),
							parse.NewStmtIf(
								parse.NewBinaryExpression(
									parse.NewVariableExpression(
										scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
									),
									scan.NewToken(scan.LESS, "<", nil, 0),
									parse.NewLiteralExpression(
										scan.NewLiteral(
											scan.NewFloatLoxValue(3.),
										),
									),
								),
								parse.NewStmtBlock(
									[]parse.Statement{
										parse.NewExecuteControlStmt(
											scan.NewToken(scan.CONTINUE, "continue", nil, 0),
										),
									},
								),
								nil,
							),
							parse.NewStmtPrint(
								parse.NewVariableExpression(
									scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
								),
							),
							parse.NewStmtIf(
								parse.NewBinaryExpression(
									parse.NewVariableExpression(
										scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
									),
									scan.NewToken(scan.GREATER, ">", nil, 0),
									parse.NewLiteralExpression(
										scan.NewLiteral(
											scan.NewFloatLoxValue(5.),
										),
									),
								),
								parse.NewStmtBlock(
									[]parse.Statement{
										parse.NewExecuteControlStmt(
											scan.NewToken(scan.BREAK, "break", nil, 0),
										),
									},
								),
								nil,
							),
						},
					),
				),
			},
			expected: "3\n4\n5\n6\n",
		}, {
			/*
				fun fib(n) {
				   if (n <= 1) return n;
				   return fib(n - 2) + fib(n - 1);
				}
				for (var i = 0; i < 15; i = i + 1) {
				   print fib(i);
				}
			*/
			stmts: []parse.Statement{
				parse.NewStmtFunction(
					scan.NewToken(scan.IDENTIFIER, "fib", nil, 0),
					[]scan.Token{
						scan.NewToken(scan.IDENTIFIER, "n", nil, 0),
					},
					[]parse.Statement{
						parse.NewStmtIf(
							parse.NewBinaryExpression(
								parse.NewVariableExpression(
									scan.NewToken(scan.IDENTIFIER, "n", nil, 0),
								),
								scan.NewToken(scan.LESS_EQUAL, "<=", nil, 0),
								parse.NewLiteralExpression(
									scan.NewLiteral(
										scan.NewFloatLoxValue(1.),
									),
								),
							),
							parse.NewStmtReturn(
								scan.NewToken(scan.RETURN, "return", nil, 0),
								parse.NewVariableExpression(
									scan.NewToken(scan.IDENTIFIER, "n", nil, 0),
								),
							),
							nil,
						),
						parse.NewStmtReturn(
							scan.NewToken(scan.RETURN, "return", nil, 0),
							parse.NewBinaryExpression(
								parse.NewCallExpression(
									parse.NewVariableExpression(
										scan.NewToken(scan.IDENTIFIER, "fib", nil, 0),
									),
									scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
									[]parse.Expression{
										parse.NewBinaryExpression(
											parse.NewVariableExpression(
												scan.NewToken(scan.IDENTIFIER, "n", nil, 0),
											),
											scan.NewToken(scan.MINUS, "-", nil, 0),
											parse.NewLiteralExpression(
												scan.NewLiteral(
													scan.NewFloatLoxValue(2.),
												),
											),
										),
									},
								),
								scan.NewToken(scan.PLUS, "+", nil, 0),
								parse.NewCallExpression(
									parse.NewVariableExpression(
										scan.NewToken(scan.IDENTIFIER, "fib", nil, 0),
									),
									scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
									[]parse.Expression{
										parse.NewBinaryExpression(
											parse.NewVariableExpression(
												scan.NewToken(scan.IDENTIFIER, "n", nil, 0),
											),
											scan.NewToken(scan.MINUS, "-", nil, 0),
											parse.NewLiteralExpression(
												scan.NewLiteral(
													scan.NewFloatLoxValue(1.),
												),
											),
										),
									},
								),
							),
						),
					},
				),
				parse.NewStmtBlock(
					[]parse.Statement{
						parse.NewStmtVar(
							scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewFloatLoxValue(0.),
								),
							),
						),
						parse.NewStmtWhile(
							parse.NewBinaryExpression(
								parse.NewVariableExpression(
									scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
								),
								scan.NewToken(scan.LESS, "<", nil, 0),
								parse.NewLiteralExpression(
									scan.NewLiteral(
										scan.NewFloatLoxValue(15.),
									),
								),
							),
							parse.NewStmtBlock(
								[]parse.Statement{
									parse.NewStmtPrint(
										parse.NewCallExpression(
											parse.NewVariableExpression(
												scan.NewToken(scan.IDENTIFIER, "fib", nil, 0),
											),
											scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
											[]parse.Expression{
												parse.NewVariableExpression(
													scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
												),
											},
										),
									),
									parse.NewStmtExpression(
										parse.NewAssignExpression(
											scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
											parse.NewBinaryExpression(
												parse.NewVariableExpression(
													scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
												),
												scan.NewToken(scan.PLUS, "+", nil, 0),
												parse.NewLiteralExpression(
													scan.NewLiteral(
														scan.NewFloatLoxValue(1.),
													),
												),
											),
										),
									),
								},
							),
						),
					},
				),
			},
			expected: "0\n1\n1\n2\n3\n5\n8\n13\n21\n34\n55\n89\n144\n233\n377\n",
		}, {
			/*
				fun makeCounter() {
				   var i = 0;
				   fun count() {
						i = i + 1;
						print i;
				   }
				   return count;
				 }
				 var counter = makeCounter();
				 counter(); // "1".
				 counter(); // "2"*/
			stmts: []parse.Statement{
				parse.NewStmtFunction(
					scan.NewToken(scan.IDENTIFIER, "makeCounter", nil, 0),
					[]scan.Token{},
					[]parse.Statement{
						parse.NewStmtVar(
							scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
							parse.NewLiteralExpression(
								scan.NewLiteral(
									scan.NewFloatLoxValue(0.),
								),
							),
						),
						parse.NewStmtFunction(
							scan.NewToken(scan.IDENTIFIER, "count", nil, 0),
							[]scan.Token{},
							[]parse.Statement{
								parse.NewStmtExpression(
									parse.NewAssignExpression(
										scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
										parse.NewBinaryExpression(
											parse.NewVariableExpression(
												scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
											),
											scan.NewToken(scan.PLUS, "+", nil, 0),
											parse.NewLiteralExpression(
												scan.NewLiteral(
													scan.NewFloatLoxValue(1.),
												),
											),
										),
									),
								),
								parse.NewStmtPrint(
									parse.NewVariableExpression(
										scan.NewToken(scan.IDENTIFIER, "i", nil, 0),
									),
								),
							},
						),
						parse.NewStmtReturn(
							scan.NewToken(scan.RETURN, "return", nil, 0),
							parse.NewVariableExpression(
								scan.NewToken(scan.IDENTIFIER, "count", nil, 0),
							),
						),
					},
				),
				parse.NewStmtVar(
					scan.NewToken(scan.IDENTIFIER, "counter", nil, 0),
					parse.NewCallExpression(
						parse.NewVariableExpression(
							scan.NewToken(scan.IDENTIFIER, "makeCounter", nil, 0),
						),
						scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
						[]parse.Expression{},
					),
				),
				parse.NewStmtExpression(
					parse.NewCallExpression(
						parse.NewVariableExpression(
							scan.NewToken(scan.IDENTIFIER, "counter", nil, 0),
						),
						scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
						[]parse.Expression{},
					),
				),
				parse.NewStmtExpression(
					parse.NewCallExpression(
						parse.NewVariableExpression(
							scan.NewToken(scan.IDENTIFIER, "counter", nil, 0),
						),
						scan.NewToken(scan.RIGHT_PAREN, ")", nil, 0),
						[]parse.Expression{},
					),
				),
			},
			expected: "1\n2\n",
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("interpreter_test_case_%d", i),
			func(t *testing.T) {
				buf := bytes.NewBufferString("")
				err := NewInterpreter(buf).Interpret(tc.stmts)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, buf.String())
			},
		)
	}

	t.Skip()
}
