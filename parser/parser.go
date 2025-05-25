package parser

import (
	"fmt"
	"github.com/fliptv97/monkey-interpreter/ast"
	"github.com/fliptv97/monkey-interpreter/lexer"
	"github.com/fliptv97/monkey-interpreter/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	/* These two calls initialise `p.currToken` and `p.peekToken`
	 * with the first and second tokens respectively */
	p.nextToken()
	p.nextToken()

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
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
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

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	for !p.isCurrToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
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
