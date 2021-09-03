package main

import (
    "fmt"
    "sort"
    "strings"
    "strconv"
    "unicode"
    "math/big"
)

const INF = 1

type Variable struct {
    blk *Block
    obj interface{}
    par *Variable
    arr []*Variable
    nom string
    sub bool
    idx int
}

type Hash map[string]interface{}

type Array []interface{}

type String []rune

type Number struct {
    val *big.Rat
    inf int
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
    if n.inf == INF {
        return "inf"
    }

    if n.inf == -INF {
        return "-inf"
    }

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

func (n Number) Array() Array {
    out := Array{ }

    for _, c := range strconv.FormatInt(int64(n.Int()), 2) {
        out = append(out, String(string(c)).Number())
    }

    return out
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
    out := Array{ }

    for _, c := range s {
        out = append(out, String(string(c)))
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

func (a Array) Hash() Hash {
    out := Hash{ }

    for _, val := range a {
        switch x := val.(type) {
        case String:
            out[string(x)] = Boolean(true)
        default:
            out[fmt.Sprintf("%v", x)] = Boolean(true)
        }
    }

    return out
}

func (h Hash) String() string {
    var out []string
    var keys []string

    for key := range h {
        keys = append(keys, key)
    }

    sort.Strings(keys)

    for _, key := range keys {
        out = append(out, fmt.Sprintf("%s: %v", String(key), h[key]))
    }

    return fmt.Sprintf("{" + strings.Join(out, ", ") + "}")
}

func (h Hash) Array() Array {
    out := Array{ }

    for _, val := range h {
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
        if len(x.args) > 0 {
            return x
        }

        return x.Run()
    case *Variable:
        switch obj := x.obj.(type) {
        case Hash:
            if x.arr != nil {
                out := Array{ }

                for _, mem := range x.arr {
                    if _, ok := obj[mem.nom]; ok {
                        out = append(out, obj[mem.nom])
                    } else {
                        out = append(out, Null{ })
                    }
                }

                return out
            }

            if _, ok := obj[x.nom]; ok {
                return b.Value(obj[x.nom])
            }

            return Null{ }
        case Array:
            if x.arr != nil {
                out := Array{ }

                for _, mem := range x.arr {
                    if mem.idx < len(obj) {
                        out = append(out, obj[mem.idx])
                    }
                }

                return out
            }

            if x.idx < len(obj) {
                return b.Value(obj[x.idx])
            }

            return Null{ }
        case String:
            if x.arr != nil {
                out := ""

                for _, mem := range x.arr {
                    if mem.idx < len(obj) {
                        out += string(obj[mem.idx])
                    }
                }

                return String(out)
            }

            if x.idx < len(obj) {
                return String(string(obj[x.idx]))
            }

            return Null{ }
        case Number:
            bin := strconv.FormatInt(int64(obj.Int()), 2)

            if x.arr != nil {
                out := ""

                for _, mem := range x.arr {
                    if mem.idx < len(bin) {
                        out += string(bin[mem.idx])
                    }
                }

                val, _ := strconv.ParseInt(out, 2, 64)

                return NewNumber(int(val))
            }

            if x.idx < len(bin) {
                return Boolean(bin[x.idx] == '1').Number()
            }

            return Null{ }
        }

        return b.FindVar(x.nom)
    default:
        return x
    }
}

func (b *Block) Eval(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        out := Hash{ }

        for key, val := range x {
            out[key] = b.Eval(val)
        }

        return out
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, b.Eval(val))
        }

        return out
    case String:
        return x
    default:
        return b.Value(a)
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
    b.cur = &Run{ idx: 0, stck: Array{ }, hash: Hash{ }, vars: Hash{ } }
    b.runs = append(b.runs, b.cur)

    if len(args) > 0 {
        if b.obj != nil {
            b.cur.vars["$_"] = b.obj
            b.cur.hash["$_"] = b.obj
        }

        for i, val := range args {
            b.cur.vars["$" + strconv.Itoa(i + 1)] = val
            b.cur.hash["$" + strconv.Itoa(i + 1)] = val

            if i < len(b.args) {
                b.cur.vars[b.args[i]] = val
                b.cur.hash[b.args[i]] = val
            }
        }
    }

    i := len(args)

    for i < len(b.args) {
        b.cur.vars[b.args[i]] = Null{ }
        b.cur.hash[b.args[i]] = Null{ }
        i = i + 1
    }

    return b.Chirp()
}

func (blk *Block) Chirp() interface{} {
    t := blk.toks[blk.cur.idx]
    blk.cur.idx = blk.cur.idx + 1

    switch t.tok {
    case OP1:
        a := blk.Deregister(t)

        if len(t.lit) > 0 && (t.opx || t.mem) && (unicode.IsLetter(rune(t.lit[0])) || t.lit[0] == '$') {
            if t.mem {
                blk.Register(t.Dot(a, String(t.lit)))
                return blk.Chirp()
            }

            switch op := blk.FindVar(t.lit).(type) {
            case Null:
            default:
                blk.Register(t.Dot(a, op))
                return blk.Chirp()
            }
        }

        switch t.lit {
        case "!", "not":
            blk.Register(t.TopBang(a))
        case "?", "boolean", "bool":
            blk.Register(t.TopWhiz(a))
        case "%", "ratio", "hash":
            blk.Register(t.TopGrep(a))
        case "~", "invert", "flip", "bnot":
            blk.Register(t.TopTwiddle(a))
        case "@", "keys", "reverse", "round":
            blk.Register(t.TopThump(a))
        case "*", "string", "str", "product", "prod", "join":
            blk.Register(t.TopSplat(a))
        case "**", "divisors", "perms":
            blk.Register(t.TopDoubleSplat(a))
        case "/", "array", "arr", "factors", "split", "words", "values", "vals":
            blk.Register(t.TopSplit(a))
        case "//", "bits", "chars", "flatten", "flat":
            blk.Register(t.TopDoubleSplit(a))
        case "+", "sum", "number", "num":
            blk.Register(t.TopCross(a))
        case "-", "negsum", "negate", "neg":
            blk.Register(t.TopDash(a))
        case ">>", "pop", "lsb", "last":
            val, obj := t.TopWakaWaka(a)

            if _, ok := a.(*Variable); ok && t.lit != "lsb" && t.lit != "last" {
                blk.Assign(a, obj, false)
            }

            blk.Register(val)
        case "<<", "shift", "msb", "first":
            val, obj := t.TopWikiWiki(a)

            if _, ok := a.(*Variable); ok && t.lit != "msb" && t.lit != "first" {
                blk.Assign(a, obj, false)
            }

            blk.Register(val)
        case "^", "sort", "squared":
            blk.Register(t.TopHat(a))
        case "<", "minimum", "min", "lowercase", "downcase", "lc", "floor", "int":
            blk.Register(t.TopWiki(a))
        case "=", "average", "avg", "eval":
            blk.Register(t.TopBars(a))
        case ">", "maxium", "max", "uppercase", "upcase", "uc", "ceiling", "ceil":
            blk.Register(t.TopWaka(a))
        case "&", "compact", "popcount":
            blk.Register(t.TopBoom(a))
        case "|", "unique", "uniq", "abs":
            blk.Register(t.TopWham(a))
        case "#", "length", "len", "bitcount":
            blk.Register(t.TopSharp(a))
        case "++", "increment", "incr":
            val := t.TopDoubleCross(a)

            if _, ok := a.(*Variable); ok {
                blk.Assign(a, val, false)
            }

            blk.Register(val)
        case "--", "decrement", "decr":
            val := t.TopDoubleDash(a)

            if _, ok := a.(*Variable); ok {
                blk.Assign(a, val, false)
            }

            blk.Register(val)
        case "digit", "letter", "consonant", "vowel", "upper", "lower", "space", "ord", "chr":
            blk.Register(t.Strings(a))
        case "rand", "sqrt", "log", "sin", "cos", "tan", "asin", "acos", "atan", "prime":
            blk.Register(t.Numbers(a))
        default:
            switch {
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                blk.Register(t.Dot(a, String(t.lit).Number()))
            default:
                blk.Register(t.Dot(a, String(t.lit)))
            }
        }
    case OP2:
        b := blk.Deregister(t)
        a := blk.Deregister(t)

        if y, ok := b.(*Block); ok && len(t.args) > 0 {
            y.args = t.args
        }

        if len(t.lit) > 0 && t.opx && (unicode.IsLetter(rune(t.lit[0])) || t.lit[0] == '$') {
            switch op := blk.FindVar(t.lit).(type) {
            case Null:
            default:
                params := Array{ blk.Value(a) }

                for _, param := range blk.Blockify(b) {
                    switch y := param.(type) {
                    case *Block:
                        params = append(params, y.Run())
                    default:
                        params = append(params, y)
                    }
                }

                blk.Register(t.Dot(params, op))

                return blk.Chirp()
            }
        }

        switch t.lit {
        case "?", "and", "switch", "redo", "while":
            blk.Register(t.Whiz(blk.Blockify(a), blk.Blockify(b)))
        case "??", "coallesce":
            blk.Register(t.DoubleWhiz(a, b))
        case "!", "or", "swap", "until":
            blk.Register(t.Bang(blk.Blockify(a), blk.Blockify(b)))
        case "@", "round", "find", "search", "indices":
            blk.Register(t.Thump(a, b))
        case "^", "power", "pow", "rotate", "rot", "sort", "zip":
            blk.Register(t.Pow(a, b))
        case "*", "multiply", "mult", "repeat", "times", "map", "each", "join", "dot":
            blk.Register(t.Splat(a, b))
        case "**", "choose", "comb":
            blk.Register(t.DoubleSplat(a, b))
        case "/", "divide", "div", "split":
            blk.Register(t.Split(a, b))
        case "//", "idiv", "partition", "group":
            blk.Register(t.DoubleSplit(a, b))
        case "%", "mod", "filter", "select", "grep", "every":
            blk.Register(t.Grep(a, b))
        case "%%", "imod", "xevery", "without":
            blk.Register(t.DoubleGrep(a, b))
        case "+", "add", "concat", "increase", "pad", "aggregate":
            blk.Register(t.Cross(a, b))
        case "-", "subtract", "sub", "remove", "delete", "del", "decrease", "trunc", "reduce":
            blk.Register(t.Dash(a, b))
        case "#", "base", "unbase", "format", "fmt":
            blk.Register(t.Sharp(a, b))
        case "<<", "push", "append", "rpad", "rtrunc", "lshift", "extend":
            val := t.WikiWiki(a, b)

            if _, ok := a.(*Variable); ok && t.lit != "lshift" {
                switch val.(type) {
                case Hash:
                    blk.Assign(a, val, false)
                case Array:
                    blk.Assign(a, val, false)
                case String:
                    blk.Assign(a, val, false)
                }
            }

            blk.Register(val)
        case ">>", "unshift", "prepend", "lpad", "ltrunc", "rshift":
            val := t.WakaWaka(a, b)

            if _, ok := a.(*Variable); ok && t.lit != "rshift" {
                switch val.(type) {
                case Hash:
                    blk.Assign(a, val, false)
                case Array:
                    blk.Assign(a, val, false)
                case String:
                    blk.Assign(a, val, false)
                }
            }

            blk.Register(val)
        case "&", "intersect", "band", "all":
            blk.Register(t.Boom(a, b))
        case "~", "exclude", "bxor", "none":
            blk.Register(t.Twiddle(a, b))
        case "|", "union", "bor", "any":
            blk.Register(t.Wham(a, b))
        case "<", "lt":
            blk.Register(t.Wiki(a, b))
        case "<=", "lte", "le":
            blk.Register(t.WikiBars(a, b))
        case ">", "gt":
            blk.Register(t.Waka(a, b))
        case ">=", "gte", "ge":
            blk.Register(t.WakaBars(a, b))
        case "==", "equals", "eq", "is":
            blk.Register(Equals(a, b))
        case "!=", "isnt", "ne":
            blk.Register(Boolean(!Equals(a, b)))
        case "&&", "then":
            blk.Register(t.DoubleBoom(a, b))
        case "~~", "xor":
            blk.Register(t.TwiddleDee(a, b))
        case "||", "else":
            blk.Register(t.DoubleWham(a, b))
        case "..", "range", "to":
            blk.Register(t.DoubleDot(a, b))
        case "=", "assign":
            blk.Register(blk.Assign(a, b, false))
        case ":", "define":
            blk.Register(blk.Assign(a, b, true))
        case "@=", "replace":
            if x, ok := a.(*Variable); ok && x.obj != nil {
                blk.Register(x.Assign(blk, b, false))
            } else {
                panic(fmt.Sprintf("Unexpected operand for member assignment \"%s\" at %s", t.lit, t.pos))
            }
        case "+=":
            blk.Register(blk.Assign(a, t.Cross(a, b), false))
        case "-=":
            blk.Register(blk.Assign(a, t.Dash(a, b), false))
        case "*=":
            blk.Register(blk.Assign(a, t.Splat(a, b), false))
        case "/=":
            blk.Register(blk.Assign(a, t.Split(a, b), false))
        case "%=":
            blk.Register(blk.Assign(a, t.Grep(a, b), false))
        case "^=":
            blk.Register(blk.Assign(a, t.Pow(a, b), false))
        case "&=":
            blk.Register(blk.Assign(a, t.Boom(a, b), false))
        case "~=":
            blk.Register(blk.Assign(a, t.Twiddle(a, b), false))
        case "|=":
            blk.Register(blk.Assign(a, t.Wham(a, b), false))
        case "?=":
            blk.Register(blk.Assign(a, t.DoubleWhiz(a, b), false))
        case "", "at", "item", "items", "subset", "call":
            blk.Register(t.Dot(a, b))
        default:
            t.UnexpectedToken()
        }
    case BLK:
        if len(blk.toks[blk.cur.idx].args) > 0 || blk.toks[blk.cur.idx].ShortCircuit() || blk.toks[blk.cur.idx].Loop() || blk.toks[blk.cur.idx].opx {
            blk.Register(t.blk)
        } else if blk.cur.idx + 1 < len(blk.toks) && (blk.toks[blk.cur.idx + 1].Assignment() || blk.toks[blk.cur.idx + 1].Loop()) {
            blk.Register(t.blk)
        } else {
            blk.Register(t.blk.Run())
        }
    case FIN:
        if len(blk.cur.stck) > 0 {
            blk.cur.stck[len(blk.cur.stck) - 1] = blk.Eval(blk.cur.stck[len(blk.cur.stck) - 1])
        }

        if blk.dim == MAP && len(blk.cur.hash) < len(blk.cur.stck) {
            switch x := blk.cur.stck[len(blk.cur.stck) - 1].(type) {
            case String:
                blk.cur.hash[string(x)] = Boolean(true)
            default:
                blk.cur.hash[fmt.Sprintf("%v", x)] = Boolean(true)
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
                    out = Null{ }
                }
            case LST:
                out = blk.Value(blk.cur.stck)
            case MAP:
                out = blk.Value(blk.cur.hash)
            default:
                out = Null{ }
            }

            blk.runs = blk.runs[:len(blk.runs) - 1]

            if len(blk.runs) > 0 {
                blk.cur = blk.runs[len(blk.runs) - 1]
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
        if len(t.args) > 0 {
            blk.Register(blk.Interpolate(t.lit))
        } else {
            blk.Register(String(t.lit))
        }
    default:
        t.UnexpectedToken()
    }

    return blk.Chirp()
}
