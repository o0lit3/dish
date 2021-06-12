
package main

import (
    "fmt"
    "bufio"
    "strings"
)

type Parser struct {
    ops []*Token
    lexr *Lexer
    blk *Block
}

type Run struct {
    idx int
    stck Array
    hash Hash
    vars Hash
}

type Block struct {
    dep int
    dim Dimension
    src *Block
    args []string
    toks []*Token
    runs []*Run
    cur *Run
}

func NewBlock() *Block {
    return &Block {
        dep: 0,
        dim: VAL,
        args: []string{ "true", "false", "null", "stdin" },
    }
}

func (blk *Block) Assign(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        blk.Assign(x.Run(), b)
    case *Variable:
        x.Assign(blk, b)
    case String:
        blk.cur.vars[string(x)] = b
    default:
        blk.cur.vars[fmt.Sprintf("%v", x)] = b
    }

    return b
}

func (v *Variable) Assign(blk *Block, b interface{}) {
    switch obj := v.obj.(type) {
    case Hash:
        obj[v.nom] = b

        if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = obj
        } else {
            v.par.Assign(blk, obj)
        }
    case Array:
        if v.idx < 0 {
            if len(obj) == 0 {
                obj = append(obj, Null { })
            }

            v.idx = len(obj) + v.idx
        }

        for v.idx + 1 > len(obj) {
            obj = append(obj, Null { })
        }

        obj[v.idx] = b

        if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = obj
        } else {
            v.par.Assign(blk, obj)
        }
    default:
        blk.cur.vars[v.nom] = b
    }
}

func (blk *Block) Define(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        blk.Define(x.Run(), b)
    case *Variable:
        switch obj := x.obj.(type) {
        case Hash:
            obj[x.nom] = b
        case Array:
            obj[x.idx] = b
        default:
            blk.cur.vars[x.nom] = b
            blk.cur.hash[x.nom] = b
        }
    case String:
        blk.cur.vars[string(x)] = b
        blk.cur.hash[string(x)] = b
    default:
        blk.cur.vars[fmt.Sprintf("%v", x)] = b
        blk.cur.hash[fmt.Sprintf("%v", x)] = b
    }

    return b
}

func (b *Block) Interpolate(s string) String {
    lexer := &Lexer {
        pos: Position { row: 1, col: 0 },
        rdr: bufio.NewReader(strings.NewReader(s)),
        opn: true,
    }

    toks := []string{ }
    ipol := false
    out := ""
    tok := ""
    dep := 0

    Interpolate:
    for {
        r := lexer.Read()

        switch r {
        case 0:
            toks = append(toks, tok)
            break Interpolate
        case '\\':
            tok += string(r)

            if n := lexer.Read(); n != 0 {
                tok += string(n)
            }
        case '$':
            n := lexer.Read()

            switch n {
            case 0:
                tok += string(r)
                break Interpolate
            case '(':
                toks = append(toks, tok)
                ipol = true
                tok = "$("
            default:
                tok += string(r) + string(n)
            }
        case '(':
            if ipol {
                dep = dep + 1
            }

            tok += string(r)
        case ')':
            switch {
            case ipol && dep > 0:
                tok += string(r)
                dep = dep - 1
            case ipol:
                tok += string(r)
                toks = append(toks, tok)
                ipol = false
                tok = ""
            default:
                tok += string(r)
            }
        default:
            tok += string(r)
        }
    }

    for _, val := range toks {
        if len(val) > 1 && val[0] == '$' && val[1] == '(' {
            reader := bufio.NewReader(strings.NewReader(val))
            parser := process(reader, b.Branch(VAL))

            switch x := parser.blk.Run().(type) {
            case Null:
            case String:
                out += string(x)
            default:
                out += fmt.Sprintf("%v", x)
            }
        } else {
            out += val
        }
    }

    return String(out)
}

func (b *Block) String() string {
    out := ""

    for _, tok := range b.toks {
        out += tok.lit
    }

    return out
}

func (b *Block) FindVar(name string) interface{} {
    val, ok := b.cur.vars[name]

    for !ok && b.src != nil {
        b = b.src

        if b.cur != nil {
            val, ok = b.cur.vars[name]
        }
    }

    if b.src == nil && !ok {
        return Null { }
    }

    return val
}

func (b *Block) Branch(d Dimension) *Block {
    return &Block { dep: b.dep + 1, dim: d, src: b }
}

func Blockify(a interface{}) *Block {
    if blk, ok := a.(*Block); ok {
        return blk
    }

    return &Block { }
}

