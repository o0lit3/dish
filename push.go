package main
import("fmt")

func Push(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Push(x.Run(), b)
    case *Variable:
        return Push(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return x.Push(y)
        case Array:
            return x.Push(y.Hash())
        case String:
            return x.Push(Hash { string(y): y })
        default:
            return x.Push(Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return x.Push(y.Array())
        case Array:
            return x.Push(y)
        default:
            return x.Push(Array { y })
        }
    case String:
        return Push(x.Number(), b)
    case Number:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return Number(int(x) << uint(len(y)))
        case Array:
            return Number(int(x) << uint(len(y)))
        case String:
            return Number(int(x) << uint(y.Number()))
        case Number:
            return Number(int(x) << uint(y))
        case Boolean:
            return Number(int(x) << uint(y.Number()))
        case Null:
            return Number(int(x) << 0)
        }
    case Boolean:
        return Push(x.Number(), b)
    case Null:
        return Push(Array { }, b)
    }

    return Number(0)
}

func (a Hash) Push(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func (a Array) Push(b Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}
