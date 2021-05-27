package main

func Or(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Or(x.Run(), b)
    case Hash:
        return Or(Boolean(len(x) != 0), b)
    case Array:
        return Or(Boolean(len(x) != 0), b)
    case String:
        return Or(Boolean(x != "" && x != "0"), b)
    case Number:
        return Or(Boolean(x != 0), b)
    case Boolean:
        if x {
            return Boolean(true)
        }

        switch y := b.(type) {
        case *Block:
            return Or(Boolean(false), y.Run())
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
