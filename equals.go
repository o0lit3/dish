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
            return Boolean(string(x) == string(y))
        case Number:
            return Equals(x.Number(), y)
        case Boolean:
            return Equals(x.Number(), y.Number())
        case Null:
            return Equals(x.Number(), NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Equals(x, NewNumber(len(y)))
        case Array:
            return Equals(x, NewNumber(len(y)))
        case String:
            return Equals(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(true)
            }

            if x.inf == INF || x.inf == -INF || y.inf == INF || y.inf == -INF {
                return Boolean(false)
            }

            return Boolean(x.val.Cmp(y.val) == 0)
        case Boolean:
            return Equals(x, y.Number())
        case Null:
            return Equals(x, NewNumber(0))
        }
    case Boolean:
        return Equals(x.Number(), b)
    case Null:
        return Equals(NewNumber(0), b)
    }

    return Boolean(false)
}
