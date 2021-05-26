package main

func Under(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Under(x.Run(), b)
    case Hash:
        return Under(Number(len(x)), b)
    case Array:
        return Under(Number(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case String:
            return Boolean(x <= y)
        case Number:
            return Boolean(x.Number() <= y)
        case Boolean:
            return Boolean(x.Number() <= y.Number())
        case Null:
            return Boolean(x.Number() <= Number(0))
        default:
            return Under(Number(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case Hash:
            return Boolean(x <= Number(len(y)))
        case Array:
            return Boolean(x <= Number(len(y)))
        case String:
            return Boolean(x <= y.Number())
        case Number:
            return Boolean(x <= y)
        case Boolean:
            return Boolean(x <= y.Number())
        case Null:
            return Boolean(x <= Number(0))
        }
    case Boolean:
        return Under(x.Number(), b)
    case Null:
        return Under(Number(0), b)
    }

    return Null { }
}
