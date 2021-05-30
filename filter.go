package main

func Filter(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Filter(x.Run(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return x.Filter(y)
        default:
            if Not(y) {
                return Hash { }
            }

            return x
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return x.Filter(y)
        default:
            if Not(y) {
                return Array { }
            }

            return x
        }
    case String:
        return Filter(x.Array(), b)
    }

    return Null { }
}

func (a Hash) Filter(b *Block) Hash {
    out := Hash { }

    for key, val := range a {
        if Not(b.Run(val, String(key))) {
            continue
        }

        out[key] = val
    }

    return out
}

func (a Array) Filter(b *Block) Array {
    out := Array { }

    for i, val := range a {
        if Not(b.Run(val, Number(i))) {
            continue
        }

        out = append(out, val)
    }

    return out
}
