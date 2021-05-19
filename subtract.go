package main
import("fmt")

func Subtract(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return Subtract(x, y.Run())
        case Array:
            return x.Subtract(y.Hash())
        case String:
            return x.Subtract(Hash { string(y): y })
        default:
            return x.Subtract(Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return Subtract(x, y.Run())
        case Array:
            return x.Subtract(y)
        default:
            return x.Subtract(Array { y })
        }
    case String:
        switch y := b.(type) {
        case Interpreter:
            return Subtract(x, y.Run())
        case Array:
            return Array{ x }.Subtract(y)
        case String:
            return x.Number() - y.Number()
        default:
            return Subtract(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case Interpreter:
            return Subtract(x, y.Run())
        case Array:
            return Array{ x }.Subtract(y)
        case String:
            return x - y.Number()
        case Number:
            return x - y
        case Boolean:
            return x - y.Number()
        case Null:
            return x - Number(0)
        }
    case Boolean:
        return Subtract(x.Number(), b)
    case Null:
        return Subtract(Number(0), b)
    }

    return Null { }
}

func (a Hash) Subtract(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, _ := range b {
        if _, ok := out[key]; ok {
            delete(out, key)
        }
    }

    return out
}

func (a Array) Subtract(b Array) Array {
    out := Array { }
    hash := Hash { }

    for _, val := range b {
        hash[fmt.Sprintf("%v", val)] = val
    }

    for _, val := range a {
        key := fmt.Sprintf("%v", val)

        if _, ok := hash[key]; ok {
            delete(hash, key)
        } else {
            out = append(out, val)
        }
    }

    return out
}