func (b *Block) Blockify(a interface{}) Array {
    switch x := a.(type) {
    case *Block:
        switch x.dim {
        case VAL:
            return Array { x }
        case LST:
            blk := b.Branch(VAL)
            out := Array { }
            i := 0

            for i < len(x.toks) {
                if i + 1 < len(x.toks) && x.toks[i].tok == BLK && x.toks[i + 1].tok == FIN {
                    x.toks[i].blk.src = b
                    out = append(out, x.toks[i].blk)
                    i = i + 2
                } else {
                    blk.toks = append(blk.toks, x.toks[i])

                    if x.toks[i].tok == FIN {
                        out = append(out, blk)
                        blk = b.Branch(VAL)
                    }

                    i = i + 1
                }
            }

            return out
        case MAP:
            out := Array { }
            blk := b.Branch(VAL)

            for _, t := range x.toks {
                blk.toks = append(blk.toks, t)

                if t.tok == FIN {
                    out = append(out, blk)
                    blk = b.Branch(VAL)
                }
            }

            return out
        }

        return Array { }
    default:
        return Array { x }
    }
}

func (p *Parser) Shift() {
    p.blk.toks = append(p.blk.toks, p.ops[len(p.ops) - 1])
    p.ops = p.ops[:len(p.ops) - 1]
}

func (p *Parser) Operator() Lexeme {
    if len(p.lexr.toks) > 0 && (p.lexr.toks[0].BlockOpen() || p.lexr.toks[0].tok == OP1) {
        if len(p.lexr.toks) > 1 && p.lexr.toks[1].BlockClose() {
            p.lexr.toks = p.lexr.toks[2:]
            return OP1
        }

        return OP2
    }

    return OP1
}

func (p *Parser) Parse() {
    t := p.lexr.toks[0]
    p.lexr.toks = p.lexr.toks[1:]

    switch {
    case t.tok == EOF:
        for len(p.ops) > 0 {
            op := p.ops[len(p.ops) - 1]

            if op.BlockOpen() || op.BlockClose() {
                op.UnmatchedBlock()
            }

            p.Shift()
        }
    case t.tok == FIN:
        for len(p.ops) > 0 && !p.ops[len(p.ops) - 1].BlockOpen() {
            if p.ops[len(p.ops) - 1].ShortCircuit() {
                p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
                p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
                p.blk = p.blk.src
            }

            p.Shift()
        }

        p.blk.toks = append(p.blk.toks, t)
    case t.tok == OP1:
        p.ops = append(p.ops, t)
    case t.tok == OP2 || t.tok == OPX:
        for len(p.ops) > 0 && p.ops[len(p.ops) - 1].Higher(t) && !p.ops[len(p.ops) - 1].BlockOpen() && !p.ops[len(p.ops) - 1].ShortCircuit() {
            p.Shift()
        }

        if len(p.ops) > 0 && p.ops[len(p.ops) - 1].ShortCircuit() && p.ops[len(p.ops) - 1].Higher(t) {
            p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
            p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
            p.blk = p.blk.src

            for len(p.ops) > 0 && p.ops[len(p.ops) - 1].Higher(t) && !p.ops[len(p.ops) - 1].BlockOpen() {
                p.Shift()
            }
        }

        if t.Redo() {
            blk := p.blk.Branch(VAL)

            for len(p.blk.toks) > 0 && p.blk.toks[len(p.blk.toks) - 1].tok != FIN {
                blk.toks = append([]*Token{ p.blk.toks[len(p.blk.toks) - 1] }, blk.toks...)
                p.blk.toks = p.blk.toks[:len(p.blk.toks) - 1]
            }

            blk.toks = append(blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
            p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: BLK, blk: blk, lit: "" })
        }

        if t.ShortCircuit() {
            p.blk = p.blk.Branch(VAL)
        }

        if t.tok == OPX {
            t.tok = p.Operator()
        }

        p.ops = append(p.ops, t)
    case t.BlockOpen():
        p.blk = p.blk.Branch(t.Dimension())
        p.ops = append(p.ops, t)
    case t.BlockClose():
        for (len(p.ops) > 0 && p.ops[len(p.ops) - 1].lit != t.BlockMatch()) {
            if p.ops[len(p.ops) - 1].ShortCircuit() {
				p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
				p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
				p.blk = p.blk.src
            }

            p.Shift()
        }

        if len(p.ops) == 0 {
            t.UnmatchedBlock()
        } else {
            p.ops = p.ops[:len(p.ops) - 1]
        }

        p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
        p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
        p.blk = p.blk.src

        if len(p.ops) > 0 && p.ops[len(p.ops) - 1].tok == OP1 {
			p.Shift()
        }
    case t.tok == STR || t.tok == NUM || t.tok == VAR:
        p.blk.toks = append(p.blk.toks, t)
    default:
        t.UnexpectedToken()
    }
}
