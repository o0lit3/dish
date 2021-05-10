package main

import (
    "os"
    "fmt"
    "bufio"
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
        dep: 0,
        toks: lexer.toks,
        blks: []Block{0: Block {
            dim: VAL,
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
        for _, term := range parser.tics {
            fmt.Println(term.String())
        }

        return
    }

    interpreter := &Interpreter {
        tics: parser.tics,
        blks: parser.blks,
    }

    for len(interpreter.tics) > 0 {
        interpreter.Interpret()
    }

    if (len(interpreter.blks) > 0 && len(interpreter.blks[0].stck) > 0) {
        fmt.Printf("%v\n", interpreter.blks[0].stck[len(interpreter.blks[0].stck) - 1])
    }
}
