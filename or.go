package main

func Or(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Or(x.Run(), b)
    case *Variable:
        return Or(x.Value(), b)
    case Hash:
        if len(x) != 0 {
            return x
        }

        return Or(Boolean(false), b)
    case Array:
        if len(x) != 0 {
            return x
        }

        return Or(Boolean(false), b)
    case String:
        if string(x) != "" && string(x) != "0" {
            return x
        }

        return Or(Boolean(false), b)
    case Number:
        if x.inf != 0 || x.val.Cmp(NewNumber(0).val) != 0 {
            return x
        }

        return Or(Boolean(false), b)
    case Boolean:
        if x {
            return Boolean(true)
        }

        switch y := b.(type) {
        case *Block:
            return Or(Boolean(false), y.Run())
        case *Variable:
            return Or(Boolean(false), y.Value())
        case Hash:
            if len(y) != 0 {
                return y
            }

            return Boolean(false)
        case Array:
            if len(y) != 0 {
                return y
            }

            return Boolean(false)
        case String:
            if string(y) != "" && string(y) != "0" {
                return y
            }

            return Boolean(false)
        case Number:
            if y.inf != 0 || y.val.Cmp(NewNumber(0).val) != 0 {
                return y
            }

            return Boolean(false)
        case Boolean:
            return y
        case Null:
            return Boolean(false)
        }
    case Null:
        return Or(Boolean(false), b)
    }

    return Boolean(false)
}
