package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "testing"
)

func main() {
    var reader *bufio.Reader

    source := os.Args[1]
    debug := false

    for i, flag := range os.Args {
        switch flag {
        case "-debug":
            debug = true
            source = os.Args[i + 1]
        case "-e":
            source = ""
            reader = bufio.NewReader(strings.NewReader(os.Args[i + 1]))
        }
    }

    if source != "" {
        reader = open(source)
    }

    parser := process(reader)

    if debug {
        for _, term := range parser.tics {
            fmt.Printf("%v\n", term)
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
            fmt.Printf("%v\n", string(val))
        default:
            fmt.Printf("%v\n", val)
        }
    }
}

func open(source string) *bufio.Reader {
    file, err := os.Open(source)

    if err != nil {
        panic(err)
    }

    return bufio.NewReader(file)
}

func process(rdr *bufio.Reader) *Parser {
    lexer := &Lexer {
        pos: Position { line: 1, column: 0},
        rdr: rdr,
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
                "true": Boolean(true),
                "false": Boolean(false),
                "null": Null { },
            },
        }},
    }

    for len(parser.toks) > 0 {
        parser.Parse()
    }

    return parser
}

func test(test *testing.T, source string) {
    r := open(source)
    p := process(r)
    c := 0
    f := 0

    i := &Interpreter {
        tics: p.tics,
        blks: p.blks,
    }

    for len(i.tics) > 0 {
        t := i.Interpret()

        if t.tok == FIN && i.blks[t.dep].com != "" {
            var val interface{}

            if t.lit == "" && t.dep > 0 && len(i.blks[t.dep - 1].stck) > 0 {
                val = i.blks[t.dep - 1].stck[len(i.blks[t.dep - 1].stck) - 1]

                if x, ok := val.(Array); ok && len(x) > 0 {
                    val = x[len(x) - 1]
                }
            } else if len(i.blks[t.dep].stck) > 0 {
                val = i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1]
            }

            if val == nil {
                test.Errorf("%s expected %s at %s; got nil", source, i.blks[t.dep].com, t.pos)
                f++
            } else {
                if fmt.Sprintf("%v", val) != i.blks[t.dep].com {
                    test.Errorf("%s expected %s at %s; got %v", source, i.blks[t.dep].com, t.pos, val)
                    f++
                } else {
                    c++
                }
            }

            i.blks[t.dep].com = ""
        }
    }

    fmt.Printf("%s passed %d of %d tests\n", source, c, c + f)
}
