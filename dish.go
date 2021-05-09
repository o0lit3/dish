package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
)

const (
    NOP = iota
    BIT // Binary register stores (not used during lexing)
    NUM // Numeric lexemes or register stores
    STR // String lexemes or register stores
    OBJ // Dynamic register stores (not used during lexing)
    VAR // Identifier lexemes or register stores (for assigment only)
    OP1 // Unary operators (!, ~, +, -, ++, --)
    OP2 // Binary operators
    OPX // Method operators (alphanumeric)
    OPA // Assignment operators (= += -= *= /= %= &= ^= |=)
    BLK // Grouping lexemes {} [] ()
    ERR // Unrecognized lexemes
    FIN // Statement ending lexemes (newline or comma)
    EOF // End of file
)

func main() {
    debug := false
    source := os.Args[1]

    for _, flag := range os.Args {
        switch flag {
        case "-debug":
            debug = true
            source = os.Args[2]
        }
    }

    file, err := os.Open(source)

    if err != nil {
        panic(err)
    }

    lexer := &Lexer {
        pos: Position { line: 1, column: 0},
        rdr: bufio.NewReader(file),
        opn: true,
    }

    for lexer.opn {
        lexer.Lexify()
    }

    parser := &Parser {
        depth: 0,
        toks: lexer.toks,
        blks: []Block{0: Block {
            style: "()",
            vars: map[string]interface{}{
                "true": true,
                "false": false,
            },
        }},
    }

    for len(parser.toks) > 0 {
        parser.Parse()
    }

    if debug {
        for _, term := range parser.terms {
            fmt.Println(strings.Repeat("\t", term.depth) + term.String())
        }

        return
    }

    interpreter := &Interpreter {
        terms: parser.terms,
        blks: parser.blks,
    }

    for len(interpreter.terms) > 0 {
        interpreter.Interpret()
    }

    if (len(interpreter.blks) > 0 && len(interpreter.blks[0].stack) > 0) {
        fmt.Printf("%v\n", interpreter.blks[0].stack[len(interpreter.blks[0].stack) - 1])
    }
}
