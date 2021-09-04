package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
	"testing"
    "math/big"
    "io/ioutil"
    "encoding/json"
)

var program *Block
var argv Array

func main() {
    var reader *bufio.Reader

    index := 0
    source := os.Args[1]
    debug := false

    for i, flag := range os.Args {
        switch flag {
        case "-debug":
            debug = true
            source = os.Args[i + 1]
            index = i + 1
        case "-e":
            source = ""
            reader = bufio.NewReader(strings.NewReader(os.Args[i + 1]))
            index = i + 1
        }
    }

    if source != "" {
        reader = open(source)
    }

    for index + 1 < len(os.Args) {
        index += 1
        argv = append(argv, String(os.Args[index])) 
    }

    program = NewBlock()

    parser := process(reader, program)

    if debug {
        for _, term := range parser.blk.toks {
            fmt.Fprintf(os.Stderr, "%v\n", term)
        }

        return
    }

    switch x := parser.blk.Run(parser.blk.def...).(type) {
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
        return Null{ }
    }

    var data interface{}

    input, _ := ioutil.ReadAll(os.Stdin)

    if err := json.Unmarshal(input, &data); err == nil {
        return parse(data)
    }

    out := Array{ }

    for _, val := range strings.Split(strings.TrimSuffix(string(input), "\n"), "\n") {
        out = append(out, String(val))
    }

    switch len(out) {
    case 1:
        return out[0]
    default:
        return out
    }
}

func parse(j interface{}) interface{} {
    switch x := j.(type) {
    case map[string]interface{}:
        out := Hash { }

        for key, val := range x {
            out[key] = parse(val)
        }

        return out
    case []interface{}:
        out := Array{ }

        for _, val := range x {
            out = append(out, parse(val))
        }

        return out
    case string:
        return String(x)
    case float64:
        return Number{ val: new(big.Rat).SetFloat64(x) }
    case bool:
        return Boolean(x)
    case nil:
        return Null{ }
    }

    return Null{ }
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
        lexer.toks = append(lexer.toks, &Token { pos: lexer.pos, tok: FIN, lit: "" })
    }

    parser := &Parser {
        lexr: lexer,
        blk: blk,
    }

    for len(parser.lexr.toks) > 0 {
        parser.Churn()
    }

    return parser
}

func test(test *testing.T, source string) {
    program = NewBlock()

    r := open(source)
    p := process(r, program)
    c := 0
    f := 0

    if val, ok := p.blk.Run(p.blk.def...).(Array); ok {
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
