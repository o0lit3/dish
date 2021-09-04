package main

import (
    "fmt"
    "bufio"
    "strconv"
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
    obj interface{}
    toks []*Token
    runs []*Run
    cur *Run
}

func NewBlock() *Block {
    stdin := stdin()
    args := []string{ "true", "false", "null", "inf", "stdin", "$_", "argv", "$0" }
    def := []interface{}{ Boolean(true), Boolean(false), Null{ }, Number{ inf: INF }, stdin, stdin, argv, argv }

    for i, arg := range argv {
        args = append(args, "$" + strconv.Itoa(i + 1))
        def = append(def, arg)
    }

    return &Block {
        dep: 0,
        dim: VAL,
        args: args,
        def: def,
    }
}

func (blk *Block) Assign(a interface{}, b interface{}, local bool) interface{} {
    switch x := a.(type) {
    case *Block:
        switch y := b.(type) {
        case *Block:
            return blk.Assign(x, y.Run(), local)
        case *Variable:
            return blk.Assign(x, y.Value(), local)
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
                    obj = v.Assign(blk, arr[i], local)
                } else {
                    obj = v.Assign(blk, Null { }, local)
                }
            case Array:
                if i < len(y) {
                   obj = v.Assign(blk, y[i], local)
                } else {
                    obj = v.Assign(blk, Null { }, local)
                }
            default:
                obj = v.Assign(blk, y, local)
            }

            if v.par == nil {
                blk.cur.stck = append(blk.cur.stck, obj)
            } else if v.par.obj == nil {
                blk.cur.vars[v.par.nom] = obj
            } else {
                v.par.Assign(blk, obj, local)
            }
        }
    case *Variable:
        x.Assign(blk, b, local)
    default:
        if local && blk.dim == MAP {
            key := ""

            if _, ok := x.(String); ok {
                key = string(x.(String))
            } else {
                key = fmt.Sprintf("%v", x)
            }

            switch y := b.(type) {
            case *Block:
                blk.cur.vars[key] = y.Run()
            case *Variable:
                blk.cur.vars[key] = y.Value()
            default:
                blk.cur.vars[key] = y
            }

            blk.cur.hash[key] = blk.cur.vars[key]
        } else {
            tok := blk.toks[len(blk.toks) - 1]

            if blk.cur != nil && blk.cur.idx - 1 < len(blk.toks) {
                tok = blk.toks[blk.cur.idx - 1]
            }

            panic(fmt.Sprintf("Assignment operator \"%s\" requires variable left-hand operand at %s", tok.lit, tok.pos))
        }
    }

    return b
}

func (v *Variable) String() string {
    return fmt.Sprintf("%v", v.Value())
}

func (v *Variable) Assign(blk *Block, b interface{}, local bool) interface{} {
    switch obj := v.obj.(type) {
    case *Block:
        v.obj = obj.Run()
        return v.Assign(blk, b, local)
    case *Variable:
        v.obj = obj.Value()
        return v.Assign(blk, b, local)
    case Hash:
        if v.arr != nil {
            switch y := b.(type) {
            case *Block:
                b = y.Run()
            case *Variable:
                b = y.Value()
            }

            for idx, mem := range v.arr {
                mem.obj = v.obj
                v.obj = v.obj.(Hash).Array().Assign(blk, mem, idx, b, local).(Hash)
            }
        } else {
            switch y := b.(type) {
            case *Variable:
                v.obj.(Hash)[v.nom] = y.Value()
            default:
                v.obj.(Hash)[v.nom] = b
            }
        }
    case Array:
        if v.arr != nil {
            switch y := b.(type) {
            case *Block:
                b = y.Run()
            case *Variable:
                b = y.Value()
            }

            for idx, mem := range v.arr {
                mem.obj = v.obj
                v.obj = v.obj.(Array).Assign(blk, mem, idx, b, local).(Array)
            }
        } else {
            if v.idx < 0 {
                if len(obj) == 0 {
                    v.obj = append(v.obj.(Array), Null { })
                }

                v.idx = len(obj) + v.idx
            }

            for v.idx + 1 > len(v.obj.(Array)) {
                v.obj = append(v.obj.(Array), Null { })
            }

            switch y := b.(type) {
            case *Variable:
                v.obj.(Array)[v.idx] = y.Value()
            default:
                v.obj.(Array)[v.idx] = b
            }
        }
    case String:
        if v.arr != nil {
            switch y := b.(type) {
            case *Block:
                b = y.Run()
            case *Variable:
                b = y.Value()
            }

            for idx, mem := range v.arr {
                mem.obj = v.obj
                v.obj = v.obj.(String).Array().Assign(blk, mem, idx, b, local).(String)
            }
        } else {
            if v.idx < 0 {
                if len(obj) == 0 {
                    v.obj = append(v.obj.(String), ' ')
                }

                v.idx = len(obj) + v.idx
            }

            for v.idx + 1 > len(obj) {
                v.obj = append(v.obj.(String), ' ')
            }

            y := ""
            switch b.(type) {
            case String:
                y = string(b.(String))
            default:
                tok := blk.toks[len(blk.toks) - 1]

                if blk.cur != nil && blk.cur.idx - 1 < len(blk.toks) {
                    tok = blk.toks[blk.cur.idx - 1]
                }

                panic(fmt.Sprintf("String member assignment near \"%s\" requires String for right-hand operand at %s", tok.lit, tok.pos))
            }

            if len(y) > 0 {
                v.obj.(String)[v.idx] = rune(y[0])
            } else {
                v.obj.(String)[v.idx] = ' '
            }
        }
    case Number:
        if v.arr != nil {
            for idx, mem := range v.arr {
                mem.obj = v.obj
                v.obj = v.obj.(Number).Array().Assign(blk, mem, idx, b, local).(Number)
            }
        } else {
            bin := []rune(strconv.FormatInt(int64(obj.Int()), 2))

            if v.idx < 0 {
                if len(bin) == 0 {
                    bin = append(bin, '0')
                }

                v.idx = len(bin) + v.idx
            }

            if v.idx > len(bin) - 1 {
                for v.idx > len(bin) - 1 {
                    bin = append(bin, '0')
                }
            }

            if Boolify(b) {
                bin[v.idx] = '1'
            } else {
                bin[v.idx] = '0'
            }

            val, _ := strconv.ParseInt(string(bin), 2, 64)
            v.obj = NewNumber(int(val))
        }
    default:
        switch y := b.(type) {
        case *Variable:
            blk.cur.vars[v.nom] = y.Value()
        default:
            blk.cur.vars[v.nom] = b
        }

        if local {
            blk.cur.hash[v.nom] = blk.cur.vars[v.nom]
        }
    }

    if !v.sub && v.par != nil {
        if v.par.obj == nil {
            blk.cur.vars[v.par.nom] = v.obj

            if local {
                blk.cur.hash[v.nom] = v.obj
            }
        } else {
            return v.par.Assign(blk, v.obj, local)
        }
    }

    return v.obj
}

