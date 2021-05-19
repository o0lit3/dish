package main

func Decrement(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Decrement()
    case Array:
        return x.Decrement()
    case String:
        return x.Decrement()
    case Number:
        return x - 1
    case Boolean:
        return Decrement(x.Number())
    }

    return Number(-1)
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
