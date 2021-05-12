package main

import (
    "os"
    "fmt"
    "bufio"
    "testing"
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

    parser := process(source)

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
        switch val := interpreter.blks[0].stck[len(interpreter.blks[0].stck) - 1].(type) {
        case String:
            fmt.Println(string(val))
        default:
            fmt.Printf("%v\n", val)
        }
    }
}

func process(source string) *Parser {
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

    return parser
}

func test(test *testing.T, source string) {
    p := process(source)
    c := 0
    f := 0

    i := &Interpreter {
        tics: p.tics,
        blks: p.blks,
    }

    for len(i.tics) > 0 {
        t := i.Interpret()

        if t.tok == FIN && i.comm != "" {
            if len(i.blks[t.dep].stck) == 0 {
                test.Errorf("%s expected %s at %s; got nil", source, i.comm, t.pos)
                f++
            } else {
                val := i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1]

                if fmt.Sprintf("%v", val) != i.comm {
                    test.Errorf("%s expected %s at %s; got %v", source, i.comm, t.pos, val)
                    f++
                } else {
                    c++
                }
            }

            i.comm = ""
        }
    }

    fmt.Printf("%s passed %d of %d tests\n", source, c, c + f)
}
