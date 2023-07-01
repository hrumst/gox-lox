package interpret

import (
	"bytes"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolver_Resolve(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		var stmtBlock *parse.StmtBlock
		testStmt1 := parse.NewVariableExpression(
			scan.NewToken(scan.IDENTIFIER, "v", nil, 0),
		)
		testStmt2 := parse.NewVariableExpression(
			scan.NewToken(scan.IDENTIFIER, "v", nil, 0),
		)

		for i := 0; i < 10; i += 1 {
			if stmtBlock == nil {
				stmtBlock = parse.NewStmtBlock(
					[]parse.Statement{
						parse.NewStmtExpression(testStmt1),
					},
				)
			} else {
				stmtBlock = parse.NewStmtBlock(
					[]parse.Statement{
						stmtBlock,
					},
				)
			}
		}
		testStmt := []parse.Statement{
			parse.NewStmtVar(
				scan.NewToken(scan.IDENTIFIER, "v", nil, 0),
				parse.NewLiteralExpression(
					scan.NewLiteral(
						scan.NewFloatLoxValue(1.),
					),
				),
			),
			stmtBlock.Stmts[0],
			parse.NewStmtPrint(testStmt2),
		}

		buf := bytes.NewBufferString("")
		interpreter := NewInterpreter(buf)
		if err := NewResolver(interpreter).Resolve(testStmt); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, interpreter.locals[testStmt1], 9)
		assert.Equal(t, interpreter.locals[testStmt2], 0)
	})

	t.Run("returnError", func(t *testing.T) {
		/*
			{
				return 1;
			}
		*/
		testStmt := []parse.Statement{
			parse.NewStmtBlock(
				[]parse.Statement{
					parse.NewStmtReturn(
						scan.NewToken(scan.IDENTIFIER, "return", nil, 0),
						parse.NewLiteralExpression(
							scan.NewLiteral(
								scan.NewFloatLoxValue(1.),
							),
						),
					),
				},
			),
		}

		buf := bytes.NewBufferString("")
		interpreter := NewInterpreter(buf)
		err := NewResolver(interpreter).Resolve(testStmt)
		assert.Errorf(t, err, "can't return from top-level code")
	})

	t.Run("doubleVarDeclaration", func(t *testing.T) {
		/*
			{
				var v1;
				var v1;
			}
		*/
		testStmt := []parse.Statement{
			parse.NewStmtBlock(
				[]parse.Statement{
					parse.NewStmtVar(
						scan.NewToken(scan.IDENTIFIER, "v1", nil, 0),
						nil,
					),
					parse.NewStmtVar(
						scan.NewToken(scan.IDENTIFIER, "v1", nil, 0),
						nil,
					),
				},
			),
		}

		buf := bytes.NewBufferString("")
		interpreter := NewInterpreter(buf)
		err := NewResolver(interpreter).Resolve(testStmt)
		assert.Errorf(t, err, "already variable with this name in this scope")
	})

	t.Run("varWithItsOwnInitializer", func(t *testing.T) {
		// var i = i;
		testStmt := []parse.Statement{
			parse.NewStmtBlock(
				[]parse.Statement{
					parse.NewStmtVar(
						scan.NewToken(scan.IDENTIFIER, "v1", nil, 0),
						parse.NewVariableExpression(
							scan.NewToken(scan.IDENTIFIER, "v1", nil, 0),
						),
					),
				},
			),
		}

		buf := bytes.NewBufferString("")
		interpreter := NewInterpreter(buf)
		err := NewResolver(interpreter).Resolve(testStmt)
		assert.Errorf(t, err, "can't read local variable in its own initializer")
	})
}
