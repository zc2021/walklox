// Package tools provides internal tools for development of GoWalkLox.
package tools

//go:generate stringer -type=TokID ../tokens
//go:generate go run gen_expressions.go
//go:generate gofmt -w ../expressions/expr_structs_ints.go
