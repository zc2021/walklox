package parser

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/tokens"
)

type Parser struct {
	tokens  []*tokens.Token
	current int
	accum   *reporters.Accumulator
	Printer *reporters.PrettyPrinter
}

func New(tks []*tokens.Token, a *reporters.Accumulator) *Parser {
	return &Parser{tokens: tks, accum: a, current: 0}
}

func (p *Parser) Parse() expressions.Expr {
	return p.expression()
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

func (p *Parser) check(tid tokens.TokID) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ID() == tid
}

func (p *Parser) match(tids ...tokens.TokID) bool {
	for _, tid := range tids {
		if p.check(tid) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expression() expressions.Expr {
	return p.equality()
}

func (p *Parser) equality() expressions.Expr {
	ex := p.comparison()
	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		op := p.previous()
		nex := p.comparison()
		ex = &expressions.Binary{
			Left:     ex,
			Operator: op,
			Right:    nex,
		}
	}
	return ex
}

func (p *Parser) comparison() expressions.Expr {
	ex := p.term()
	for p.match(tokens.GREATER, tokens.GREATER_EQUAL, tokens.LESS, tokens.LESS_EQUAL) {
		op := p.previous()
		nex := p.term()
		ex = &expressions.Binary{
			Left:     ex,
			Operator: op,
			Right:    nex,
		}
	}
	return ex
}

func (p *Parser) term() expressions.Expr {
	ex := p.factor()
	for p.match(tokens.MINUS, tokens.PLUS) {
		op := p.previous()
		nex := p.factor()
		ex = &expressions.Binary{
			Left:     ex,
			Operator: op,
			Right:    nex,
		}
	}
	return ex
}

func (p *Parser) factor() expressions.Expr {
	ex := p.unary()
	for p.match(tokens.SLASH, tokens.STAR) {
		op := p.previous()
		nex := p.factor()
		ex = &expressions.Binary{
			Left:     ex,
			Operator: op,
			Right:    nex,
		}
	}
	return ex
}

func (p *Parser) unary() expressions.Expr {
	if p.match(tokens.BANG, tokens.MINUS) {
		op := p.previous()
		right := p.unary()
		return &expressions.Unary{
			Operator: op,
			Right:    right,
		}
	}
	return p.primary()
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.FALSE) {
		return &expressions.Literal{Value: false}
	}
	if p.match(tokens.TRUE) {
		return &expressions.Literal{Value: true}
	}
	if p.match(tokens.NIL) {
		return &expressions.Literal{Value: nil}
	}
	if p.match(tokens.NUMBER, tokens.STRING) {
		return &expressions.Literal{Value: p.previous().Literal()}
	}
	if p.match(tokens.LEFT_PAREN) {
		ex := p.expression()
		p.consume(tokens.RIGHT_PAREN, "Expect ')' after expression.")
		return ex
	}
	p.error("Expect expression.")
	return nil
}

func (p *Parser) consume(tid tokens.TokID, msg string) *tokens.Token {
	if p.check(tid) {
		return p.advance()
	}
	p.error(msg)
	return nil
}

func (p *Parser) error(msg string) {
	p.accum.AddError(p.peek().Line(), msg)
	if !p.isAtEnd() {
		p.synchronize()
	}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().ID() == tokens.SEMICOLON {
			return
		}
		switch p.peek().ID() {
		case tokens.CLASS, tokens.FUN, tokens.VAR, tokens.FOR, tokens.IF, tokens.WHILE, tokens.PRINT, tokens.RETURN:
			return
		}
		p.advance()
	}
}
