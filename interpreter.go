package main

import (
    "fmt"
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
    stck Set
    vars Map
    hash []string
}

type Map map[string]interface{}

type Set []interface{}

type String string

type Number float64

type Boolean bool

type Interpreter struct {
    tics []Token
    blks []Block
    comm string
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

func (s Set) String() string {
    var out []string

    for _, val := range s {
        out = append(out, fmt.Sprintf("%v", val))
    }

    return fmt.Sprintf("[" + strings.Join(out, ", ") + "]")
}

func (s Set) Map() Map {
    out := Map { }

    for _, val := range s {
        out[fmt.Sprintf("%v", val)] = val
    }

    return out
}

func (m Map) String() string {
    var out []string

    for key, val := range m {
        out = append(out, fmt.Sprintf("\"%s\": %v", key, val))
    }

    return fmt.Sprintf("{" + strings.Join(out, ", ") + "}")
}

func (m Map) Set() Set {
    out := Set { }

    for _, val := range m {
        out = append(out, val)
    }

    return out
}

func (i *Interpreter) Bind(op Token, val interface{}) interface{} {
    switch x := val.(type) {
    case Variable:
        if op.tok == OPA {
            return x
        } else {
            return i.Bind(op, i.blks[x.dep].vars[x.nom])
        }
    case string:
        return String(x)
    case float64:
        return Number(x)
    case int:
        return Number(x)
    case bool:
        return Boolean(x)
    default:
        return val
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

func (i *Interpreter) Interpret() Token {
    t := i.tics[0]
    i.tics = i.tics[1:]

    switch t.tok {
    case OP1:
        a := i.Deregister(t)

        switch t.lit {
        case "!":
            i.Register(t, Not(t, a))
        case "~":
        case "+":
        case "-":
        case "++":
        case "--":
        default:
            t.UnexpectedToken()
        }
    case OP2:
        b := i.Deregister(t)
        a := i.Deregister(t)

        switch t.lit {
        case "**":
        case "*":
            i.Register(t, Times(t, a, b))
        case "/":
        case "%":
        case "+":
            i.Register(t, Plus(t, a, b))
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
        case "..":
        case "??":
        default:
            t.UnexpectedToken()
        }
    case OPA:
        b := i.Deregister(t)
        a := i.Deregister(t)

        dep := t.dep
        nom := ""
        val := b

        switch x := a.(type) {
        case Variable:
            nom = x.nom
        default:
            if t.lit == ":" {
                nom = fmt.Sprintf("%v", x)
            } else {
                panic(fmt.Sprintf("Assignment operator \"%s\" requires left-hand variable operand at %s", t.lit, t.pos))
            }
        }

        switch y := b.(type) {
        case Variable:
            val = i.blks[y.dep].vars[y.nom]
        }

        switch t.lit {
        case "=":
            i.blks[dep].vars[nom] = val
        case ":":
            if t.dim == MAP {
                i.blks[dep].vars[nom] = val

                if len(i.blks[dep].hash) > len(i.blks[dep].stck) {
                    i.blks[dep].hash[len(i.blks[dep].hash) - 1] = nom
                } else {
                    i.blks[dep].hash = append(i.blks[dep].hash, nom)
                }
            } else {
                panic(fmt.Sprintf("Assigment operator \"%s\" can only be used in hashes at %s", t.lit, t.pos))
            }
        case "+=":
        case "-=":
        case "**=":
        case "*=":
        case "/=":
        case "%=":
        case "&=":
        case "|=":
        default:
            t.UnexpectedToken()
        }

        i.Register(t, i.blks[dep].vars[nom])
    case OPX:
        var a interface{}
        var b interface{}

        switch t.lit {
        case "not", "bnot":
            a = i.Deregister(t)
        default:
            b = i.Deregister(t)
            a = i.Deregister(t)
        }

        switch t.lit {
        case "not":
            i.Register(t, Not(t, a))
        case "plus":
            i.Register(t, Plus(t, a, b))
        case "times":
            i.Register(t, Times(t, a, b))
        case "join":
            i.Register(t, Times(t, a, b))
        default:
            panic(fmt.Sprintf("Unrecognized method \"%s\" at %s", t.lit, t.pos))
        }
    case FIN:
        if len(i.blks[t.dep].stck) > 0 {
            val := i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1]

            switch x := val.(type) {
            case Variable:
                i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1] = i.blks[x.dep].vars[x.nom]
            case string:
                i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1] = String(x)
            }

            if t.dim == MAP && len(i.blks[t.dep].stck) != len(i.blks[t.dep].hash) {
                i.blks[t.dep].hash = append(i.blks[t.dep].hash, fmt.Sprintf("%v", val))
            }
        }

        if t.lit == "" {
            blk := Token { dep: t.dep - 1 }

            switch t.dim {
            case VAL:
                for key, val := range i.blks[t.dep].vars {
                    i.blks[blk.dep].vars[key] = val
                }

                if len(i.blks[t.dep].stck) > 0 {
                    switch x := i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1].(type) {
                    case String:
                        i.Register(blk, String(x))
                    default:
                        i.Register(blk, x)
                    }
                }
            case SET:
                for key, val := range i.blks[t.dep].vars {
                    i.blks[blk.dep].vars[key] = val
                }

                i.Register(blk, i.blks[t.dep].stck)
            case MAP:
                hash := Map { }

                for idx, key := range i.blks[t.dep].hash {
                    hash[key] = i.blks[t.dep].stck[idx]
                }

                i.Register(blk, hash)
            }

            i.comm = ""
            i.blks[t.dep].stck = nil
            i.blks[t.dep].hash = nil
            i.blks[t.dep].vars = make(map[string]interface{})
        }
    case VAR:
        n := t.dep
        _, ok := i.blks[n].vars[t.lit]

        for (n > 0 && !ok) {
            n--
            _, ok = i.blks[n].vars[t.lit]
        }

        if n == 0 && !ok {
            i.blks[t.dep].vars[t.lit] = ""
            n = t.dep
        }

        i.Register(t, Variable { dep: n, nom: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(t, val)
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic(fmt.Sprintf("Malformed number at %s", t.pos))
            }

            i.Register(t, val)
        }
    case STR:
        i.Register(t, t.lit)
    case COM:
        i.comm = t.lit
    default:
        t.UnexpectedToken()
    }

    return t
}
