package main

func Add(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Add(x.Run(), b)
    case *Variable:
        return Add(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return x.Add(y)
        default:
            return x.Add(Hashify(b))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return x.Add(y.Array())
        case Array:
            return x.Add(y)
        default:
            return x.Add(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return Add(x, NewNumber(len(y)))
        case Array:
            return Add(x, NewNumber(len(y)))
        case String:
            return String(string(x) + string(y))
        case Number:
            return x.Add(y)
        case Boolean:
            return x.Add(y.Number())
        case Null:
            return x.Add(NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Add(x, y.Run())
        case *Variable:
            return Add(x, y.Value())
        case Hash:
            return Add(x, NewNumber(len(y)))
        case Array:
            return Add(x, NewNumber(len(y)))
        case String:
            return Add(x, y.Number())
        case Number:
            if x.inf == INF || y.inf == INF {
                return Number{ inf: INF }
            }

            if x.inf == -INF || y.inf == -INF {
                return Number{ inf: -INF }
            }

            if (x.inf == INF && y.inf == -INF) || (x.inf == -INF && y.inf == INF) {
                return Null { }
            }

            return Number{ val: NewNumber(0).val.Add(x.val, y.val) }
        case Boolean:
            return Add(x, y.Number())
        case Null:
            return x
        }
    case Boolean:
        return Add(x.Number(), b)
    case Null:
        return Add(NewNumber(0), b)
    }

    return Null { }
}

func (a Hash) Add(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func (a Array) Add(b Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}

func (a String) Add(b Number) String {
    out := ""

    for _, c := range a {
        out += string(int(c) + b.Int())
    }

    return String(out)
}
