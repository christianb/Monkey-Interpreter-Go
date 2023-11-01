package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(testing *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(testing, parser)
	if program == nil {
		testing.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		testing.Fatalf("expected program.Statements to contain 3 statements, but was %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(testing, stmt, test.expectedIdentifier) {
			return
		}
	}
}

func TestParseErrors(testing *testing.T) {
	input := `
	let x 5;
	let = 10;
	let 838383;
	`
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if len(p.errors) != 3 {
		testing.Fatalf("expected p.errors to contain 3 erros, but was %d", len(p.errors))
	}
}

func checkParseErrors(testing *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	testing.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		testing.Errorf("parser error: %q", msg)
	}
	testing.FailNow()
}

func testLetStatement(testint *testing.T, stmt ast.Statement, name string) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		testint.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		testint.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	return true
}

func TestReturnStatements(testing *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322,;
	`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(testing, parser)

	if len(program.Statements) != 3 {
		testing.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			testing.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
	}
}

func TestIdentifierExpression(testing *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(testing, parser)

	if len(program.Statements) != 1 {
		testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testing.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	identifier, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		testing.Fatalf("stmt.Expression is not *ast.Identifier. got=%T", stmt.Expression)
	}

	if identifier.Value != "foobar" {
		testing.Errorf("identifier.Value is not foobar. got=%s", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		testing.Errorf("identifier.TokenLiteral is not foobar. got=%s", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(testing *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(testing, parser)

	if len(program.Statements) != 1 {
		testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testing.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		testing.Fatalf("stmt.Expression is not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		testing.Errorf("literal.Value is not 5. got=%d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		testing.Errorf("literal.TokenLiteral is not 5. got=%s", literal.TokenLiteral())
	}
}
