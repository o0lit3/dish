package main

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
    default:
        return x
    }
}

func (a Hash) Compact() Hash {
    out := Hash { }

    for key, val := range a {
        if !Not(val) {
            out[key] = val
        }
    }

    return out
}

func (a Array) Compact() Array {
    out := Array { }

    for _, val := range a {
        if !Not(val) {
            out = append(out, val)
        }
    }

    return out
}
