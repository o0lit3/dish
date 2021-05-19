package main

func Sum(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Array().Sum()
    case Array:
        return x.Sum()
    case String:
        return x.Sum()
    case Number:
        return x
    case Boolean:
        return x.Number()
    case Null:
        return Number(0)
    }

    return Null { }
}

func (a Array) Sum() Number {
    out := Number(0)

    for _, val := range a {
        if x, ok := Sum(val).(Number); ok {
            out += x
        }
    }

    return out
}

func (a String) Sum() Number {
    out := 0

    for _, c := range a {
        out += int(c)
    }

    return Number(out)
}
