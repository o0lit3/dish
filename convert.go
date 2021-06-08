package main
import("strconv")

func Convert(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Convert(x.Run(), b)
    case *Variable:
        return Convert(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Convert(x, y.Run())
        case *Variable:
            return Convert(x, y.Value())
        case Hash:
            return x.Convert(NewNumber(len(y)))
        case Array:
            return x.Convert(NewNumber(len(y)))
        case String:
            return x.Convert(y.Number())
        case Number:
            return x.Convert(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Convert(x, y.Run())
        case *Variable:
            return Convert(x, y.Value())
        case Hash:
            return x.Convert(NewNumber(len(y)))
        case Array:
            return x.Convert(NewNumber(len(y)))
        case String:
            return x.Convert(y.Number())
        case Number:
            return x.Convert(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
           return Convert(x, y.Run())
        case *Variable:
            return Convert(x, y.Value())
        case Hash:
            return x.Convert(NewNumber(len(y)))
        case Array:
            return x.Convert(NewNumber(len(y)))
        case String:
            return x.Convert(y.Number())
        case Number:
            return x.Convert(y)
        }
    case Number:
        return Convert(String(x.String()), b)
    case Boolean:
        return Base(x.Number(), b)
    case Null:
        return NewNumber(0)
    }

    return Null { }
}

func (a Hash) Convert(b Number) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Convert(val, b)
    }

    return out
}

func (a Array) Convert(b Number) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Convert(val, b))
    }

    return out
}

func (a String) Convert(b Number) Number {
    if b.val.Cmp(NewNumber(2).val) == -1 || b.val.Cmp(NewNumber(36).val) == 1 {
        return NewNumber(0)
    }

    out, _ := strconv.ParseInt(string(a), b.Int(), 64)

    return NewNumber(int(out))
}
