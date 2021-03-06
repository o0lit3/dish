package main

import (
    "io"
    "fmt"
    "bufio"
    "strings"
    "unicode"
)

type Lexeme int
type Dimension int

const (
    NOP Lexeme = iota
    NUM // Numeric lexemes
    STR // String lexemes
    VAR // Identifier lexemes
    OP1 // Unary operators (!, ~, +, -, #, ++, --)
    OP2 // Binary operators
    OPX // Method operators (alphanumeric)
    MEM // Member access (as string)
    BLK // Grouping lexemes {} [] ()
    FIN // Statement ending lexemes (newline or comma)
    EOF // End of file
)

const (
    NIL Dimension = iota
    VAL // () Scalar
    LST // [] Array
    MAP // {} Hash or Logic
)

type Position struct {
    row int
    col int
}

type Token struct {
    pos Position
    tok Lexeme
    opx bool
    mem bool
    lit string
    blk *Block
    args []string
}

type Lexer struct {
    pos Position
    rdr *bufio.Reader
    opn bool
    toks []*Token
    buf rune
    coms []string
}

func (l Lexeme) String() string {
    return [...]string{"NOP", "NUM", "STR", "VAR", "OP1", "OP2", "OPX", "MEM", "BLK", "FIN", "EOF"}[l]
}

func (d Dimension) String() string {
    return [...]string{"", "()", "[]", "{}"}[d]
}

func (p Position) String() string {
    return fmt.Sprintf("%d:%d", p.row, p.col)
}

func (p Position) UnexpectedToken(r rune) {
    panic(fmt.Sprintf("Unexpected token \"%v\" at %s", string(r), p))
}

func (t *Token) DivideByZero() {
    panic(fmt.Sprintf("Divide by 0 error near \"%s\" at %s", t.lit, t.pos))
}

func (t *Token) TypeMismatch(a interface{}, b interface{}) Boolean {
    if b == nil {
        panic(fmt.Sprintf("Type mismatch on unary operator \"%s\" near %s on type %T", t.lit, t.pos, a))
    } else {
        panic(fmt.Sprintf("Type mismatch on binary operator \"%s\" near %s on types %T and %T", t.lit, t.pos, a, b))
    }

    return Boolean(false)
}

func (t *Token) String() string {
    if t.tok == BLK {
        out := []string{ }

        for _, term := range t.blk.toks {
            if term.tok == BLK {
                out = append(out, fmt.Sprintf("%v", term))
            } else {
                out = append(out, fmt.Sprintf(strings.Repeat("  ", t.blk.dep) + "%s %v", t.blk.dim, term))
            }
        }

        return strings.Join(out, "\n")
    }

    return fmt.Sprintf("%s%v %s %s", t.pos, t.args, t.tok, t.lit)
}

func (t *Token) UnexpectedToken() {
    panic(fmt.Sprintf("Unexpected token \"%s\" at %s", t.lit, t.pos))
}

func (t *Token) UnexpectedOperand() {
    panic(fmt.Sprintf("Unexpected operand for \"%s\" at %s", t.lit, t.pos))
}

func (t *Token) UnmatchedBlock() {
    panic(fmt.Sprintf("Unmatched \"%s\" at %s", t.lit, t.pos))
}

func (t *Token) Continuator() bool {
    if t.tok == OP2 {
        return true
    }

    switch t.lit {
    case "{", "[", "(", ",":
        return true
    }

    return false
}

func (t *Token) Assignment() bool {
    if t.tok == OP1 && (t.lit == "=" || t.lit == ":") {
        return false
    }

    switch t.lit {
    case "=", ":", "+=", "-=", "*=", "/=", "%=", "^=", "&=", "|=", "~=", "?=", "<<", ">>":
        return true
    }

    return false
}

func (t *Token) Term() bool {
    if t.BlockClose() || t.tok == MEM {
        return true
    }

    if t.tok == OPX && t.lit == "" {
        return false
    }

    switch t.tok {
    case STR, NUM, VAR, OPX:
        return true
    default:
        return false
    }
}

