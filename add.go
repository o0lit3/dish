package main
import("fmt")

func Add(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return Add(x, y.Run())
        case Array:
            return x.Add(y.Hash())
        case String:
            return x.Add(Hash { string(y): y })
        default:
            return x.Add(Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return Add(x, y.Run())
        case Array:
            return x.Add(y)
        default:
            return x.Add(Array { y })
        }
    case String:
        switch y := b.(type) {
        case Interpreter:
            return Add(x, y.Run())
        case Array:
            return Array{ x }.Add(y)
        case String:
            return x.Number() + y.Number()
        default:
            return Add(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case Interpreter:
            return Add(x, y.Run())
        case Array:
            return Array{ x }.Add(y)
        case String:
            return x + y.Number()
        case Number:
            return x + y
        case Boolean:
            return x + y.Number()
        case Null:
            return x + Number(0)
        }
    case Boolean:
        return Add(x.Number(), b)
    case Null:
        return Add(Number(0), b)
    }

    return Null { }
}

func (a Hash) Add(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func (a Array) Add(b Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}
