package main

func Decrement(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Decrement(x.Run())
    case *Variable:
        return Decrement(x.Value())
    case Hash:
        return x.Decrement()
    case Array:
        return x.Decrement()
    case String:
        return x.Decrement()
    case Number:
        return Number{ val: NewNumber(0).val.Sub(x.val, NewNumber(1).val) }
    case Boolean:
        return Decrement(x.Number())
    }

    return NewNumber(-1)
}

func (a Hash) Decrement() Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Decrement(val)
    }

    return out
}

func (a Array) Decrement() Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Decrement(val))
    }

    return out
}

func (a String) Decrement() String {
    out := ""

    for _, c := range a {
        out += string(c - 1)
    }

    return String(out)
}
