package main
import("fmt")

func Gte(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Gte(x.Run(), b)
    case *Variable:
        return Gte(x.Value(), b)
    case Hash:
        return Gte(NewNumber(len(x)), b)
    case Array:
        return Gte(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Gte(x, y.Run())
        case *Variable:
            return Gte(x, y.Value())
        case String:
            return Boolean(string(x) >= string(y))
        case Null:
            return Boolean(string(x) >= "")
        default:
            return Gte(x, String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Gte(x, y.Run())
        case *Variable:
            return Gte(x, y.Value())
        case Hash:
            return Gte(x, NewNumber(len(y)))
        case Array:
            return Gte(x, NewNumber(len(y)))
        case String:
            return Gte(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(true)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(false)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(y.val) >= 0)
        case Boolean:
            return Gte(x, y.Number())
        case Null:
            return Gte(x, NewNumber(0))
        }
    case Boolean:
        return Gte(x.Number(), b)
    case Null:
        return Gte(NewNumber(0), b)
    }

    return Boolean(false)
}
