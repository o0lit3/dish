package main
import("fmt")

func Subtract(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Subtract(x.Run(), b)
    case *Variable:
        return Subtract(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return x.Subtract(y)
        case Array:
            return x.Subtract(y.Hash())
        case String:
            return x.Subtract(Hash { string(y): y })
        default:
            return x.Subtract(Hash { fmt.Sprintf("%v", y): y})
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return x.Subtract(y.Array())
        case Array:
            return x.Subtract(y)
        default:
            return x.Subtract(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return Hash{ string(x): x }.Subtract(y)
        case Array:
            return Array{ x }.Subtract(y)
        case String:
            return Number{ val: NewNumber(0).val.Sub(x.Number().val, y.Number().val) }
        default:
            return Subtract(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return Hash { fmt.Sprintf("%v", x): x }.Subtract(y)
        case Array:
            return Array{ x }.Subtract(y)
        case String:
            return Number{ val: NewNumber(0).val.Sub(x.val, y.Number().val) }
        case Number:
            return Number{ val: NewNumber(0).val.Sub(x.val, y.val) }
        case Boolean:
            return Number{ val: NewNumber(0).val.Sub(x.val, y.Number().val) }
        case Null:
            return x
        }
    case Boolean:
        return Subtract(x.Number(), b)
    case Null:
        return Subtract(NewNumber(0), b)
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
