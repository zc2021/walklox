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
	printer *reporters.PrettyPrinter
}

func New(tks []*tokens.Token, a *reporters.Accumulator) *Parser {
	return &Parser{tokens: tks, accum: a}
}

func (p *Parser) AddPrinter(pp *reporters.PrettyPrinter) {
	p.printer = pp
}

func (p *Parser) isAtEnd() bool {
	return p.peek().ID() == tokens.EOF
}

func (p *Parser) peek() *tokens.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() *tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(tid tokens.TokID) bool {
	if p.isAtEnd() {
		return false
	}
	return p.previous().ID() == tid
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
