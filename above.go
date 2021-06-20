package main

func Above(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Above(x.Run(), b)
    case *Variable:
        return Above(x.Value(), b)
    case Hash:
        return Above(NewNumber(len(x)), b)
    case Array:
        return Above(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Above(x, y.Run())
        case *Variable:
            return Above(x, y.Value())
        case Hash:
            return Above(x.Number(), NewNumber(len(y)))
        case Array:
            return Above(x.Number(), NewNumber(len(y)))
        case String:
            return Boolean(string(x) > string(y))
        case Number:
            return Above(x.Number(), y)
        case Boolean:
            return Above(x.Number(), y.Number())
        case Null:
            return Above(x.Number(), NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Above(x, y.Run())
        case *Variable:
            return Above(x, y.Value())
        case Hash:
            return Above(x, NewNumber(len(y)))
        case Array:
            return Above(x, NewNumber(len(y)))
        case String:
            return Above(x, y.Number())
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
            return Above(x, y.Number())
        case Null:
            return Above(x, NewNumber(0))
        }
    case Boolean:
        return Above(x.Number(), b)
    case Null:
        return Above(NewNumber(0), b)
    }

    return Boolean(false)
}
