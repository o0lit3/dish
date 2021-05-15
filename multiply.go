package main

func Multiply(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Hash:
            return MultiplyArray(x.Array(), Number(len(y)))
        case Array:
            return MultiplyArray(x.Array(), Number(len(y)))
        case String:
            return Null { }
        case Number:
            return MultiplyArray(x.Array(), y)
        case Boolean:
            return Identity(x, y)
        }
    case Array:
        switch y := b.(type) {
        case Hash:
            return MultiplyArray(x, Number(len(y)))
        case Array:
            return MultiplyArray(x, Number(len(y)))
        case String:
            return Null { }
        case Number:
            return MultiplyArray(x, y)
        case Boolean:
            return Identity(x, y)
        }
    case String:
        switch y := b.(type) {
        case Hash:
            return Null { }
        case Array:
            return Null { }
        case String:
            return MultiplyString(x, Number(len(y)))
        case Number:
            return MultiplyString(x, y)
        case Boolean:
            return Identity(x, y)
        }
    case Number:
        switch y := b.(type) {
        case Hash:
            return MultiplyArray(y.Array(), x)
        case Array:
            return MultiplyArray(y, x)
        case String:
            return MultiplyString(y, x)
        case Number:
            return x * y
        case Boolean:
            return Identity(x, y)
        }
    case Boolean:
        return Identity(b, x)
    }

    return Null { }
}

func MultiplyArray(a Array, b Number) Array {
    out := Array { }

    for n := 0; n < int(b); n++ {
        for _, val := range a {
            out = append(out, val)
        }
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, val := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out = append(out, val)
            } else {
                break
            }
        }
    }

    return out
}

func MultiplyString(a String, b Number) String {
    out := ""

    for n := 0; n < int(b); n++ {
        out += string(a)
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, c := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return String(out)
}

func Identity(a interface{}, b Boolean) interface{} {
    if b {
        return a
    } else {
       return Null { }
    }
}
