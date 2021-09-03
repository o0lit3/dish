package main
import("math/big")

func (t *Token) Thump(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Thump(x.Run(), b)
    case *Variable:
        return t.Thump(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                    t.TypeMismatch(x, y)
                }

                return t.FindInHash(x, y)
            }

            return t.Thump(x, y.Run())
        case *Variable:
            return t.Thump(x, y.Value())
        default:
            if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                t.TypeMismatch(x, y)
            }

            return t.SearchHash(x, y)
        }

    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                    t.TypeMismatch(x, y)
                }

                return t.FindInArray(x, y)
            }

            return t.Thump(x, y.Run())
        case *Variable:
            return t.Thump(x, y.Value())
        default:
            if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                t.TypeMismatch(x, y)
            }

            return t.SearchArray(x, y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                    t.TypeMismatch(x, y)
                }

                return t.FindInString(x, y)
            }

            return t.Thump(x, y.Run())
        case *Variable:
            return t.Thump(x, y.Value())
        case String:
            if t.lit != "@" && t.lit != "find" && t.lit != "search" && t.lit != "indices" {
                t.TypeMismatch(x, y)
            }

            return t.SearchString(x, y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Thump(x, y.Run())
        case *Variable:
            return t.Thump(x, y.Value())
        case Number:
            if t.lit != "@" && t.lit != "round" {
                t.TypeMismatch(x, y)
            }

            return t.RoundNumber(x, y)
        case Boolean:
            return t.Thump(x, y.Number())
        case Null:
            return t.Thump(x, NewNumber(0))
        }
    case Boolean:
       return t.Thump(x.Number(), b)
    case Null:
        switch b.(type) {
        case Number:
            return t.Thump(NewNumber(0), b)
        default:
            return t.Thump(Array{ }, b)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopThump(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopThump(x.Run())
    case *Variable:
        return t.TopThump(x.Value())
    case Hash:
        if t.lit != "@" && t.lit != "keys" {
            t.TypeMismatch(a, nil)
        }

        return t.HashKeys(x)
    case Array:
        if t.lit != "@" && t.lit != "reverse" {
            t.TypeMismatch(a, nil)
        }

        return t.ReverseArray(x)
    case String:
        if t.lit != "@" && t.lit != "reverse" {
            t.TypeMismatch(a, nil)
        }

        return t.ReverseString(x)
    case Number:
        if t.lit != "@" && t.lit != "round" {
            t.TypeMismatch(a, nil)
        }

        return t.RoundNumber(x, NewNumber(0))
    case Boolean:
        return t.TopThump(x.Number())
    case Null:
        return t.TopThump(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) FindInHash(x Hash, y *Block) Array {
    out := Array { }

    for key := range x {
        if Boolify(y.Run(x[key], String(key))) {
            out = append(out, String(key))
        }
    }

    return out
}

func (t *Token) FindInArray(x Array, y *Block) Array {
    out := Array { }

    for i := range x {
        if Boolify(y.Run(x[i], NewNumber(i))) {
            out = append(out, NewNumber(i))
        }
    }

    return out
}

func (t *Token) FindInString(x String, y *Block) Array {
    out := Array { }

    for i := range x {
        if Boolify(y.Run(String(string(x[i])), NewNumber(i))) {
            out = append(out, NewNumber(i))
        }
    }

    return out
}

func (t *Token) SearchHash(x Hash, y interface{}) Array {
    out := Array { }

    for key := range x {
        if Equals(x[key], y) {
            out = append(out, String(key))
        }
    }

    return out
}

func (t *Token) SearchArray(x Array, y interface{}) Array {
    out := Array { }

    for i := range x {
        if Equals(x[i], y) {
            out = append(out, NewNumber(i))
        }
    }

    return out
}

func (t *Token) SearchString(x String, y String) Array {
    out := Array { }

    for i := range x {
        if i + len(y) <= len(x) && string(x[i:i + len(y)]) == string(y) {
            out = append(out, NewNumber(i))
        }
    }

    return out
}

func (t *Token) RoundNumber(x Number, y Number) Number {
    if y.inf == INF || y.inf == -INF {
        return x
    }

    if x.inf == INF || x.inf == -INF {
        return x
    }

    if pow, ok := t.PowerNumber(NewNumber(10), y).(Number); ok {
        o := new(big.Rat).Mul(x.val, pow.val)
        i := new(big.Int).Quo(o.Num(), o.Denom())
        j := new(big.Rat).SetInt(i)
        d := new(big.Rat).Sub(o, j)

        if d.Cmp(big.NewRat(1, 2)) >= 0 {
            j = j.Add(j, big.NewRat(1, 1))
        }

        if d.Cmp(big.NewRat(-1, 2)) <= 0 {
            j = j.Sub(j, big.NewRat(1, 1))
        }

        return Number{ val: j.Quo(j, pow.val) }
    }

    return NewNumber(0)
}

func (t *Token) HashKeys(x Hash) Array {
    out := Array { }

    for key := range x {
        out = append(out, String(key))
    }

    return out
}

func (t *Token) ReverseArray(x Array) Array {
    out := Array { }

    for i := range x {
        out = append(out, x[len(x) - i - 1])
    }

    return out
}

func (t *Token) ReverseString(x String) String {
    out := ""

    for i := range x {
        out += string(x[len(x) - i - 1])
    }

    return String(out)
}
