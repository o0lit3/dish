package main
import("fmt")

func Add(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Add(x.Run(), b)
    case *Variable:
        return Add(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return x.Add(y)
        case Array:
            return x.Add(y.Hash())
        case String:
            return x.Add(Hash { string(y): y })
        default:
            return x.Add(Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return x.Add(y.Array())
        case Array:
            return x.Add(y)
        default:
            return x.Add(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return Add(x.Number(), NewNumber(len(y)))
        case Array:
            return Add(x.Number(), NewNumber(len(y)))
        case String:
            return Number{ val: NewNumber(0).val.Add(x.Number().val, y.Number().val) }
        default:
            return Add(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return Add(x, NewNumber(len(y)))
        case Array:
            return Add(x, NewNumber(len(y)))
        case String:
            return Number{ val: NewNumber(0).val.Add(x.val, y.Number().val) }
        case Number:
            return Number{ val: NewNumber(0).val.Add(x.val, y.val) }
        case Boolean:
            return Number{ val: NewNumber(0).val.Add(x.val, y.Number().val) }
        case Null:
            return x
        }
    case Boolean:
        return Add(x.Number(), b)
    case Null:
        return Add(NewNumber(0), b)
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
