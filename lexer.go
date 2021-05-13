package main

import (
    "io"
    "fmt"
    "bufio"
    "strings"
    "unicode"
    "unicode/utf8"
)

type Lexeme int
type Dimension int

const (
    NOP Lexeme = iota
    NUM // Numeric lexemes
    STR // String lexemes
    VAR // Identifier lexemes
    OP1 // Unary operators (!, ~, +, -, ++, --)
    OP2 // Binary operators
    OPX // Method operators (alphanumeric)
    OPA // Assignment operators (= += -= *= /= %= &= ^= |=)
    BLK // Grouping lexemes {} [] ()
    COM // Comment lexemes
    FIN // Statement ending lexemes (newline or comma)
    EOF // End of file
)

const (
    NIL Dimension = iota
    VAL // () Scalar
    SET // [] Array
    MAP // {} Hash
)

type Position struct {
    line int
    column int
}

type Token struct {
    dep int
    dim Dimension
    pos Position
    tok Lexeme
    lit string
}

type Lexer struct {
    pos Position
    rdr *bufio.Reader
    opn bool
    toks []Token
}

func (l Lexeme) String() string {
    return [...]string{"NOP", "NUM", "STR", "VAR", "OP1", "OP2", "OPX", "OPA", "BLK", "COM", "FIN", "EOF"}[l]
}

func (d Dimension) String() string {
    return [...]string{"", "()", "[]", "{}"}[d]
}

func (p Position) String() string {
    return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (p Position) UnexpectedToken(r rune) {
    panic(fmt.Sprintf("Unexpected token \"%v\" at %s", r, p))
}

func (t Token) String() string {
    return fmt.Sprintf(strings.Repeat("\t", t.dep) + "%s %s %s %s", t.pos, t.dim, t.tok, t.lit)
}

func (t Token) UnexpectedToken() {
    panic(fmt.Sprintf("Unexpected token \"%s\" at %s", t.lit, t.pos))
}

func (t Token) UnexpectedOperand() {
    panic(fmt.Sprintf("Unexpected operand for \"%s\" at %s", t.lit, t.pos))
}

func (t Token) UnmatchedBlock() {
    panic(fmt.Sprintf("Unmatched \"%s\" at %s", t.lit, t.pos))
}

func (t Token) Continuator() bool {
    if t.tok == OP2 || t.tok == OPA {
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
    return a.Precedence() > b.Precedence() || (a.Precedence() == b.Precedence() && b.tok != OPA && b.tok != OP1)
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

func (t Token) Dimension() Dimension {
    switch t.lit {
    case "{", "}":
        return MAP
    case "[", "]":
        return SET
    case "(", ")":
        return VAL
    default:
        return NIL
    }
}

func (l *Lexer) Read() rune {
    r, _, err := l.rdr.ReadRune()

    if err != nil {
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

func (l *Lexer) Tokenize(pos Position, tok Lexeme, lit string) Token {
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
        case '#':
            return l.Tokenize(l.pos, COM, l.LexCom())
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
            n := l.Read()

            switch {
            case n == 0:
                s.UnexpectedToken(r)
            case n == '.':
                return l.Tokenize(s, OP2, string(r) + string(n))
            case unicode.IsDigit(n):
                return l.Tokenize(l.Backup(), NUM, l.LexNum())
            default:
                return l.Tokenize(l.Backup(), OPX, l.LexVar())
            }
        case '+', '-', '*', '/', '%', '!', '~', '&', '|', '^', ':', '=', '<', '>':
            s := l.pos
            n := l.Read()

            switch n {
            case 0:
                s.UnexpectedToken(r)
            case '=':
                if len(l.toks) == 0 || !l.toks[len(l.toks ) - 1].Term() {
                    s.UnexpectedToken(r)
                }

                if r == '<' || r == '>' {
                    return l.Tokenize(s, OP2, string(r) + string(n))
                }

                return l.Tokenize(s, OPA, string(r) + string(n))
            case r:
                if r == '+' || r == '-' {
                    return l.Tokenize(s, OP1, string(r) + string(n))
                }

                return l.Tokenize(s, OP2, string(r) + string(n))
            default:
                l.Backup()

                if r == '=' || r == ':' {
                    return l.Tokenize(l.pos, OPA, string(r))
                }

                if len(l.toks) > 0 && l.toks[len(l.toks) - 1].Term() {
                    return l.Tokenize(l.pos, OP2, string(r))
                }

                return l.Tokenize(l.pos, OP1, string(r))
            }
        case '{', '[', '(':
            return l.Tokenize(l.pos, BLK, string(r))
        case '}', ']', ')':
            if len(l.toks) > 0 && l.toks[len(l.toks) - 1].tok == FIN {
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
                l.pos.UnexpectedToken(r)
            }
        }
    }
}

func (l *Lexer) LexCom() string {
    var lit string

    r := l.Read()

    for r == ' ' {
        r = l.Read()
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

        r = l.Read()
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

            if unicode.IsDigit(n) {
                lit = lit + string(r)
            } else {
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
        r := l.Read()

        switch r {
        case 0:
            return lit
        case quote:
            return lit
        case '\n':
            l.Reset()
            lit = lit + string(r)
        case '\\':
            n := l.Read()

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
