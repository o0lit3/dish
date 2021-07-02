package main

func Push(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Push(x.Run(), b)
    case *Variable:
        return Push(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return x.Push(y)
        default:
            return x.Push(Hashify(b))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return x.Push(Array{ y })
        case Array:
            return x.Push(Array { y })
        default:
            return x.Push(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return Push(x, NewNumber(len(y)))
        case Array:
            return Push(x, NewNumber(len(y)))
        case String:
            return Push(x, y.Number())
        case Number:
            return Join(x, Multiply(String(" "), y))
        case Boolean:
            return Push(x, y.Number())
        case Null:
            return x
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Push(x, y.Run())
        case *Variable:
            return Push(x, y.Value())
        case Hash:
            return Push(x, NewNumber(len(y)))
        case Array:
            return Push(x, NewNumber(len(y)))
        case String:
            return Push(x, y.Number())
        case Number:
            if x.inf == INF || x.inf == -INF {
                return NewNumber(0)
            }

            if y.inf == INF || y.inf == -INF {
                return x
            }

            return NewNumber(x.Int() << uint(y.Int()))
        case Boolean:
            return Push(x, y.Number())
        case Null:
            return Push(x, NewNumber(0))
        }
    case Boolean:
        return Push(x.Number(), b)
    case Null:
        return Push(Array { }, b)
    }

    return NewNumber(0)
}

func (a Hash) Push(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, val := range b {
        out[key] = val
    }

    return out
}

func (a Array) Push(b Array) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, val)
    }

    for _, val := range b {
        out = append(out, val)
    }

    return out
}
