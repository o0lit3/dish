package main

import (
    "fmt"
    "math"
    "strings"
    "strconv"
    "unicode"
)

type Variable string

type Hash map[string]interface{}

type Array []interface{}

type String string

type Number float64

type Boolean bool

type Null struct { }

type Interpreter struct {
    parser *Parser
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

func (i *Interpreter) Bind(op *Token, val interface{}) interface{} {
    switch x := val.(type) {
    case Variable:
        if op.Assignment() {
            return x
        } else {
            return i.parser.blk.vars[string(x)]
        }
    default:
        return x
    }
}

func (i *Interpreter) Value(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return x.Run()
    case Variable:
        return i.Value(i.parser.blk.vars[string(x)])
    case Hash:
        out := Hash { }

        for key, val := range x {
            out[key] = i.Value(val)
        }

        return out
    case Array:
        out := Array { }

        for _, val := range x {
            out = append(out, i.Value(val))
        }

        return out
    default:
        return x
    }
}

func (i *Interpreter) Register(val interface{}) {
    i.parser.blk.stck = append(i.parser.blk.stck, val)
}

func (i *Interpreter) Deregister(op *Token) interface{} {
    if len(i.parser.blk.stck) < 1 {
        panic(fmt.Sprintf("Missing operand for \"%s\" at %s", op.lit, op.pos))
    }

    val := i.parser.blk.stck[len(i.parser.blk.stck) - 1]
    i.parser.blk.stck = i.parser.blk.stck[:len(i.parser.blk.stck) - 1]

    return i.Bind(op, val)
}

func (b *Block) Run(args ...interface{}) interface{} {
    b.idx = 0
    b.stck = Array { }
    b.vars = Hash { "$0": b }
    b.hash = Hash { }

    for i, val := range args {
        b.vars["$" + strconv.Itoa(i + 1)] = val

        if i < len(b.args) {
            b.vars[b.args[i]] = val
        }
    }

    return (&Interpreter { parser: &Parser { blk: b.Reset() } }).Interpret()
}

func (i *Interpreter) Interpret() interface{} {
    t := i.parser.blk.toks[i.parser.blk.idx]
    i.parser.blk.idx = i.parser.blk.idx + 1

    switch t.tok {
    case OP1:
        a := i.Deregister(t)

        if blk, ok := i.parser.blk.FindVar(t.lit).(*Block); ok {
            i.Register(blk.Invoke(a))
            return i.Interpret()
        }

        switch t.lit {
        case "!", "not":
            i.Register(Not(a))
        case "^", "invert":
            i.Register(Invert(a))
        case "*", "product":
            i.Register(Product(a))
        case "+", "number", "num", "sum":
            i.Register(Sum(a))
        case "-", "negative", "negate":
            i.Register(Negate(a))
        case "~", "string", "str":
            i.Register(Stringify(a))
        case "<", "minimum", "min", "floor":
            i.Register(Min(a))
        case "=", "average", "avg", "mean":
            i.Register(Average(a))
        case ">", "maxium", "max", "ceiling", "ceil":
            i.Register(Max(a))
        case "|", "unique", "uniq":
            i.Register(Unique(a))
        case "#", "size", "length", "len":
            i.Register(Length(a))
        case "@", "switch", "which":
            i.Register(Switch(a))
        case "++", "increment", "incr":
            val := Increment(i.Value(a))
            i.parser.blk.vars[t.VarName(a)] = val
            i.Register(val)
        case "--", "decrement", "decr":
            val := Decrement(i.Value(a))
            i.parser.blk.vars[t.VarName(a)] = val
            i.Register(val)
        default:
            switch {
            case len(t.lit) > 0 && unicode.IsDigit(rune(t.lit[0])):
                i.Register(Member(a, String(t.lit).Number()))
            default:
                i.Register(Member(a, String(t.lit)))
            }
        }
    case OP2:
        b := i.Deregister(t)
        a := i.Deregister(t)

        if y, ok := b.(*Block); ok && len(t.args) > 0 {
            y.args = t.args
        }

        if blk, ok := i.parser.blk.FindVar(t.lit).(*Block); ok {
            i.Register(blk.Invoke(Array{ i.Value(a), i.Value(b) }))
            return i.Interpret()
        }

        switch t.lit {
        case "@", "find":
            i.Register(Find(a, b))
        case "**", "pow", "power":
        case "*", "map", "multiply":
            i.Register(Multiply(a, b))
        case "/", "divide", "split":
            i.Register(Divide(a, b))
        case "%", "mod":
        case "+", "add":
            i.Register(Add(a, b))
        case "-", "subtract":
            i.Register(Subtract(a, b))
        case "~", "join":
            i.Register(Join(a, b))
        case "~~", "base":
        case "<<", "shovel":
        case ">>", "shift":
        case "<", "below":
            i.Register(Below(a, b))
        case "<=", "under":
            i.Register(Under(a, b))
        case ">", "above":
            i.Register(Above(a, b))
        case ">=", "over":
            i.Register(Over(a, b))
        case "==", "is":
        case "!=", "isnt":
        case "&", "intersect":
        case "^", "exclude":
        case "|", "union":
        case "&&", "and":
        case "||", "or":
        case "..", "range", "to":
            i.Register(Range(a, b))
        case "??":
        case "=", "assign":
            i.parser.blk.vars[t.VarName(a)] = b
            i.Register(b)
        case ":", "define":
            i.parser.blk.vars[t.VarName(a)] = b
            i.parser.blk.hash[t.VarName(a)] = b
            i.Register(b)
        case "+=":
            val := Add(i.Value(a), i.Value(b))
            i.parser.blk.vars[t.VarName(a)] = val
            i.Register(val)
        case "-=":
        case "*=":
        case "/=":
        case "%=":
        case "&=":
        case "^=":
        case "|=":
        case "", "member":
            i.Register(Member(a, b))
        default:
            t.UnexpectedToken()
        }
    case BLK:
        i.Register(t.blk)
    case FIN:
        if i.parser.blk.dim == MAP && len(i.parser.blk.hash) < len(i.parser.blk.stck) {
            switch x := i.parser.blk.stck[len(i.parser.blk.stck) - 1].(type) {
            case Variable:
                i.parser.blk.hash[string(x)] = String(x)
            case String:
                i.parser.blk.hash[string(x)] = x
            default:
                i.parser.blk.hash[fmt.Sprintf("%v", x)] = x
            }
        }

        if i.parser.blk.idx == len(i.parser.blk.toks) {
            if i.parser.blk.src != nil {
                for key, val := range i.parser.blk.vars {
                    if _, ok := i.parser.blk.hash[key]; !ok {
                        i.parser.blk.src.vars[key] = val
                    }
                }
            }

            switch i.parser.blk.dim {
            case VAL:
                if len(i.parser.blk.stck) > 0 {
                    return i.Value(i.parser.blk.stck[len(i.parser.blk.stck) - 1])
                } else {
                    return Null { }
                }
            case LST:
                return i.Value(i.parser.blk.stck)
            case MAP:
                return i.Value(i.parser.blk.hash)
            default:
                return Null { }
            }
        }
    case VAR:
        i.parser.blk.vars[t.lit] = i.parser.blk.FindVar(t.lit)
        i.Register(Variable(t.lit))
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(Number(val))
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(Number(val))
        }
    case STR:
        i.Register(String(t.lit))
    case COM:
        if i.parser.blk.src != nil {
            i.parser.blk.src.coms = append(i.parser.blk.src.coms, t.lit)
        }
    default:
        t.UnexpectedToken()
    }

    return i.Interpret()
}
