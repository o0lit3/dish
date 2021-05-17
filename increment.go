package main

func Increment(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return IncrementHash(x)
    case Array:
        return IncrementArray(x)
    case String:
        return IncrementString(x)
    case Number:
        return x + 1
    case Boolean:
        return Boolean(true)
    }

    return Number(1)
}

func IncrementHash(a Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Increment(val)
    }

    return out
}

func IncrementArray(a Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Increment(val))
    }

    return out
}

func IncrementString(a String) String {
    out := ""

    for _, c := range a {
        out += string(c + 1)
    }

    return String(out)
}
