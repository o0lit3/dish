package main

func Under(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Under(x.Run(), b)
    case *Variable:
        return Under(x.Value(), b)
    case Hash:
        return Under(NewNumber(len(x)), b)
    case Array:
        return Under(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case *Variable:
            return Under(x, y.Value())
        case String:
            return Boolean(x <= y)
        case Number:
            return Boolean(x.Number().val.Cmp(y.val) <= 0)
        case Boolean:
            return Boolean(x.Number().val.Cmp(y.Number().val) <= 0)
        case Null:
            return Boolean(x.Number().val.Cmp(NewNumber(0).val) <= 0)
        default:
            return Under(NewNumber(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case *Variable:
            return Under(x, y.Value())
        case Hash:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) <= 0)
        case Array:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) <= 0)
        case String:
            return Boolean(x.val.Cmp(y.Number().val) <= 0)
        case Number:
            return Boolean(x.val.Cmp(y.val) <= 0)
        case Boolean:
            return Boolean(x.val.Cmp(y.Number().val) <= 0)
        case Null:
            return Boolean(x.val.Cmp(NewNumber(0).val) <= 0)
        }
    case Boolean:
        return Under(x.Number(), b)
    case Null:
        return Under(NewNumber(0), b)
    }

    return Boolean(false)
}
