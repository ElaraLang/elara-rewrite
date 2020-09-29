package parser

import "C"
import (
	lexer "elara/lexer"
	"fmt"
)

type Scanner = lexer.Scanner
type Token = lexer.Token
type TokenType = lexer.TokenType

type ParseError struct {
	token   Token
	message string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("Parse Error: %s at %s", pe.message, pe.token.Text)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens *[]Token) *Parser {
	return &Parser{
		tokens: *tokens,
	}
}

func (p *Parser) Parse() (result []Stmt, error []ParseError) {
	p.current = 0
	result = make([]Stmt, 0)
	error = make([]ParseError, 0)
	for !p.isAtEnd() {
		p.parseLine(&result, &error)
	}
	return
}

func (p *Parser) parseLine(result *[]Stmt, error *[]ParseError) {
	defer p.handleError(error)
	stmt := p.declaration()
	*result = append(*result, stmt)
	if !(p.match(lexer.NEWLINE) || p.isAtEnd()) {
		panic(ParseError{
			token:   p.peek(),
			message: "Expected new line",
		})
	}
}

func (p *Parser) handleError(error *[]ParseError) {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ParseError:
			*error = append(*error, err)
			break
		case []ParseError:
			*error = append(*error, err...)
		default:
			*error = append(*error, ParseError{
				token:   p.previous(),
				message: "Invalid error thrown by Parser",
			})
			break
		}
		p.syncError()
	}
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == lexer.EOF
}

func (p *Parser) check(tokenType TokenType) bool {
	return !p.isAtEnd() && p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, msg string) (token Token) {
	if p.check(tokenType) {
		token = p.advance()
		return
	}
	panic(ParseError{token: p.peek(), message: msg})
}

func (p *Parser) cleanNewLines() {
	for p.match(lexer.NEWLINE) {
	}
}
func (p *Parser) insert(index int, value ...Token) {
	if len(p.tokens) == index {
		p.tokens = append(p.tokens, value...)
	}
	p.tokens = append(p.tokens[:index+len(value)], p.tokens[index:]...)
	for i := 0; i < len(value); i++ {
		p.tokens[index+i] = value[i]
	}
}

func (p *Parser) syncError() {
	for !p.isAtEnd() && !p.check(lexer.NEWLINE) && !p.check(lexer.EOF) {
		p.advance()
	}
	p.cleanNewLines()
}
