package parser

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/statements"
	"devZ/lox/internal/tokens"
)

const CTX = reporters.PARSING

type Parser struct {
	tokens   []*tokens.Token
	current  int
	loop_lvl int
	accum    *reporters.Accumulator
	Printer  *reporters.PrettyPrinter
}

func New(tks []*tokens.Token, a *reporters.Accumulator) *Parser {
	return &Parser{tokens: tks, accum: a, current: 0}
}

func (p *Parser) Parse() []statements.Stmt {
	stmts := make([]statements.Stmt, 0)
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().ID() == tokens.EOF
}

func (p *Parser) peek() *tokens.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() *tokens.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) check(tid tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ID() == tid
}

func (p *Parser) match(tids ...tokens.TokenType) bool {
	for _, tid := range tids {
		if p.check(tid) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) declaration() statements.Stmt {
	var stmt statements.Stmt
	err := p.accum.LastError()
	if p.match(tokens.VAR) {
		stmt = p.varDeclaration()
	} else {
		stmt = p.statement()
	}
	if p.accum.LastError() != err {
		p.synchronize()
		return nil
	}
	return stmt
}

func (p *Parser) varDeclaration() statements.Stmt {
	name := p.consume(tokens.IDENTIFIER, reporters.Expectation("a variable name", "after", "var"))
	var initializer expressions.Expr
	if p.match(tokens.EQUAL) {
		initializer = p.expression()
	}
	p.consume(tokens.SEMICOLON, reporters.Expectation("';'", "after", "variable declaration"))
	return statements.NewVarStmt(name, initializer)
}

func (p *Parser) loopStatement() statements.Stmt {
	p.loop_lvl += 1
	stmt := p.statement()
	p.loop_lvl -= 1
	return stmt
}

func (p *Parser) statement() statements.Stmt {
	if p.match(tokens.BREAK) {
		if p.loop_lvl > 0 {
			return p.breakStatement()
		}
		p.error("'break' outside of loop body")
		return nil
	}
	if p.match(tokens.FOR) {
		return p.forStatement()
	}
	if p.match(tokens.IF) {
		return p.ifStatement()
	}
	if p.match(tokens.WHILE) {
		return p.whileStatement()
	}
	if p.match(tokens.PRINT) {
		return p.printStatement()
	}
	if p.match(tokens.LEFT_BRACE) {
		return p.blockStatement()
	}
	return p.exprStatement()
}

func (p *Parser) breakStatement() statements.Stmt {
	tok := p.previous()
	p.consume(tokens.SEMICOLON, reporters.Expectation("';'", "after", "'break'"))
	return statements.NewBreak(tok)
}

func (p *Parser) forStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, reporters.Expectation("'('", "after", "'for'"))
	var init statements.Stmt
	if p.match(tokens.SEMICOLON) {
		init = nil
	} else if p.match(tokens.VAR) {
		init = p.varDeclaration()
	} else {
		init = p.exprStatement()
	}
	var cnd, incr expressions.Expr
	if !p.check(tokens.SEMICOLON) {
		cnd = p.expression()
	}
	p.consume(tokens.SEMICOLON, reporters.Expectation("';'", "after", "loop condition"))
	if !p.check(tokens.RIGHT_PAREN) {
		incr = p.expression()
	}
	p.consume(tokens.RIGHT_PAREN, reporters.Expectation("')'", "after", "for clauses"))
	bd := p.loopStatement()
	if incr != nil {
		bd = statements.NewBlock([]statements.Stmt{bd, statements.NewExpression(incr)})
	}
	if cnd == nil {
		cnd = expressions.NewLiteral(true)
	}
	bd = statements.NewWhile(cnd, bd)
	if init != nil {
		bd = statements.NewBlock([]statements.Stmt{init, bd})
	}
	return bd
}

func (p *Parser) printStatement() statements.Stmt {
	expr := p.expression()
	p.consume(tokens.SEMICOLON, reporters.Expectation("';'", "after", "value"))
	return statements.NewPrint(expr)
}

func (p *Parser) blockStatement() statements.Stmt {
	stmts := make([]statements.Stmt, 0)
	for !p.check(tokens.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(tokens.RIGHT_BRACE, reporters.Expectation("'}'", "after", "block"))
	return statements.NewBlock(stmts)
}

func (p *Parser) ifStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, reporters.Expectation("'('", "after", "if"))
	condition := p.expression()
	p.consume(tokens.RIGHT_PAREN, reporters.Expectation("')'", "after", "if condition"))
	thenBr := p.statement()
	var elseBr statements.Stmt
	if p.match(tokens.ELSE) {
		elseBr = p.statement()
	}
	return statements.NewIf(condition, thenBr, elseBr)
}

