package main

import (
    "fmt"
    "bufio"
    "strings"
    "unicode"
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
    def []interface{}
    toks []*Token
    runs []*Run
    cur *Run
}

func NewBlock() *Block {
    stdin := stdin()

    return &Block {
        dep: 0,
        dim: VAL,
        args: []string{ "null", "true", "false", "inf", "stdin", "$_" },
        def: []interface{}{ Null{ }, Boolean(true), Boolean(false), Number{ inf: INF }, stdin, stdin },
    }
}

func (blk *Block) Assign(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        switch y := b.(type) {
        case *Block:
            return blk.Assign(x, y.Run())
        case *Variable:
            return blk.Assign(x, y.Value())
        }

        var obj interface{}
        x.cur = blk.cur

        for i, v := range x.Variate() {
            if v.par != nil && v.par.obj == nil {
                v.obj = blk.cur.vars[v.par.nom]
            }

            switch y := b.(type) {
            case Hash:
                arr := y.Array()

                if i < len(arr) {
                    obj = v.Assign(blk, arr[i])
                } else {
                    obj = v.Assign(blk, Null { })
                }
            case Array:
                if i < len(y) {
                   obj = v.Assign(blk, y[i])
                } else {
                    obj = v.Assign(blk, Null { })
                }
            default:
                obj = v.Assign(blk, y)
            }

            if v.par == nil {
                blk.cur.stck = append(blk.cur.stck, obj)
            } else if v.par.obj == nil {
                blk.cur.vars[v.par.nom] = obj
            } else {
                v.par.Assign(blk, obj)
            }
        }
    case *Variable:
        x.Assign(blk, b)
    case String:
        blk.cur.vars[string(x)] = b
    default:
        blk.cur.vars[fmt.Sprintf("%v", x)] = b
    }

    return b
}

func (v *Variable) Assign(blk *Block, b interface{}) interface{} {
    switch obj := v.obj.(type) {
    case *Block:
        v.obj = obj.Run()
        return v.Assign(blk, b)
    case *Variable:
        v.obj = obj.Value()
        return v.Assign(blk, b)
    case Hash:
        obj[v.nom] = b

        if v.par == nil {
            blk.cur.stck = append(blk.cur.stck, obj)
        } else if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = obj
        } else {
            v.par.Assign(blk, obj)
        }

        return obj
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

        if v.par == nil {
            blk.cur.stck = append(blk.cur.stck, obj)
        } else if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = obj
        } else {
            v.par.Assign(blk, obj)
        }

        return obj
    case String:
        if v.idx < 0 {
            if len(obj) == 0 {
                obj = append(obj, ' ')
            }

            v.idx = len(obj) + v.idx
        }

        for v.idx + 1 > len(obj) {
            obj = append(obj, ' ')
        }

        y := ""
        switch b.(type) {
        case String:
            y = string(b.(String))
        default:
            y = fmt.Sprintf("%v", b)
        }

        if len(y) > 0 {
            obj[v.idx] = rune(y[0])
        } else {
            obj[v.idx] = ' '
        }

        if v.par == nil {
            blk.cur.stck = append(blk.cur.stck, obj)
        } else if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = obj
        } else {
            v.par.Assign(blk, obj)
        }

        return obj
    default:
        blk.cur.vars[v.nom] = b
        return obj
    }
}

func (blk *Block) Define(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        switch y := b.(type) {
        case *Block:
            return blk.Define(x, y.Run())
        case *Variable:
            return blk.Define(x, y.Value())
        }

        x.cur = blk.cur

        for i, item := range x.Variate() {
            switch y := b.(type) {
            case Hash:
                arr := y.Array()

                if i < len(arr) {
                    blk.Define(item, arr[i])
                } else {
                    blk.Define(item, Null { })
                }
            case Array:
                if i < len(y) {
                    blk.Define(item, y[i])
                } else {
                    blk.Define(item, Null { })
                }
            default:
                blk.Define(item, y)
            }
        }
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

func (b *Block) Variate() []*Variable {
    var reg *Variable = nil
    out := []*Variable{}

    for _, t := range b.toks {
        switch t.tok {
        case OP1:
            switch {
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                if v, ok := Member(reg, String(t.lit).Number()).(*Variable); ok {
                    v.par = reg
                    reg = v
                } else {
                    reg = &Variable{ blk: reg.blk, par: reg, obj: Array{ }, idx: String(t.lit).Number().Int() }
                }
            default:
                res := false

                if len(t.lit) > 0 && t.opx && unicode.IsLetter(rune(t.lit[0])) {
                    switch op := b.FindVar(t.lit).(type) {
                    case Null:
                    default:
                        if v, ok := Member(reg, op).(*Variable); ok {
                            v.par = reg
                            reg = v
                            res = true
                        }
                    }
                }

                if !res {
                    if v, ok := Member(reg, String(t.lit)).(*Variable); ok {
                        v.par = reg
                        reg = v
                    } else {
                        reg = &Variable{ blk: reg.blk, par: reg, obj: Hash{ }, nom: t.lit }
                    }
                }
            }
        case VAR:
            b.cur.vars[t.lit] = b.FindVar(t.lit)
            reg = &Variable{ blk: b, nom: t.lit }
        case FIN:
            if reg != nil {
                out = append(out, reg)
                reg = nil
            }
        default:
            panic(fmt.Sprintf("Assigment block may only contain variables near %s at %s", t.lit, t.pos))
        }
    }

    return out
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
                    if x.toks[i].tok == BLK {
                        x.toks[i].blk.src = b
                    }

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

func (p *Parser) Branch(t *Token) {
    p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
    p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
    p.blk = p.blk.src
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

            if p.ops[len(p.ops) - 1].ShortCircuit() {
                p.Branch(t)
            }

            p.Shift()
        }
    case t.tok == FIN:
        for len(p.ops) > 0 && !p.ops[len(p.ops) - 1].BlockOpen() {
            if p.ops[len(p.ops) - 1].ShortCircuit() {
                p.Branch(t)
            }

            p.Shift()
        }

        p.blk.toks = append(p.blk.toks, t)
    case t.tok == OP1:
        p.ops = append(p.ops, t)
    case t.tok == OP2 || t.tok == OPX || t.tok == MEM:
        for len(p.ops) > 0 && p.ops[len(p.ops) - 1].Higher(t) && !p.ops[len(p.ops) - 1].BlockOpen() && !p.ops[len(p.ops) - 1].ShortCircuit() {
            p.Shift()
        }

        if len(p.ops) > 0 && p.ops[len(p.ops) - 1].ShortCircuit() && p.ops[len(p.ops) - 1].Higher(t) {
            p.Branch(t)

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

        if t.tok == OPX || t.tok == MEM {
            t.tok = p.Operator()
        }

        p.ops = append(p.ops, t)
    case t.BlockOpen() && t.tok != STR:
        p.blk = p.blk.Branch(t.Dimension())
        p.ops = append(p.ops, t)
    case t.BlockClose() && t.tok != STR:
        for (len(p.ops) > 0 && p.ops[len(p.ops) - 1].lit != t.BlockMatch()) {
            if p.ops[len(p.ops) - 1].ShortCircuit() {
                p.Branch(t)
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
