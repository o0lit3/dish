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
        if i < len(out) {
            out[i] = append(out[i].(Array), val)
        } else {
            out = append(out, Array{ val })
        }
    }

    return out
}

func (t *Token) SortArray(x Array, y *Block) Array {
    named := 0

    if y != nil {
        named = y.Named()
    }

    sort.Slice(x, func(i, j int) bool {
        if y == nil {
            if b, ok := t.Wiki(x[i], x[j]).(Boolean); ok {
                return bool(b)
            }

            return false
        }

        if named == 1 {
            if b, ok := t.Wiki(y.Context(x).Run(x[i]), y.Context(x).Run(x[j])).(Boolean); ok {
                return bool(b)
            }

            return false
        }

        val := y.Context(x).Run(x[i], x[j])

        if b, ok := val.(Boolean); ok {
            return bool(b)
        }

        if named == 0 {
            if b, ok := t.Wiki(val, y.Context(x).Run(x[j], x[i])).(Boolean); ok {
                return bool(b)
            }
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

    if len(x) == 0 {
        return out
    }

    e := -(y.Int() % len(x))
    i := -(y.Int() % len(x))

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
        switch y.Cmp(NewNumber(0)) {
        case -1:
            return NewNumber(0)
        case 0:
            return NewNumber(1)
        case 1:
            return Number{ inf: x.inf }
        }

        return Null { }
    }

    if !y.Rat().IsInt() {
        x, _ := x.Rat().Float64()
        y, _ := y.Rat().Float64()
        val := math.Pow(x, y)

        if math.IsInf(val, 1) {
            return Number{ inf: INF }
        }

        if math.IsInf(val, -1) {
            return Number{ inf: -INF }
        }

        if math.IsNaN(val) {
            return Null{ }
        }

        return NewRat(new(big.Rat).SetFloat64(val))
    }

    neg := y.Cmp(NewNumber(0)) == -1
    mag := NewRat(new(big.Rat).Abs(y.Rat()))
    exp := int64(0)

    if num := mag.Rat().Num(); num.IsInt64() {
        exp = num.Int64()
    } else {
        out := NewNumber(1)
        idx := NewNumber(0)

        for idx.Cmp(mag) == -1 {
            out = NewRat(out.Rat().Mul(out.Rat(), x.Rat()))
            idx = NewRat(idx.Rat().Add(idx.Rat(), NewNumber(1).Rat()))
        }

        return Inverted(out, neg)
    }

    if x.Fits() {
        val := int64(1)
        fits := true

        for i := int64(0); i < exp; i++ {
            if x.num == 0 {
                val = 0
                break
            }

            prod := val * x.num

            if prod / x.num != val {
                fits = false
                break
            }

            val = prod
        }

        if fits {
            return Inverted(Number{ num: val }, neg)
        }
    }

    out := NewNumber(1)

    for i := int64(0); i < exp; i++ {
        out = NewRat(new(big.Rat).Mul(out.Rat(), x.Rat()))
    }

    return Inverted(out, neg)
}

func Inverted(n Number, neg bool) interface{} {
    if !neg {
        return n
    }

    if n.Cmp(NewNumber(0)) == 0 {
        return Number{ inf: INF }
    }

    return NewRat(new(big.Rat).Inv(n.Rat()))
}

