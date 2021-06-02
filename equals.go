package main
import("fmt")

func Equals(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Equals(x.Run(), b)
    case *Variable:
        return Equals(x.Value(), b)
    case Hash:
        return Equals(String(fmt.Sprintf("%v", x)), b)
    case Array:
        return Equals(String(fmt.Sprintf("%v", x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Equals(x, String(fmt.Sprintf("%v", y)))
        case Array:
            return Equals(x, String(fmt.Sprintf("%v", y)))
        case String:
            return Boolean(x == y)
        case Number:
            return Boolean(x.Number() == y)
        case Boolean:
            return Boolean(x.Number() == y.Number())
        case Null:
            return Boolean(x.Number() == Number(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Boolean(x == Number(len(y)))
        case Array:
            return Boolean(x == Number(len(y)))
        case String:
            return Boolean(x == y.Number())
        case Number:
            return Boolean(x == y)
        case Boolean:
            return Boolean(x == y.Number())
        case Null:
            return Boolean(x == Number(0))
        }
    case Boolean:
        return Equals(x.Number(), b)
    case Null:
        return Equals(Number(0), b)
    }

    return Boolean(false)
}
