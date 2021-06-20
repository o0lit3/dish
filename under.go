package main

func Under(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Under(x.Run(), b)
    case *Variable:
        return Under(x.Value(), b)
    case Hash:
        return Under(NewNumber(len(x)), b)
    case Array:
        return Under(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case *Variable:
            return Under(x, y.Value())
        case Hash:
            return Under(x.Number(), NewNumber(len(y)))
        case String:
            return Boolean(string(x) <= string(y))
        case Number:
            return Under(x.Number(), y)
        case Boolean:
            return Under(x.Number(), y.Number())
        case Null:
            return Under(x.Number(), NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Under(x, y.Run())
        case *Variable:
            return Under(x, y.Value())
        case Hash:
            return Under(x, NewNumber(len(y)))
        case Array:
            return Under(x, NewNumber(len(y)))
        case String:
            return Under(x, y.Number())
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
            return Under(x, y.Number())
        case Null:
            return Under(x, NewNumber(0))
        }
    case Boolean:
        return Under(x.Number(), b)
    case Null:
        return Under(NewNumber(0), b)
    }

    return Boolean(false)
}
