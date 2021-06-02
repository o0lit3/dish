package main

func Negate(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Negate(x.Run())
    case *Variable:
        return Negate(x.Value())
    case Hash:
        return x.Array().Negate()
    case Array:
        return x.Negate()
    case String:
        return -x.Number()
    case Number:
        return -x
    case Boolean:
        return -x.Number()
    case Null:
        return Number(0)
    }

    return Null { }
}

func (a Array) Negate() Number {
    out := Number(0)

    for i, val := range a {
        if x, ok := Negate(val).(Number); ok {
            if i == 0 {
                out = x
            } else {
                out += x
            }
        }
    }

    return out
}
