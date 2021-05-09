package main

import ("fmt")

type Position struct {
    line int
    column int
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


type Token struct {
    depth int
    style string
    pos Position
    tok int
    lit string
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
