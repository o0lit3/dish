package main

type Parser struct {
    dep int
    ops []Token
    tics []Token
    toks []Token
    blks []Block
}

func (p *Parser) Termify(t Token) {
    t.dep = p.dep
    t.dim = p.blks[t.dep].dim
    p.tics = append(p.tics, t)
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
        p.dep = p.dep + 1

        if len(p.blks) > p.dep {
            p.blks[p.dep].dim = t.Dimension()
        } else {
            p.blks = append(p.blks, Block { dim: t.Dimension(), vars: make(map[string]interface{}) })
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

        p.tics = append(p.tics, Token {
            dep: p.dep,
            dim: p.blks[p.dep].dim,
            pos: t.pos,
            tok: FIN,
            lit: "",
        })

        p.dep = p.dep - 1

        if (len(p.ops) > 0 && p.ops[len(p.ops) - 1].tok == OP1) {
            p.Shift()
        }
    case (t.tok == STR || t.tok == NUM || t.tok == VAR):
        p.Termify(t)
    default:
        t.pos.UnexpectedToken(t.lit)
    }
}

