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

    parser := process(reader, NewBlock())

    if debug {
        for _, term := range parser.blk.toks {
            fmt.Printf("%v\n", term)
        }

        return
    }

    switch x := parser.blk.Run(Boolean(true), Boolean(false), Null { }, stdin()).(type) {
    case Null:
    case String:
        fmt.Printf("%v\n", string(x))
    default:
        fmt.Printf("%v\n", x)
    }
}

func stdin() interface{} {
    stat, _ := os.Stdin.Stat()

    if (stat.Mode() & os.ModeCharDevice) != 0 {
        return Null { }
    }

    out := Array { }
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        out = append(out, String(scanner.Text()))
    }

    switch len(out) {
    case 1:
        return out[0]
    default:
        return out
    }
}

func open(source string) *bufio.Reader {
    file, err := os.Open(source)

    if err != nil {
        panic(err)
    }

    return bufio.NewReader(file)
}

func process(rdr *bufio.Reader, blk *Block) *Parser {
    lexer := &Lexer {
        pos: Position { row: 1, col: 0 },
        rdr: rdr,
        opn: true,
    }

    for lexer.opn {
        lexer.Lexify()
    }

    if len(lexer.toks) > 0 && lexer.toks[len(lexer.toks) - 1].tok != FIN {
        lexer.toks = append(lexer.toks, &Token { tok: FIN, lit: "" })
    }

    parser := &Parser {
        lexr: lexer,
        blk: blk,
    }

    for len(parser.lexr.toks) > 0 {
        parser.Parse()
    }

    return parser
}

func test(test *testing.T, source string) {
    r := open(source)
    p := process(r, NewBlock())
    c := 0
    f := 0

    if val, ok := p.blk.Run(Boolean(true), Boolean(false), Null { }, stdin()).(Array); ok {
        for i, _ := range val {
            if i < len(p.lexr.coms) && fmt.Sprintf("%v", val[i]) != p.lexr.coms[i] {
                test.Errorf("%s expected %s at index %d; got %v", source, p.lexr.coms[i], i, val[i])
                f++
            } else {
                c++
            }
        }
	} else {
		test.Errorf("%s is a malformed test file", source)
	}

    fmt.Printf("%s passed %d of %d tests\n", source, c, c + f)
}
