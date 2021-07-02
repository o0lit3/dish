package main

import(
    "fmt"
    "strings"
)

func Subtract(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Subtract(x.Run(), b)
    case *Variable:
        return Subtract(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return x.Subtract(y)
        default:
            return x.Subtract(Hashify(b))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return x.Subtract(y.Array())
        case Array:
            return x.Subtract(y)
        default:
            return x.Subtract(Array { y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return Subtract(x, NewNumber(len(y)))
        case Array:
            return Subtract(x, NewNumber(len(y)))
        case String:
            return x.Remove(y)
        case Number:
            return x.Subtract(y)
        case Boolean:
            return x.Subtract(y.Number())
        case Null:
            return x.Subtract(NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Subtract(x, y.Run())
        case *Variable:
            return Subtract(x, y.Value())
        case Hash:
            return Subtract(x, NewNumber(len(y)))
        case Array:
            return Subtract(x, NewNumber(len(y)))
        case String:
            return Subtract(x, y.Number())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Null { }
            }

            if x.inf == INF && y.inf == -INF {
                return Number{ inf: INF }
            }

            if x.inf == -INF && y.inf == INF {
                return Number{ inf: -INF }
            }

            return Number{ val: NewNumber(0).val.Sub(x.val, y.val) }
        case Boolean:
            return Subtract(x, y.Number())
        case Null:
            return x
        }
    case Boolean:
        return Subtract(x.Number(), b)
    case Null:
        return Subtract(NewNumber(0), b)
    }

    return Null { }
}

func (a Hash) Subtract(b Hash) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = val
    }

    for key, _ := range b {
        if _, ok := out[key]; ok {
            delete(out, key)
        }
    }

    return out
}

func (a Array) Subtract(b Array) Array {
    out := Array { }
    hash := Hash { }

    for _, val := range b {
        hash[fmt.Sprintf("%v", val)] = val
    }

    for _, val := range a {
        key := fmt.Sprintf("%v", val)

        if _, ok := hash[key]; ok {
            delete(hash, key)
        } else {
            out = append(out, val)
        }
    }

    return out
}

func (a String) Subtract(b Number) String {
    out := ""

    for _, c := range a {
        out += string(int(c) - b.Int())
    }

    return String(out)
}

func (a String) Remove(b String) String {
    return String(strings.Replace(string(a), string(b), "", 1))
}
