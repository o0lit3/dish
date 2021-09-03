package main
import("sort"; "math"; "math/big")

func (t *Token) Pow(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Pow(x.Run(), b)
    case *Variable:
        return t.Pow(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.Pow(x, y.Run())
        case *Variable:
            return t.Pow(x, y.Value())
        case Hash:
            if t.lit != "^" && t.lit != "^=" && t.lit != "zip" {
                t.TypeMismatch(x, y)
            }

            return t.ZipHash(x, y)
        case Null:
            return t.Pow(x, Hash{ })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "^" && t.lit != "^=" && t.lit != "sort" {
                    t.TypeMismatch(x, y)
                }

                return t.SortArray(x ,y)
            }

            return t.Pow(x, y.Run())
        case *Variable:
            return t.Pow(x, y.Value())
        case Array:
            if t.lit != "^" && t.lit != "^=" && t.lit != "zip" {
                t.TypeMismatch(x, y)
            }

            return t.ZipArray(x, y)
        case Number:
            if t.lit != "^" && t.lit != "^=" && t.lit != "rotate" && t.lit != "rot" {
                t.TypeMismatch(x, y)
            }

            return t.RotateArray(x, y)
        case Boolean:
            return t.Pow(x, y.Number())
        case Null:
            switch t.lit {
            case "rotate", "rot":
                return t.Pow(x, NewNumber(0))
            default:
                return t.Pow(x, Array{ })
            }
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "^" && t.lit != "^=" && t.lit != "sort" {
                    t.TypeMismatch(x, y)
                }

                return t.SortString(x, y)
            }

            return t.Pow(x, y.Run())
        case *Variable:
            return t.Pow(x, y.Value())
        case String:
            if t.lit != "^" && t.lit != "^=" && t.lit != "zip" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(t.FlattenArray(t.ZipArray(x.Array(), y.Array())), String(""))
        case Number:
            if t.lit != "^" && t.lit != "^=" && t.lit != "rotate" && t.lit != "rot" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(t.RotateArray(x.Array(), y), String(""))
        case Boolean:
            return t.Pow(x, y.Number())
        case Null:
            switch t.lit {
            case "rotate", "rot":
                return t.Pow(x, NewNumber(0))
            default:
                return t.Pow(x, String(""))
            }
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Pow(x, y.Run())
        case *Variable:
            return t.Pow(x, y.Value())
        case Number:
            if t.lit != "^" && t.lit != "^=" && t.lit != "power" && t.lit != "pow" {
                t.TypeMismatch(x, y)
            }

            return t.PowerNumber(x, y)
        case Boolean:
            return t.Pow(x, y.Number())
        case Null:
            return t.Pow(x, NewNumber(0))
        }
    case Boolean:
        return t.Pow(x.Number(), b)
    case Null:
        switch b.(type) {
        case Hash:
            return t.Pow(Hash{ }, b)
        case Array:
            return t.Pow(Array{ }, b)
        case String:
            return t.Pow(String(""), b)
        default:
            return t.Pow(NewNumber(0), b)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopHat(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopHat(x.Run())
    case *Variable:
        return t.TopHat(x.Value())
    case Array:
        if t.lit != "^" && t.lit != "sort" {
            t.TypeMismatch(x, nil)
        }

        return t.SortArray(x, nil)
    case String:
        if t.lit != "^" && t.lit != "sort" {
            t.TypeMismatch(x, nil)
        }

        return t.SortString(x, nil)
    case Number:
        if t.lit != "^" && t.lit != "squared" {
            t.TypeMismatch(x, nil)
        }

        return t.PowerNumber(x, NewNumber(2))
    case Boolean:
        return t.TopHat(x.Number())
    case Null:
        return t.TopHat(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) ZipHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key, val := range x {
        if _, ok := val.(Array); ok {
            out[key] = val
        } else {
            out[key] = Array{ val }
        }
    }

    for key, val := range y {
        if _, ok := out[key]; ok {
            out[key] = append(out[key].(Array), val)
        } else {
            if _, ok := val.(Array); ok {
                out[key] = val
            } else {
                out[key] = Array{ val }
            }

        }
    }

    return out
}

func (t *Token) ZipArray(x Array, y Array) Array {
    out := Array { }

    for _, val := range x {
        if _, ok := val.(Array); ok {
            out = append(out, val)
        } else {
            out = append(out, Array{ val })
        }
    }

    for i, val := range y {
        out[i] = append(out[i].(Array), val)
    }

    return out
}

func (t *Token) SortArray(x Array, y *Block) Array {
    sort.Slice(x, func(i, j int) bool {
        if y == nil {
            if b, ok := t.Wiki(x[i], x[j]).(Boolean); ok {
                return bool(b)
            }

            return false
        }

        if b, ok := y.Context(x).Run(x[i], x[j]).(Boolean); ok {
            return bool(b)
        }

        if b, ok := t.Wiki(x[i], x[j]).(Boolean); ok {
            return bool(b)
        }

        return false
    })

    return x
}

func (t *Token) SortString(x String, y *Block) String {
    return t.JoinArray(t.SortArray(x.Array(), y), String(""))
}

func (t *Token) RotateArray(x Array, y Number) Array {
    out := Array { }

    e := -y.Int()
    i := -y.Int()

    if i < 0 {
        e = len(x) + i
        i = len(x) + i
    }

    for i < len(x) {
        out = append(out, x[i])
        i = i + 1
    }

    i = 0

    for i < e && i < len(x) {
        out = append(out, x[i])
        i = i + 1
    }

    return out
}

func (t *Token) PowerNumber(x Number, y Number) interface{} {
    switch y.inf {
    case INF:
        return Number{ inf: INF }
    case -INF:
        return NewNumber(0)
    }

    if x.inf == INF || x.inf == -INF {
        switch y.val.Cmp(NewNumber(0).val) {
        case -1:
            return NewNumber(0)
        case 0:
            return NewNumber(1)
        case 1:
            return Number{ inf: x.inf }
        }

        return Null { }
    }

    if y.val.Cmp(NewNumber(0).val) == -1 || !y.val.IsInt() {
        x, _ := x.val.Float64()
        y, _ := y.val.Float64()

        return Number{ val: new(big.Rat).SetFloat64(math.Pow(x, y)) }
    }

    out := NewNumber(1)
    idx := NewNumber(0)

    for idx.val.Cmp(y.val) == -1 {
        out = Number{ val: out.val.Mul(out.val, x.val) }
        idx = Number{ val: idx.val.Add(idx.val, NewNumber(1).val) }
    }

    return out
}

