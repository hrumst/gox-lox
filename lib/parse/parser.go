package parse

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/scan"
)

type Parser struct {
	tokens  []scan.Token
	current int
}

// todo add comma separated expressions, add ternar Operator
func NewParser(tokens []scan.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) assignment() (Expression, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(scan.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if varExpr, ok := expr.(*VariableExpression); ok {
			return NewAssignExpression(varExpr.Name, value), nil
		}
		return nil, NewParseError(equals, fmt.Errorf("invalid assignment target"))
	}

	return expr, nil
}

func (p *Parser) or() (Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(scan.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpression(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(scan.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpression(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) expression() (Expression, error) {
	return p.assignment()
}

func (p *Parser) match(tokenTypes ...scan.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType scan.TokenType) bool {
	return !p.isAtEnd() && p.peek().Type == tokenType
}

func (p *Parser) advance() scan.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scan.EOF
}

func (p *Parser) peek() scan.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scan.Token {
	return p.tokens[p.current-1]
}

// equality → comparison ( ( "!=" | "==" ) comparison )
func (p *Parser) equality() (Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scan.BANG_EQUAL, scan.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, right)
	}
	return expr, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scan.GREATER, scan.GREATER_EQUAL, scan.LESS, scan.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, right)
	}
	return expr, nil
}

// term → factor ( ( "-" | "+" ) factor )
func (p *Parser) term() (Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scan.MINUS, scan.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, right)
	}
	return expr, nil
}

// factor → unary ( ( "/" | "*" ) unary )
func (p *Parser) factor() (Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scan.SLASH, scan.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, right)
	}
	return expr, nil
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (Expression, error) {
	if p.match(scan.BANG, scan.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpression(operator, right), nil
	}
	return p.call()
}

func (p *Parser) call() (Expression, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scan.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expression) (Expression, error) {
	arguments := make([]Expression, 0)
	if !p.check(scan.RIGHT_PAREN) {
		for {
			argExpr, err := p.expression()
			if err != nil {
				return nil, err
			}
			if len(arguments) >= 255 {
				return nil, NewParseError(p.peek(), fmt.Errorf("can't have more then 255 call arguments"))
			}
			arguments = append(arguments, argExpr)
			if !p.match(scan.COMMA) {
				break
			}
		}
	}
	closeParen, err := p.consume(scan.RIGHT_PAREN, "expect ')' after arguments")
	if err != nil {
		return nil, err
	}
	return NewCallExpression(callee, closeParen, arguments), nil
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (Expression, error) {
	if p.match(scan.FALSE) {
		return NewLiteralExpression(scan.NewLiteral(scan.NewBooleanLoxValue(false))), nil
	} else if p.match(scan.TRUE) {
		return NewLiteralExpression(scan.NewLiteral(scan.NewBooleanLoxValue(true))), nil
	} else if p.match(scan.NIL) {
		return NewLiteralExpression(scan.NewLiteral(scan.NewNilLoxValue())), nil
	} else if p.match(scan.NUMBER, scan.STRING) {
		return NewLiteralExpression(p.previous().Literal), nil
	} else if p.match(scan.IDENTIFIER) {
		return NewVariableExpression(p.previous()), nil
	} else if p.match(scan.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(scan.RIGHT_PAREN, "expect ')' after expression"); err != nil {
			return nil, err
		}
		return NewGroupingExpression(expr), nil
	}

	return nil, NewParseError(p.peek(), fmt.Errorf("unexpected token type"))
}

func (p *Parser) consume(tokenType scan.TokenType, message string) (scan.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return scan.Token{}, NewParseError(p.peek(), fmt.Errorf(message))
}

func (p *Parser) printStmt() (*StmtPrint, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scan.SEMICOLON, "expect ';' after value"); err != nil {
		return nil, err
	}
	return NewStmtPrint(expr), nil
}

func (p *Parser) expressionStmt() (*StmtExpression, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scan.SEMICOLON, "expect ';' after value"); err != nil {
		return nil, err
	}
	return NewStmtExpression(expr), nil
}

func (p *Parser) statement() (Statement, error) {
	if p.match(scan.FOR) {
		return p.forStatement()
	} else if p.match(scan.IF) {
		return p.ifStatement()
	} else if p.match(scan.PRINT) {
		return p.printStmt()
	} else if p.match(scan.RETURN) {
		return p.returnStatement()
	} else if p.match(scan.WHILE) {
		return p.whileStatement()
	} else if p.match(scan.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewStmtBlock(stmts), nil
	} else if p.match(scan.CONTINUE) || p.match(scan.BREAK) {
		return p.breakContinueStatement()
	}
	return p.expressionStmt()
}

