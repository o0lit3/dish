package main

import (
    "strconv"
    "strings"
)

type Variable struct {
    dep int
    nom string
}

type Block struct {
    dim Dimension
    vars map[string]interface{}
    regs []interface{}
    stck []interface{}
}

type Interpreter struct {
    tics []Token
    blks []Block
}

func (i *Interpreter) Stack(op Token, val interface{}) {
    i.blks[op.dep].stck = append(i.blks[op.dep].stck, val)
}

func (i *Interpreter) Register(op Token, val interface{}) {
    i.blks[op.dep].regs = append(i.blks[op.dep].regs, val)
}

func (i *Interpreter) Deregister(op Token) interface{} {
    if len(i.blks[op.dep].regs) < 1 {
        panic("Missing operand for " +  op.lit + " at " + op.pos.String())
    }

    val := i.blks[op.dep].regs[0]
    i.blks[op.dep].regs = i.blks[op.dep].regs[1:]

    if op.tok == OPA {
        return val
    }

    switch x := val.(type) {
    case Variable:
        return i.blks[x.dep].vars[x.nom]
    default:
        return val
    }
}

func (i *Interpreter) Interpret() {
    t := i.tics[0]
    i.tics = i.tics[1:]

    switch t.tok {
    case OP1:
        a := i.Deregister(t)

        switch t.lit {
        case "!":
            switch x := a.(type) {
            case map[string]interface{}:
                i.Register(t, len(x) == 0)
            case []interface{}:
                i.Register(t, len(x) == 0)
            case string:
                i.Register(t, x == "" || x == "0")
            case float64:
                i.Register(t, x == 0)
            case int:
                i.Register(t, x == 0)
            case bool:
                i.Register(t, !x)
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
        case "~":
        case "+":
        case "-":
        case "++":
        case "--":
        default:
            t.pos.UnexpectedToken(t.lit)
        }
    case OP2:
        a := i.Deregister(t)
        b := i.Deregister(t)

        switch t.lit {
        case "**":
        case "*":
            switch x := a.(type) {
            case []interface{}:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                default:
                    t.pos.UnexpectedOperand(t.lit)
                }

                switch y {
                case 0:
                    i.Register(t, nil)
                case 1:
                    i.Register(t, x)
                default:
                    for n:= 1; n < y; n++ {
                        for _, val := range x {
                            x = append(x, val)
                        }
                    }

                    i.Register(t, x)
                }
            case string:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                i.Register(t, strings.Repeat(x, y))
            case float64:
                switch y := b.(type) {
                case int:
                    i.Register(t, x * float64(y))
                case float64:
                    i.Register(t, x * y)
                }
            case int:
                switch y := b.(type) {
                case int:
                    i.Register(t, x * y)
                case float64:
                    i.Register(t, float64(x) * y)
                }
            case bool:
                if x {
                    i.Register(t, b)
                } else {
                    i.Register(t, 0)
                }
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
        case "/":
        case "%":
        case "+":
            a := i.Deregister(t)
            b := i.Deregister(t)

            switch x := a.(type) {
            case map[string]interface{}:
                switch y := b.(type) {
                case map[string]interface{}:
                    for key, val := range y {
                        x[key] = val
                    }

                    i.Register(t, x)
                case []interface{}:
                    for idx, val := range y {
                        x[strconv.Itoa(idx)] = val
                    }

                    i.Register(t, x)
                default:
                    t.pos.UnexpectedOperand(t.lit)
                }
            case []interface{}:
                switch y := b.(type) {
                case map[string]interface{}:
                    for _, val := range y {
                        x = append(x, val)
                    }

                    i.Register(t, x)
                case []interface{}:
                    for _, val := range y {
                        x = append(x, val)
                    }

                    i.Register(t, x)
                default:
                    t.pos.UnexpectedOperand(t.lit)
                }
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
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
            t.pos.UnexpectedToken(t.lit)
        }
    case OPA:
        a := i.Deregister(t)

        switch x := a.(type) {
        case Variable:
            switch t.lit {
            case "=", ":":
                b := i.Deregister(t)
                i.blks[x.dep].vars[x.nom] = b
                i.Register(t, b)
            case "+=":
            case "-=":
            case "**=":
            case "*=":
            case "/=":
            case "%=":
            case "&=":
            case "|=":
            default:
                t.pos.UnexpectedToken(t.lit)
            }
        default:
            panic("Assignment operator requires variable operand near " + t.lit + " at " + t.pos.String())
        }
    case FIN:
        switch len(i.blks[t.dep].regs) {
        case 0:
        case 1:
            i.Stack(t, i.blks[t.dep].regs[0])
        default:
            panic("Too many items in register near " + t.lit + " at " + t.pos.String())
        }

        i.blks[t.dep].regs = nil

        if (len(i.tics) > 0 && i.tics[0].dep < t.dep) {
            switch len(i.blks[t.dep].stck) {
            case 0:
                switch t.dim {
                case VAL:
                case SET:
                    var x []interface{}
                    i.Register(i.tics[0], x)
                case MAP:
                    var x map[string]interface{}
                    i.Register(i.tics[0], x)
                }
            default:
                switch t.dim {
                case VAL:
                    i.Register(i.tics[0], i.blks[t.dep].stck[len(i.blks[t.dep].stck) - 1])
                case SET:
                    i.Register(i.tics[0], i.blks[t.dep].stck)
                case MAP:
                    i.Register(i.tics[0], i.blks[t.dep].vars)
                }
            }

            i.blks[t.dep].stck = nil
            i.blks[t.dep].vars = make(map[string]interface{})
        }
    case VAR:
        n := t.dep
        _, ok := i.blks[n].vars[t.lit]

        for (n > 0 && !ok) {
            n--
            _, ok = i.blks[n].vars[t.lit]
        }

        if (n == 0 && !ok) {
            n = t.dep
        }

        i.Register(t, Variable { dep: n, nom: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic("Malformed float at " + t.pos.String())
            }

            i.Register(t, val)
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic("Malformed int at " + t.pos.String())
            }

            i.Register(t, val)
        }
    case STR:
        i.Register(t, t.lit)
    default:
        t.pos.UnexpectedToken(t.lit)
    }
}

