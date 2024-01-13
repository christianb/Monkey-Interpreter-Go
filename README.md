# Writing An Interpreter In Go
This is an interpreter written in GoLang for the Monkey programming language. 

All code and samples are based on the book ["Writing An Interpreter in Go"](https://interpreterbook.com) by Thorsten Ball.

## Lexer
The lexer transforms meaningless string into a (flat) list of things like "number literal", "string literal", "identifier", or "operator", called Tokens. It also recognizing reserved identifiers (keywords) and discarding whitespace. 

## Parser
The parser is turning a stream of Tokens, produced by the lexer, into a parse tree (Abstract Syntax Tree) representing the structure of the parsed language.

### Run Tests
To execute all tests in all packages: `go test ./...`

### Run the REPL
To start the REPL: `go run main.go`

### About The Monkey Programming Language
* [https://interpreterbook.com](https//interpreterbook.com)
* [The Lost Chapter](https//interpreterbook.com/lost)

### Links
* [Writing A Compiler In Go](https://compilerbook.com)
* [What is a Lexer, Anyway?](https://dev.to/cad97/what-is-a-lexer-anyway-4kdo)