func (p *Parser) breakContinueStatement() (Statement, error) {
	stmt := NewExecuteControlStmt(p.previous())
	if _, err := p.consume(scan.SEMICOLON, "expect ';' after value"); err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) whileStatement() (Statement, error) {
	if _, err := p.consume(scan.LEFT_PAREN, "expect '(' after 'while'"); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scan.RIGHT_PAREN, "expect ')' after 'while' condition"); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return NewStmtWhile(condition, body), nil
}

func (p *Parser) forStatement() (Statement, error) {
	_, err := p.consume(scan.LEFT_PAREN, "expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer Statement
	if p.match(scan.SEMICOLON) {
		initializer = nil
	} else if p.match(scan.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStmt()
		if err != nil {
			return nil, err
		}
	}

	var condition Expression
	if !p.check(scan.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(scan.SEMICOLON, "expect ';' after loop condition"); err != nil {
		return nil, err
	}

	var increment Expression
	if !p.check(scan.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(scan.RIGHT_PAREN, "expect ')' after for clause"); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewStmtBlock(
			[]Statement{
				body,
				NewStmtExpression(increment),
			},
		)
	}

	if condition == nil {
		condition = NewLiteralExpression(
			scan.NewLiteral(
				scan.NewBooleanLoxValue(true),
			),
		)
	}
	body = NewStmtWhile(condition, body)

	if initializer != nil {
		body = NewStmtBlock(
			[]Statement{
				initializer,
				body,
			},
		)
	}
	return body, nil
}

func (p *Parser) ifStatement() (Statement, error) {
	if _, err := p.consume(scan.LEFT_PAREN, "expect '(' after 'if'"); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scan.RIGHT_PAREN, "expect ')' after 'if' condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Statement
	if p.match(scan.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewStmtIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) block() ([]Statement, error) {
	stmts := make([]Statement, 0)
	for !p.check(scan.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if _, err := p.consume(scan.RIGHT_BRACE, "expect '}' after block"); err != nil {
		return nil, err
	}
	return stmts, nil
}

func (p *Parser) varDeclaration() (Statement, error) {
	name, err := p.consume(scan.IDENTIFIER, "expect variable name")
	if err != nil {
		return nil, err
	}

	var initializer Expression
	if p.match(scan.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(scan.SEMICOLON, "expect ';' after variable declaration"); err != nil {
		return nil, err
	}
	return NewStmtVar(name, initializer), nil
}

func (p *Parser) declaration() (Statement, error) {
	if p.match(scan.FUN) {
		return p.function("function")
	} else if p.match(scan.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string) (Statement, error) {
	name, err := p.consume(scan.IDENTIFIER, fmt.Sprintf("exect %s name", kind))
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scan.LEFT_PAREN, fmt.Sprintf("exect '(' after %s name", kind)); err != nil {
		return nil, err
	}

	parameters := make([]scan.Token, 0)
	if !p.check(scan.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, NewParseError(p.peek(), fmt.Errorf("can't have more than 255 parameters"))
			}
			param, err := p.consume(scan.IDENTIFIER, "expect parameter name")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, param)
			if !p.match(scan.COMMA) {
				break
			}
		}
	}

	if _, err := p.consume(scan.RIGHT_PAREN, fmt.Sprintf("exect ')' after parameters")); err != nil {
		return nil, err
	}
	if _, err := p.consume(scan.LEFT_BRACE, fmt.Sprintf("exect '{' before %s name", kind)); err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return NewStmtFunction(name, parameters, body), nil
}

func (p *Parser) Parse() ([]Statement, error) {
	stmts := make([]Statement, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scan.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case scan.CLASS, scan.FUN, scan.VAR, scan.FOR, scan.IF, scan.WHILE, scan.PRINT, scan.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) returnStatement() (Statement, error) {
	keyword := p.previous()
	var value Expression
	if !p.check(scan.SEMICOLON) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(scan.SEMICOLON, "expect ';' after return value"); err != nil {
		return nil, err
	}
	return NewStmtReturn(keyword, value), nil
}
