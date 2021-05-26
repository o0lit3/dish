package main

type Parser struct {
    ops []*Token
    lexr *Lexer
    blk *Block
}

type Block struct {
    dep int
    dim Dimension
    src *Block
    idx int
    toks []*Token
    stck Array
    hash Hash
    vars Hash
    args []string
    coms []string
}

func (b *Block) Reset() *Block {
    b.idx = 0
    return b
}

func (b *Block) FindVar(name string) interface{} {
    val, ok := b.vars[name]

    for !ok && b.src != nil {
        b = b.src
        val, ok = b.vars[name]
    }

    if b.src == nil && !ok {
        return Null { }
    }

    return val
}

func (b *Block) Invoke(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return b.Invoke(x.Run())
    case Hash:
        return b.Run(x.Array()...)
    case Array:
        return b.Run(x...)
    default:
        return b.Run(x)
    }
}

func (p *Parser) Shift() {
    p.blk.toks = append(p.blk.toks, p.ops[len(p.ops) - 1])
    p.ops = p.ops[:len(p.ops) - 1]
}

func (p *Parser) Operator() Lexeme {
    if len(p.lexr.toks) > 0 && p.lexr.toks[0].BlockOpen() {
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
        for (len(p.ops) > 0 && !p.ops[len(p.ops) - 1].BlockOpen()) {
            p.Shift()
        }

        p.blk.toks = append(p.blk.toks, t)
    case t.tok == OP1:
        p.ops = append(p.ops, t)
    case t.tok == OP2 || t.tok == OPX:
        for (len(p.ops) > 0 && p.ops[len(p.ops) - 1].Higher(t) && !p.ops[len(p.ops) - 1].BlockOpen()) {
            p.Shift()
        }

        if t.tok == OPX {
            t.tok = p.Operator()
        }

        p.ops = append(p.ops, t)
    case t.BlockOpen():
        p.blk = &Block {
            dep: p.blk.dep + 1,
            dim: t.Dimension(),
            src: p.blk,
            vars: Hash { },
        }

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

        p.blk.toks = append(p.blk.toks, &Token { pos: t.pos, tok: FIN, lit: "" })
        p.blk.src.toks = append(p.blk.src.toks, &Token { pos: t.pos, tok: BLK, blk: p.blk, lit: "" })
        p.blk = p.blk.src

        if len(p.ops) > 0 && p.ops[len(p.ops) - 1].tok == OP1 {
            p.Shift()
        }
    case t.tok == STR || t.tok == NUM || t.tok == VAR || t.tok == COM:
        p.blk.toks = append(p.blk.toks, t)
    default:
        t.UnexpectedToken()
    }
}
