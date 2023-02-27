## Go Walk Lox
This is a Go implementation of the tree-walk interpreter for Lox, following Robert Nystrom's (excellent!) [__Crafting Interpreters__](https://www.craftinginterpreters.com)

If you are interested in a discussion of general architecture choices and motivational theory, please refer to [section III of __Crafting Interpreters__](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html). For an overview of the Lox language, refer to [section II of the same](https://www.craftinginterpreters.com/the-lox-language.html).

The canonical Lox tree-walk interpreter is written in Java, and (of course) the canonical implementation is designed around Java's class system and idiosyncracies. Below are some notes on implementing in Go.

## General Approach
Go is not straightforwardly object-oriented in the same way as Java; tangibly, there is no such construct as a Go class. Structs and interfaces can be used to achieve something resembling an instantiable Java Class, but attempting to exactly reproduce an inheritance hierarchy quickly becomes cumbersome. Thus, there are context-sensitive decisions to be made anywhere the canonical implementation leans on Java's Object or Class behavior. 

There is one other notable quirk of Go which precipitates significant departure from the Lox canon: error handling. Go's "errors are normal data" philosophy is not _necessarily_ at odds with JLox's use of exceptions for control flow in parsing and interpreting, but does once again call for some retooling in the translation.

In facing these crossroads, I have tried to hew as closely as comfortable to the spirit of the original, while producing somewhat idiomatic Go code. Major departures from the canonical implementation are discussed in the relevant section of the notes.

## Main (and the Accumulator)
JLox's Lox class, `lox/Lox.java` containing the entrypoint `Lox.main()`, handles the CLI, error reporting, and code input via Lox script or REPL. In GLox, this is replaced by package `main`, which contains `main()`, the CLI, and code input functions, and the `reporters` package[^rptnote], which contains the `Accumulator` struct with associated types and helpers. The `reporters` package additionally provides the `ErrCtx` integer type, which is used as an exit value and by other Lox packages to communicate the context for any errors encountered.[^errnote]

Fundamentally, of course, the behavior of `main()` is unchanged.[^cmdnote] Instead of maintaining a global error state (`hadError`), an `Accumulator` tracks errors in all stages of Lox processing[^stgnote]. This accumulator is checked at the end of each stage, and processing stops if any errors were encountered during the previous stage.



## Tokens
An individual token is implemented as a struct with fields for token type, source line location, represented lexeme, and literal value.

Go's answer to Java's `enum` type is the `const iota` identifier; its use, along with the aliased `int` type `TokenType`, as a value for token "type" in GWLox's package tokens is directly inspired by the Go source.

A Token has no exported fields and minimal setters, to reflect the use `final` of in the canonical class definition.



[^cmdnote]: When invoked from the command line with: 
    - more than one argument, the program prints a usage hint to stdout and exits; 
    - a single argument, expected to be a path to a readable Lox script file, the program reads the file and executes the script;
    - no arguments, the program starts an in-shell REPL.

[^rptnote]: Currently, package reporters contains exit codes, the Accumulator, and the PrettyPrinter. It could easily be extended to include other reporting and reflective structs.

[^errnote]: Error codes as given in UNIX's [sysexits.h](https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html), as linked by Nystrom.

[^stgnote]: Scanning, parsing, interpreting