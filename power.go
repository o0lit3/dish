package main

import(
    "sort"
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
            if len(y.args) > 0 {
                return y.Sort(x.Array())
            }

            return Power(x, y.Run())
        case *Variable:
            return Power(x, y.Value())
        default:
            return Power(x.Array(), y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Sort(x)
            }

            return Power(x, y.Run())
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
            if len(y.args) > 0 {
                return Join(y.Sort(x.Array()), String(""))
            }

            return Power(x, y.Run())
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
            return Power(x, y.Value())
        case Hash:
            return Power(x, NewNumber(len(y)))
        case Array:
            return Power(x, NewNumber(len(y)))
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

func (b *Block) Sort(a Array) Array {
    sort.Slice(a, func(i, j int) bool {
        if b, ok := b.Run(a[i], a[j]).(Boolean); ok {
            return bool(b)
        }

        return bool(Below(a[i], a[j]))
    })

    return a
}
