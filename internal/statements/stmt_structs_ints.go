// Code generated by walklox/internal/tools. DO NOT EDIT.

package statements

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/tokens"
)

type Stmt interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitExpression(expst *Expression)
	VisitPrint(prnst *Print)
	VisitVarStmt(varst *VarStmt)
	VisitBlock(blk *Block)
}

type Expression struct {
	expression expressions.Expr
}

type Print struct {
	expression expressions.Expr
}

type VarStmt struct {
	name        *tokens.Token
	initializer expressions.Expr
}

type Block struct {
	statements []Stmt
}

func (expst *Expression) Accept(v Visitor) {
	v.VisitExpression(expst)
}

func (expst *Expression) Expression() expressions.Expr {
	return expst.expression
}

func (expst *Expression) SetExpression(ex expressions.Expr) {
	expst.expression = ex
}

func (prnst *Print) Accept(v Visitor) {
	v.VisitPrint(prnst)
}

func (prnst *Print) Expression() expressions.Expr {
	return prnst.expression
}

func (prnst *Print) SetExpression(ex expressions.Expr) {
	prnst.expression = ex
}

func (varst *VarStmt) Accept(v Visitor) {
	v.VisitVarStmt(varst)
}

func (varst *VarStmt) Name() *tokens.Token {
	return varst.name
}

func (varst *VarStmt) Initializer() expressions.Expr {
	return varst.initializer
}

func (varst *VarStmt) SetName(nm *tokens.Token) {
	varst.name = nm
}

func (varst *VarStmt) SetInitializer(init expressions.Expr) {
	varst.initializer = init
}

func (blk *Block) Accept(v Visitor) {
	v.VisitBlock(blk)
}

func (blk *Block) Statements() []Stmt {
	return blk.statements
}

func (blk *Block) SetStatements(sts []Stmt) {
	blk.statements = sts
}

func NewExpression(ex expressions.Expr) *Expression {
	return &Expression{
		expression: ex,
	}
}

func NewPrint(ex expressions.Expr) *Print {
	return &Print{
		expression: ex,
	}
}

func NewVarStmt(nm *tokens.Token, init expressions.Expr) *VarStmt {
	return &VarStmt{
		name:        nm,
		initializer: init,
	}
}

func NewBlock(sts []Stmt) *Block {
	return &Block{
		statements: sts,
	}
}