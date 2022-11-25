// Code generated by walklox/internal/tools. DO NOT EDIT.

package expressions

import (
	"devZ/lox/internal/tokens"
)

type Expr interface {
	Accept(v Visitor) interface{}
}

type Visitor interface {
	VisitBinary(bi *Binary) interface{}
	VisitGrouping(gr *Grouping) interface{}
	VisitLiteral(li *Literal) interface{}
	VisitUnary(un *Unary) interface{}
	VisitVarExpr(vr *VarExpr) interface{}
	VisitAssignment(as *Assignment) interface{}
}

type Binary struct {
	left     Expr
	operator *tokens.Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value interface{}
}

type Unary struct {
	operator *tokens.Token
	right    Expr
}

type VarExpr struct {
	name *tokens.Token
}

type Assignment struct {
	name  *tokens.Token
	value Expr
}

func (bi *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinary(bi)
}

func (bi *Binary) Left() Expr {
	return bi.left
}

func (bi *Binary) Operator() *tokens.Token {
	return bi.operator
}

func (bi *Binary) Right() Expr {
	return bi.right
}

func (bi *Binary) SetLeft(lf Expr) {
	bi.left = lf
}

func (bi *Binary) SetOperator(op *tokens.Token) {
	bi.operator = op
}

func (bi *Binary) SetRight(rt Expr) {
	bi.right = rt
}

func (gr *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGrouping(gr)
}

func (gr *Grouping) Expression() Expr {
	return gr.expression
}

func (gr *Grouping) SetExpression(ex Expr) {
	gr.expression = ex
}

func (li *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteral(li)
}

func (li *Literal) Value() interface{} {
	return li.value
}

func (li *Literal) SetValue(val interface{}) {
	li.value = val
}

func (un *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnary(un)
}

func (un *Unary) Operator() *tokens.Token {
	return un.operator
}

func (un *Unary) Right() Expr {
	return un.right
}

func (un *Unary) SetOperator(op *tokens.Token) {
	un.operator = op
}

func (un *Unary) SetRight(rt Expr) {
	un.right = rt
}

func (vr *VarExpr) Accept(v Visitor) interface{} {
	return v.VisitVarExpr(vr)
}

func (vr *VarExpr) Name() *tokens.Token {
	return vr.name
}

func (vr *VarExpr) SetName(nm *tokens.Token) {
	vr.name = nm
}

func (as *Assignment) Accept(v Visitor) interface{} {
	return v.VisitAssignment(as)
}

func (as *Assignment) Name() *tokens.Token {
	return as.name
}

func (as *Assignment) Value() Expr {
	return as.value
}

func (as *Assignment) SetName(nm *tokens.Token) {
	as.name = nm
}

func (as *Assignment) SetValue(vl Expr) {
	as.value = vl
}

func NewBinary(lf Expr, op *tokens.Token, rt Expr) *Binary {
	return &Binary{
		left:     lf,
		operator: op,
		right:    rt,
	}
}

func NewGrouping(ex Expr) *Grouping {
	return &Grouping{
		expression: ex,
	}
}

func NewLiteral(val interface{}) *Literal {
	return &Literal{
		value: val,
	}
}

func NewUnary(op *tokens.Token, rt Expr) *Unary {
	return &Unary{
		operator: op,
		right:    rt,
	}
}

func NewVarExpr(nm *tokens.Token) *VarExpr {
	return &VarExpr{
		name: nm,
	}
}

func NewAssignment(nm *tokens.Token, vl Expr) *Assignment {
	return &Assignment{
		name:  nm,
		value: vl,
	}
}
