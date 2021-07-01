package main
import("fmt")

func Gt(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Gt(x.Run(), b)
    case *Variable:
        return Gt(x.Value(), b)
    case Hash:
        return Gt(NewNumber(len(x)), b)
    case Array:
        return Gt(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Gt(x, y.Run())
        case *Variable:
            return Gt(x, y.Value())
        case String:
            return Boolean(string(x) > string(y))
        case Null:
            return Boolean(string(x) > "")
        default:
            return Gt(x, String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Gt(x, y.Run())
        case *Variable:
            return Gt(x, y.Value())
        case Hash:
            return Gt(x, NewNumber(len(y)))
        case Array:
            return Gt(x, NewNumber(len(y)))
        case String:
            return Gt(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(false)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(false)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(y.val) > 0)
        case Boolean:
            return Gt(x, y.Number())
        case Null:
            return Gt(x, NewNumber(0))
        }
    case Boolean:
        return Gt(x.Number(), b)
    case Null:
        return Gt(NewNumber(0), b)
    }

    return Boolean(false)
}
