package main
import("fmt")

func Add(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Hash:
            return AddHash(x, y)
        case Array:
            return AddHash(x, y.Hash())
        case String:
            return AddHash(x, Hash { string(y): y })
        default:
            return AddHash(x, Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case Hash:
            return AddArray(x, y.Array())
        case Array:
            return AddArray(x, y)
        default:
            return AddArray(x, Array { y })
        }
    case String:
        switch y := b.(type) {
        case Hash:
            return x.Number() + Number(len(y))
        case Array:
            return x.Number() + Number(len(y))
        case String:
            return x.Number() + y.Number()
        case Number:
            return x.Number() + y
        case Boolean:
            return x.Number() + y.Number()
        }
    case Number:
        switch y := b.(type) {
        case Hash:
            return x + Number(len(y))
        case Array:
            return x + Number(len(y))
        case String:
            return x + y.Number()
        case Number:
            return x + y
        case Boolean:
            return x + y.Number()
        }
    case Boolean:
        switch y := b.(type) {
        case Hash:
            return x || Boolean(len(y) > 0)
        case Array:
            return x || Boolean(len(y) > 0)
        case String:
            return x || Boolean(y != "")
        case Number:
            return x || Boolean(y != 0)
        case Boolean:
            return x || y
        }
    }

    return Null { }
}

func AddHash(a Hash, b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func AddArray(a Array, b Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}
