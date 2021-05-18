package main

import (
    "fmt"
    "math"
    "strings"
    "strconv"
    "unicode"
)

type Variable struct {
    dep int
    nom string
}

type Block struct {
    dim Dimension
    stck Array
    vars Hash
    hash []string
    com string
}

type Hash map[string]interface{}

type Array []interface{}

type String string

type Number float64

type Boolean bool

type Null struct { }

type Interpreter struct {
    tics []Token
    blks []Block
}

func (n Null) String() string {
    return "null"
}

func (b Boolean) Number() Number {
    if b {
        return Number(1)
    }

    return Number(0)
}

func (n Number) String() string {
    p := math.Pow(10, float64(12))
    r := int(float64(n) * p + math.Copysign(0.5, float64(n) * p))

    return strconv.FormatFloat(float64(r) / p, 'f', -1, 64)
}

func (s String) String() string {
    out := ""

    for _, c := range s {
        switch c {
        case '\b':
            out += "\\b"
        case '\f':
            out += "\\f"
        case '\n':
            out += "\\n"
        case '\r':
            out += "\\r"
        case '\t':
            out += "\\t"
        case '"':
            out += "\\\""
        case '\\':
            out += "\\\\"
        default:
            out += string(c)
        }
    }

    return "\"" + out + "\""
}

func (s String) Number() Number {
    num := ""
    dec := false;

    for _, c := range s {
        if unicode.IsDigit(c) || (c == '.' && !dec) {
            num += string(c)

            if c == '.' {
                dec = true
            }
        } else {
            break
        }
    }

    if num != "" {
        if dec {
            val, err := strconv.ParseFloat(num, 64)

            if err == nil {
                return Number(val)
            }
        } else {
            val, err := strconv.Atoi(num)

            if err == nil {
                return Number(val)
            }
        }
    }

    return Number(0)
}

func (s String) Array() Array {
    out := Array { }

    for _, c := range s {
        out = append(out, String(c))
    }

    return out
}

func (s Array) String() string {
    var out []string

    for _, val := range s {
        out = append(out, fmt.Sprintf("%v", val))
    }

    return fmt.Sprintf("[" + strings.Join(out, ", ") + "]")
}

func NewArray(n int) Array {
    out := Array { }

    for n > 0 {
        out = append(out, Boolean(true))
        n = n - 1
    }

    return out
}

func (s Array) Hash() Hash {
    out := Hash { }

    for _, val := range s {
        switch x := val.(type) {
        case String:
            out[string(x)] = x
        default:
            out[fmt.Sprintf("%v", x)] = x
        }
    }

    return out
}

func (m Hash) String() string {
    var out []string

    for key, val := range m {
        out = append(out, fmt.Sprintf("%s: %v", String(key), val))
    }

    return fmt.Sprintf("{" + strings.Join(out, ", ") + "}")
}

func (m Hash) Array() Array {
    out := Array { }

    for _, val := range m {
        out = append(out, val)
    }

    return out
}

func (i *Interpreter) Bind(op Token, val interface{}) interface{} {
    switch x := val.(type) {
    case Variable:
        if op.Assignment() {
            return x
        } else {
            return i.blks[x.dep].vars[x.nom]
        }
    default:
        return x
    }
}

func (i *Interpreter) Value(a interface{}) interface{} {
    switch x := a.(type) {
    case Variable:
        return i.blks[x.dep].vars[x.nom]
    default:
        return x
    }
}

func (i *Interpreter) Register(op Token, val interface{}) {
    i.blks[op.dep].stck = append(i.blks[op.dep].stck, val)
}

func (i *Interpreter) Deregister(op Token) interface{} {
    if len(i.blks[op.dep].stck) < 1 {
        panic(fmt.Sprintf("Missing operand for \"%s\" at %s", op.lit, op.pos))
    }

    val := i.blks[op.dep].stck[len(i.blks[op.dep].stck) - 1]
    i.blks[op.dep].stck = i.blks[op.dep].stck[:len(i.blks[op.dep].stck) - 1]

    return i.Bind(op, val)
}

func (i *Interpreter) Scope(depth int) Hash {
    vars := Hash { }
    n := 0

    for n < depth {
        for key, val := range i.blks[n].vars {
            vars[key] = val
        }

        n = n + 1
    }

    return vars
}

func (i *Interpreter) LogicBlock(depth int) ([]Token, []Block, int) {
    out := []Token{ }
    blks := []Block {0: Block { dim: VAL, vars: i.Scope(depth) }}
    n := 0

    for i.tics[n].dep > depth - 1 {
        if i.tics[n].tok == FIN && i.tics[n].lit == "" && i.tics[n].dep == i.tics[n + 1].dep {
            break
        }

        dim := i.tics[n].dim

        if i.tics[n].dep - depth == 0 {
            dim = VAL
        }

        if i.tics[n].dep - depth > len(blks) {
            blks = append(blks, Block { dim: dim, vars: Hash { }})
        }

        out = append(out, Token {
            dep: i.tics[n].dep - depth,
            dim: dim,
            pos: i.tics[n].pos,
            tok: i.tics[n].tok,
            lit: i.tics[n].lit,
        })

        if i.tics[n + 1].tok == COM {
            n = n + 2
        } else {
            for i.tics[n].tok == FIN && i.tics[n].lit == "" {
                n = n + 1
            }

            n = n + 1
        }
    }

    return out, blks, n
}

