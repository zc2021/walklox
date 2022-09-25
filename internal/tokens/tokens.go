// Package tokens defines the Token struct and associated methods, representing
// Lox tokens. Tokens are created by the scanner, and contain source-code
// context information.
// Elements of the implementation (TokID as an iota with string-slice,
// Token.String() method, etc) are taken almost directly from the Go language
// token implementation (thanks!).
package tokens

import "strconv"

type TokID int

const (
	EOF TokID = iota

	begin_single_char_toks
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // [
	RIGHT_BRACE // ]
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
	LEFT_BRACE:  "[",
	RIGHT_BRACE: "]",
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

type Token struct {
	ID TokID
}

func (t Token) String() string {
	s := ""
	if 0 <= t.ID && t.ID < TokID(len(tokens)) {
		s = tokens[t.ID]
	}
	if s == "" {
		s = "token (" + strconv.Itoa(int(t.ID)) + ")"
	}
	return s
}
