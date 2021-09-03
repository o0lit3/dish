package main

import(
    "time"
    "math"
    "math/big"
    "math/rand"
)

func (t *Token) Numbers(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Numbers(x.Run())
    case *Variable:
        return t.Numbers(x.Value())
    case Number:
        val, _ := x.val.Float64()

        switch t.lit {
        case "rand":
            rand.Seed(time.Now().UnixNano())
            return Number{ val: new(big.Rat).SetFloat64(rand.Float64() * val) }
        case "prime":
            if x.inf == INF || x.inf == -INF {
                return Boolean(false)
            }

            return Boolean(new(big.Int).Quo(x.val.Num(), x.val.Denom()).ProbablyPrime(0))
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
        return t.Numbers(x.Number())
    case Null:
        return t.Numbers(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}
