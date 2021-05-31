package main
import("strconv")

func Base(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Base(x.Run(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case Hash:
            return x.Base(Number(len(y)))
        case Array:
            return x.Base(Number(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case Hash:
            return x.Base(Number(len(y)))
        case Array:
            return x.Base(Number(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case String:
        return Base(x.Number(), b)
    case Number:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case Hash:
            return x.Base(Number(len(y)))
        case Array:
            return x.Base(Number(len(y)))
        case String:
            return x.Base(y.Number())
        case Number:
            return x.Base(y)
        }
    case Boolean:
        return Base(x.Number(), b)
    case Null:
        return String("0")
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

func (a Number) Base(b Number) String {
    if b < 2 || b > 36 {
        return String("")
    }

    return String(strconv.FormatInt(int64(a), int(b)))
}
