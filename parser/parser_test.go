package parser

import (
	"fmt"
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
	checkParserErrors(testing, parser)
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

	for _, err := range p.errors {
		testing.Logf(err)
	}

	if len(p.errors) != 4 {
		testing.Fatalf("expected p.errors to contain 3 erros, but was %d", len(p.errors))
	}
}

func checkParserErrors(testing *testing.T, parser *Parser) {
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

func testLetStatement(testing *testing.T, stmt ast.Statement, name string) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		testing.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		testing.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
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
	checkParserErrors(testing, parser)

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
	checkParserErrors(testing, parser)

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
	checkParserErrors(testing, parser)

	if len(program.Statements) != 1 {
		testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testing.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIntegerLiteral(testing, stmt.Expression, 5)
}

func TestParsingPrefixExpressions(testing *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, test := range prefixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(testing, parser)

		if len(program.Statements) != 1 {
			testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			testing.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			testing.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != test.operator {
			testing.Fatalf("exp.Operator is not %s. got=%s", test.operator, exp.Operator)
		}

		testLiteralExpression(testing, exp.Right, test.value)
	}
}

func testIntegerLiteral(testing *testing.T, il ast.Expression, expected int64) bool {
	integerLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		testing.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integerLiteral.Value != expected {
		testing.Errorf("integerLiteral.Value is not %d. got=%d", expected, integerLiteral.Value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", expected) {
		testing.Errorf("integerLiteral.TokenLiteral is not %d. got=%s", expected, integerLiteral.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(testing *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false != false", false, "!=", false},
	}

	for _, test := range infixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(testing, parser)

		if len(program.Statements) != 1 {
			testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			testing.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testInfixExpression(testing, stmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func TestOperatorPrecedenceParsing(testing *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(testing, parser)

		actual := program.String()
		if actual != test.expected {
			testing.Errorf("expected=%q got=%q", test.expected, actual)
		}
	}
}

func testIdentifier(testing *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		testing.Errorf("expression is not *ast.Identifier. got=%T", expression)
		return false
	}

	if identifier.Value != value {
		testing.Errorf("identifier.Value is not %s. got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		testing.Errorf("identifier.TokenLiteral is not %s. got=%s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(testing *testing.T, expression ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(testing, expression, int64(value))
	case int64:
		return testIntegerLiteral(testing, expression, value)
	case string:
		return testIdentifier(testing, expression, value)
	case bool:
		return testBooleanLiteral(testing, expression, value)
	}

	testing.Errorf("type of expression not handled. got=%T", expression)
	return false
}

func testInfixExpression(
	testing *testing.T,
	expression ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	operatorExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		testing.Errorf("expression is not ast.InfixExpression. got=%T(%s)", expression, expression)
		return false
	}

	if !testLiteralExpression(testing, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		testing.Errorf("expression.Operator is not %s. got=%q", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(testing, operatorExpression.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(testing *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(testing, parser)

		if len(program.Statements) != 1 {
			testing.Fatalf("program should have 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			testing.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			testing.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}

		if boolean.Value != test.expectedBoolean {
			testing.Errorf("boolean.Value not %t. got=%t", test.expectedBoolean, boolean.Value)
		}
	}
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	expressionBoolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expressionBoolean not *ast.Boolean. got=%T", expression)
		return false
	}

	if expressionBoolean.Value != value {
		t.Errorf("expressionBoolean.Value not %t. got=%t", value, expressionBoolean.Value)
		return false
	}

	if expressionBoolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("expressionBoolean.TokenLiteral not %t. got=%s",
			value, expressionBoolean.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statments does not contain exactly 1 statment. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative was not nil. got=%+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}
