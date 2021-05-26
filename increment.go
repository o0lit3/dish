package main

func Increment(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Increment(x.Run())
    case Hash:
        return x.Increment()
    case Array:
        return x.Increment()
    case String:
        return x.Increment()
    case Number:
        return x + 1
    case Boolean:
        return Increment(x.Number())
    }

    return Number(1)
}

func (a Hash) Increment() Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Increment(val)
    }

    return out
}

func (a Array) Increment() Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Increment(val))
    }

    return out
}

func (a String) Increment() String {
    out := ""

    for _, c := range a {
        out += string(c + 1)
    }

    return String(out)
}
