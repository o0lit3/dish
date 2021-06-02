package main

func Keys(a interface{}) Array {
    switch x := a.(type) {
    case *Block:
        return Keys(x.Run())
    case *Variable:
        return Keys(x.Value())
    case Hash:
        return x.Keys()
    case Array:
        return x.Keys()
    case String:
        return x.Array().Keys()
    default:
        return Array { }
    }
}

func (a Hash) Keys() Array {
    out := Array { }

    for key := range a {
        out = append(out, String(key))
    }

    return out
}

func (a Array) Keys() Array {
    out := Array { }

    for i := range a {
        out = append(out, Number(i))
    }

    return out
}
