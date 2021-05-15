package main

func Invert(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return InvertHash(x)
    case Array:
        return InvertArray(x)
    case String:
        return InvertString(x)
    case Number:
        return Number(^int(x))
    case Boolean:
        return Boolean(!x)
    }

    return Null { }
}

func InvertHash(a Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Invert(val)
    }

    return out
}

func InvertArray(a Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Invert(val))
    }

    return out
}

func InvertString(a String) String {
    out := ""

    for _, c := range a {
        out += string(^c)
    }

    return String(out)
}
