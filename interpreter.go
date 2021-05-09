package main

import (
    "fmt"
    "strconv"
    "unicode"
    "strings"
)

type Variable struct {
    depth int
    name string
}

type Block struct {
    style string
    vars map[string]interface{}
    regs []interface{}
    stack []interface{}
}

type Interpreter struct {
    terms []Token
    blks []Block
}

func (i *Interpreter) Stack(depth int, val interface{}) {
    i.blks[depth].stack = append(i.blks[depth].stack, val)
}

func (i *Interpreter) Register(depth int, val interface{}) {
    i.blks[depth].regs = append(i.blks[depth].regs, val)
}

func (i *Interpreter) Cast(val interface{}, t int) interface{} {
    switch x := val.(type) {
    case Variable:
        if t == VAR {
            return x
        }

        v, ok := i.blks[x.depth].vars[x.name]

        if ok {
            return i.Cast(v, t)
        } else {
            switch t {
            case BIT:
                return false
            case NUM:
                return 0
            case STR:
                return ""
            default:
                return nil
            }
        }
    case map[string]interface{}:
        switch t {
        case BIT:
            return len(x) > 0
        case NUM:
            return len(x)
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case []interface{}:
        switch t {
        case BIT:
            return len(x) > 0
        case NUM:
            return len(x)
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case bool:
        switch t {
        case NUM:
            if x {
                return 1
            } else {
                return 0
            }
        case STR:
            if x {
                return "true"
            } else {
                return "false"
            }
        default:
            return x
        }
    case int:
        switch t {
        case BIT:
            return x != 0
        case STR:
            return strconv.Itoa(x)
        default:
            return x
        }
    case float64:
        switch t {
        case BIT:
            return x != 0
        case STR:
            return fmt.Sprintf("%v", x)
        default:
            return x
        }
    case string:
        switch t {
        case BIT:
            return x != "" && x != "0"
        case NUM:
            v := ""
            decimal := false

            for _, c := range x {
                if unicode.IsDigit(c) {
                    v += string(c)
                } else if (c == '.' && !decimal) {
                    v += string(c)
                    decimal = true
                }
            }

            if decimal {
                w, _ := strconv.ParseFloat(v, 64)
                return w
            }

            w, _ := strconv.Atoi(v)
            return w
        default:
            return x
        }
    default:
        panic("Could not determine operand type.")
    }
}

func (i *Interpreter) Bind(op Token, t int) interface{} {
    if len(i.blks[op.depth].regs) < 1 {
        panic("Missing operand for " +  op.lit + " at " + op.pos.String())
    }

    val := i.blks[op.depth].regs[0]
    i.blks[op.depth].regs = i.blks[op.depth].regs[1:]

    return i.Cast(val, t)
}

func (i *Interpreter) Interpret() {
    t := i.terms[0]
    i.terms = i.terms[1:]

    switch t.tok {
    case OP1:
        switch t.lit {
        case "!":
            i.Register(t.depth, !i.Bind(t, BIT).(bool))
        case "~":
        case "+":
        case "-":
        case "++":
        case "--":
        default:
            t.pos.UnexpectedToken(t.lit)
        }
    case OP2:
        switch t.lit {
        case "**":
        case "*":
            a := i.Bind(t, OBJ)
            b := i.Bind(t, NUM)

            switch x := a.(type) {
            case []interface{}:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                switch y {
                case 0:
                    i.Register(t.depth, nil)
                case 1:
                    i.Register(t.depth, x)
                default:
                    for n:= 1; n < y; n++ {
                        for _, val := range x {
                            x = append(x, val)
                        }
                    }

                    i.Register(t.depth, x)
                }
            case string:
                y := 0

                switch b.(type) {
                case int:
                    y = b.(int)
                case float64:
                    y = int(b.(float64))
                }

                i.Register(t.depth, strings.Repeat(x, y))
            case float64:
                switch y := b.(type) {
                case int:
                    i.Register(t.depth, x * float64(y))
                case float64:
                    i.Register(t.depth, x * y)
                }
            case int:
                switch y := b.(type) {
                case int:
                    i.Register(t.depth, x * y)
                case float64:
                    i.Register(t.depth, float64(x) * y)
                }
            case bool:
                if x {
                    i.Register(t.depth, b)
                } else {
                    i.Register(t.depth, 0)
                }
            default:
                t.pos.UnexpectedOperand(t.lit)
            }
        case "/":
        case "%":
        case "+":
            a := i.Bind(t, OBJ)
            b := i.Bind(t, OBJ)

            switch x := a.(type) {
            case map[string]interface{}:
                switch y := b.(type) {
                case map[string]interface{}:
                    for key, val := range y {
                        x[key] = val
                    }

                    i.Register(t.depth, x)
                case []interface{}:
                    for idx, val := range y {
                        x[strconv.Itoa(idx)] = val
                    }

                    i.Register(t.depth, x)
                default:
                    t.pos.UnexpectedOperand(t.lit)
                }
            case []interface{}:
                switch y := b.(type) {
                case map[string]interface{}:
                    for _, val := range y {
                        x = append(x, val)
                    }

                    i.Register(t.depth, x)
                case []interface{}:
                    for _, val := range y {
                        x = append(x, val)
                    }

                    i.Register(t.depth, x)
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
        a := i.Bind(t, VAR)

        switch x := a.(type) {
        case Variable:
            switch t.lit {
            case "=":
                b := i.Bind(t, OBJ)
                i.blks[x.depth].vars[x.name] = b
                i.Register(t.depth, b)
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
            fmt.Printf("%T %v\n", x, x)
            panic("Assignment operator requires variable operand near " + t.lit + " at " + t.pos.String())
        }
    case FIN:
        switch len(i.blks[t.depth].regs) {
        case 0:
            i.Stack(t.depth, "")
        case 1:
            i.Stack(t.depth, i.blks[t.depth].regs[0])
        default:
            panic("Too many items in register near " + t.lit + " at " + t.pos.String())
        }

        i.blks[t.depth].regs = nil

        if (len(i.terms) > 0 && i.terms[0].depth < t.depth) {
            switch len(i.blks[t.depth].stack) {
            case 0:
                switch t.style {
                case "()":
                    i.Register(t.depth - 1, "")
                case "[]":
                    var x []interface{}
                    i.Register(t.depth - 1, x)
                case "{}":
                    var x map[string]interface{}
                    i.Register(t.depth - 1, x)
                }
            default:
                switch t.style {
                case "()":
                    i.Register(t.depth - 1, i.blks[t.depth].stack[len(i.blks[t.depth].stack) - 1])
                case "[]", "{}":
                    i.Register(t.depth - 1, i.blks[t.depth].stack)
                }
            }

            i.blks[t.depth].stack = nil
            i.blks[t.depth].vars = nil
        }
    case VAR:
        n := t.depth
        _, ok := i.blks[n].vars[t.lit]

        for (n > 0 && !ok) {
            n--
            _, ok = i.blks[n].vars[t.lit]
        }

        if (n == 0 && !ok) {
            n = t.depth
        }

        i.Register(t.depth, Variable { depth: n, name: t.lit })
    case NUM:
        if strings.Contains(t.lit, ".") {
            val, err := strconv.ParseFloat(t.lit, 64)

            if err != nil {
                panic("Malformed float at " + t.pos.String())
            }

            i.Register(t.depth, val)
        } else {
            val, err := strconv.Atoi(t.lit)

            if err != nil {
                panic("Malformed int at " + t.pos.String())
            }

            i.Register(t.depth, val)
        }
    case STR:
        i.Register(t.depth, t.lit)
    default:
        t.pos.UnexpectedToken(t.lit)
    }
}

