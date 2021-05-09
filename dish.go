package main

import (
    "os"
    "io"
    "fmt"
    "bufio"
    "strings"
    "strconv"
    "unicode"
    "unicode/utf8"
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

type Position struct {
    line int
    column int
}

type Token struct {
    depth int
    style string
    pos Position
    tok int
    lit string
}

type Lexer struct {
    pos Position
    rdr *bufio.Reader
    opn bool
    toks []Token
}

type Block struct {
    style string
    vars map[string]interface{}
    regs []interface{}
    stack []interface{}
}

type Variable struct {
    depth int
    name string
}

type Parser struct {
    depth int
    ops []Token
    terms []Token
    toks []Token
    blks []Block
}

type Interpreter struct {
    terms []Token
    blks []Block
}

func (p Position) String() string {
    return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (p Position) UnexpectedToken(s string) Position {
    panic("Unexpected token " + s + " at " + p.String())
    return p
}

func (p Position) UnexpectedOperand(s string) Position {
    panic("Unexpected operand for " + s + " at " + p.String())
    return p
}

func (t Token) String() string {
    return fmt.Sprintf("%s\t%s\t%s\t%s", t.pos.String(), t.style, t.tok, t.lit)
}

func (t Token) UnmatchedBlock() Token {
    panic("Unmatched " + t.lit + " at " + t.pos.String())
    return t
}

func (t Token) Continuator() bool {
    if (t.tok == OP2 || t.tok == OPA) {
        return true
    }

    switch t.lit {
    case "{", "[", "(", ",":
        return true
    default:
        return false
    }
}

func (t Token) Term() bool {
    if t.BlockClose() {
        return true
    }

    switch t.tok {
    case STR, NUM, VAR:
        return true
    default:
        return false
    }
}

func (t Token) Precedence() int {
    if t.tok == OP1 {
        return 17
    }

    switch t.lit {
    case "**":
        return 16
    case "*", "/", "%":
        return 15
    case "+", "-":
        return 14
    case "<<", ">>":
        return 13
    case "<", "<=", ">", ">=":
        return 12
    case "==", "!=":
        return 11
    case "&":
        return 10
    case "^":
        return 9
    case "|":
        return 8
    case "&&":
        return 7
    case "||":
        return 6
    case "??":
        return 5
    case "..":
        return 4
    case "=", "+=", "-=", "*=", "/=", "%=", "&=", "^=", "|=":
        return 3
    case ",":
        return 2
    default:
        return 0
    }
}

func (a Token) Higher(b Token) bool {
    return a.Precedence() >= b.Precedence()
}

func (t Token) BlockOpen() bool {
    switch t.lit {
    case "{", "[", "(":
        return true
    default:
        return false
    }
}

func (t Token) BlockClose() bool {
    switch t.lit {
    case "}", "]", ")":
        return true
    default:
        return false
    }
}

func (t Token) BlockMatch() string {
    switch t.lit {
    case "{":
        return "}"
    case "}":
        return "{"
    case "[":
        return "]"
    case "]":
        return "["
    case "(":
        return ")"
    case ")":
        return "("
    default:
    return "="
    }
}

func (l *Lexer) Read() rune {
    r, _, err := l.rdr.ReadRune()

    if (err != nil && err == io.EOF) {
        return 0
    }

    l.pos.column++

    return r
}

func (l *Lexer) Peek() rune {
    for bytes := 4; bytes > 0; bytes-- {
        b, err := l.rdr.Peek(bytes)

        if err == nil {
            r, _ := utf8.DecodeRune(b)
            return r
        }
    }

    return 0
}


func (l *Lexer) Backup() Position {
    last := l.pos
    err := l.rdr.UnreadRune()

    if err != nil {
        panic(err)
    }

    l.pos.column--
    return last
}

func (l *Lexer) Reset() Position {
    last := l.pos
    l.pos.line++
    l.pos.column = 0
    return last
}

func (l *Lexer) Tokenize(pos Position, tok int, lit string) Token {
    token := Token { pos: pos, tok: tok, lit: lit }
    l.toks = append(l.toks, token)
    return token
}

func (l *Lexer) Lexify() Token {
    for {
        r, _, err := l.rdr.ReadRune()

        if err != nil {
            if err == io.EOF {
                l.opn = false
                return l.Tokenize(l.pos, EOF, "")
            }

            panic(err)
        }

        l.pos.column++

        switch r {
        case '"', '\'':
            return l.Tokenize(l.pos, STR, l.LexStr(r))
        case ' ', '\t', '\r':
            continue
        case '\n':
            if (len(l.toks) > 0 && !l.toks[len(l.toks) - 1].Continuator()) {
                return l.Tokenize(l.Reset(), FIN, string(r))
            } else {
                l.Reset()
                continue
            }
        case ',':
            return l.Tokenize(l.pos, FIN, string(r))
        case '.':
            s := l.pos
            n := l.Read()

            switch n {
            case 0:
                return l.Tokenize(s.UnexpectedToken(string(r)), ERR, string(r))
            case '.':
                return l.Tokenize(s, OP2, string(r) + string(n))
            default:
                return l.Tokenize(l.Backup(), OPX, l.LexVar())
            }
        case '+', '-', '*', '/', '%', '!', '~', '&', '|', '^', ':', '=', '<', '>':
            s := l.pos
            n := l.Read()

            switch n {
            case 0:
                return l.Tokenize(s.UnexpectedToken(string(r)), ERR, string(r))
            case '=':
                if (len(l.toks) == 0 || !l.toks[len(l.toks ) - 1].Term()) {
                    return l.Tokenize(s.UnexpectedToken(string(r)), ERR, string(r))
                }

                if (r == '<' || r == '>') {
                    return l.Tokenize(s, OP2, string(r) + string(n))
                }

                return l.Tokenize(s, OPA, string(r) + string(n))
            case r:
                if (r == '+' || r == '-') {
                    return l.Tokenize(s, OP1, string(r) + string(n))
                }

                return l.Tokenize(s, OP2, string(r) + string(n))
            default:
                l.Backup()

                if (r == '=' || r == ':') {
                    return l.Tokenize(l.pos, OPA, string(r))
                }

                if (len(l.toks) > 0 && l.toks[len(l.toks) - 1].Term()) {
                    return l.Tokenize(l.pos, OP2, string(r))
                }

                return l.Tokenize(l.pos, OP1, string(r))
            }
        case '{', '[', '(':
            return l.Tokenize(l.pos, BLK, string(r))
        case '}', ']', ')':
            if (len(l.toks) > 0 && l.toks[len(l.toks) - 1].tok == FIN) {
                l.toks = l.toks[:len(l.toks) - 1]
            }

            return l.Tokenize(l.pos, BLK, string(r))
        default:
            switch {
            case unicode.IsDigit(r):
                return l.Tokenize(l.Backup(), NUM, l.LexNum())
            case unicode.IsLetter(r):
                return l.Tokenize(l.Backup(), VAR, l.LexVar())
            default:
                return l.Tokenize(l.pos, ERR, string(r))
            }
        }
    }
}

func (l *Lexer) LexNum() string {
    var lit string

    for {
        r := l.Read()

        switch {
        case r == 0:
            return lit
        case r == '.':
            n := l.Peek()

            switch {
            case n == 0:
                return lit + string(r)
            case unicode.IsDigit(n):
                lit = lit + string(r)
            default:
                l.Backup()
                return lit
            }
        case unicode.IsDigit(r):
            lit = lit + string(r)
        default:
            l.Backup()
            return lit
        }
    }
}

func (l *Lexer) LexStr(quote rune) string {
    var lit string

    for {
        r := l.Read();

        switch r {
        case 0:
            return lit
        case quote:
            return lit
        case '\n':
            l.Reset()
            lit = lit + string(r)
        case '\\':
            n := l.Read();

            switch n {
            case 0:
                return lit
            case 'n':
                lit = lit + "\n"
            case 'r':
                lit = lit + "\r"
            case 't':
                lit = lit + "\t"
            case 'b':
                lit = lit + "\b"
            case 'f':
                lit = lit + "\f"
            case 'v':
                lit = lit + "\v"
            default:
                lit = lit + string(n)
            }
        default:
            lit = lit + string(r)
        }
    }
}

func (l *Lexer) LexVar() string {
    var lit string

    for {
        r := l.Read()

        switch {
        case r == 0:
            return lit
        case unicode.IsLetter(r), unicode.IsDigit(r):
            lit = lit + string(r)
        default:
            l.Backup()
            return lit
        }
    }
}

func (p *Parser) Termify(t Token) {
    t.depth = p.depth
    t.style = p.blks[t.depth].style
    p.terms = append(p.terms, t)
}

func (p *Parser) Shift() {
    p.Termify(p.ops[len(p.ops) - 1])
    p.ops = p.ops[:len(p.ops) - 1]
}

func (p *Parser) Parse() {
    t := p.toks[0]
    p.toks = p.toks[1:]

    switch {
    case t.tok == EOF:
        for len(p.ops) > 0 {
            op := p.ops[len(p.ops) - 1]

            if (op.BlockOpen() || op.BlockClose()) {
                op.UnmatchedBlock()
            }

            p.Shift()
        }
    case t.tok == FIN:
        for (len(p.ops) > 0 && !p.ops[len(p.ops) - 1].BlockOpen()) {
            p.Shift()
        }

        p.Termify(t)
    case t.tok == OP1:
        p.ops = append(p.ops, t)
    case (t.tok == OP2 || t.tok == OPA || t.tok == OPX):
        for (len(p.ops) > 0 && p.ops[len(p.ops) - 1].Higher(t) && !p.ops[len(p.ops) - 1].BlockOpen()) {
            p.Shift()
        }

        p.ops = append(p.ops, t)
    case t.BlockOpen():
        p.depth = p.depth + 1
        p.blks = append(p.blks, Block { style: t.lit + t.BlockMatch(), vars: make(map[string]interface{}) })
        p.ops = append(p.ops, t)
    case t.BlockClose():
        for (len(p.ops) > 0 && p.ops[len(p.ops) - 1].lit != t.BlockMatch()) {
            p.Shift()
        }

        if len(p.ops) == 0 {
            t.UnmatchedBlock()
        } else {
            p.ops = p.ops[:len(p.ops) - 1]
        }

        p.terms = append(p.terms, Token {
            depth: p.depth,
            style: p.blks[p.depth].style,
            pos: t.pos,
            tok: FIN,
            lit: "",
        })

        p.depth = p.depth - 1

        if (len(p.ops) > 0 && p.ops[len(p.ops) - 1].tok == OP1) {
            p.Shift()
        }
    case (t.tok == STR || t.tok == NUM || t.tok == VAR):
        p.Termify(t)
    default:
        t.pos.UnexpectedToken(t.lit)
    }
}

func (i *Interpreter) Stack(depth int, val interface{}) {
    i.blks[depth].stack = append(i.blks[depth].stack, val)
}

func (i *Interpreter) Register(depth int, val interface{}) {
    i.blks[depth].regs = append(i.blks[depth].regs, val)
}

func (i *Interpreter) Cast(val interface{}, t int) interface{} {
    switch x := val.(type) {
    case Variable:
        if t == VAR {
            return x
        }

        v, ok := i.blks[x.depth].vars[x.name]

        if ok {
            return i.Cast(v, t)
        } else {
            switch t {
            case BIT:
                return false
            case NUM:
                return 0
            case STR:
                return ""
            default:
                return nil
            }
        }
    case map[string]interface{}:
        switch t {
        case BIT:
            return len(x) > 0
        case NUM:
            return len(x)
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case []interface{}:
        switch t {
        case BIT:
            return len(x) > 0
        case NUM:
            return len(x)
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case bool:
        switch t {
        case NUM:
            if x {
                return 1
            } else {
                return 0
            }
        case STR:
            if x {
                return "true"
            } else {
                return "false"
            }
        default:
            return x
        }
    case int:
        switch t {
        case BIT:
            return x != 0
        case STR:
            return strconv.Itoa(x)
        default:
            return x
        }
    case float64:
        switch t {
        case BIT:
            return x != 0
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case string:
        switch t {
        case BIT:
            return x != "" && x != "0"
        case NUM:
            v := ""
            decimal := false

            for _, c := range x {
                if unicode.IsDigit(c) {
                    v += string(c)
                } else if (c == '.' && !decimal) {
                    v += string(c)
                    decimal = true
                }
            }

            if decimal {
                w, _ := strconv.ParseFloat(v, 64)
                return w
            }

            w, _ := strconv.Atoi(v)
            return w
        default:
            return x
        }
    default:
        panic("Could not determine operand type.")
    }
}

func (i *Interpreter) Bind(op Token, t int) interface{} {
    if len(i.blks[op.depth].regs) < 1 {
        panic("Missing operand for " +  op.lit + " at " + op.pos.String())
    }

    val := i.blks[op.depth].regs[0]
    i.blks[op.depth].regs = i.blks[op.depth].regs[1:]

    return i.Cast(val, t)
}

func (i *Interpreter) Interpret() {
    t := i.terms[0]
    i.terms = i.terms[1:]

    switch t.tok {
    case OP1:
        switch t.lit {
        case "!":
            i.Register(t.depth, !i.Bind(t, BIT).(bool))
        case "~":
        case "+":
        case "-":
        case "++":
        case "--":
        default:
            t.pos.UnexpectedToken(t.lit)
        }
    case OP2:
        switch t.lit {
        case "**":
        case "*":
            a := i.Bind(t, OBJ)
            b := i.Bind(t, NUM)

            switch x := a.(type) {
            case []interface{}:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                switch y {
                case 0:
                    i.Register(t.depth, nil)
                case 1:
                    i.Register(t.depth, x)
                default:
                    for n:= 1; n < y; n++ {
                        for _, w := range x {
                            x = append(x, w)
                        }
                    }

                    i.Register(t.depth, x)
                }
            case bool:
                if x {
                    i.Register(t.depth, b)
                } else {
                    i.Register(t.depth, 0)
                }
            case int:
                switch b.(type) {
                case int:
                    i.Register(t.depth, a.(int) * b.(int))
                case float64:
                    i.Register(t.depth, float64(a.(int)) * b.(float64))
                }
            case float64:
                switch b.(type) {
                case int:
                    i.Register(t.depth, a.(float64) * float64(b.(int)))
                case float64:
                    i.Register(t.depth, a.(float64) * b.(float64))
                }
            case string:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                i.Register(t.depth, strings.Repeat(x, y))
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
        case "/":
        case "%":
        case "+":
            a := i.Bind(t, OBJ)
            b := i.Bind(t, OBJ)

            switch x := a.(type) {
            case []interface{}:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                switch y {
                case 0:
                    i.Register(t.depth, nil)
                case 1:
                    i.Register(t.depth, x)
                default:
                    for n:= 1; n < y; n++ {
                        for _, w := range x {
                            x = append(x, w)
                        }
                    }

                    i.Register(t.depth, x)
                }
            case bool:
                if x {
                    i.Register(t.depth, b)
                } else {
                    i.Register(t.depth, 0)
                }
            case int:
                switch b.(type) {
                case int:
                    i.Register(t.depth, a.(int) * b.(int))
                case float64:
                    i.Register(t.depth, float64(a.(int)) * b.(float64))
                }
            case float64:
                switch b.(type) {
                case int:
                    i.Register(t.depth, a.(float64) * float64(b.(int)))
                case float64:
                    i.Register(t.depth, a.(float64) * b.(float64))
                }
            case string:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                i.Register(t.depth, strings.Repeat(x, y))
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
        case "-":
        case "<<":
        case ">>":
        case "<":
        case "<=":
        case ">":
        case ">=":
        case "==":
        case "!=":
        case "&":
        case "^":
        case "|":
        case "&&":
        case "||":
        case "..":
        case "??":
        default:
            t.pos.UnexpectedToken(t.lit)
        }
    case OPA:
        a := i.Bind(t, VAR)

        switch x := a.(type) {
        case Variable:
            switch t.lit {
            case "=":
                b := i.Bind(t, OBJ)
                i.blks[x.depth].vars[x.name] = b
                i.Register(t.depth, b)
            case "+=":
            case "-=":
            case "**=":
            case "*=":
            case "/=":
            case "%=":
            case "&=":
            case "|=":
            default:
                t.pos.UnexpectedToken(t.lit)
            }
        default:
            fmt.Printf("%T %v\n", x, x)
            panic("Assignment operator requires variable operand near " + t.lit + " at " + t.pos.String())
        }
    case FIN:
        switch len(i.blks[t.depth].regs) {
        case 0:
            i.Stack(t.depth, "")
        case 1:
            i.Stack(t.depth, i.blks[t.depth].regs[0])
        default:
            panic("Too many items in register near " + t.lit + " at " + t.pos.String())
        }

        i.blks[t.depth].regs = nil

        if (len(i.terms) > 0 && i.terms[0].depth < t.depth) {
            switch len(i.blks[t.depth].stack) {
            case 0:
                switch t.style {
                case "()":
                    i.Register(t.depth - 1, "")
                case "[]":
                    var x []interface{}
                    i.Register(t.depth - 1, x)
                case "{}":
                    var x map[string]interface{}
                    i.Register(t.depth - 1, x)
                }
            default:
                switch t.style {
                case "()":
                    i.Register(t.depth - 1, i.blks[t.depth].stack[len(i.blks[t.depth].stack) - 1])
                case "[]", "{}":
                    i.Register(t.depth - 1, i.blks[t.depth].stack)
                }
            }

            i.blks[t.depth].stack = nil
            i.blks[t.depth].vars = nil
        }
    case VAR:
        n := t.depth
        _, ok := i.blks[n].vars[t.lit]

        for (n > 0 && !ok) {
            n--
            _, ok = i.blks[n].vars[t.lit]
        }

        if (n == 0 && !ok) {
            n = t.depth
        }

        i.Register(t.depth, Variable { depth: n, name: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic("Malformed float at " + t.pos.String())
            }

            i.Register(t.depth, val)
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic("Malformed int at " + t.pos.String())
            }

            i.Register(t.depth, val)
        }
    case STR:
        i.Register(t.depth, t.lit)
    default:
        t.pos.UnexpectedToken(t.lit)
    }
}

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
