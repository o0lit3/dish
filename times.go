package main
import("fmt")

func Times(t Token, a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        switch y := b.(type) {
        case Map:
            return MultiplySet(x.Set(), Number(len(y)))
        case Set:
            return MultiplySet(x.Set(), Number(len(y)))
        case String:
            return JoinSet(x.Set(), y)
        case Number:
            return MultiplySet(x.Set(), y)
        case Boolean:
            return Identity(x, y)
        }
    case Set:
        switch y := b.(type) {
        case Map:
            return MultiplySet(x, Number(len(y)))
        case Set:
            return MultiplySet(x, Number(len(y)))
        case String:
            return JoinSet(x, y)
        case Number:
            return MultiplySet(x, y)
        case Boolean:
            return Identity(x, y)
        }
    case String:
        switch y := b.(type) {
        case Map:
            return JoinSet(y.Set(), x)
        case Set:
            return JoinSet(y, x)
        case String:
            return MultiplyString(x, Number(len(y)))
        case Number:
            return MultiplyString(x, y)
        case Boolean:
            return Identity(x, y)
        }
    case Number:
        switch y := b.(type) {
        case Map:
            return MultiplySet(y.Set(), x)
        case Set:
            return MultiplySet(y, x)
        case String:
            return MultiplyString(y, x)
        case Number:
            return x * y
        case Boolean:
            return Identity(x, y)
        }
    case Boolean:
        return Identity(b, x)
    }

    return String("")
}

func JoinSet(a Set, b String) String {
    out := ""

    for i, val := range a {
        switch x := val.(type) {
        case String:
            out += map[bool]string{true: string(b), false: ""}[i > 0] + string(x)
        default:
            out += map[bool]string{true: string(b), false: ""}[i > 0] + fmt.Sprintf("%v", x)
        }
    }

    return String(out)
}

func MultiplySet(a Set, b Number) Set {
    out := Set { }

    for n := 0; n < int(b); n++ {
        for _, val := range a {
            out = append(out, val)
        }
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, val := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out = append(out, val)
            } else {
                break
            }
        }
    }

    return out
}

func MultiplyString(a String, b Number) String {
    out := ""

    for n := 0; n < int(b); n++ {
        out += string(a)
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, c := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return String(out)
}

func Identity(a interface{}, b Boolean) interface{} {
    if b {
        return a
    } else {
       return String("")
    }
}
