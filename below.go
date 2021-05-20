package main

func Below(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return Below(Number(len(x)), b)
    case Array:
        return Below(Number(len(x)), b)
    case String:
        switch y := b.(type) {
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
        case Interpreter:
            return Below(x, y.Run())
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

    return Null { }
}
