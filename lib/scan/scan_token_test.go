package scan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanner_ScanTokensOk(t *testing.T) {
	assert.Equal(t, 1, 1)

	type testCase struct {
		source       string
		expectTokens []Token
	}

	testCases := []testCase{
		{
			`!false; // true.`,
			[]Token{
				{BANG, "!", nil, 0},
				{FALSE, "false", nil, 0},
				{SEMICOLON, ";", nil, 0},
				{EOF, "", nil, 0},
			},
		}, {
			`var average = (min + max) / 2;`,
			[]Token{
				{VAR, "var", nil, 0},
				{IDENTIFIER, "average", nil, 0},
				{EQUAL, "=", nil, 0},
				{LEFT_PAREN, "(", nil, 0},
				{IDENTIFIER, "min", nil, 0},
				{PLUS, "+", nil, 0},
				{IDENTIFIER, "max", nil, 0},
				{RIGHT_PAREN, ")", nil, 0},
				{SLASH, "/", nil, 0},
				{NUMBER, "2", NewLiteral(NewFloatLoxValue(2.)), 0},
				{SEMICOLON, ";", nil, 0},
				{EOF, "", nil, 0},
			},
		}, {
			`for (var a = 1; a < 10; a = a + 1) {
						print a;
					}`,
			[]Token{
				{FOR, "for", nil, 0},
				{LEFT_PAREN, "(", nil, 0},
				{VAR, "var", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{EQUAL, "=", nil, 0},
				{NUMBER, "1", NewLiteral(NewFloatLoxValue(1.)), 0},
				{SEMICOLON, ";", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{LESS, "<", nil, 0},
				{NUMBER, "10", NewLiteral(NewFloatLoxValue(10.)), 0},
				{SEMICOLON, ";", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{EQUAL, "=", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{PLUS, "+", nil, 0},
				{NUMBER, "1", NewLiteral(NewFloatLoxValue(1.)), 0},
				{RIGHT_PAREN, ")", nil, 0},
				{LEFT_BRACE, "{", nil, 0},
				{PRINT, "print", nil, 1},
				{IDENTIFIER, "a", nil, 1},
				{SEMICOLON, ";", nil, 1},
				{RIGHT_BRACE, "}", nil, 2},
				{EOF, "", nil, 2},
			},
		}, {
			` var a = 1;
					  while (a < 10) {
						print a;
						a = a + 1; 
					  }`,
			[]Token{
				{VAR, "var", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{EQUAL, "=", nil, 0},
				{NUMBER, "1", NewLiteral(NewFloatLoxValue(1.)), 0},
				{SEMICOLON, ";", nil, 0},
				{WHILE, "while", nil, 1},
				{LEFT_PAREN, "(", nil, 1},
				{IDENTIFIER, "a", nil, 1},
				{LESS, "<", nil, 1},
				{NUMBER, "10", NewLiteral(NewFloatLoxValue(10.)), 1},
				{RIGHT_PAREN, ")", nil, 1},
				{LEFT_BRACE, "{", nil, 1},
				{PRINT, "print", nil, 2},
				{IDENTIFIER, "a", nil, 2},
				{SEMICOLON, ";", nil, 2},
				{IDENTIFIER, "a", nil, 3},
				{EQUAL, "=", nil, 3},
				{IDENTIFIER, "a", nil, 3},
				{PLUS, "+", nil, 3},
				{NUMBER, "1", NewLiteral(NewFloatLoxValue(1.)), 3},
				{SEMICOLON, ";", nil, 3},
				{RIGHT_BRACE, "}", nil, 4},
				{EOF, "", nil, 4},
			},
		}, {
			`if (condition) {
						print "yes";
					  } else {
						print "no";
					}`,
			[]Token{
				{IF, "if", nil, 0},
				{LEFT_PAREN, "(", nil, 0},
				{IDENTIFIER, "condition", nil, 0},
				{RIGHT_PAREN, ")", nil, 0},
				{LEFT_BRACE, "{", nil, 0},
				{PRINT, "print", nil, 1},
				{STRING, "\"yes\"", NewLiteral(NewStringLoxValue("yes")), 1},
				{SEMICOLON, ";", nil, 1},
				{RIGHT_BRACE, "}", nil, 2},
				{ELSE, "else", nil, 2},
				{LEFT_BRACE, "{", nil, 2},
				{PRINT, "print", nil, 3},
				{STRING, "\"no\"", NewLiteral(NewStringLoxValue("no")), 3},
				{SEMICOLON, ";", nil, 3},
				{RIGHT_BRACE, "}", nil, 4},
				{EOF, "", nil, 4},
			},
		}, {
			`fun calculation(arg1, arg2) { 
                        return (arg1+45.6)*arg2/3; // parameters calculation
                    }`,
			[]Token{
				{FUN, "fun", nil, 0},
				{IDENTIFIER, "calculation", nil, 0},
				{LEFT_PAREN, "(", nil, 0},
				{IDENTIFIER, "arg1", nil, 0},
				{COMMA, ",", nil, 0},
				{IDENTIFIER, "arg2", nil, 0},
				{RIGHT_PAREN, ")", nil, 0},
				{LEFT_BRACE, "{", nil, 0},
				{RETURN, "return", nil, 1},
				{LEFT_PAREN, "(", nil, 1},
				{IDENTIFIER, "arg1", nil, 1},
				{PLUS, "+", nil, 1},
				{NUMBER, "45.6", NewLiteral(NewFloatLoxValue(45.6)), 1},
				{RIGHT_PAREN, ")", nil, 1},
				{STAR, "*", nil, 1},
				{IDENTIFIER, "arg2", nil, 1},
				{SLASH, "/", nil, 1},
				{NUMBER, "3", NewLiteral(NewFloatLoxValue(3)), 1},
				{SEMICOLON, ";", nil, 1},
				{RIGHT_BRACE, "}", nil, 2},
				{EOF, "", nil, 2},
			},
		}, {
			`class Breakfast {
						init(meat, bread) {
						  this.meat = meat;
						  this.bread = bread;
						}
					    // ...
					}
				    var baconAndToast = Breakfast("bacon", "toast");
				    baconAndToast.serve("Dear Reader");

					class Brunch < Breakfast {
						drink() {
						  print "How about a Bloody Mary?";
						}
					}
				
					var benedict = Brunch("ham", "English muffin");
				    benedict.serve("Noble Reader");
				`,

			[]Token{
				{CLASS, "class", nil, 0},
				{IDENTIFIER, "Breakfast", nil, 0},
				{LEFT_BRACE, "{", nil, 0},
				{IDENTIFIER, "init", nil, 1},
				{LEFT_PAREN, "(", nil, 1},
				{IDENTIFIER, "meat", nil, 1},
				{COMMA, ",", nil, 1},
				{IDENTIFIER, "bread", nil, 1},
				{RIGHT_PAREN, ")", nil, 1},
				{LEFT_BRACE, "{", nil, 1},
				{THIS, "this", nil, 2},
				{DOT, ".", nil, 2},
				{IDENTIFIER, "meat", nil, 2},
				{EQUAL, "=", nil, 2},
				{IDENTIFIER, "meat", nil, 2},
				{SEMICOLON, ";", nil, 2},
				{THIS, "this", nil, 3},
				{DOT, ".", nil, 3},
				{IDENTIFIER, "bread", nil, 3},
				{EQUAL, "=", nil, 3},
				{IDENTIFIER, "bread", nil, 3},
				{SEMICOLON, ";", nil, 3},
				{RIGHT_BRACE, "}", nil, 4},
				{RIGHT_BRACE, "}", nil, 6},
				{VAR, "var", nil, 7},
				{IDENTIFIER, "baconAndToast", nil, 7},
				{EQUAL, "=", nil, 7},
				{IDENTIFIER, "Breakfast", nil, 7},
				{LEFT_PAREN, "(", nil, 7},
				{STRING, "\"bacon\"", NewLiteral(NewStringLoxValue("bacon")), 7},
				{COMMA, ",", nil, 7},
				{STRING, "\"toast\"", NewLiteral(NewStringLoxValue("toast")), 7},
				{RIGHT_PAREN, ")", nil, 7},
				{SEMICOLON, ";", nil, 7},
				{IDENTIFIER, "baconAndToast", nil, 8},
				{DOT, ".", nil, 8},
				{IDENTIFIER, "serve", nil, 8},
				{LEFT_PAREN, "(", nil, 8},
				{STRING, "\"Dear Reader\"", NewLiteral(NewStringLoxValue("Dear Reader")), 8},
				{RIGHT_PAREN, ")", nil, 8},
				{SEMICOLON, ";", nil, 8},
				{CLASS, "class", nil, 10},
				{IDENTIFIER, "Brunch", nil, 10},
				{LESS, "<", nil, 10},
				{IDENTIFIER, "Breakfast", nil, 10},
				{LEFT_BRACE, "{", nil, 10},
				{IDENTIFIER, "drink", nil, 11},
				{LEFT_PAREN, "(", nil, 11},
				{RIGHT_PAREN, ")", nil, 11},
				{LEFT_BRACE, "{", nil, 11},
				{PRINT, "print", nil, 12},
				{STRING, "\"How about a Bloody Mary?\"", NewLiteral(NewStringLoxValue("How about a Bloody Mary?")), 12},
				{SEMICOLON, ";", nil, 12},
				{RIGHT_BRACE, "}", nil, 13},
				{RIGHT_BRACE, "}", nil, 14},
				{VAR, "var", nil, 16},
				{IDENTIFIER, "benedict", nil, 16},
				{EQUAL, "=", nil, 16},
				{IDENTIFIER, "Brunch", nil, 16},
				{LEFT_PAREN, "(", nil, 16},
				{STRING, "\"ham\"", NewLiteral(NewStringLoxValue("ham")), 16},
				{COMMA, ",", nil, 16},
				{STRING, "\"English muffin\"", NewLiteral(NewStringLoxValue("English muffin")), 16},
				{RIGHT_PAREN, ")", nil, 16},
				{SEMICOLON, ";", nil, 16},
				{IDENTIFIER, "benedict", nil, 17},
				{DOT, ".", nil, 17},
				{IDENTIFIER, "serve", nil, 17},
				{LEFT_PAREN, "(", nil, 17},
				{STRING, "\"Noble Reader\"", NewLiteral(NewStringLoxValue("Noble Reader")), 17},
				{RIGHT_PAREN, ")", nil, 17},
				{SEMICOLON, ";", nil, 17},
				{EOF, "", nil, 18},
			},
		}, {
			`var a = 1;
					/*
					Here is multi-line
					comment
					*/
					var b = 2;`,
			[]Token{
				{VAR, "var", nil, 0},
				{IDENTIFIER, "a", nil, 0},
				{EQUAL, "=", nil, 0},
				{NUMBER, "1", NewLiteral(NewFloatLoxValue(1.)), 0},
				{SEMICOLON, ";", nil, 0},
				{VAR, "var", nil, 5},
				{IDENTIFIER, "b", nil, 5},
				{EQUAL, "=", nil, 5},
				{NUMBER, "2", NewLiteral(NewFloatLoxValue(2.)), 5},
				{SEMICOLON, ";", nil, 5},
				{EOF, "", nil, 5},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(
			fmt.Sprintf("test_case_%d", i),
			func(t *testing.T) {
				sc := NewScanner(tc.source)
				tokens, err := sc.ScanTokens()
				assert.NoError(t, err)
				assert.Equal(t, tc.expectTokens, tokens)
			},
		)
	}
}

func TestScanner_ScanTokensFail(t *testing.T) {
	t.Skip() // todo
}
