package main

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
