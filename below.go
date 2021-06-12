package main

func Below(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Below(x.Run(), b)
    case *Variable:
        return Below(x.blk.Value(x), b)
    case Hash:
        return Below(NewNumber(len(x)), b)
    case Array:
        return Below(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Below(x, y.Run())
        case *Variable:
            return Below(x, y.blk.Value(y))
        case String:
            return Boolean(string(x) < string(y))
        case Number:
            return Boolean(x.Number().val.Cmp(y.val) == -1)
        case Boolean:
            return Boolean(x.Number().val.Cmp(y.Number().val) == -1)
        case Null:
            return Boolean(x.Number().val.Cmp(NewNumber(0).val) == -1)
        default:
            return Below(NewNumber(len(x)), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Below(x, y.Run())
        case *Variable:
            return Below(x, y.blk.Value(y))
        case Hash:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == -1)
        case Array:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == -1)
        case String:
            return Boolean(x.val.Cmp(y.Number().val) == -1)
        case Number:
            return Boolean(x.val.Cmp(y.val) == -1)
        case Boolean:
            return Boolean(x.val.Cmp(y.Number().val) == -1)
        case Null:
            return Boolean(x.val.Cmp(NewNumber(0).val) == -1)
        }
    case Boolean:
        return Below(x.Number(), b)
    case Null:
        return Below(NewNumber(0), b)
    }

    return Boolean(false)
}
