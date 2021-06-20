package main

func Ascii(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Ascii(x.Run())
    case *Variable:
        return Ascii(x.Value())
    case Hash:
        return x.Ascii()
    case Array:
        return x.Ascii()
    case String:
        if len(x) == 1 {
            return x.Ascii()
        }

        return x.Array().Ascii()
    case Number:
        if x.inf == INF || x.inf == -INF {
            return Null { }
        }

        return x.Ascii()
    case Boolean:
        return Ascii(x.Number())
    }

    return Null { }
}

func (a Hash) Ascii() Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Ascii(val)
    }

    return out
}

func (a Array) Ascii() Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Ascii(val))
    }

    return out
}

func (a String) Ascii() Number {
    if len(a) > 0 {
        return NewNumber(int(a[0]))
    }

    return NewNumber(0)
}

func (a Number) Ascii() String {
    return String(string(rune(a.Int())))
}
