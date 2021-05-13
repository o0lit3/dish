package main
import("fmt")

func Plus(t Token, a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        switch y := b.(type) {
        case Map:
            return AddMap(x, y)
        case Set:
            return AddMap(x, y.Map())
        case String:
            return AddMap(x, Map { string(y): y })
        default:
            return AddMap(x, Map { fmt.Sprintf("%v", y): y})
        }
    case Set:
        switch y := b.(type) {
        case Map:
            return AddSet(x, y.Set())
        case Set:
            return AddSet(x, y)
        case String:
            return AddSet(x, Set { y })
        default:
            return AddSet(x, Set { y })
        }
    case String:
        switch y := b.(type) {
        case String:
            return String(x + y)
        default:
            return String(string(x) + fmt.Sprintf("%v", y))
        }
    case Number:
        switch y := b.(type) {
        case Map:
            return x + Number(len(y))
        case Set:
            return x + Number(len(y))
        case String:
            return x + y.Number()
        case Number:
            return x + y
        case Boolean:
            return map[bool]Number{true: Number(x + 1), false: x}[bool(y)]
        }
    case Boolean:
        switch y := b.(type) {
        case Map:
            return x || Boolean(len(y) > 0)
        case Set:
            return x || Boolean(len(y) > 0)
        case String:
            return x || Boolean(y != "")
        case Number:
            return x || Boolean(y != 0)
        case Boolean:
            return x || y
        }
    }

    return String("")
}

func AddMap(a Map, b Map) Map {
    out := Map { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func AddSet(a Set, b Set) Set {
    out := Set { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}
