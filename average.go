package main

func Average(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Average(x.Run())
    case Hash:
        return x.Array().Average()
    case Array:
        return x.Average()
    case String:
        return x.Average()
    case Number:
        return x
    case Boolean:
        return x.Number()
    }

    return Number(0)
}

func (a Array) Average() interface{} {
    return Divide(Sum(a), Length(a))
}

func (a String) Average() interface{} {
    out := 0

    for _, c := range a {
        out += int(c)
    }

    return String(string(out / len(a)))
}
