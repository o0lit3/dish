package main
import("strconv")

func Base(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Base(x.Run(), b)
    case *Variable:
        return Base(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
           return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case Boolean:
        return Base(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Base(b Number) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Base(val, b)
    }

    return out
}

func (a Array) Base(b Number) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Base(val, b))
    }

    return out
}

func (a String) Base(b Number) Number {
    if b.val.Cmp(NewNumber(2).val) == -1 || b.val.Cmp(NewNumber(36).val) == 1 {
        return NewNumber(0)
    }

    out, _ := strconv.ParseInt(string(a), b.Int(), 64)

    return NewNumber(int(out))
}

func (a Number) Base(b Number) String {
    if b.val.Cmp(NewNumber(2).val) == -1 || b.val.Cmp(NewNumber(36).val) == 1 {
        return String("")
    }

    return String(strconv.FormatInt(int64(a.Int()), b.Int()))
}
