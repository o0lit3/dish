package main

import (
    "fmt"
    "math"
    "sort"
    "strings"
    "strconv"
    "unicode"
)

type Variable struct {
    blk *Block
    nom string
}

type Hash map[string]interface{}

type Array []interface{}

type String string

type Number float64

type Boolean bool

type Null struct { }

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
    if n == 0 {
        return "0"
    }

    if float64(n) == float64(int64(n)) {
        return strconv.FormatFloat(float64(n), 'f', -1, 64)
    }

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
    var keys []string

    for key := range m {
        keys = append(keys, key)
    }

    sort.Strings(keys)

    for _, key := range keys {
        out = append(out, fmt.Sprintf("%s: %v", String(key), m[key]))
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

func (v *Variable) Value() interface{} {
    return v.blk.Value(v)
}

func (b *Block) Value(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return x.Run()
    case *Variable:
        return b.Value(b.cur.vars[x.nom])
    case Hash:
        out := Hash { }

        for key, val := range x {
            out[key] = b.Value(val)
        }

        return out
    case Array:
        out := Array { }

        for _, val := range x {
            out = append(out, b.Value(val))
        }

        return out
    default:
        return x
    }
}

func (b *Block) Register(val interface{}) {
    b.cur.stck = append(b.cur.stck, val)
}

func (b *Block) Deregister(op *Token) interface{} {
    if len(b.cur.stck) < 1 {
        panic(fmt.Sprintf("Missing operand for \"%s\" at %s", op.lit, op.pos))
    }

    val := b.cur.stck[len(b.cur.stck) - 1]
    b.cur.stck = b.cur.stck[:len(b.cur.stck) - 1]

    return val
}

func (b *Block) Run(args ...interface{}) interface{} {
    b.cur = &Run { idx: 0, stck: Array { }, hash: Hash { }, vars: Hash { } }
    b.runs = append(b.runs, b.cur)

    if len(args) > 0 {
        b.cur.vars["$0"] = b
        b.cur.vars["$_"] = args[0]

        for i, val := range args {
            b.cur.vars["$" + strconv.Itoa(i + 1)] = val

            if i < len(b.args) {
                b.cur.vars[b.args[i]] = val
            }
        }
    }

    return b.Interpret()
}

