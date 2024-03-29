// Package tokens defines the Token struct and associated methods, representing
// Lox tokens. Tokens are created by the scanner, and contain source-code
// context information.
package tokens

import (
	"errors"
	"strconv"
)

type TokenType int

const (
	EOF TokenType = iota

	begin_delimiters
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // {
	RIGHT_BRACE // }
	COMMA       // ,
	DOT         // .
	SEMICOLON   // ;
	end_delimiters

	begin_literals
	IDENTIFIER // a user-supplied (nonreserved) identifier
	STRING     // a literal string
	NUMBER     // a literal number
	end_literals

	begin_binary_operators
	PLUS          // +
	SLASH         // /
	STAR          // *
	EQUAL_EQUAL   // ==
	EQUAL         // =
	GREATER_EQUAL // >=
	GREATER       // >
	LESS_EQUAL    // <=
	LESS          // <
	BANG_EQUAL    // !=
	end_binary_operators

	begin_unary_operators
	BANG // !
	end_undary_operators

	begin_variadic_operators
	MINUS // -
	end_variadic_operators

	begin_reserved_keywords
	NIL    // empty value
	AND    // logical and
	OR     // logical or
	TRUE   // boolean true
	FALSE  // boolean false
	VAR    // begin variable declaration statement
	CLASS  // opens class definition block
	SUPER  // call parent class method on child class instance
	FUN    // opens function definition block
	RETURN // begins return statement
	THIS   // refers to current object within a method definition
	IF     // opens control flow if block
	ELSE   // opens control flow if-block else branch
	FOR    // opens control flow for loop block
	WHILE  // opens control flow while loop block
	BREAK  // exits current control flow structure
	PRINT  // call to fundamental native print function
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
	BREAK:  "BREAK",
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
	"break":  BREAK,
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

func (tt TokenType) Lexeme() string {
	return tokens[tt]
}
