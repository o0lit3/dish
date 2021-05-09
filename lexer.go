package main

import (
    "io"
    "bufio"
    "unicode"
    "unicode/utf8"
)

type Lexer struct {
    pos Position
    rdr *bufio.Reader
    opn bool
    toks []Token
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
