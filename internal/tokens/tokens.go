// Package tokens defines the Token struct and associated methods, representing
// Lox tokens. Tokens are created by the scanner, and contain source-code
// context information.
// Elements of the implementation (TokID as an iota with string-slice,
// Token.String() method, etc) are taken almost directly from the Go language
// token implementation (thanks!).
package tokens

import (
	"errors"
	"fmt"
	"strconv"
)

type TokenType int

const (
	EOF TokenType = iota

	begin_single_char_toks
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // {
	RIGHT_BRACE // }
	COMMA       // ,
	DOT         // .
	MINUS       // -
	PLUS        // +
	SEMICOLON   // ;
	SLASH       // /
	STAR        // *
	end_single_char_toks

	begin_one_two_toks
	BANG          // !
	BANG_EQUAL    // !=
	EQUAL         // =
	EQUAL_EQUAL   // ==
	GREATER       // >
	GREATER_EQUAL // >=
	LESS          // <
	LESS_EQUAL    // <=
	end_one_two_toks

	begin_literal_toks
	IDENTIFIER // a user-supplied (nonreserved) identifier
	STRING     // a literal string
	NUMBER     // a literal number
	end_literal_toks

	begin_reserved_keywords
	AND    // logical and
	CLASS  // opens class definition block
	ELSE   // opens control flow if-block else branch
	FALSE  // boolean false
	FUN    // opens function definition block
	FOR    // opens control flow for loop block
	IF     // opens control flow if block
	NIL    // empty value
	OR     // logical or
	PRINT  // call to fundamental native print function
	RETURN // begins return statement
	SUPER  // call parent class method on child class instance
	THIS   // refers to current object within a method definition
	TRUE   // boolean true
	VAR    // begin variable declaration statement
	WHILE  // opens control flow while loop block
	end_reserved_keywords
)

var tokens = [...]string{
	EOF: "EOF",

	LEFT_PAREN:  "(",
	RIGHT_PAREN: ")",
	LEFT_BRACE:  "{",
	RIGHT_BRACE: "}",
	COMMA:       ",",
	DOT:         ".",
	MINUS:       "-",
	PLUS:        "+",
	SEMICOLON:   ";",
	SLASH:       "/",
	STAR:        "*",

	BANG:          "!",
	BANG_EQUAL:    "!=",
	EQUAL:         "=",
	EQUAL_EQUAL:   "==",
	GREATER:       ">",
	GREATER_EQUAL: ">=",
	LESS:          "<",
	LESS_EQUAL:    "<=",

	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",

	AND:    "AND",
	CLASS:  "CLASS",
	ELSE:   "ELSE",
	FALSE:  "FALSE",
	FUN:    "FUN",
	FOR:    "FOR",
	IF:     "IF",
	NIL:    "NIL",
	OR:     "OR",
	PRINT:  "PRINT",
	RETURN: "RETURN",
	SUPER:  "SUPER",
	THIS:   "THIS",
	TRUE:   "TRUE",
	VAR:    "VAR",
	WHILE:  "WHILE",
}

var Keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (tid TokenType) Valid() bool {
	return tid >= 0 && tid < TokenType(len(tokens))
}

type Token struct {
	tp      TokenType
	lexeme  string
	line    int
	literal interface{}
}

func New(tid TokenType, lex string, loc int, obj interface{}) (*Token, error) {
	if !tid.Valid() {
		return nil, errors.New("invalid TokID passed to NewToken")
	}
	return &Token{tp: tid, lexeme: lex, line: loc, literal: obj}, nil
}

func (t *Token) ID() TokenType {
	return t.tp
}

func (t *Token) Lexeme() string {
	return t.lexeme
}

func (t *Token) Line() int {
	return t.line
}

func (t *Token) Literal() interface{} {
	return t.literal
}

func (t *Token) String() string {
	s := ""
	if t.ID().Valid() {
		s = tokens[t.ID()]
	}
	if s == "" {
		s = "token (" + strconv.Itoa(int(t.ID())) + ")"
	}
	return s
}

func (t *Token) SetBiFunc() error {
	fn := binaryFuncs[t.tp]
	if fn == nil {
		return errors.New("Unrecognized token type for binary function.")
	}
	t.literal = fn
	return nil
}

func (t *Token) BiFunc(x, y interface{}) interface{} {
	return t.literal.(biFunc)(x, y)
}

func (t *Token) SetUnFunc() error {
	fn := unaryFuncs[t.tp]
	if fn == nil {
		return errors.New("Unrecognized token type for unary function.")
	}
	t.literal = fn
	return nil
}

func (t *Token) UnFunc(x interface{}) interface{} {
	return t.literal.(unFunc)(x)
}

type biFunc func(x, y interface{}) interface{}

func add(x, y interface{}) interface{} {
	n, isNum := x.(float64)
	if !isNum {
		return fmt.Sprintf("%s%s", x.(string), y.(string))
	}
	return n + y.(float64)
}

func subtract(x, y interface{}) interface{} {
	return x.(float64) - y.(float64)
}

func divide(x, y interface{}) interface{} {
	return x.(float64) / y.(float64)
}

func multiply(x, y interface{}) interface{} {
	return x.(float64) * y.(float64)
}

func greater(x, y interface{}) interface{} {
	return x.(float64) > y.(float64)
}

func less(x, y interface{}) interface{} {
	return x.(float64) < y.(float64)
}

func greaterEq(x, y interface{}) interface{} {
	return x.(float64) >= y.(float64)
}

func lessEq(x, y interface{}) interface{} {
	return x.(float64) <= y.(float64)
}

func checkVal[T any](y interface{}) interface{} {
	yt, ok := y.(T)
	if !ok {
		return nil
	}
	return yt
}

func equal(x, y interface{}) interface{} {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	var yt interface{}
	switch x.(type) {
	case bool:
		yt = checkVal[bool](y)
	case float64:
		yt = checkVal[float64](y)
	case string:
		yt = checkVal[string](y)
	default:
		return false
	}
	if yt == nil {
		return false
	}
	return yt == x
}

func notEqual(x, y interface{}) interface{} {
	return !equal(x, y).(bool)
}

var binaryFuncs = [...]biFunc{
	PLUS:          add,
	MINUS:         subtract,
	SLASH:         divide,
	STAR:          multiply,
	GREATER:       greater,
	LESS:          less,
	GREATER_EQUAL: greaterEq,
	LESS_EQUAL:    lessEq,
	EQUAL_EQUAL:   equal,
	BANG_EQUAL:    notEqual,
}

type unFunc func(x interface{}) interface{}

func isTruthy(x interface{}) interface{} {
	if x == nil {
		return false
	}
	xb, ok := x.(bool)
	if !ok {
		return true
	}
	return xb
}

func isNotTruthy(x interface{}) interface{} {
	return !isTruthy(x).(bool)
}

func negate(x interface{}) interface{} {
	return -1 * x.(float64)
}

var unaryFuncs = [...]unFunc{
	MINUS: negate,
	BANG:  isNotTruthy,
}
