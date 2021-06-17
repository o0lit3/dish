package main
import("math/big")

func Min(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Min(x.Run())
    case *Variable:
        return Min(x.Value())
    case Hash:
        return x.Array().Min()
    case Array:
        return x.Min()
    case String:
        return Min(x.Number())
    case Number:
        v := new(big.Int).Quo(x.val.Num(), x.val.Denom())
        return Number{ val: new(big.Rat).SetInt(v) }
    case Boolean:
        return x.Number()
    }

    return NewNumber(0)
}

func (a Array) Min() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if Below(val, out) {
            out = val
        }
    }

    return out
}
