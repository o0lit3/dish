package main

func Invert(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Invert()
    case Array:
        return x.Invert()
    case String:
        return x.Invert()
    case Number:
        return Number(^int(x))
    case Boolean:
        return Boolean(!x)
    }

    return Number(^0)
}

func (a Hash) Invert() Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Invert(val)
    }

    return out
}

func (a Array) Invert() Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Invert(val))
    }

    return out
}

func (a String) Invert() String {
    out := ""

    for _, c := range a {
        out += string(^c)
    }

    return String(out)
}
