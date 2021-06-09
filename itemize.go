package main
import("math/big")

func Itemize(a interface{}) Array {
    switch x := a.(type) {
    case *Block:
        return Itemize(x.Run())
    case *Variable:
        return Itemize(x.Value())
    case Hash:
        return x.Array()
    case Array:
        return x
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
