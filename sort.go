package main

import(
    "sort"
    "math/big"
)

func Sort(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Sort(x.Run())
    case *Variable:
        return Sort(x.Value())
    case Hash:
        return x.Array().Sort()
    case Array:
        return x.Sort()
    case String:
        return Join(x.Array().Sort(), String(""))
    case Number:
        return x.Divisors()
    case Boolean:
        return x.Number().Divisors()
    }

    return Array { }
}

func (a Array) Sort() Array {
    sort.Slice(a, func(i, j int) bool {
        return bool(Below(a[i], a[j]))
    })

    return a
}

func (a Number) Divisors() interface{} {
    i := new(big.Int).Quo(a.val.Num(), a.val.Denom())
    j := big.NewInt(1)

    if i.Cmp(big.NewInt(0)) < 0 {
        i = i.Neg(i)
    }

    if i.Cmp(big.NewInt(2)) < 0 {
        return Array { }
    }

    out := Array { NewNumber(1) }
    sqrt := new(big.Int).Sqrt(i)

    for j.Cmp(sqrt) <= 0 {
        if j.Cmp(big.NewInt(1)) > 0 && new(big.Int).Rem(i, j).Cmp(big.NewInt(0)) == 0 {
            out = append(out, Number{ val: new(big.Rat).SetInt(j) })
            div := new(big.Int).Div(i, j)

            if j.Cmp(div) != 0 {
                out = append(out, Number{ val: new(big.Rat).SetInt(div) })
            }
        }

        j = new(big.Int).Add(j, big.NewInt(1))
    }

    return Sort(out)
}
