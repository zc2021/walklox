// Package tools provides internal tools for development of GoWalkLox.
package tools

//go:generate stringer -type=TokenType ../tokens
//go:generate go run gen_expressions.go
//go:generate gofmt -w ../expressions/expr_structs_ints.go
//go:generate go run gen_statements.go
//go:generate gofmt -w ../statements/stmt_structs_ints.go
