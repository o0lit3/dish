package main

import(
    "bufio"
    "strings"
    "math/big"
)

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
        return x.Eval()
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

func (a String) Eval() interface{} {
    reader := bufio.NewReader(strings.NewReader(string(a)))
    parser := process(reader, program.Branch(VAL))
    return parser.blk.Run()
}

func (a Number) Prime() Boolean {
    if a.inf == INF || a.inf == -INF {
        return Boolean(false)
    }

    return Boolean(new(big.Int).Quo(a.val.Num(), a.val.Denom()).ProbablyPrime(0))
}
