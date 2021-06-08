package main

func Product(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Product(x.Run())
    case *Variable:
        return Product(x.Value())
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
        return NewNumber(0)
    }

    return Null { }
}

func (a Array) Product() Number {
    out := NewNumber(1)

    for _, val := range a {
        if x, ok := Product(val).(Number); ok {
            out = Number{ val: out.val.Mul(out.val, x.val) }
        }
    }

    return out
}
