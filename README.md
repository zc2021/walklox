## Go Walk Lox
This is a Go implementation of the tree-walk interpreter for Lox, following Robert Nystrom's (excellent!) [__Crafting Interpreters__](https://www.craftinginterpreters.com)

If you are interested in a discussion of general architecture choices and motivational theory, please refer to [section III of __Crafting Interpreters__](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html). For an overview of the Lox language, refer to [section II of the same](https://www.craftinginterpreters.com/the-lox-language.html).

The canonical Lox tree-walk interpreter is written in Java, and (of course) the canonical implementation is designed around Java's class system and idiosyncracies. Below are some notes on implementing in Go.

## General Approach
Go is not straightforwardly object-oriented in the same way as Java; tangibly, there is no such construct as a Go class. Thus, there are decisions to be made anywhere the canonical implementation leans on Java's Object or Class behavior. 
I have tried to hew as closely as possible to the spirit of the original, while producing idiomatic (to a newb's eye) Go code. Major discrepancies are discussed in the relevant section of the notes.

There is one other notable quirk of Go: error handling. 

## Tokens
