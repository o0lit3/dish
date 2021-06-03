package main
import("fmt")

func Intersect(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Intersect(x.Run(), b)
    case *Variable:
        return Intersect(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Intersect(x, y.Run())
        case *Variable:
            return Intersect(x, y.Value())
        case Hash:
            return x.Intersect(y)
        case Array:
            return x.Intersect(y.Hash())
        case String:
            return x.Intersect(Hash { string(y): y })
        default:
            return x.Intersect(Hash { fmt.Sprintf("%v", y): y })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Intersect(x, y.Run())
        case *Variable:
            return Intersect(x, y.Value())
        case Hash:
            return x.Intersect(y.Array())
        case Array:
            return x.Intersect(y)
        default:
            return x.Intersect(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Intersect(x, y.Run())
        case *Variable:
            return Intersect(x, y.Value())
        case String:
            return x.Intersect(y)
        default:
            return Intersect(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Intersect(x, y.Run())
        case *Variable:
            return Intersect(x, y.Value())
        case Hash:
            return int(x) & len(y)
        case Array:
            return int(x) & len(y)
        case String:
            return int(x) & int(y.Number())
        case Number:
            return int(x) & int(y)
        case Boolean:
            return int(x) & int(y.Number())
        case Null:
            return int(x) & 0
        }
    case Boolean:
        return Intersect(x.Number(), b)
    case Null:
        return Intersect(Number(0), b)
    }

    return Null { }
}

func (a Hash) Intersect(b Hash) Hash {
    out := Hash { }

    for key := range a {
        if _, ok := b[key]; ok {
            out[key] = b[key]
        }
    }

    return out
}

func (a Array) Intersect(b Array) Array {
    out := Array { }

    for _, aval := range a {
        for _, bval := range b {
            if Equals(aval, bval) {
                out = append(out, aval)
                break
            }
        }
    }

    return out
}

func (a String) Intersect(b String) String {
    out := ""

    if len(b) < len(a) {
        for i := range b {
            out += string(a[i] & b[i])
        }
    } else {
        for i := range a {
            out += string(a[i] & b[i])
        }
    }

    return String(out)
}
