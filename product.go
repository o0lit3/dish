package main

func Product(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Array().Product()
    case Array:
        return x.Product()
    case String:
        return x.Number()
    case Number:
        return x
    case Boolean:
        return x.Number()
    case Null:
        return Number(0)
    }

    return Null { }
}

func (a Array) Product() Number {
    out := Number(1)

    for _, val := range a {
        if x, ok := Product(val).(Number); ok {
            out *= x
        }
    }

    return out
}
