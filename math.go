package main

import(
    "math"
    "math/big"
)

func Math(a interface{}, tok *Token) interface{} {
    switch x := a.(type) {
    case *Block:
        return Math(x.Run(), tok)
    case *Variable:
        return Math(x.Value(), tok)
    case Hash:
        return Math(NewNumber(len(x)), tok)
    case Array:
        return Math(NewNumber(len(x)), tok)
    case String:
        return Math(x.Number(), tok)
    case Number:
        val, _ := x.val.Float64()

        switch tok.lit {
        case "sqrt":
            return Number{ val: new(big.Rat).SetFloat64(math.Sqrt(val)) }
        case "log":
            return Number{ val: new(big.Rat).SetFloat64(math.Log(val)) }
        case "sin":
            return Number{ val: new(big.Rat).SetFloat64(math.Sin(val)) }
        case "cos":
            return Number{ val: new(big.Rat).SetFloat64(math.Cos(val)) }
        case "tan":
            return Number{ val: new(big.Rat).SetFloat64(math.Tan(val)) }
        case "asin":
            return Number{ val: new(big.Rat).SetFloat64(math.Asin(val)) }
        case "acos":
            return Number{ val: new(big.Rat).SetFloat64(math.Acos(val)) }
        case "atan":
            return Number{ val: new(big.Rat).SetFloat64(math.Atan(val)) }
        }
    case Boolean:
        return Math(x.Number(), tok)
    case Null:
        return Math(NewNumber(0), tok)
    }

    return Null { }
}
