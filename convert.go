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
            return x.Convert(Number(len(y)))
        case Array:
            return x.Convert(Number(len(y)))
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
            return x.Convert(Number(len(y)))
        case Array:
            return x.Convert(Number(len(y)))
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
            return x.Convert(Number(len(y)))
        case Array:
            return x.Convert(Number(len(y)))
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
        return Number(0)
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
    if b < 2 || b > 36 {
        return Number(0)
    }

    out, _ := strconv.ParseInt(string(a), int(b), 64)

    return Number(out)
}
