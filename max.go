package main
import("math/big")

func Max(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Max(x.Run())
    case *Variable:
        return Max(x.Value())
    case Hash:
        return x.Array().Max()
    case Array:
        return x.Max()
    case String:
        return Max(x.Number())
    case Number:
        if x.val.Denom().Cmp(big.NewInt(1)) == 0 {
            return x
        }

        return Add(Min(x), NewNumber(1))
    case Boolean:
        return x.Number()
    }

    return NewNumber(0)
}

func (a Array) Max() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if Above(val, out) {
            out = val
        }
    }

    return out
}
