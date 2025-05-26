package parser

import (
	"fmt"
	"github.com/fliptv97/monkey-interpreter/ast"
	"github.com/fliptv97/monkey-interpreter/lexer"
	"github.com/fliptv97/monkey-interpreter/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(lvalue ast.Expression) ast.Expression
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

type OperationPrecedence int

const (
	_ OperationPrecedence = iota
	PrecedenceLowest
	PrecedenceEquals      // ==
	PrecedenceLessGreater // > or <
	PrecedenceSum         // +
	PrecedenceProduct     // *
	PrecedencePrefix      // -X or !X
	PrecedenceCall        // myFunction(X)
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	/* These two calls initialise `p.currToken` and `p.peekToken`
	 * with the first and second tokens respectively */
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefixParseFn(token.Ident, p.parseIdentifier)
	p.registerPrefixParseFn(token.Int, p.parseIntegerLiteral)
	p.registerPrefixParseFn(token.Bang, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.Minus, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfixParseFn(token.Equal, p.parseInfixExpression)
	p.registerInfixParseFn(token.NotEqual, p.parseInfixExpression)
	p.registerInfixParseFn(token.LessThan, p.parseInfixExpression)
	p.registerInfixParseFn(token.GreaterThan, p.parseInfixExpression)
	p.registerInfixParseFn(token.Plus, p.parseInfixExpression)
	p.registerInfixParseFn(token.Minus, p.parseInfixExpression)
	p.registerInfixParseFn(token.Asterisk, p.parseInfixExpression)
	p.registerInfixParseFn(token.Slash, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isCurrToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.consumeSpecific(token.Ident) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.consumeSpecific(token.Assign) {
		return nil
	}
	for !p.isCurrToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	for !p.isCurrToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(PrecedenceLowest)

	if p.isPeekToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence OperationPrecedence) ast.Expression {
	prefixParser := p.prefixParseFns[p.currToken.Type]
	if prefixParser == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpr := prefixParser()

	for !p.isPeekToken(token.Semicolon) && precedence < p.getTokenPrecedence(p.peekToken) {
		infixParser := p.infixParseFns[p.peekToken.Type]
		if infixParser == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixParser(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	p.nextToken()

	expr.Right = p.parseExpression(PrecedencePrefix)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.getTokenPrecedence(p.currToken)
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) isCurrToken(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) isPeekToken(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) consumeSpecific(t token.Type) bool {
	if !p.isPeekToken(t) {
		p.registerPeekError(t)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPeekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefixParseFn(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFn(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

var precedences = map[token.Type]OperationPrecedence{
	token.Equal:       PrecedenceEquals,
	token.NotEqual:    PrecedenceEquals,
	token.LessThan:    PrecedenceLessGreater,
	token.GreaterThan: PrecedenceLessGreater,
	token.Plus:        PrecedenceSum,
	token.Minus:       PrecedenceSum,
	token.Asterisk:    PrecedenceProduct,
	token.Slash:       PrecedenceProduct,
}

func (p *Parser) getTokenPrecedence(token token.Token) OperationPrecedence {
	if precedence, ok := precedences[token.Type]; ok {
		return precedence
	}

	return PrecedenceLowest
}
