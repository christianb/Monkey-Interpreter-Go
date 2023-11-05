package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFn map[token.TokenType]prefixParseFn
	infixParseFn  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	parser.nextToken() // curToken is still nil, only peekToken is set now
	parser.nextToken() // sets also curToken, and peekToken again

	parser.prefixParseFn = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefixFn(token.IDENT, parser.parseIdentifier)
	parser.registerPrefixFn(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixFn(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFn(token.MINUS, parser.parsePrefixExpression)

	parser.infixParseFn = make(map[token.TokenType]infixParseFn)
	parser.registerInfixFn(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixFn(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixFn(token.EQ, parser.parseInfixExpression)
	parser.registerInfixFn(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfixFn(token.LT, parser.parseInfixExpression)
	parser.registerInfixFn(token.GT, parser.parseInfixExpression)

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !parser.curTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	stmt := &ast.LetStatement{Name: &ast.Identifier{Value: parser.curToken.Literal}}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We are skipping the expressions until we encounter a semicolon
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) expectPeek(t token.TokenType) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(t)
		return false
	}
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	parser.nextToken()

	// TODO: we are skipping the expressions until we encounter a semicolon
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return &ast.ReturnStatement{}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.curToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFn[parser.curToken.Type]
	if prefix == nil {
		parser.noPrefixPerseFnErrror(parser.curToken.Type)
		return nil
	}

	leftExpression := prefix()
	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFn[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		parser.nextToken()
		leftExpression = infix(leftExpression)

	}

	return leftExpression
}

func (parser *Parser) noPrefixPerseFnErrror(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
	}

	integerLiteral.Value = value

	return integerLiteral
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (parser *Parser) registerPrefixFn(tokkenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFn[tokkenType] = fn
}

func (parser *Parser) registerInfixFn(tokkenType token.TokenType, fn infixParseFn) {
	parser.infixParseFn[tokkenType] = fn
}

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (parser *Parser) getPrecedence(tokenType token.TokenType) int {
	precedence, ok := precedences[tokenType]

	if ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) peekPrecedence() int {
	return parser.getPrecedence(parser.peekToken.Type)
}

func (parser *Parser) curPrecendence() int {
	return parser.getPrecedence(parser.curToken.Type)
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}

	precedence := parser.curPrecendence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) curTokenIs(t token.TokenType) bool {
	return parser.curToken.Type == t
}

func (parser *Parser) peekTokenIs(t token.TokenType) bool {
	return parser.peekToken.Type == t
}
