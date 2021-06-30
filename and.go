package main

func And(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return And(x.Run(), b)
    case *Variable:
        return And(x.Value(), b)
    case Hash:
        if len(x) == 0 {
            return x
        }

        return And(Boolean(true), b)
    case Array:
        if len(x) == 0 {
            return x
        }

        return And(Boolean(true), b)
    case String:
        if string(x) == "" {
            return x
        }

        return And(Boolean(true), b)
    case Number:
        if x.inf == 0 && x.val.Cmp(NewNumber(0).val) == 0 {
            return x
        }

        return And(Boolean(true), b)
    case Boolean:
        if !x {
            return x
        }

        switch y := b.(type) {
        case *Block:
            return And(Boolean(true), y.Run())
        case *Variable:
            return And(Boolean(true), y.Value())
        default:
            return y
        }
    case Null:
        return x
    }

    return Boolean(false)
}
