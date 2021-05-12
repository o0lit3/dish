package main

import("fmt")

func Multiply(t Token, a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        switch y := b.(type) {
        case string:
            return JoinSet(x.ToSet(), y)
        case float64:
            return MultiplySet(x.ToSet(), y)
        case int:
            return MultiplySet(x.ToSet(), float64(y))
        case bool:
            return Identity(x, y)
        default:
            t.UnexpectedOperand()
        }
    case Set:
        switch y := b.(type) {
        case string:
            return JoinSet(x, y)
        case float64:
            return MultiplySet(x, y)
        case int:
            return MultiplySet(x, float64(y))
        case bool:
            return Identity(x, y)
        default:
            t.UnexpectedOperand()
        }
    case string:
        switch y := b.(type) {
        case Map:
            return JoinSet(y.ToSet(), x)
        case Set:
            return JoinSet(y, x)
        case float64:
            return MultiplyString(x, y)
        case int:
            return MultiplyString(x, float64(y))
        case bool:
            return Identity(x, y)
        default:
            t.UnexpectedOperand()
        }
    case float64:
        switch y := b.(type) {
        case Map:
            return MultiplySet(y.ToSet(), x)
        case Set:
            return MultiplySet(y, x)
        case string:
            return MultiplyString(y, x)
        case float64:
            return x * y
        case int:
            return x * float64(y)
        case bool:
            return Identity(x, y)
        }
    case int:
        switch y := b.(type) {
        case Map:
            return MultiplySet(y.ToSet(), float64(x))
        case Set:
            return MultiplySet(y, float64(x))
        case string:
            return MultiplyString(y, float64(x))
        case float64:
            return float64(x) * y
        case int:
            return x * y
        case bool:
            return Identity(x, y)
        }
    case bool:
        return Identity(b, x)
    default:
        t.UnexpectedOperand()
    }

    return 0
}

func JoinSet(a Set, b string) string {
    out := ""

    for i, val := range a {
        switch x := val.(type) {
        case string:
            out += map[bool]string{true: b, false: ""}[i > 0] + x
        default:
            out += map[bool]string{true: b, false: ""}[i > 0] + fmt.Sprintf("%v", x)
        }
    }

    return out
}

func MultiplySet(a Set, b float64) Set {
    out := Set { }

    for n := 0; n < int(b); n++ {
        for _, val := range a {
            out = append(out, val)
        }
    }

    if b != float64(int(b)) {
        rem := float64(len(a)) * (b - float64(int(b)))

        for _, val := range a {
            if float64(len(out)) < float64(int(b) * len(a)) + rem {
                out = append(out, val)
            } else {
                break
            }
        }
    }

    return out
}

func MultiplyString(a string, b float64) string {
    out := ""

    for n := 0; n < int(b); n++ {
        out += a
    }

    if b != float64(int(b)) {
        rem := float64(len(a)) * (b - float64(int(b)))

        for _, c := range a {
            if float64(len(out)) < float64(int(b) * len(a)) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return out
}

func Identity(a interface{}, b bool) interface{} {
    return map[bool]interface{}{true: a, false: 0}[b]
}
