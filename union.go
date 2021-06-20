package main
import("fmt")

func Union(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Union(x.Run(), b)
    case *Variable:
        return Union(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Union(x, y.Run())
        case *Variable:
            return Union(x, y.Value())
        case Hash:
            return x.Union(y)
        case Array:
            return x.Union(y.Hash())
        case String:
            return x.Union(Hash { string(y): y })
        default:
            return x.Union(Hash { fmt.Sprintf("%v", y): y })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Union(x, y.Run())
        case *Variable:
            return Union(x, y.Value())
        case Hash:
            return x.Union(y.Array())
        case Array:
            return x.Union(y)
        default:
            return x.Union(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Union(x, y.Run())
        case *Variable:
            return Union(x, y.Value())
        case String:
            return x.Union(y)
        default:
            return Union(x.Number(), y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Union(x, y.Run())
        case *Variable:
            return Union(x, y.Value())
        case Hash:
            return Union(x, NewNumber(len(y)))
        case Array:
            return Union(x, NewNumber(len(y)))
        case String:
            return Union(x, y.Number())
        case Number:
            if (x.inf == INF || x.inf == -INF) && (y.inf == INF || y.inf == -INF) {
                return NewNumber(0)
            }

            if x.inf == INF || x.inf == -INF {
                return y
            }

            if y.inf == INF || y.inf == -INF {
                return x
            }

            return NewNumber(x.Int() | y.Int())
        case Boolean:
            return Union(x, y.Number())
        case Null:
            return Union(x, NewNumber(0))
        }
    case Boolean:
        return Union(x.Number(), b)
    case Null:
        return Union(NewNumber(0), b)
    }

    return Null { }
}

func (a Hash) Union(b Hash) interface{} {
    out := Hash { }

    for key := range a {
        out[key] = a[key]
    }

    for key := range b {
        out[key] = b[key]
    }

    return out
}

func (a Array) Union(b Array) interface{} {
    out := Array { }

    for i := range a {
        out = append(out, a[i])
    }

    for i := range b {
        out = append(out, b[i])
    }

    return out.Unique()
}

func (a String) Union(b String) interface{} {
    out := ""

    if len(b) > len(a) {
        for i := range b {
            if i < len(a) {
                out += string(a[i] | b[i])
            } else {
                out += string(b[i])
            }
        }
    } else {
        for i := range a {
            if i < len(b) {
                out += string(a[i] | b[i])
            } else {
                out += string(a[i])
            }
        }

    }

    return String(out)
}
