package main
import("fmt")

func Lt(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Lt(x.Run(), b)
    case *Variable:
        return Lt(x.blk.Value(x), b)
    case Hash:
        return Lt(NewNumber(len(x)), b)
    case Array:
        return Lt(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Lt(x, y.Run())
        case *Variable:
            return Lt(x, y.blk.Value(y))
        case String:
            return Boolean(string(x) < string(y))
        case Null:
            return Boolean(string(x) < "")
        default:
            return Lt(x, String(fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Lt(x, y.Run())
        case *Variable:
            return Lt(x, y.blk.Value(y))
        case Hash:
            return Lt(x, NewNumber(len(y)))
        case Array:
            return Lt(x, NewNumber(len(y)))
        case String:
            return Lt(x, y.Number())
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
            return Lt(x, y.Number())
        case Null:
            return Lt(x, NewNumber(0))
        }
    case Boolean:
        return Lt(x.Number(), b)
    case Null:
        return Lt(NewNumber(0), b)
    }

    return Boolean(false)
}
