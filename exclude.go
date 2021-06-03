package main
import("fmt")

func Exclude(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Exclude(x.Run(), b)
    case *Variable:
        return Exclude(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Exclude(x, y.Run())
        case *Variable:
            return Exclude(x, y.Value())
        case Hash:
            return x.Exclude(y)
        case Array:
            return x.Exclude(y.Hash())
        case String:
            return x.Exclude(Hash { string(y): y })
        default:
            return x.Exclude(Hash { fmt.Sprintf("%v", y): y })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Exclude(x, y.Run())
        case *Variable:
            return Exclude(x, y.Value())
        case Hash:
            return x.Exclude(y.Array())
        case Array:
            return x.Exclude(y)
        default:
            return x.Exclude(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Exclude(x, y.Run())
        case *Variable:
            return Exclude(x, y.Value())
        case String:
            return x.Exclude(y)
        default:
            return Exclude(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Exclude(x, y.Run())
        case *Variable:
            return Exclude(x, y.Value())
        case Hash:
            return int(x) ^ len(y)
        case Array:
            return int(x) ^ len(y)
        case String:
            return int(x) ^ int(y.Number())
        case Number:
            return int(x) ^ int(y)
        case Boolean:
            return int(x) ^ int(y.Number())
        case Null:
            return int(x) ^ 0
        }
    case Boolean:
        return Exclude(x.Number(), b)
    case Null:
        return Exclude(Number(0), b)
    }

    return Null { }
}

func (a Hash) Exclude(b Hash) Hash {
    out := Hash { }

    for key := range a {
        if _, ok := b[key]; !ok {
            out[key] = a[key]
        }
    }

    for key := range b {
        if _, ok := a[key]; !ok {
            out[key] = b[key]
        }
    }

    return out
}

func (a Array) Exclude(b Array) Array {
    out := Array { }

    for _, aval := range a {
        found := false

        for _, bval := range b {
            if Equals(aval, bval) {
                found = true
                break
            }
        }

        if !found {
            out = append(out, aval)
        }
    }

    for _, bval := range b {
        found := false

        for _, aval := range a {
            if Equals(bval, aval) {
                found = true
                break
            }
        }

        if !found {
            out = append(out, bval)
        }
    }

    return out
}

func (a String) Exclude(b String) String {
    out := ""

    if len(b) > len(a) {
        for i := range b {
            if i < len(a) {
                out += string(a[i] ^ b[i])
            } else {
                out += string(b[i])
            }
        }
    } else {
        for i := range a {
            if i < len(b) {
                out += string(a[i] ^ b[i])
            } else {
                out += string(a[i])
            }
        }
    }

    return String(out)
}