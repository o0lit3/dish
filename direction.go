package main
import("fmt")

func Direction(a interface{}, b interface{}) Number {
    switch x := a.(type) {
    case *Block:
        return Direction(x.Run(), b)
    case *Variable:
        return Direction(x.Value(), b)
    case Hash:
        return Direction(NewNumber(len(x)), b)
    case Array:
        return Direction(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Direction(x, y.Run())
        case *Variable:
            return Direction(x, y.Value())
        case String:
            return x.Direction(y)
        case Null:
            return x.Direction(String(""))
        default:
            return x.Direction(String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Direction(x, y.Run())
        case *Variable:
            return Direction(x, y.Value())
        case Hash:
            return Direction(x, NewNumber(len(y)))
        case Array:
            return Direction(x, NewNumber(len(y)))
        case String:
            return x.Direction(y.Number())
        case Number:
            return x.Direction(y)
        case Boolean:
            return x.Direction(y.Number())
        case Null:
            return x.Direction(NewNumber(0))
        }
    case Boolean:
        return Direction(x.Number(), b)
    case Null:
        return Direction(NewNumber(0), b)
    }

    return NewNumber(0)
}

func (a String) Direction(b String) Number {
    if (string(a) < string(b)) {
        return NewNumber(-1)
    }

    if (string(a) > string(b)) {
        return NewNumber(1)
    }

    return NewNumber(0)
}

func (a Number) Direction(b Number) Number {
    return NewNumber(a.val.Cmp(b.val))
}