func (a Array) Assign(blk *Block, mem *Variable, idx int, b interface{}, local bool) interface{} {
    switch w := b.(type) {
    case Hash:
        if idx < len(w) {
            return mem.Assign(blk, w.Array()[idx], local)
        }
    case Array:
        if idx < len(w) {
            return mem.Assign(blk, w[idx], local)
        }
    case String:
        if _, ok := mem.obj.(String); ok && idx < len(w) {
            return mem.Assign(blk, String(string(w[idx])), local)
        } else {
            return mem.Assign(blk, b, local)
        }
    case Number:
        bin := []rune(strconv.FormatInt(int64(w.Int()), 2))

        if _, ok := mem.obj.(Number); ok {
            if bin[idx % len(bin)] == '1' {
                return mem.Assign(blk, Boolean(true), local)
            } else {
                return mem.Assign(blk, Boolean(false), local)
            }
        } else {
            return mem.Assign(blk, b, local)
        }
    default:
        return mem.Assign(blk, b, local)
    }

    return mem.obj
}

func (b *Block) Interpolate(s string) String {
    lexer := &Lexer {
        pos: Position { row: 1, col: 0 },
        rdr: bufio.NewReader(strings.NewReader(s)),
        opn: true,
    }

    toks := []string{ }
    ipol := false
    ivar := false
    out := ""
    tok := ""
    dep := 0

    Interpolate:
    for {
        r := lexer.Chomp()

        switch {
        case r == 0:
            toks = append(toks, tok)
            break Interpolate
        case r == '\\':
            tok += string(r)

            if n := lexer.Chomp(); n != 0 {
                tok += string(n)
            }
        case r == '$':
            n := lexer.Chomp()

            switch {
            case ipol:
                tok += string(r) + string(n)
            case n == 0:
                tok += string(r)
                break Interpolate
            case n == '(', n == '[', n == '{':
                toks = append(toks, tok)
                ipol = true
                tok = string(r) + string(n)
            case unicode.IsLetter(n), unicode.IsDigit(n), r == '_':
                toks = append(toks, tok)
                ivar = true
                tok = string(r) + string(n)
            default:
                tok += string(r) + string(n)
            }
        case r == '(', r == '[', r == '{':
            if ipol {
                dep = dep + 1
            }

            tok += string(r)
        case r == ')', r == ']', r == '}':
            tok += string(r)

            switch {
            case ipol && dep > 0:
                dep = dep - 1
            case ipol:
                toks = append(toks, tok)
                ipol = false
                tok = ""
            }
        case !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_':
            if ivar {
                toks = append(toks, tok)
                ivar = false
                tok = string(r)
            } else {
                tok += string(r)
            }
        default:
            tok += string(r)
        }
    }

    for _, val := range toks {
        if len(val) > 1 && val[0] == '$' {
            if unicode.IsLetter(rune(val[1])) {
                val = val[1:]
            }

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
    return fmt.Sprintf("%v", b.Run())
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

func (b *Block) Context(x interface{}) *Block {
    b.obj = x
    return b;
}

func (b *Block) Variate() []*Variable {
    var reg *Variable = nil
    out := []*Variable{}

    for _, t := range b.toks {
        switch t.tok {
        case OP1:
            switch {
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                if v, ok := t.Dot(reg, String(t.lit).Number()).(*Variable); ok {
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
                        if v, ok := t.Dot(reg, op).(*Variable); ok {
                            v.par = reg
                            reg = v
                            res = true
                        }
                    }
                }

                if !res {
                    if v, ok := t.Dot(reg, String(t.lit)).(*Variable); ok {
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

func (p *Parser) Churn() {
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
