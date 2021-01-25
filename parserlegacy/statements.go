package parserlegacy

import "github.com/ElaraLang/elara/lexer"

type Stmt interface {
	stmtNode()
}

type ExpressionStmt struct {
	Expr Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type VarDefStmt struct {
	Mutable    bool
	Lazy       bool
	Open       bool
	Identifier string
	Type       Type
	Value      Expr
}

type StructDefStmt struct {
	Identifier   string
	StructFields []StructField
}

type IfElseStmt struct {
	Condition  Expr
	MainBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type ExtendStmt struct {
	Identifier string
	Body       BlockStmt
	Alias      string
}
type TypeStmt struct {
	Identifier string
	Contract   Type
}
type GenerifiedStmt struct {
	Contracts []GenericContract
	Statement Stmt
}

type ReturnStmt struct {
	Returning Expr
}

func (ExpressionStmt) stmtNode() {}
func (BlockStmt) stmtNode()      {}
func (VarDefStmt) stmtNode()     {}
func (StructDefStmt) stmtNode()  {}
func (IfElseStmt) stmtNode()     {}
func (WhileStmt) stmtNode()      {}
func (ExtendStmt) stmtNode()     {}
func (GenerifiedStmt) stmtNode() {}
func (TypeStmt) stmtNode()       {}
func (ReturnStmt) stmtNode()     {}

func (p *Parser) declaration() (stmt Stmt) {
	if p.check(lexer.Let) {
		return p.varDefStatement()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	switch p.peek().TokenType {
	case lexer.While:
		return p.whileStatement()
	case lexer.If:
		return p.ifStatement()
	case lexer.LBrace:
		return p.blockStatement()
	case lexer.Struct:
		return p.structStatement()
	case lexer.Type:
		return p.typeStatement()
	case lexer.LAngle:
		return p.genericStatement()
	case lexer.Return:
		return p.returnStatement()
	case lexer.Extend:
		return p.extendStatement()
	default:
		return p.exprStatement()
	}
}

func (p *Parser) varDefStatement() Stmt {
	p.consume(lexer.Let, "Expected variable declaration to start with let")

	properties := p.parseProperties(lexer.Mut, lexer.Lazy, lexer.Open)
	mut := properties[0]
	lazy := properties[1]
	open := properties[2]

	id := p.consume(lexer.Identifier, "Expected identifier for variable declaration")
	var typ Type
	if p.match(lexer.Colon) {
		typ = p.typeContract()
	}

	switch p.peek().TokenType {
	case lexer.LParen:
		p.insertBlankType(p.current, lexer.Equal)
		break
	case lexer.Arrow:
		p.insertBlankType(p.current, lexer.Equal, lexer.LParen, lexer.RParen)
		break
	}
	p.consume(lexer.Equal, "Expected Equal on variable declaration")
	expr := p.expression()

	return VarDefStmt{
		Mutable:    mut,
		Lazy:       lazy,
		Open:       open,
		Identifier: string(id.Data),
		Type:       typ,
		Value:      expr,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(lexer.While, "Expected while at beginning of while loop")
	expr := p.expression()
	body := p.blockStatement()
	return WhileStmt{
		Condition: expr,
		Body:      body,
	}
}

func (p *Parser) ifStatement() (stmt Stmt) {
	p.consume(lexer.If, "Expected if at beginning of if statement")
	condition := p.logicalOr()
	p.cleanNewLines()
	mainBranch := p.blockStatement()
	p.cleanNewLines()

	var elseBranch Stmt
	if p.match(lexer.Else) {
		if p.check(lexer.If) {
			elseBranch = p.ifStatement()
		} else {
			elseBranch = p.blockStatement()
		}
	}
	stmt = IfElseStmt{
		Condition:  condition,
		MainBranch: mainBranch,
		ElseBranch: elseBranch,
	}
	return
}

func (p *Parser) blockStatement() BlockStmt {
	result := make([]Stmt, 0)
	errors := make([]ParseError, 0)
	p.consume(lexer.LBrace, "Expected { at beginning of block")
	p.cleanNewLines()
	for !p.check(lexer.RBrace) {
		declaration := p.blockedDeclaration(&errors)
		result = append(result, declaration)
	}
	p.consume(lexer.RBrace, "Expected } at end of block")
	if len(errors) > 0 {
		panic(errors)
	}
	return BlockStmt{Stmts: result}
}

func (p *Parser) blockedDeclaration(errors *[]ParseError) (s Stmt) {
	defer p.handleError(errors)
	s = p.declaration()

	//This is no longer guaranteed as if statements clean any new lines while looking for an else branch
	//nxt := p.peek()
	//if nxt.TokenType != lexer.NEWLINE && nxt.TokenType != lexer.RBrace {
	//	panic("Expected newline after declaration in block")
	//}
	p.cleanNewLines()

	return s
}

func (p *Parser) structStatement() Stmt {
	p.consume(lexer.Struct, "Expected struct start to begin with `struct` keyword")
	return StructDefStmt{
		Identifier:   string(p.consume(lexer.Identifier, "Expected identifier after `struct` keyword").Data),
		StructFields: p.structFields(),
	}
}

func (p *Parser) returnStatement() Stmt {
	p.consume(lexer.Return, "Expected return")
	var expr Expr
	if p.peek().TokenType != lexer.NEWLINE {
		expr = p.expression()
	}
	return ReturnStmt{Returning: expr}
}

func (p *Parser) exprStatement() Stmt {
	return ExpressionStmt{Expr: p.expression()}
}

func (p *Parser) extendStatement() Stmt {
	p.consume(lexer.Extend, "Expected 'extend'")
	id := p.consumeValidIdentifier("Expected struct name to extend")
	alias := "this" //
	next := p.peek()
	if next.TokenType == lexer.As {
		p.advance()
		alias = string(p.consume(lexer.Identifier, "Expected identifier for extend alias").Data)
	}
	return ExtendStmt{
		Identifier: string(id.Data),
		Body:       p.blockStatement(),
		Alias:      alias,
	}
}
