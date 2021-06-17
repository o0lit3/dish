package main
import("math/big")

func Itemize(a interface{}) Array {
    switch x := a.(type) {
    case *Block:
        return Itemize(x.Run())
    case *Variable:
        return Itemize(x.Value())
    case Hash:
        return Flatten(x)
    case Array:
        return Flatten(x)
    case String:
        return x.Array()
    case Number:
        return Array{
            Number{ val: new(big.Rat).SetInt(x.val.Num()) },
            Number{ val: new(big.Rat).SetInt(x.val.Denom()) },
        }
    case Boolean:
        return Itemize(x.Number())
    }

    return Array{ }
}

func Flatten (a interface{}) Array {
    out := Array { }

    switch x := a.(type) {
    case *Block:
        return Flatten(x.Run())
    case *Variable:
        return Flatten(x.Value())
    case Hash:
        return Flatten(x.Array())
    case Array:
        for _, item := range x {
            for _, val := range Flatten(item) {
                out = append(out, val)
            }
        }
    default:
        out = append(out, x)
    }

    return out
}
