package main
import("fmt")

func Unshift(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Unshift(x.Run(), b)
    case *Variable:
        return Unshift(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Unshift(x, y.Run())
        case *Variable:
            return Unshift(x, y.Value())
        case Hash:
            return x.Unshift(y)
        case Array:
            return x.Unshift(y.Hash())
        case String:
            return x.Unshift(Hash { string(y): y })
        default:
            return x.Unshift(Hash { fmt.Sprintf("%v", y): y })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Unshift(x, y.Run())
        case *Variable:
            return Unshift(x, y.Value())
        case Hash:
            return x.Unshift(y.Array())
        case Array:
            return x.Unshift(y)
        default:
            return x.Unshift(Array { y })
        }
    case String:
        return Unshift(x.Number(), b)
    case Number:
        switch y := b.(type) {
        case *Block:
            return Unshift(x, y.Run())
        case *Variable:
            return Unshift(x, y.Value())
        case Hash:
            return int(x) >> uint(len(y))
        case Array:
            return int(x) >> uint(len(y))
        case String:
            return int(x) >> uint(y.Number())
        case Number:
            return int(x) >> uint(y)
        case Boolean:
            return int(x) >> uint(y.Number())
        case Null:
            return int(x) << 0
        }
    case Boolean:
        return Unshift(x.Number(), b)
    case Null:
        return Unshift(Number(0), b)
    }

    return Number(0)
}

func (a Hash) Unshift(b Hash) Hash {
    out := Hash { }

    for key, val := range b {
        out[key] = val
    }

    for key, val := range a {
        out[key] = val
    }

    return out
}

func (a Array) Unshift(b Array) Array {
    out := Array { }

    for _, val := range b {
        out = append(out, val)
    }

    for _, val := range a {
        out = append(out, val)
    }

    return out
}
