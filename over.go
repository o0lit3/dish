package main

func Over(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Over(x.Run(), b)
    case *Variable:
        return Over(x.Value(), b)
    case Hash:
        return Over(NewNumber(len(x)), b)
    case Array:
        return Over(NewNumber(len(x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Over(x, y.Run())
        case *Variable:
            return Over(x, y.Value())
        case Hash:
            return Over(x.Number(), NewNumber(len(y)))
        case Array:
            return Over(x.Number(), NewNumber(len(y)))
        case String:
            return Boolean(string(x) >= string(y))
        case Number:
            return Over(x.Number(), y)
        case Boolean:
            return Over(x.Number(), y.Number())
        case Null:
            return Over(x.Number(), NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Over(x, y.Run())
        case *Variable:
            return Over(x, y.Value())
        case Hash:
            return Over(x, NewNumber(len(y)))
        case Array:
            return Over(x, NewNumber(len(y)))
        case String:
            return Over(x, y.Number())
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
            return Over(x, y.Number())
        case Null:
            return Over(x, NewNumber(0))
        }
    case Boolean:
        return Over(x.Number(), b)
    case Null:
        return Over(NewNumber(0), b)
    }

    return Boolean(false)
}
