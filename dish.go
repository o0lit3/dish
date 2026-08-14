package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "testing"
    "math/big"
    "runtime"
    "io/ioutil"
    "encoding/json"
)

var program *Block
var argv Array

func main() {
    var reader *bufio.Reader

    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: dish [-dfpt] [-e expression | file] [arguments...]")
        os.Exit(1)
    }

    index := 0
    source := os.Args[1]
    pretty := false
    format := false
    tokens := false
    debug := false

    for i, flag := range os.Args {
        if i == 0 || reader != nil {
            continue
        }

        if len(flag) == 0 || flag[0:1] != "-" {
            break
        }

        switch flag {
        case "-d", "-debug":
            debug = true
            source = arg(i + 1)
            index = i + 1
        case "-t", "-tokens":
            tokens = true
            source = arg(i + 1)
            index = i + 1
        case "-p", "-pretty":
            pretty = true
            source = arg(i + 1)
            index = i + 1
        case "-f", "-format":
            format = true
            source = arg(i + 1)
            index = i + 1
        case "-e", "-exec":
            source = ""
            reader = bufio.NewReader(strings.NewReader(arg(i + 1)))
            index = i + 1
        default:
            if flag[0:1] == "-" {
                source = arg(i + 1)
                index = i + 1

                if strings.Contains(flag, "d") {
                    debug = true
                }

                if strings.Contains(flag, "t") {
                    tokens = true
                }

                if strings.Contains(flag, "p") {
                    pretty = true
                }

                if strings.Contains(flag, "f") {
                    format = true
                }

                if flag[len(flag) - 1:] == "e" {
                    source = ""
                    reader = bufio.NewReader(strings.NewReader(arg(i + 1)))
                    index = i + 1
                }
            }
        }
    }

    if !debug {
        defer func() {
            if err := recover(); err != nil {
                fail(describe(err))
            }
        }()
    }

    if source != "" {
        reader = open(source)
    }

    for index + 1 < len(os.Args) {
        index += 1
        argv = append(argv, bind([]byte(os.Args[index])))
    }

    program = NewBlock()

    parser := process(reader, program)

    if tokens {
        for _, term := range parser.blk.toks {
            fmt.Fprintf(os.Stderr, "%v\n", term)
        }

        return
    }

    if parser.blk.ReadStdin() {
        parser.blk.BindStdin(stdin())
    }

    if format || pretty {
        print(parser.blk.Run(parser.blk.def...), pretty, 0, true, true)
    } else {
        switch x := parser.blk.Run(parser.blk.def...).(type) {
            case Null:
        case String:
            fmt.Printf("%v\n", string(x))
        default:
            fmt.Printf("%v\n", x)
        }
    }
}

func print(v interface{}, pretty bool, depth int, line bool, indent bool) {
    if indent {
        fmt.Print(strings.Repeat("  ", depth))
    }

    switch x := v.(type) {
    case *Block:
        print(x.Run(), pretty, depth, line, indent)
    case *Variable:
        print(x.Value(), pretty, depth, line, indent)
    case Hash:
        i := 0

        if pretty {
            fmt.Print("\033[0m\033[37;1m{\n")
        } else {
            fmt.Print("{\n")
        }

        for key, val := range x {
            if pretty {
                fmt.Print("\033[0m\033[34;1m" + strings.Repeat("  ", depth + 1) + "\"" + key + "\"")
                fmt.Print("\033[0m\033[37;1m" + ": ")
            } else {
                fmt.Print(strings.Repeat("  ", depth + 1) + "\"" + key + "\"")
                fmt.Print(": ")
            }

            if i < len(x) - 1 {
                print(val, pretty, depth + 1, false, false)

                if pretty {
                    fmt.Print("\033[0m\033[37;1m,\n")
                } else {
                    fmt.Print(",\n")
                }
            } else {
                print(val, pretty, depth + 1, true, false)
            }

            i++
        }

        if pretty {
            fmt.Print("\033[0m\033[37;1m" + strings.Repeat("  ", depth) + "}")
        } else {
            fmt.Print(strings.Repeat("  ", depth) + "}")
        }
    case Array:
        if pretty {
            fmt.Print("\033[0m\033[37;1m[\n")
        } else {
            fmt.Print("[\n")
        }

        for i, val := range x {
            if i < len(x) - 1 {
                print(val, pretty, depth + 1, false, true)

                if pretty {
                    fmt.Print("\033[0m\033[37;1m,\n")
                } else {
                    fmt.Print(",\n")
                }
            } else {
                print(val, pretty, depth + 1, true, true)
            }
       }

        if pretty {
            fmt.Print("\033[0m\033[37;1m" + strings.Repeat("  ", depth) + "]")
        } else {
            fmt.Print(strings.Repeat("  ", depth) + "]")
        }
    case String:
        if (depth > 0) {
            if pretty {
                fmt.Printf("\033[0m\033[32m%v", x)
            } else {
                fmt.Printf("%v", x)
            }
        } else {
            if pretty {
                fmt.Printf("\033[0m\033[32m%v", string(x))
            } else {
                fmt.Printf("%v", string(x))
            }
        }
    case Number:
        if pretty {
            fmt.Printf("\033[0m\033[30m%v", x)
        } else {
            fmt.Printf("%v", x)
        }
    case Boolean:
        if pretty {
            fmt.Printf("\033[0m\033[33m%v", x)
        } else {
            fmt.Printf("%v", x)
        }
    case Null:
        if (depth > 0) {
            if pretty {
                fmt.Printf("\033[0m\033[31m%v", x)
            } else {
                fmt.Printf("%v", x)
            }
        }
    }

    if line {
        fmt.Print("\n")
    }
}

func stdin() interface{} {
    stat, _ := os.Stdin.Stat()

    if (stat.Mode() & os.ModeCharDevice) != 0 {
        return Null{ }
    }

    input, _ := ioutil.ReadAll(os.Stdin)

    return bind(input)
}

func bind(input []byte) interface{} {
    var data interface{}

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
        if path, ok := err.(*os.PathError); ok {
            panic(fmt.Sprintf("cannot open \"%s\": %v", source, path.Err))
        }

        panic(fmt.Sprintf("cannot open \"%s\": %v", source, err))
    }

    return bufio.NewReader(file)
}

func fail(msg string) {
    fmt.Fprintf(os.Stderr, "dish: %s\n", msg)
    os.Exit(1)
}

func arg(i int) string {
    if i >= len(os.Args) {
        fail(fmt.Sprintf("missing argument after \"%s\"", os.Args[i - 1]))
    }

    return os.Args[i]
}

func describe(err interface{}) string {
    switch x := err.(type) {
    case string:
        return x
    case runtime.Error:
        return fmt.Sprintf("internal error: %v (rerun with -d for a trace)", x)
    case error:
        return x.Error()
    }

    return fmt.Sprintf("%v", err)
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
    defer func() {
        if err := recover(); err != nil {
            test.Errorf("%s panicked: %s", source, describe(err))
        }
    }()

    program = NewBlock()

    r := open(source)
    p := process(r, program)
    c := 0
    f := 0

    if p.blk.ReadStdin() {
        p.blk.BindStdin(stdin())
    }

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