func (p *Parser) whileStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, reporters.Expectation("'('", "after", "while"))
	cnd := p.expression()
	p.consume(tokens.RIGHT_PAREN, reporters.Expectation("')'", "after", "while condition"))
	bd := p.loopStatement()
	return statements.NewWhile(cnd, bd)
}

func (p *Parser) exprStatement() statements.Stmt {
	expr := p.expression()
	p.consume(tokens.SEMICOLON, reporters.Expectation("';'", "after", "expression"))
	return statements.NewExpression(expr)
}

func (p *Parser) expression() expressions.Expr {
	return p.assignment()
}

func (p *Parser) assignment() expressions.Expr {
	expr := p.logical()
	if p.match(tokens.EQUAL) {
		eqs := p.previous()
		val := p.assignment()
		exVar, ok := expr.(*expressions.VarExpr)
		if ok {
			return expressions.NewAssignment(exVar.Name(), val)
		}
		p.errorAt("invalid assignment target", eqs.Line())
	}
	return expr
}

func (p *Parser) logical() expressions.Expr {
	expr := p.equality()
	for p.match(tokens.OR, tokens.AND) {
		op := p.previous()
		right := p.equality()
		expr = expressions.NewLogical(expr, op, right)
	}
	return expr
}

func (p *Parser) equality() expressions.Expr {
	ex := p.comparison()
	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		op := p.previous()
		nex := p.comparison()
		ex = expressions.NewBinary(ex, op, nex)
	}
	return ex
}

func (p *Parser) comparison() expressions.Expr {
	ex := p.term()
	for p.match(tokens.GREATER, tokens.GREATER_EQUAL, tokens.LESS, tokens.LESS_EQUAL) {
		op := p.previous()
		nex := p.term()
		ex = expressions.NewBinary(ex, op, nex)
	}
	return ex
}

func (p *Parser) term() expressions.Expr {
	ex := p.factor()
	for p.match(tokens.MINUS, tokens.PLUS) {
		op := p.previous()
		nex := p.factor()
		ex = expressions.NewBinary(ex, op, nex)
	}
	return ex
}

func (p *Parser) factor() expressions.Expr {
	ex := p.unary()
	for p.match(tokens.SLASH, tokens.STAR) {
		op := p.previous()
		nex := p.factor()
		ex = expressions.NewBinary(ex, op, nex)
	}
	return ex
}

func (p *Parser) unary() expressions.Expr {
	if p.match(tokens.BANG, tokens.MINUS) {
		op := p.previous()
		right := p.unary()
		return expressions.NewUnary(op, right)
	}
	return p.call()
}

func (p *Parser) call() expressions.Expr {
	expr := p.primary()
	for true {
		if p.match(tokens.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(cle expressions.Expr) expressions.Expr {
	args := make([]expressions.Expr, 0)
	if !p.check(tokens.RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(tokens.COMMA) {
			args = append(args, p.expression())
			if len(args) >= 255 {
				p.error("call cannot have more than 255 arguments")
			}
		}
	}
	paren := p.consume(tokens.LEFT_PAREN, reporters.Expectation("')", "after", "arguments"))
	return expressions.NewCall(cle, paren, args)
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.FALSE) {
		return expressions.NewLiteral(false)
	}
	if p.match(tokens.TRUE) {
		return expressions.NewLiteral(true)
	}
	if p.match(tokens.NIL) {
		return expressions.NewLiteral(nil)
	}
	if p.match(tokens.NUMBER, tokens.STRING) {
		return expressions.NewLiteral(p.previous().Literal())
	}
	if p.match(tokens.IDENTIFIER) {
		return expressions.NewVarExpr(p.previous())
	}
	if p.match(tokens.LEFT_PAREN) {
		ex := p.expression()
		p.consume(tokens.RIGHT_PAREN, reporters.Expectation("')'", "after", "expression"))
		return ex
	}
	p.error(reporters.Expectation("expression", "", ""))
	return nil
}

func (p *Parser) consume(tid tokens.TokenType, msg string) *tokens.Token {
	if p.check(tid) {
		return p.advance()
	}
	p.error(msg)
	return nil
}

func (p *Parser) error(msg string) {
	p.errorAt(msg, p.peek().Line())
}

func (p *Parser) errorAt(msg string, loc int) {
	p.accum.AddError(loc, msg, CTX)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().ID() == tokens.SEMICOLON {
			return
		}
		switch p.peek().ID() {
		case tokens.CLASS, tokens.FUN, tokens.VAR, tokens.FOR, tokens.IF,
			tokens.WHILE, tokens.PRINT, tokens.RETURN:
			return
		}
		p.advance()
	}
}
