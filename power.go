package main
import("math")

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
            return x.Rotate(Number(0))
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
            return Number(math.Pow(float64(x), float64(y.Number())))
        case Number:
            return Number(math.Pow(float64(x), float64(y)))
        case Boolean:
            return Number(math.Pow(float64(x), float64(y.Number())))
        case Null:
            return Number(math.Pow(float64(x), 0))
        }
    case Boolean:
        return Power(x.Number(), b)
    case Null:
        return Number(0)
    }

    return Number(0)
}

func (a Array) Rotate(b Number) Array {
    out := Array { }

    e := -int(b)
    i := -int(b)

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
