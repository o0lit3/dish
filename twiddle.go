package main

func (t *Token) Twiddle(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Twiddle(x.Run(), b)
    case *Variable:
        return t.Twiddle(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "~" && t.lit != "~=" && t.lit != "none" {
                    t.TypeMismatch(x, y)
                }

                for key, val := range(x) {
                    if Boolify(y.Run(val, String(key))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Twiddle(x, y.Run())
        case *Variable:
            return t.Twiddle(x, y.Value())
        case Hash:
            if t.lit != "~" && t.lit != "~=" && t.lit != "exclusion" {
                t.TypeMismatch(x, y)
            }

            return t.ExclusionHash(x, y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "~" && t.lit != "~=" && t.lit != "none" {
                    t.TypeMismatch(x, y)
                }

                for i, val := range(x) {
                    if Boolify(y.Run(val, NewNumber(i))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Twiddle(x, y.Run())
        case *Variable:
            return t.Twiddle(x, y.Value())
        case Array:
            if t.lit != "~" && t.lit != "~=" && t.lit != "exclusion" {
                t.TypeMismatch(x, y)
            }

            return t.ExclusionArray(x, y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "~" && t.lit != "~=" && t.lit != "none" {
                    t.TypeMismatch(x, y)
                }

                for i, val := range(x) {
                    if Boolify(y.Run(String(string(val)), NewNumber(i))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Twiddle(x, y.Run())
        case *Variable:
            return t.Twiddle(x, y.Value())
        case String:
            if t.lit != "~" && t.lit != "~=" && t.lit != "exclusion" {
                t.TypeMismatch(x, y)
            }

            return t.ExclusionString(x, y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Twiddle(x, y.Run())
        case *Variable:
            return t.Twiddle(x, y.Value())
        case Number:
            if t.lit != "~" && t.lit != "~=" && t.lit != "bxor" {
                t.TypeMismatch(x, y)
            }

            return t.BxorNumber(x, y)
        case Boolean:
            return t.Twiddle(x, y.Number())
        case Null:
            return t.Twiddle(x, NewNumber(0))
        }
    case Boolean:
        return t.Twiddle(x.Number(), b)
    case Null:
        return t.Twiddle(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TwiddleDee(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TwiddleDee(x.Run(), b)
    case *Variable:
        return t.TwiddleDee(x.Value(), b)
    default:
        if Boolify(x) {
            if !Boolify(b) {
                return x
            }

            return Null { }
        }

        switch y := b.(type) {
        case *Block:
            return y.Run()
        case *Variable:
            return y.Value()
        default:
            return y
        }
    }
}

func (t *Token) TopTwiddle(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopTwiddle(x.Run())
    case *Variable:
        return t.TopTwiddle(x.Value())
    case Hash:
        if t.lit != "~" && t.lit != "invert" && t.lit != "flip" {
            t.TypeMismatch(x, nil)
        }

        return t.InvertHash(x)
    case Array:
        if t.lit != "~" && t.lit != "invert" && t.lit != "flip" {
            t.TypeMismatch(x, nil)
        }

        return t.TransposeArray(x)
    case String:
        if t.lit != "~" && t.lit != "flip" {
            t.TypeMismatch(x, nil)
        }

        return t.FlipString(x)
    case Number:
        if t.lit != "~" && t.lit != "bnot" {
            t.TypeMismatch(x, nil)
        }

        return t.BnotNumber(x)
    case Boolean:
        return t.TopTwiddle(x.Number())
    case Null:
        return t.TopTwiddle(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) InvertHash(x Hash) Hash {
    out := Hash { }

    for key, val := range x {
        out[string(Stringify(val))] = String(key)
    }

    return out
}

func (t *Token) TransposeArray(x Array) Array {
    out := Array{ }

    for i, row := range x {
        switch row := row.(type) {
        case Array:
            for j, val := range row {
                if len(out) <= j {
                    out = append(out, Array { })
                }

                out[j] = append(out[j].(Array), val)
            }
        default:
            if len(out) <= i {
                out = append(out, Array { })
            }

            out[i] = append(out[i].(Array), row)
        }
    }

    return out
}

func (t *Token) FlipString(x String) String {
    out := ""

    for _, c := range x {
        out += string(c ^ ' ')
    }

    return String(out)
}

func (t *Token) BnotNumber(x Number) Number {
    if x.inf == INF || x.inf == -INF {
        return NewNumber(-1)
    }

    return NewNumber(^x.Int())
}

func (t *Token) ExclusionHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key := range x {
        if _, ok := y[key]; !ok {
            out[key] = x[key]
        }
    }

    for key := range y {
        if _, ok := x[key]; !ok {
            out[key] = y[key]
        }
    }

    return out
}

func (t *Token) ExclusionArray(x Array, y Array) Array {
    out := Array { }

    for _, a := range x {
        found := false

        for _, b := range y {
            if Equals(a, b) {
                found = true
                break
            }
        }

        if !found {
            out = append(out, a)
        }
    }

    for _, b := range y {
        found := false

        for _, a := range x {
            if Equals(b, a) {
                found = true
                break
            }
        }

        if !found {
            out = append(out, b)
        }
    }

    return t.UniqueArray(out)
}

func (t *Token) ExclusionString(x String, y String) String {
    return t.JoinArray(t.ExclusionArray(x.Array(), y.Array()), String(""))
}

func (t *Token) BxorNumber(x Number, y Number) Number {
    if (x.inf == INF || x.inf == -INF) && (y.inf == INF || y.inf == -INF) {
        return NewNumber(0)
    }

    if x.inf == INF || x.inf == -INF {
        return y
    }

    if y.inf == INF || y.inf == -INF {
        return x
    }

    return NewNumber(x.Int() ^ y.Int())
}
