package main

func And(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return And(x.Run(), b)
    case *Variable:
        return And(x.Value(), b)
    case Hash:
        if len(x) == 0 {
            return Boolean(false)
        }

        return And(Boolean(true), b)
    case Array:
        if len(x) == 0 {
            return Boolean(false)
        }

        return And(Boolean(true), b)
    case String:
        if x == "" || x == "0" {
            return Boolean(false)
        }

        return And(Boolean(true), b)
    case Number:
        if x == 0 {
            return Boolean(false)
        }

        return And(Boolean(true), b)
    case Boolean:
        if !x {
            return Boolean(false)
        }

        switch y := b.(type) {
        case *Block:
            return And(Boolean(true), y.Run())
        case *Variable:
            return And(Boolean(true), y.Value())
        case Hash:
            if len(y) == 0 {
                return Boolean(false)
            }

            return y
        case Array:
            if len(y) == 0 {
                return Boolean(false)
            }

            return y
        case String:
            if y == "" || y == "0" {
                return Boolean(false)
            }

            return y
        case Number:
            if y == 0 {
                return Boolean(false)
            }

            return y
        case Boolean:
            return y
        case Null:
            return Boolean(false)
        }
    case Null:
        return Boolean(false)
    }

    return Boolean(false)
}
