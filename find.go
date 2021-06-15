package main
import("math/big")

func Find(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Find(x.Run(), b)
    case *Variable:
        return Find(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Find(x)
            }

            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        default:
            return x.Find(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Find(x)
            }

            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        default:
            return x.Find(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Find(x.Array())
            }

            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        case String:
            return x.Find(y)
        default:
            return x.Array().Find(y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        case String:
            return x.Round(y.Number())
        case Number:
            return x.Round(y)
        case Boolean:
            return x.Round(y.Number())
        case Null:
            return x.Round(NewNumber(0))
        }
    case Boolean:
       return Find(x.Number(), b)
    }

    return NewNumber(0)
}

func (a Hash) Find(b interface{}) Array {
    out := Array { }

    for key, val := range a {
        if Equals(val, b) {
            out = append(out, String(key))
        }
    }

    return out
}

func (a Array) Find(b interface{}) Array {
    out := Array { }

    for i, val := range a {
        if Equals(val, b) {
            out = append(out, NewNumber(i))
        }
    }

    return out
}

func (a String) Find(b String) Array {
    out := Array { }

    for i := range a {
        if i + len(b) <= len(a) && string(a[i:i + len(b)]) == string(b) {
            out = append(out, i)
        }
    }

    return out
}

func (a Number) Round(b Number) Number {
    if pow, ok := Power(NewNumber(10), b).(Number); ok {
        o := new(big.Rat).Mul(a.val, pow.val)
        i := new(big.Int).Quo(o.Num(), o.Denom())

        o = o.SetInt(i)
        return Number{ val: o.Quo(o, pow.val) }
    }

    return NewNumber(0)
}

func (b *Block) Find(a interface{}) Array {
    out := Array { }

    switch x := a.(type) {
    case Hash:
        for key, val := range x {
            if b, ok := b.Run(val).(Boolean); ok && bool(b) {
                out = append(out, key)
            }
        }
    case Array:
        for i, val := range x {
            if b, ok := b.Run(val).(Boolean); ok && bool(b) {
                out = append(out, i)
            }
        }
    }

    return out
}
