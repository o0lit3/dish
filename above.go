package main

func Above(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Above(x.Run(), b)
    case *Variable:
        return Above(x.Value(), b)
    case Hash:
        return Above(NewNumber(len(x)), b)
    case Array:
        return Above(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Above(x, y.Run())
        case *Variable:
            return Above(x, y.Value())
        case String:
            return Boolean(x > y)
        case Number:
            return Boolean(x.Number().val.Cmp(y.val) == 1)
        case Boolean:
            return Boolean(x.Number().val.Cmp(y.Number().val) == 1)
        case Null:
            return Boolean(x.Number().val.Cmp(NewNumber(0).val) == 1)
        default:
            return Above(NewNumber(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Above(x, y.Run())
        case *Variable:
            return Above(x, y.Value())
        case Hash:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == 1)
        case Array:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == 1)
        case String:
            return Boolean(x.val.Cmp(y.Number().val) == 1)
        case Number:
            return Boolean(x.val.Cmp(y.val) == 1)
        case Boolean:
            return Boolean(x.val.Cmp(y.Number().val) == 1)
        case Null:
            return Boolean(x.val.Cmp(NewNumber(0).val) == 1)
        }
    case Boolean:
        return Above(x.Number(), b)
    case Null:
        return Above(NewNumber(0), b)
    }

    return Boolean(false)
}
