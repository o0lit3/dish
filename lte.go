package main
import("fmt")

func Lte(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Lte(x.Run(), b)
    case *Variable:
        return Lte(x.Value(), b)
    case Hash:
        return Lte(NewNumber(len(x)), b)
    case Array:
        return Lte(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Lte(x, y.Run())
        case String:
            return Boolean(string(x) <= string(y))
        case Null:
            return Boolean(string(x) <= "")
        default:
            return Lte(x, String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Lte(x, y.Run())
        case *Variable:
            return Lte(x, y.Value())
        case Hash:
            return Lte(x, NewNumber(len(y)))
        case Array:
            return Lte(x, NewNumber(len(y)))
        case String:
            return Lte(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(true)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(true)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(false)
            }

            return Boolean(x.val.Cmp(y.val) <= 0)
        case Boolean:
            return Lte(x, y.Number())
        case Null:
            return Lte(x, NewNumber(0))
        }
    case Boolean:
        return Lte(x.Number(), b)
    case Null:
        return Lte(NewNumber(0), b)
    }

    return Boolean(false)
}
