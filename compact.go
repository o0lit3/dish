package main
import("unicode")

func Compact(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Compact(x.Run())
    case *Variable:
        return Compact(x.Value())
    case Hash:
        return x.Compact()
    case Array:
        return x.Compact()
    case String:
        return Boolean(true)
    case Number:
        return Boolean(true)
    case Boolean:
        return Boolean(true)
    default:
        return Boolean(false)
    }
}

func (a Hash) Compact() Hash {
    out := Hash { }

    for key, val := range a {
        if Boolify(val) {
            out[key] = val
        }
    }

    return out
}

func (a Array) Compact() Array {
    out := Array { }

    for _, val := range a {
        if Boolify(val) {
            out = append(out, val)
        }
    }

    return out
}

func (a String) Letter() Boolean {
    if len(a) > 0 && unicode.IsLetter(a[0]) {
        return Boolean(true)
    }

    return Boolean(false)
}
