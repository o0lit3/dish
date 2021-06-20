package main
import("math/big")

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
        return x.Vowel()
    case Number:
        return x.Prime()
    case Boolean:
        return Product(x.Number())
    case Null:
        return Boolean(false)
    }

    return Null { }
}

func (a Array) Product() Number {
    out := NewNumber(1)

    for _, val := range a {
        out = Multiply(out, Numberize(val)).(Number)
    }

    return out
}

func (a String) Vowel() Boolean {
    if len(a) > 0 {
        switch a[0] {
        case 'A', 'a', 'E', 'e', 'I', 'i', 'O', 'o', 'U', 'u':
            return Boolean(true)
        }
    }

    return Boolean(false)
}

func (a Number) Prime() Boolean {
    if a.inf == INF || a.inf == -INF {
        return Boolean(false)
    }

    return Boolean(new(big.Int).Quo(a.val.Num(), a.val.Denom()).ProbablyPrime(0))
}