func (blk *Block) Interpret() interface{} {
    t := blk.toks[blk.cur.idx]
    blk.cur.idx = blk.cur.idx + 1

    switch t.tok {
    case OP1:
        a := blk.Deregister(t)

        if vblk, ok := blk.FindVar(t.lit).(*Block); ok {
            blk.Register(Member(a, vblk))
            return blk.Interpret()
        }

        switch t.lit {
        case "!", "not":
            blk.Register(Not(a))
        case "^", "invert":
            blk.Register(Invert(a))
        case "*", "product":
            blk.Register(Product(a))
        case "@", "keys":
            blk.Register(Keys(a))
        case "/", "itemize", "array", "values":
            blk.Register(Itemize(a))
        case "+", "number", "num", "sum":
            blk.Register(Sum(a))
        case "-", "negative", "negate":
            blk.Register(Negate(a))
        case ">>", "pop":
            val, obj := Pop(a)
            blk.cur.vars[t.VarName(a)] = obj
            blk.Register(val)
        case "<<", "shift":
            val, obj := Shift(a)
            blk.cur.vars[t.VarName(a)] = obj
            blk.Register(val)
        case "~", "stringify", "string", "str":
            blk.Register(Stringify(a))
        case "<", "minimum", "min", "floor":
            blk.Register(Min(a))
        case "=", "average", "avg", "mean":
            blk.Register(Average(a))
        case ">", "maxium", "max", "ceiling", "ceil":
            blk.Register(Max(a))
        case "|", "unique", "uniq":
            blk.Register(Unique(a))
        case "#", "size", "length", "len":
            blk.Register(Length(a))
        case "++", "increment", "incr":
            val := Increment(blk.Value(a))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "--", "decrement", "decr":
            val := Decrement(blk.Value(a))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        default:
            switch {
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                blk.Register(Member(a, String(t.lit).Number()))
            default:
                blk.Register(Member(a, String(t.lit)))
            }
        }
    case OP2:
        b := blk.Deregister(t)
        a := blk.Deregister(t)

        if y, ok := b.(*Block); ok && len(t.args) > 0 {
            y.args = t.args
        }

        if vblk, ok := blk.FindVar(t.lit).(*Block); ok {
            blk.Register(Member(Array { blk.Value(a), blk.Value(b) }, vblk))
            return blk.Interpret()
        }

        switch t.lit {
        case "?", "switch":
            blk.Register(Switch(blk.Blockify(a), blk.Blockify(b)))
        case "??", "redo":
            blk.Register(Redo(Blockify(a), Blockify(b)))
        case "@", "find", "index":
            blk.Register(Find(a, b))
        case "**", "power", "pow", "rotate":
            blk.Register(Power(a, b))
        case "*", "multiply", "repeat", "map":
            blk.Register(Multiply(a, b))
        case "/", "divide", "split":
            blk.Register(Divide(a, b))
        case "%", "mod", "filter", "select", "grep":
            blk.Register(Mod(a, b))
        case "+", "add":
            blk.Register(Add(a, b))
        case "-", "subtract":
            blk.Register(Subtract(a, b))
        case "~", "join":
            blk.Register(Join(a, b))
        case "~~", "base":
            blk.Register(Base(a, b))
        case "++", "convert":
            blk.Register(Convert(a, b))
        case "<<", "push":
        case ">>", "unshift":
        case "&", "intersect":
        case "^", "exclude":
        case "|", "union":
        case "<", "below":
            blk.Register(Below(a, b))
        case "<=", "under":
            blk.Register(Under(a, b))
        case ">", "above":
            blk.Register(Above(a, b))
        case ">=", "over":
            blk.Register(Over(a, b))
        case "==", "equals", "is":
            blk.Register(Equals(a, b))
        case "!=", "isnt":
            blk.Register(Not(Equals(a, b)))
        case "&&", "and":
            blk.Register(And(a, b))
        case "^^", "xor":
            blk.Register(Xor(a, b))
        case "||", "or":
            blk.Register(Or(a, b))
        case "..", "range", "to":
            blk.Register(Range(a, b))
        case "=", "assign":
            blk.cur.vars[t.VarName(a)] = b

            if _, ok := b.(*Block); ok {
                blk.Register(Null { })
            } else {
                blk.Register(b)
            }
        case ":", "define":
            blk.cur.vars[t.VarName(a)] = b
            blk.cur.hash[t.VarName(a)] = b
            blk.Register(b)
        case "+=":
            val := Add(blk.Value(a), blk.Value(b))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "-=":
            val := Subtract(blk.Value(a), blk.Value(b))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "*=":
            val := Multiply(blk.Value(a), blk.Value(b))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "/=":
            val := Divide(blk.Value(a), blk.Value(b))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "%=":
            val := Mod(blk.Value(a), blk.Value(b))
            blk.cur.vars[t.VarName(a)] = val
            blk.Register(val)
        case "&=":
        case "^=":
        case "|=":
        case "", "member", "item":
            blk.Register(Member(a, b))
        default:
            t.UnexpectedToken()
        }
    case BLK:
        blk.Register(t.blk)
    case FIN:
        if len(blk.cur.stck) > 0 {
            blk.cur.stck[len(blk.cur.stck) - 1] = blk.Value(blk.cur.stck[len(blk.cur.stck) - 1])
        }

        if blk.dim == MAP && len(blk.cur.hash) < len(blk.cur.stck) {
            switch x := blk.cur.stck[len(blk.cur.stck) - 1].(type) {
            case String:
                blk.cur.hash[string(x)] = x
            default:
                blk.cur.hash[fmt.Sprintf("%v", x)] = x
            }
        }

        if blk.cur.idx == len(blk.toks) {
            if blk.src != nil && blk.src.cur != nil {
                for key, val := range blk.cur.vars {
                    if _, ok := blk.cur.hash[key]; !ok {
                        blk.src.cur.vars[key] = val
                    }
                }
            }

            var out interface{}

            switch blk.dim {
            case VAL:
                if len(blk.cur.stck) > 0 {
                    out = blk.Value(blk.cur.stck[len(blk.cur.stck) - 1])
                } else {
                    out = Null { }
                }
            case LST:
                out = blk.Value(blk.cur.stck)
            case MAP:
                out = blk.Value(blk.cur.hash)
            default:
                out = Null { }
            }

            blk.runs = blk.runs[:len(blk.runs) - 1]

            if len(blk.runs) > 0 {
                blk.cur = blk.runs[len(blk.runs) - 1]
            } else {
                blk.cur = nil
            }

            return out
        }
    case VAR:
        blk.cur.vars[t.lit] = blk.FindVar(t.lit)
        blk.Register(&Variable { blk: blk, nom: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            blk.Register(Number(val))
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            blk.Register(Number(val))
        }
    case STR:
        blk.Register(String(t.lit))
    default:
        t.UnexpectedToken()
    }

    return blk.Interpret()
}