func (t *Token) MemberAccessible() bool {
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

func (t *Token) Precedence() int {
    if t.tok == OP1 || t.tok == MEM {
        return 15
    }

    switch t.lit {
    case "^":
        return 14
    case "*", "/", "//", "%", "%%":
        return 13
    case "+", "-":
        return 12
    case "<<", ">>":
        return 11
    case "&":
        return 10
    case "~":
        return 9
    case "|":
        return 8
    case "<", "<=", ">", ">=", "?=":
        return 7
    case "==", "!=":
        return 6
    case "&&":
        return 5
    case "~~":
        return 4
    case "||":
        return 3
    case "..":
        return 2
    case "?", "??", "!", "#", "@":
        return 1
    case ":", "=", "@=", ":=", "+=", "-=", "*=", "/=", "%=", "^=", "&=", "~=", "|=":
        return 0
    default:
        return 20
    }
}

func (a *Token) Higher(b *Token) bool {
    if a.opx && (b.opx || b.tok == MEM) {
        return true
    }

    if a.Precedence() > b.Precedence() {
        return true
    }

    return a.Precedence() == b.Precedence() && !b.Assignment() && b.tok != OP1 && b.tok != MEM
}

func (t *Token) ShortCircuit() bool {
    switch t.lit {
    case "&&", "||", "and", "or":
        return true
    default:
        return false
    }
}

func (t *Token) Loop() bool {
    switch t.lit {
    case "?", "!", "switch", "swap", "redo", "while", "until":
        return true
    default:
        return false
    }
}

func (t *Token) BlockOpen() bool {
    switch t.lit {
    case "{", "[", "(":
        return true
    default:
        return false
    }
}

func (t *Token) BlockClose() bool {
    switch t.lit {
    case "}", "]", ")":
        return true
    default:
        return false
    }
}

func (t *Token) BlockMatch() string {
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

func (t *Token) Dimension() Dimension {
    switch t.lit {
    case "{", "}":
        return MAP
    case "[", "]":
        return LST
    case "(", ")":
        return VAL
    default:
        return NIL
    }
}

func (l *Lexer) Chomp() rune {
    r, _, err := l.rdr.ReadRune()

    if err != nil {
        return 0
    }

    l.pos.col++

    return r
}

func (l *Lexer) Backup() Position {
    last := l.pos
    err := l.rdr.UnreadRune()

    if err != nil {
        panic(err)
    }

    l.pos.col--
    return last
}

func (l *Lexer) Reset() Position {
    last := l.pos
    l.pos.row++
    l.pos.col = 0
    return last
}

func (l *Lexer) Tokenize(pos Position, tok Lexeme, lit string) *Token {
    token := &Token { pos: pos, tok: tok, opx: tok == OPX || tok == MEM, mem: tok == MEM, lit: lit }
    l.toks = append(l.toks, token)
    return token
}

func (l *Lexer) Lexify() *Token {
    for {
        var r rune;
        var err error;

        if l.buf != 0 {
            r = l.buf
            l.buf = 0
        } else {
            r, _, err = l.rdr.ReadRune()
        }

        if err != nil {
            if err == io.EOF {
                l.opn = false
                return l.Tokenize(l.pos, EOF, "")
            }

            panic(err)
        }

        l.pos.col++

        switch r {
        case '"', '\'':
            out := l.Tokenize(l.pos, STR, l.LexStr(r))

            if r == '"' {
                out.args = append(out.args, "interpolate")
            }

            return out
        case ' ', '\t', '\r':
            continue
        case '\n':
            if len(l.toks) > 0 && !l.toks[len(l.toks) - 1].Continuator() {
                return l.Tokenize(l.Reset(), FIN, string(r))
            } else {
                l.Reset()
                continue
            }
        case ',', ';':
            return l.Tokenize(l.pos, FIN, string(r))
        case '.':
            s := l.pos
            n := l.Chomp()

            switch {
            case n == 0:
                s.UnexpectedToken(r)
            case n == '.':
                return l.Tokenize(s, OP2, string(r) + string(n))
            case unicode.IsDigit(n):
                if len(l.toks) == 0 || !l.toks[len(l.toks) - 1].Term() {
                    return l.Tokenize(l.Backup(), NUM, l.LexNum(false))
                }

                return l.Tokenize(l.Backup(), OPX, l.LexNum(true))
            case unicode.IsLetter(n), n == '$':
                return l.Tokenize(l.Backup(), OPX, l.LexVar()).LexArgs(l)
            case n == '"', n == '\'':
                return l.Tokenize(l.pos, MEM, l.LexStr(n))
            default:
                return l.Tokenize(l.Backup(), OPX, "").LexArgs(l)
            }
        case '+', '-', '*', '/', '%', '!', '~', '#', '@', '?', '&', '|', '^', ':', '=', '<', '>':
            if r == ':' && (len(l.toks) == 0 || !l.toks[len(l.toks) - 1].Term()) {
                l.Backup()
                l.Tokenize(l.pos, VAR, "$" + fmt.Sprintf("%v", l.pos))
                return l.Tokenize(l.pos, OP2, "=").LexArgs(l)
            }

            s := l.pos
            n := l.Chomp()

            if r == '#' && n == '#' {
                l.coms = append(l.coms, l.LexCom())
                continue
            }

            switch n {
            case 0:
                s.UnexpectedToken(r)
            case '=':
                if len(l.toks) == 0 || !l.toks[len(l.toks) - 1].Term() {
                    s.UnexpectedToken(r)
                }

                return l.Tokenize(l.pos, OP2, string(r) + string(n)).LexArgs(l)
            case r:
                switch string(r) + string(n) {
                case "++", "--": // Unary double-ops
                    if len(l.toks) == 0 || !l.toks[len(l.toks) - 1].Term() {
                        return l.Tokenize(s, OP1, string(r) + string(n))
                    }
                case "~~", "&&", "||", "??", "%%": // Binary double-ops
                    if len(l.toks) > 0 && l.toks[len(l.toks) - 1].Term() {
                        return l.Tokenize(s, OP2, string(r) + string(n)).LexArgs(l)
                    }
                case ">>", "<<", "//", "**":  // Unary/binary double-ops
                    if len(l.toks) == 0 || !l.toks[len(l.toks) - 1].Term() {
                        return l.Tokenize(s, OP1, string(r) + string(n))
                    }

                    if len(l.toks) > 0 && l.toks[len(l.toks) - 1].Term() {
                        return l.Tokenize(s, OP2, string(r) + string(n)).LexArgs(l)
                    }
                }
            }

            l.Backup()

            if len(l.toks) > 0 && l.toks[len(l.toks) - 1].Term() {
                return l.Tokenize(l.pos, OP2, string(r)).LexArgs(l)
            }

            return l.Tokenize(l.pos, OP1, string(r))
        case '{', '[', '(':
            if len(l.toks) > 0 && l.toks[len(l.toks) - 1].MemberAccessible() && r == '[' {
                l.Tokenize(l.pos, OPX, "")
            }

            return l.Tokenize(l.pos, BLK, string(r))
        case '}', ']', ')':
            if len(l.toks) > 0 && l.toks[len(l.toks) - 1].tok == FIN {
                l.toks = l.toks[:len(l.toks) - 1]
            }

            return l.Tokenize(l.pos, BLK, string(r))
        case '$':
            return l.Tokenize(l.Backup(), VAR, l.LexVar())
        default:
            switch {
            case unicode.IsDigit(r):
                return l.Tokenize(l.Backup(), NUM, l.LexNum(false))
            case unicode.IsLetter(r):
                return l.Tokenize(l.Backup(), VAR, l.LexVar())
            default:
                l.pos.UnexpectedToken(r)
            }
        }
    }
}

func (l *Lexer) LexCom() string {
    var lit string

    r := l.Chomp()

    for r == ' ' {
        r = l.Chomp()
    }

    for {
        switch {
        case r == 0:
            return lit
        case r != '\n':
            lit = lit + string(r)
        default:
            l.Backup()
            return lit
        }

        r = l.Chomp()
    }
}

func (l *Lexer) LexNum(idx bool) string {
    var lit string

    for {
        r := l.Chomp()

        switch {
        case r == 0:
            return lit
        case r == '.':
            if idx {
                l.Backup()
                return lit
            }

            l.buf = r
            n := l.Chomp()

            if unicode.IsDigit(n) {
                l.buf = 0
                l.Backup()
                lit = lit + string(r)
            } else {
                l.Backup()
                return lit
            }
        case r == 'e':
            lit = lit + string(r)
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
        r := l.Chomp()

        switch r {
        case 0:
            return lit
        case quote:
            return lit
        case '\n':
            l.Reset()
            lit = lit + string(r)
        case '\\':
            if quote == '"' {
                n := l.Chomp()

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
                case '$':
                    lit = lit + "\\$"
                default:
                    lit = lit + string(n)
                }
            } else {
                lit = lit + string(r)
            }
        default:
            lit = lit + string(r)
        }
    }
}

func (l *Lexer) LexVar() string {
    var lit string

    for {
        r := l.Chomp()

        switch {
        case r == 0:
            return lit
        case r == '$' && len(lit) == 0:
            lit = lit + string(r)
        case r == '_', unicode.IsLetter(r), unicode.IsDigit(r):
            lit = lit + string(r)
        default:
            l.Backup()
            return lit
        }
    }
}

func (t *Token) LexArgs(l *Lexer) *Token {
    r := l.Chomp()

    switch r {
    case 0:
    case ' ', '\t', '\r':
        t = t.LexArgs(l)
    case ':':
        n := l.Chomp()

        switch {
        case unicode.IsLetter(n), n == '$', n == '{', n == '[', n == '(':
            l.Backup()
            t.args = append(t.args, l.LexVar())
            t = t.LexArgs(l)
        default:
            panic(fmt.Sprintf("Logic block arguments must start with a letter or a dollar sign near \"%v\" at %s", string(n), l.pos))
        }
    default:
        l.Backup()
    }

    return t
}