func (i *Interpreter) Run() interface{} {
    for len(i.tics) > 0 {
        i.Interpret()
    }

    if (len(i.blks) > 0 && len(i.blks[0].stck) > 0) {
        return i.blks[0].stck[len(i.blks[0].stck) - 1]
    }

    return Null { }
}

func (i *Interpreter) Interpret() Token {
    t := i.tics[0]

    if t.dep > 0 && t.dim == MAP {
        tics, blks, n := i.LogicBlock(t.dep)
        logic := Interpreter { tics: tics, blks: blks }

        if i.tics[n].tok == OP2 {
            t = i.tics[n]
            i.Register(t, logic)
            i.tics = i.tics[n + 1:]
        } else {
            logic.tics = nil
            logic.blks = nil
            i.tics = i.tics[1:]
        }
    } else {
        i.tics = i.tics[1:]
    }

    switch t.tok {
    case OP1:
        a := i.Deregister(t)

        switch t.lit {
        case "!", "not":
            i.Register(t, Not(a))
        case "~", "invert":
            i.Register(t, Invert(a))
        case "+":
        case "-":
        case "#", "length":
            i.Register(t, Length(a))
        case "++", "increment":
            var val interface{}

            switch x := a.(type) {
            case Variable:
                val = Increment(i.blks[x.dep].vars[x.nom])
                i.blks[x.dep].vars[x.nom] = val
            default:
                val = Increment(a)
            }

            i.Register(t, val)
        case "--":
        default:
            switch {
            case t.lit == "true":
                i.Register(t, Member(a, Boolean(true)))
            case t.lit == "false":
                i.Register(t, Member(a, Boolean(false)))
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                i.Register(t, Member(a, String(t.lit).Number()))
            default:
                i.Register(t, Member(a, String(t.lit)))
            }
        }
    case OP2:
        b := i.Deregister(t)
        a := i.Deregister(t)

        switch t.lit {
        case "**":
        case "*", "map", "multiply":
            i.Register(t, Multiply(a, b))
        case "/":
        case "%":
        case "+", "add":
            i.Register(t, Add(a, b))
        case "-":
        case "<<":
        case ">>":
        case "<":
        case "<=":
        case ">":
        case ">=":
        case "==":
        case "!=":
        case "&":
        case "^":
        case "|":
        case "&&":
        case "||":
        case "..", "range":
        case "??":
        case "=", "assign":
            val := i.Value(b)
            i.blks[t.VarDepth(a)].vars[t.VarName(a)] = val
            i.Register(t, val)
        case ":", "define":
            if t.dim == MAP {
                dep := t.VarDepth(a)
                nom := t.VarName(a)
                val := i.Value(b)

                i.blks[dep].vars[nom] = val

                if len(i.blks[dep].hash) > len(i.blks[dep].stck) {
                    i.blks[dep].hash[len(i.blks[dep].hash) - 1] = nom
                } else {
                    i.blks[dep].hash = append(i.blks[dep].hash, nom)
                }

                i.Register(t, val)
            } else {
                panic(fmt.Sprintf("Assigment operator \"%s\" can only be used in hashes at %s", t.lit, t.pos))
            }
        case "+=":
        case "-=":
        case "*=":
        case "/=":
        case "%=":
        case "&=":
        case "^=":
        case "|=":
        case "", "member":
            i.Register(t, Member(a, b))
        default:
            t.UnexpectedToken()
        }
    case FIN:
        if len(i.blks[t.dep].stck) > 0 {
            val := i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1]

            switch x := val.(type) {
            case Variable:
                i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1] = i.blks[x.dep].vars[x.nom]
            }

            if t.dim == MAP && len(i.blks[t.dep].stck) != len(i.blks[t.dep].hash) {
                switch x := val.(type) {
                case Variable:
                    i.blks[t.dep].hash = append(i.blks[t.dep].hash, x.nom)
                case String:
                    i.blks[t.dep].hash = append(i.blks[t.dep].hash, string(x))
                default:
                    i.blks[t.dep].hash = append(i.blks[t.dep].hash, fmt.Sprintf("%v", x))
                }
            }
        }

        if t.lit == "" && t.dep > 0 {
            blk := Token { dep: t.dep - 1 }

            switch t.dim {
            case VAL:
                for key, val := range i.blks[t.dep].vars {
                    i.blks[blk.dep].vars[key] = val
                }

                if len(i.blks[t.dep].stck) > 0 {
                    i.Register(blk, i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1])
                }
            case LST:
                for key, val := range i.blks[t.dep].vars {
                    i.blks[blk.dep].vars[key] = val
                }

                i.Register(blk, i.blks[t.dep].stck)
            case MAP:
                hash := Hash { }

                for idx, key := range i.blks[t.dep].hash {
                    hash[key] = i.blks[t.dep].stck[idx]
                }

                i.Register(blk, hash)
            }

            i.blks[t.dep].stck = nil
            i.blks[t.dep].hash = nil
            i.blks[t.dep].vars = make(map[string]interface{})
        }
    case VAR:
        n := t.dep
        _, ok := i.blks[n].vars[t.lit]

        for n > 0 && !ok {
            n = n - 1
            _, ok = i.blks[n].vars[t.lit]
        }

        if n == 0 && !ok {
            i.blks[t.dep].vars[t.lit] = Null { }
            n = t.dep
        }

        i.Register(t, Variable { dep: n, nom: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(t, Number(val))
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(t, Number(val))
        }
    case STR:
        i.Register(t, String(t.lit))
    case COM:
        i.blks[t.dep].com = t.lit
    default:
        t.UnexpectedToken()
    }

    return t
}
