package main

import (
    "fmt"
    "sort"
    "strings"
    "strconv"
    "unicode"
    "math/big"
)

type Variable struct {
    blk *Block
    obj interface{}
    par *Variable
    nom string
    idx int
}

type Hash map[string]interface{}

type Array []interface{}

type String string

type Number struct {
    val *big.Rat
}

type Boolean bool

type Null struct { }

func (n Null) String() string {
    return "null"
}

func (b Boolean) Number() Number {
    if b {
        return NewNumber(1)
    }

    return NewNumber(0)
}

func NewNumber(num int) Number {
    return Number{ val: big.NewRat(int64(num), 1) }
}

func (n Number) Int() int {
    return int(new(big.Int).Quo(n.val.Num(), n.val.Denom()).Int64())
}

func (n Number) String() string {
    if n.val.IsInt() {
        return n.val.RatString()
    }

    prec := 100
    val := n.val.FloatString(prec)

    for !strings.HasSuffix(val, "0000000000") && prec < 1000 {
        prec += 100
        val = n.val.FloatString(prec)
    }

    return strings.TrimRight(val, "0")
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
    out := NewNumber(0)

    if val, ok := out.val.SetString(string(s)); ok {
        return Number{ val: val }
    }

    return out
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

func (b *Block) Eval(a interface{}) interface{} {
    switch x := a.(type) {
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
        return b.Value(a)
    }
}

func (b *Block) Value(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        if len(x.args) > 0 {
            return x
        }

        return x.Run()
    case *Variable:
        switch obj := x.obj.(type) {
        case Hash:
            if _, ok := obj[x.nom]; ok {
                return b.Value(obj[x.nom])
            }

            return Null { }
        case Array:
            if x.idx < len(obj) {
                return b.Value(obj[x.idx])
            }

            return Null { }
        case String:
            if x.idx < len(obj) {
                return String(string(obj[x.idx]))
            }

            return Null { }
        }

        return b.Value(b.cur.vars[x.nom])
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
                b.cur.hash[b.args[i]] = val
            }
        }
    }

    i := len(args)

    for i < len(b.args) {
        b.cur.vars[b.args[i]] = Null { }
        b.cur.hash[b.args[i]] = Null { }
        i = i + 1
    }

    return b.Interpret()
}

