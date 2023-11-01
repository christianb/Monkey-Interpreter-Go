package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
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

	for parser.curToken.Type != token.EOF {
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
	for parser.curToken.Type != token.SEMICOLON {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) expectPeek(t token.TokenType) bool {
	if parser.peekToken.Type == t {
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
	for parser.curToken.Type != token.SEMICOLON {
		parser.nextToken()
	}

	return &ast.ReturnStatement{}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.curToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekToken.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFn[parser.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	return leftExpression
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
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
