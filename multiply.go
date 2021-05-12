package main

import ("fmt")

func Multiply(t Token, a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        switch y := b.(type) {
        case string:
            return JoinSet(x.ToSet(), y)
        case float64:
            return MultiplySet(x.ToSet(), int(y))
        case int:
            return MultiplySet(x.ToSet(), y)
        default:
            t.UnexpectedOperand()
        }
    case Set:
        switch y := b.(type) {
        case string:
            return JoinSet(x, y)
        case float64:
            return MultiplySet(x, int(y))
        case int:
            return MultiplySet(x, y)
        default:
            t.UnexpectedOperand()
        }
    case string:
        switch y := b.(type) {
        case float64:
            return MultiplyString(x, int(y))
        case int:
            return MultiplyString(x, y)
        default:
            t.UnexpectedOperand()
        }
    case float64:
        switch y := b.(type) {
        case float64:
            return x * y
        case int:
            return x * float64(y)
        }
    case int:
        switch y := b.(type) {
        case float64:
            return float64(x) * y
        case int:
            return x * y
        }
    case bool:
        if x {
            return b
        } else {
            return 0
        }
    default:
        t.UnexpectedOperand()
    }

    return 0
}

func JoinSet(a Set, b string) string {
    out := ""

    for i, val := range a {
        if i == 0 {
            out += fmt.Sprintf("%v", val)
        } else {
            out += b + fmt.Sprintf("%v", val)
        }
    }

    return out
}

func MultiplySet(a Set, b int) Set {
    out := Set { }

    for n := 0; n < b; n++ {
        for _, val := range a {
            out = append(out, val)
        }
    }

    return out
}

func MultiplyString(a string, b int) string {
    out := ""

    for n := 0; n < b; n++ {
        out += a
    }

    return out
}
