## Go Walk Lox
This is a Go implementation of the tree-walk interpreter for Lox, following Robert Nystrom's (excellent!) [__Crafting Interpreters__](https://www.craftinginterpreters.com)

 If you are interested in a discussion of general architecture choices and motivational theory, please refer to [section III of __Crafting Interpreters__](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html). For an overview of the Lox language, refer to [section II of the same](https://www.craftinginterpreters.com/the-lox-language.html).

The canonical Lox tree-walk interpreter is written in Java, and (of course) the canonical implementation is designed around Java's class system and idiosyncracies. Below are some notes on implementing in Go.

## General Approach
 Go is not straightforwardly object-oriented in the same way as Java; tangibly, there is no such construct as a Go class. This implies some consideration to be made in any place the canonical interpreter implements a Class. In some cases - such as the "pure data" Expr - Go's structs are a better fit than the Java Class. Other areas, like Tokens, feel a bit awkward without the help of inheritance. The aim on this first pass is to reflect the canonical implementation wherever possible.

There is one notable exception: error handling. 

## Tokens
