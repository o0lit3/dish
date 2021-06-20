package main
import("math/big")

func Remainder(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Remainder(x.Run(), b)
    case *Variable:
        return Remainder(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return x.Filter(y)
            }

            return Remainder(x, y.Run())
        case *Variable:
            return Remainder(x, y.Value())
        default:
            return Remainder(x.Array(), b)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return x.Filter(y)
            }

            return Remainder(x, y.Run())
        case *Variable:
            return Remainder(x, y.Value())
        case String:
            return x.Remainder(y.Number())
        case Number:
            return x.Remainder(y)
        case Boolean:
            return x.Remainder(y.Number())
        case Null:
            return x.Remainder(NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return Join(x.Array().Filter(y), String(""))
            }

            return Remainder(x, y.Run())
        case *Variable:
            return Remainder(x, y.Value())
        case String:
            return x.Remainder(y.Number())
        case Number:
            return x.Remainder(y)
        case Boolean:
            return x.Remainder(y.Number())
        case Null:
            return x.Remainder(NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Remainder(x, y.Run())
        case *Variable:
            return Remainder(x, y.Value())
        case Hash:
            return Remainder(x, NewNumber(len(y)))
        case Array:
            return Remainder(x, NewNumber(len(y)))
        case String:
            return Remainder(x, y.Number())
        case Number:
            if x.inf == INF || x.inf == -INF {
                return Null { }
            }

            if y.inf == INF || y.inf == -INF {
                return x
            }

            if y.val.Cmp(NewNumber(0).val) != 0 {
                i := new(big.Int).Quo(x.val.Num(), x.val.Denom())
                j := new(big.Int).Quo(y.val.Num(), y.val.Denom())
                return Number{ val: new(big.Rat).SetInt(new(big.Int).Rem(i, j)) }
            }
        case Boolean:
            return Remainder(x, y.Number())
        }
    case Boolean:
        return Remainder(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Filter(b *Block) Hash {
    out := Hash { }

    for key, val := range a {
        if Not(b.Run(val, String(key))) {
            continue
        }

        out[key] = val
    }

    return out
}

func (a Array) Filter(b *Block) Array {
    out := Array { }

    for i, val := range a {
        if Not(b.Run(val, NewNumber(i))) {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (a Array) Remainder(b Number) interface{} {
    if b.val.Cmp(NewNumber(0).val) == 0 {
        return Null { }
    }

    out := Array { }

    for i, val := range a{
        if i % b.Int() > 0 {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (a String) Remainder(b Number) interface{} {
    if b.val.Cmp(NewNumber(0).val) == 0 {
        return Null { }
    }

    out := ""

    for i, c := range a {
        if i % b.Int() > 0 {
            continue
        }

        out += string(c)
    }

    return String(out)
}
