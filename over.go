package main

func Over(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Over(x.Run(), b)
    case *Variable:
        return Over(x.Value(), b)
    case Hash:
        return Over(Number(len(x)), b)
    case Array:
        return Over(Number(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Over(x, y.Run())
        case *Variable:
            return Over(x, y.Value())
        case String:
            return Boolean(x >= y)
        case Number:
            return Boolean(x.Number() >= y)
        case Boolean:
            return Boolean(x.Number() >= y.Number())
        case Null:
            return Boolean(x.Number() >= Number(0))
        default:
            return Over(Number(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Over(x, y.Run())
        case *Variable:
            return Over(x, y.Value())
        case Hash:
            return Boolean(x >= Number(len(y)))
        case Array:
            return Boolean(x >= Number(len(y)))
        case String:
            return Boolean(x >= y.Number())
        case Number:
            return Boolean(x >= y)
        case Boolean:
            return Boolean(x >= y.Number())
        case Null:
            return Boolean(x >= Number(0))
        }
    case Boolean:
        return Over(x.Number(), b)
    case Null:
        return Over(Number(0), b)
    }

    return Boolean(false)
}
