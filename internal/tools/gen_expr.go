// The following build directive is necessary to make the package coherent:

//go:build ignore

//go:generate go run gen_expr.go

// This program generates Expr structs. It can be invoked using go generate.
package main

import (
	"os"
)

func main() {
	os.Create("../expressions/expr_structs.go")
}
