package main

func Equals(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Equals(x.Run(), b)
    case Hash:
        return Equals(Number(len(x)), b)
    case Array:
        return Equals(Number(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case String:
            return Boolean(x == y)
        case Number:
            return Boolean(x.Number() == y)
        case Boolean:
            return Boolean(x.Number() == y.Number())
        case Null:
            return Boolean(x.Number() == Number(0))
        default:
            return Equals(Number(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case Hash:
            return Boolean(x == Number(len(y)))
        case Array:
            return Boolean(x == Number(len(y)))
        case String:
            return Boolean(x == y.Number())
        case Number:
            return Boolean(x == y)
        case Boolean:
            return Boolean(x == y.Number())
        case Null:
            return Boolean(x == Number(0))
        }
    case Boolean:
        return Equals(x.Number(), b)
    case Null:
        return Equals(Number(0), b)
    }

    return Boolean(false)
}
