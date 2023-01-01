## Go Walk Lox
This is a Go implementation of the tree-walk interpreter for Lox, following Robert Nystrom's (excellent!) [__Crafting Interpreters__](https://www.craftinginterpreters.com)

If you are interested in a discussion of general architecture choices and motivational theory, please refer to [section III of __Crafting Interpreters__](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html). For an overview of the Lox language, refer to [section II of the same](https://www.craftinginterpreters.com/the-lox-language.html).

The canonical Lox tree-walk interpreter is written in Java, and (of course) the canonical implementation is designed around Java's class system and idiosyncracies. Below are some notes on implementing in Go.

## General Approach
Go is not straightforwardly object-oriented in the same way as Java; tangibly, there is no such construct as a Go class. Structs and interfaces can be used to achieve something resembling an instantiable Java Class, but attempting to exactly reproduce an inheritance hierarchy quickly becomes cumbersome. Thus, there are context-sensitive decisions to be made anywhere the canonical implementation leans on Java's Object or Class behavior. 

There is one other notable quirk of Go which precipitates significant departure from the Lox canon: error handling. Go's "errors are normal data" philosophy is not _necessarily_ at odds with JLox's use of exceptions for control flow in parsing and interpreting, but does once again call for some retooling in the translation.

In facing these crossroads, I have tried to hew as closely as comfortable to the spirit of the original, while producing idiomatic (to a newb's eye) Go code. Major departures from the canonical implementation are discussed in the relevant section of the notes.

## Main (and the Accumulator)
Fundamentally, of course, the behavior of `main()` is unchanged.[^cmdnote] I've introduced a reporters package[^rptnote] containing an Accumulator struct, which tracks errors as they are encountered. When `run(script []byte)` is called, a new Accumulator is created within the function body. At the end of each processing stage, the Accumulator is checked for errors; any accumulated errors are printed and the Accumulator is reset. If any errors were found, the program exits using a code provided by the package for the previous processing stage's CTX value (itself taken from the package `reporters`).[^errnote]

Introducing the accumulator served two immediate purposes: first, to more thoroughly formalize the separation of concerns between ["the code that _generates_ the errors [and] the code that _reports_ them"](https://www.craftinginterpreters.com/scanning.html#error-handling); and second, to avoid having the component packages attempting to access (and set!) global state.[^statenote] 

Aside from being (imho) poor form for a more formal project[^donthateme], it is not immediately clear that an exact replica of the canonical error state would be possible in Go. Package main cannot be imported; any value or structure that is accessed from both main and a second package must belong to either the second or a third package.[^pkgnote] This technical requirement provided the perfect excuse to fill in the sketch of error handling provided by the canonical implementation.

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

[^statenote]: It is worth noting that the canonical error tracking variable, `hadError`, is not _truly_ global as it is a public static member of the Lox class. This is as close as Java gets to a true global constant, and I still dislike the pattern.

[^donthateme]: The choice of a global error state completely makes sense in a book focused on introducing concerns of language design over precise technical implementation, similar to the use of static imports.

[^pkgnote]: Which is then, of course, imported to main for use.