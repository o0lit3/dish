package main

func Below(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Below(x.Run(), b)
    case *Variable:
        return Below(x.blk.Value(x), b)
    case Hash:
        return Below(Number(len(x)), b)
    case Array:
        return Below(Number(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Below(x, y.Run())
        case *Variable:
            return Below(x, y.blk.Value(y))
        case String:
            return Boolean(x < y)
        case Number:
            return Boolean(x.Number() < y)
        case Boolean:
            return Boolean(x.Number() < y.Number())
        case Null:
            return Boolean(x.Number() < Number(0))
        default:
            return Below(Number(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Below(x, y.Run())
        case *Variable:
            return Below(x, y.blk.Value(y))
        case Hash:
            return Boolean(x < Number(len(y)))
        case Array:
            return Boolean(x < Number(len(y)))
        case String:
            return Boolean(x < y.Number())
        case Number:
            return Boolean(x < y)
        case Boolean:
            return Boolean(x < y.Number())
        case Null:
            return Boolean(x < Number(0))
        }
    case Boolean:
        return Below(x.Number(), b)
    case Null:
        return Below(Number(0), b)
    }

    return Boolean(false)
}