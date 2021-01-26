package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Parser struct {
	Output             chan ast.Statement
	Tape               *TokenTape
	statementParslets  map[lexer.TokenType]statementParslet
	prefixParslets     map[lexer.TokenType]prefixParslet
	infixParslets      map[lexer.TokenType]infixParslet
	prefixTypeParslets map[lexer.TokenType]prefixTypeParslet
	infixTypeParslets  map[lexer.TokenType]infixTypeParslet
}

func NewParser(inputChannel chan lexer.Token, outputChannel chan ast.Statement) Parser {
	tape := NewTokenTape(inputChannel)
	p := Parser{Output: outputChannel, Tape: &tape}
	p.initPrefixParselets()
	p.initInfixParselets()
	p.initStatementParselets()
	return p
}

func (p *Parser) Parse() {
	if p.Tape.isClosed() {
		p.Tape.unwind()
	}
	for !p.Tape.ValidateHead(lexer.EOF) {
		p.Output <- p.parseStatement()
		p.Tape.Expect(lexer.EOF, lexer.NEWLINE)
	}
}

type (
	statementParslet  func() ast.Statement
	prefixParslet     func() ast.Expression
	infixParslet      func(ast.Expression) ast.Expression
	prefixTypeParslet func() ast.Type
	infixTypeParslet  func(ast.Type) ast.Type
)

func (p *Parser) registerPrefix(tokenType lexer.TokenType, function prefixParslet) {
	p.prefixParslets[tokenType] = function
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, function infixParslet) {
	p.infixParslets[tokenType] = function
}

func (p *Parser) registerTypePrefix(tokenType lexer.TokenType, function prefixTypeParslet) {
	p.prefixTypeParslets[tokenType] = function
}
func (p *Parser) registerTypeInfix(tokenType lexer.TokenType, function infixTypeParslet) {
	p.infixTypeParslets[tokenType] = function
}

func (p *Parser) registerStatement(tokenType lexer.TokenType, function statementParslet) {
	p.statementParslets[tokenType] = function
}

func (p *Parser) parseStatement() ast.Statement {
	parseStmt := p.statementParslets[p.Tape.Current().TokenType]
	if parseStmt == nil {
		return p.parseExpressionStatement()
	}
	return parseStmt()
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.prefixParslets[p.Tape.Current().TokenType]
	if parsePrefix == nil {
		// panic
		return nil
	}
	expr := parsePrefix()
	for !p.Tape.ValidationPeek(0, lexer.NEWLINE) && precedence < precedenceOf(p.Tape.Current().TokenType) {
		infix := p.infixParslets[p.Tape.Current().TokenType]
		if infix == nil {
			return expr
		}
		expr = infix(expr)
	}
	return expr
}

func (p *Parser) parseType(precedence int) ast.Type {
	parsePrefixType := p.prefixTypeParslets[p.Tape.Current().TokenType]
	if parsePrefixType == nil {
		// panic
		return nil
	}
	typ := parsePrefixType()
	for !p.Tape.ValidationPeek(0, lexer.NEWLINE) && precedence < typePrecedenceOf(p.Tape.Current().TokenType) {
		infix := p.infixTypeParslets[p.Tape.Current().TokenType]
		if infix == nil {
			return typ
		}
		typ = infix(typ)
	}
	return typ
}
