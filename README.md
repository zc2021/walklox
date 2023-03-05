## Go Walk Lox
This is a Go implementation of the tree-walk interpreter for Lox, following Robert Nystrom's (excellent!) [__Crafting Interpreters__](https://www.craftinginterpreters.com)

If you are interested in a discussion of general architecture choices and motivational theory, please refer to [section III of __Crafting Interpreters__](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html). For an overview of the Lox language, refer to [section II of the same](https://www.craftinginterpreters.com/the-lox-language.html).

The canonical Lox tree-walk interpreter is written in Java, and (of course) the canonical implementation is designed around Java's class system and idiosyncracies. Below are some notes on implementing in Go.

## General Approach
Go is not straightforwardly object-oriented in the same way as Java; tangibly, there is no such construct as a Go class. Structs and interfaces can be used to achieve something resembling an instantiable Java Class, but attempting to exactly reproduce an inheritance hierarchy quickly becomes cumbersome. Thus, there are context-sensitive decisions to be made anywhere the canonical implementation leans on Java's Object or Class behavior. 

There is one other notable quirk of Go which precipitates significant departure from the Lox canon: error handling. Go's "errors as data" philosophy is not _necessarily_ at odds with JLox's use of exceptions for control flow in parsing and interpreting, but missing throw/catch exception mechanics does once again call for some retooling in the translation.

In facing these crossroads, I have tried to hew as closely as comfortable to the spirit of the original, while producing somewhat idiomatic Go code. Major departures from the canonical implementation are discussed in the relevant section of the notes.

## Project Structure
```
walklox
├── cmd
│   └── walklox
│       └── walklox.go
└── internal
    ├── environment
    │   ├── core_ops.go
    │   └── env.go
    ├── expressions
    │   └── expr_structs_ints.go
    ├── interpreter
    │   ├── callable.go
    │   ├── interpreter.go
    │   ├── native_fns.go
    │   └── operators.go
    ├── parser
    │   └── parser.go
    ├── reporters
    │   ├── accumulator.go
    │   └── pretty_printer.go
    ├── scanner
    │   └── scanner.go
    ├── statements
    │   └── stmt_structs_ints.go
    ├── tokens
    │   ├── tokens.go
    │   └── tokentype_string.go
    └── tools
        ├── gen_expressions.go
        ├── gen_statements.go
        ├── generator.go
        ├── helpers.go
        ├── meta_structs_methods.go
        └── pkg_templates.go
```

## Main (and the Accumulator)
JLox's Lox class, `lox/Lox.java` containing the entrypoint `Lox.main()`, handles the CLI, error reporting, and code input via Lox script or REPL. 
In GLox, this is replaced by two packages: `main` in `/cmd/walklox`, which contains the CLI entrypoint `main()` and code input functions; and the `reporters` package[^rptnote], which contains the `Accumulator` struct with associated types and helpers. The `reporters` package additionally provides the `ErrCtx` integer type, which is used as an exit value and by other Lox packages to communicate the context for any errors encountered.[^errnote]

Fundamentally, of course, the behavior of `main()` is unchanged.[^cmdnote] Instead of maintaining a global error state (`Lox.hadError`), an `Accumulator` tracks errors in all stages of Lox processing[^stgnote]. This accumulator is checked at the end of each stage, and processing stops if any errors were encountered during the previous stage.

Structurally, JLox consists of various classes within a single Java package; this means each processor stage has access to the main Lox class without issue. Error handling takes advantage of this fact by using methods on the main `Lox` class to manage a `Lox.hadError` boolean. 
GLox roughly replicates the JLox class structure using Go packages. This meant that any code shared across stages, such as error handling, would have to be extracted to an independent package.[^imptnote]

## Tokens
An individual token is implemented as a struct with fields for token type, source line location, represented lexeme, and literal value.

Go's answer to Java's `enum` type is the `const iota` identifier; its use, along with the aliased `int` type `TokenType`, as a value for token "type" in GLox's package tokens is directly inspired by the Go source.

## Operators
Throughout GLox, there are minor changes to the implementation of operator expressions as a set. In all packages, these changes are confined to files named `operators.go`. 

[^cmdnote]: When invoked from the command line with: 
    - more than one argument, the program prints a usage hint to stdout and exits; 
    - a single argument, expected to be a path to a readable Lox script file, the program reads the file and executes the script;
    - no arguments, the program starts an in-shell REPL.

[^rptnote]: Currently, package reporters contains exit codes, the Accumulator, and the PrettyPrinter. It could easily be extended to include other reporting and reflective structs.

[^errnote]: Error codes as given in UNIX's [sysexits.h](https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html), as linked by Nystrom.

[^stgnote]: Scanning, parsing, interpreting

[^imptnote]: `main` cannot be imported in Go; even if it could, doing so in a package imported by `main`, like any of the processing stages, would cause an illegal import cycle.