func (blk *Block) Interpret() interface{} {
    t := blk.toks[blk.cur.idx]
    blk.cur.idx = blk.cur.idx + 1

    switch t.tok {
    case OP1:
        a := blk.Deregister(t)

        if len(t.lit) > 0 && t.opx && unicode.IsLetter(rune(t.lit[0])) {
            switch op := blk.FindVar(t.lit).(type) {
            case Null:
            default:
                blk.Register(Member(a, op))
                return blk.Interpret()
            }
        }

        switch t.lit {
        case "!", "not":
            blk.Register(Not(a))
        case "^", "invert":
            blk.Register(Invert(a))
        case "%", "hashify", "hash":
            blk.Register(Hashify(a))
        case "@", "keys":
            blk.Register(Keys(a))
        case "**", "sort":
            blk.Register(Sort(a))
        case "*", "product":
            blk.Register(Product(a))
        case "/", "itemize", "array", "arr", "values", "ratio":
            blk.Register(Itemize(a))
        case "+", "number", "num", "sum":
            blk.Register(Sum(a))
        case "-", "negative", "negate":
            blk.Register(Negate(a))
        case ">>", "pop":
            val, obj := Pop(a)
            blk.Assign(a, obj)
            blk.Register(val)
        case "<<", "shift":
            val, obj := Shift(a)
            blk.Assign(a, obj)
            blk.Register(val)
        case "~", "stringify", "string", "str":
            blk.Register(Stringify(a))
        case "<", "minimum", "min", "floor", "int":
            blk.Register(Min(a))
        case "=", "average", "avg", "mean":
            blk.Register(Average(a))
        case ">", "maxium", "max", "ceiling", "ceil":
            blk.Register(Max(a))
        case "&", "compact":
            blk.Register(Compact(a))
        case "|", "unique", "uniq", "abs":
            blk.Register(Unique(a))
        case "#", "size", "length", "len":
            blk.Register(Length(a))
        case "++", "increment", "incr":
            blk.Register(blk.Assign(a, Increment(a)))
        case "--", "decrement", "decr":
            blk.Register(blk.Assign(a, Decrement(a)))
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

        if len(t.lit) > 0 && t.opx && unicode.IsLetter(rune(t.lit[0])) {
            switch op := blk.FindVar(t.lit).(type) {
            case Null:
            default:
                blk.Register(Member(Array { blk.Value(a), blk.Value(b) }, op))
                return blk.Interpret()
            }
        }

        switch t.lit {
        case "?", "switch":
            blk.Register(Switch(blk.Blockify(a), blk.Blockify(b)))
        case "??", "redo":
            blk.Register(Redo(Blockify(a), Blockify(b)))
        case "@", "find", "index", "format", "round":
            blk.Register(Find(a, b))
        case "**", "power", "pow", "sort":
            blk.Register(Power(a, b))
        case "*", "multiply", "repeat", "map", "each", "dot":
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
        case "<<", "push", "append", "lshift":
            val := Push(a, b)

            switch val.(type) {
            case Hash:
                blk.Register(blk.Assign(a, val))
            case Array:
                blk.Register(blk.Assign(a, val))
            default:
                blk.Register(val)
            }
        case ">>", "unshift", "prepend", "rshift":
            val := Unshift(a, b)

            switch val.(type) {
            case Hash:
                blk.Register(blk.Assign(a, val))
            case Array:
                blk.Register(blk.Assign(a, val))
            default:
                blk.Register(val)
            }
        case "&", "intersect":
            blk.Register(Intersect(a, b))
        case "^", "exclude":
            blk.Register(Exclude(a, b))
        case "|", "union":
            blk.Register(Union(a, b))
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
        case "&&", "and", "then":
            blk.Register(And(a, b))
        case "^^", "xor":
            blk.Register(Xor(a, b))
        case "||", "or", "else":
            blk.Register(Or(a, b))
        case "..", "range", "to":
            blk.Register(Range(a, b))
        case "=", "assign":
            switch y := b.(type) {
            case *Variable:
                blk.Register(blk.Assign(a, y.Value()))
            default:
                blk.Register(blk.Assign(a, b))
            }
        case ":", "define":
            switch y := b.(type) {
            case *Variable:
                blk.Register(blk.Define(a, y.Value()))
            default:
                blk.Register(blk.Define(a, b))
            }
        case "+=":
            blk.Register(blk.Assign(a, Add(a, b)))
        case "-=":
            blk.Register(blk.Assign(a, Subtract(a, b)))
        case "*=":
            blk.Register(blk.Assign(a, Multiply(a, b)))
        case "/=":
            blk.Register(blk.Assign(a, Divide(a, b)))
        case "%=":
            blk.Register(blk.Assign(a, Mod(a, b)))
        case "&=":
            blk.Register(blk.Assign(a, Intersect(a, b)))
        case "^=":
            blk.Register(blk.Assign(a, Exclude(a, b)))
        case "|=":
            blk.Register(blk.Assign(a, Union(a, b)))
        case "", "member", "item":
            blk.Register(Member(a, b))
        default:
            t.UnexpectedToken()
        }
    case BLK:
        blk.Register(t.blk)
    case FIN:
        if len(blk.cur.stck) > 0 {
            blk.cur.stck[len(blk.cur.stck) - 1] = blk.Eval(blk.cur.stck[len(blk.cur.stck) - 1])
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
                    out = blk.cur.stck[len(blk.cur.stck) - 1]
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
        if val, ok := new(big.Rat).SetString(t.lit); ok {
            blk.Register(Number{ val: val })
        } else {
            blk.Register(NewNumber(0))
        }
    case STR:
        blk.Register(blk.Interpolate(t.lit))
    default:
        t.UnexpectedToken()
    }

    return blk.Interpret()
}
