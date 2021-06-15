package main

func Keys(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Keys(x.Run())
    case *Variable:
        return Keys(x.Value())
    case Hash:
        return x.Keys()
    case Array:
        return x.Reverse()
    case String:
        return x.Reverse()
    case Number:
        return String(x.String()).Reverse().Number()
    case Boolean:
        return Not(x)
    default:
        return Null { }
    }
}

func (a Hash) Keys() Array {
    out := Array { }

    for key := range a {
        out = append(out, String(key))
    }

    return out
}

func (a Array) Reverse() Array {
    out := Array { }

    for i := range a {
        out = append(out, a[len(a) - i - 1])
    }

    return out
}

func (a String) Reverse() String {
    out := ""

    for i := range a {
        out += string(a[len(a) - i - 1])
    }

    return String(out)
}
