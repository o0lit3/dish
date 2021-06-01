package main
import("math")

func Power(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        switch y := b.(type) {
        case *Block:
            return Redo(x, y)
        default:
            return Power(x.Run(), y)
        }
    case Hash:
        return Power(x.Array(), b)
    case Array:
        switch y := b.(type) {
        case String:
            return x.Power(y.Number())
        case Number:
            return x.Power(y)
        case Boolean:
            return x.Power(y.Number())
        case Null:
            return x.Power(Number(0))
        }
    case String:
        return Power(x.Number(), b)
    case Number:
        switch y := b.(type) {
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

func Redo(a *Block, b *Block) interface{} {
    var val interface{}

    i := 0
    val = Null { }

    for !Not(a.Run()) {
        val = b.Run(a, Number(i))
        i = i + 1
    }

    return val
}

func (a Array) Power(b Number) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Power(val, b))
    }

    return out
}
