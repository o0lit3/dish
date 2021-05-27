package main

func And(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return And(x.Run(), b)
    case Hash:
        return And(Boolean(len(x) != 0), b)
    case Array:
        return And(Boolean(len(x) != 0), b)
    case String:
        return And(Boolean(x != "" && x != "0"), b)
    case Number:
        return And(Boolean(x != 0), b)
    case Boolean:
        if !x {
            return Boolean(false)
        }

        switch y := b.(type) {
        case *Block:
            return And(Boolean(true), y.Run())
        case Hash:
            return Boolean(len(y) != 0)
        case Array:
            return Boolean(len(y) != 0)
        case String:
            return Boolean(y != "" && y != "0")
        case Number:
            return Boolean(y != 0)
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
