package main
import("fmt")

func Compare(a interface{}, b interface{}) Number {
    switch x := a.(type) {
    case *Block:
        return Compare(x.Run(), b)
    case *Variable:
        return Compare(x.Value(), b)
    case Hash:
        return Compare(NewNumber(len(x)), b)
    case Array:
        return Compare(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Compare(x, y.Run())
        case *Variable:
            return Compare(x, y.Value())
        case String:
            return x.Compare(y)
        case Null:
            return x.Compare(String(""))
        default:
            return x.Compare(String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Compare(x, y.Run())
        case *Variable:
            return Compare(x, y.Value())
        case Hash:
            return Compare(x, NewNumber(len(y)))
        case Array:
            return Compare(x, NewNumber(len(y)))
        case String:
            return x.Compare(y.Number())
        case Number:
            return x.Compare(y)
        case Boolean:
            return x.Compare(y.Number())
        case Null:
            return x.Compare(NewNumber(0))
        }
    case Boolean:
        return Compare(x.Number(), b)
    case Null:
        return Compare(NewNumber(0), b)
    }

    return NewNumber(0)
}

func (a String) Compare(b String) Number {
    if (string(a) < string(b)) {
        return NewNumber(-1)
    }

    if (string(a) > string(b)) {
        return NewNumber(1)
    }

    return NewNumber(0)
}

func (a Number) Compare(b Number) Number {
    return NewNumber(a.val.Cmp(b.val))
}
