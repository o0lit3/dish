package main

import(
    "math"
    "math/big"
)

func Power(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Power(x.Run(), b)
    case *Variable:
        return Power(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return x.Array().UserSort(y)
        case *Variable:
            return Power(x, y.Value())
        default:
            return Power(x.Array(), y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return x.UserSort(y)
        case *Variable:
            return Power(x, y.Value())
        case String:
            return x.Rotate(y.Number())
        case Number:
            return x.Rotate(y)
        case Boolean:
            return x.Rotate(y.Number())
        case Null:
            return x.Rotate(NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Join(x.Array().UserSort(y), String(""))
        case *Variable:
            return Power(x, y.Value())
        default:
            return Power(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Power(x, y.Run())
        case *Variable:
            return Power(x, y.Value)
        case Hash:
            return y.Array().Rotate(x)
        case Array:
            return y.Rotate(x)
        case String:
            return x.Power(y.Number())
        case Number:
            return x.Power(y)
        case Boolean:
            return x.Power(y.Number())
        case Null:
            return NewNumber(1)
        }
    case Boolean:
        return Power(x.Number(), b)
    case Null:
        return NewNumber(0)
    }

    return NewNumber(0)
}

func (a Array) Rotate(b Number) Array {
    out := Array { }

    e := -b.Int()
    i := -b.Int()

    if i < 0 {
        e = len(a) + i
        i = len(a) + i
    }

    for i < len(a) {
        out = append(out, a[i])
        i = i + 1
    }

    i = 0

    for i < e && i < len(a) {
        out = append(out, a[i])
        i = i + 1
    }

    return out
}

func (a Number) Power(b Number) Number {
    if b.val.Cmp(NewNumber(0).val) == -1 || !b.val.IsInt() {
        x, _ := a.val.Float64()
        y, _ := b.val.Float64()

        return Number{ val: new(big.Rat).SetFloat64(math.Pow(x, y)) }
    }

    out := NewNumber(1)
    idx := NewNumber(0)

    for idx.val.Cmp(b.val) == -1 {
        out = Number{ val: out.val.Mul(out.val, a.val) }
        idx = Number{ val: idx.val.Add(idx.val, NewNumber(1).val) }
    }

    return out
}
