package main
import("fmt")

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
        default:
            return Below(x, String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Below(x, y.Run())
        case *Variable:
            return Below(x, y.blk.Value(y))
        case Hash:
            return Below(x, NewNumber(len(y)))
        case Array:
            return Below(x, NewNumber(len(y)))
        case String:
            return Below(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(false)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(true)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(false)
            }

            return Boolean(x.val.Cmp(y.val) < 0)
        case Boolean:
            return Below(x, y.Number())
        case Null:
            return Below(x, NewNumber(0))
        }
    case Boolean:
        return Below(x.Number(), b)
    case Null:
        return Below(NewNumber(0), b)
    }

    return Boolean(false)
}
