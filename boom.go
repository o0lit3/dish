package main

func (t *Token) Boom(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Boom(x.Run(), b)
    case *Variable:
        return t.Boom(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "&" && t.lit != "&=" && t.lit != "all" {
                    t.TypeMismatch(x, y)
                }

                for key, val := range(x) {
                    if !Boolify(y.Run(val, String(key))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Boom(x, y.Run())
        case *Variable:
            return t.Boom(x, y.Value())
        case Hash:
            if t.lit != "&" && t.lit != "&=" && t.lit != "intersect" {
                t.TypeMismatch(x, y)
            }

            return t.IntersectHash(x, y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "&" && t.lit != "&=" && t.lit != "all" {
                    t.TypeMismatch(x, y)
                }

                for i, val := range(x) {
                    if !Boolify(y.Run(val, NewNumber(i))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Boom(x, y.Run())
        case *Variable:
            return t.Boom(x, y.Value())
        case Array:
            if t.lit != "&" && t.lit != "&=" && t.lit != "intersect" {
                t.TypeMismatch(x, y)
            }

            return t.IntersectArray(x, y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "&" && t.lit != "&=" && t.lit != "all" {
                    t.TypeMismatch(x, y)
                }

                for i, val := range(x) {
                    if !Boolify(y.Run(String(string(val)), NewNumber(i))) {
                        return Boolean(false)
                    }
                }

                return Boolean(true)
            }

            return t.Boom(x, y.Run())
        case *Variable:
            return t.Boom(x, y.Value())
        case String:
            if t.lit != "&" && t.lit != "&=" && t.lit != "intersect" {
                t.TypeMismatch(x, y)
            }

            return t.IntersectString(x, y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Boom(x, y.Run())
        case *Variable:
            return t.Boom(x, y.Value())
        case Number:
            if t.lit != "&" && t.lit != "&=" && t.lit != "band" {
                t.TypeMismatch(x, y)
            }

            return t.BandNumber(x, y)
        case Boolean:
            return t.Boom(x, y.Number())
        case Null:
            return t.Boom(x, NewNumber(0))
        }
    case Boolean:
        return t.Boom(x.Number(), b)
    case Null:
        return t.Boom(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) DoubleBoom(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleBoom(x.Run(), b)
    case *Variable:
        return t.DoubleBoom(x.Value(), b)
    case Boolean:
        if !x {
            return x
        }

        switch y := b.(type) {
        case *Block:
            return t.DoubleBoom(Boolean(true), y.Run())
        case *Variable:
            return t.DoubleBoom(Boolean(true), y.Value())
        default:
            return y
        }
    default:
        if Boolify(x) {
            return t.DoubleBoom(Boolean(true), b)
        }

        return x
    }
}

func (t *Token) TopBoom(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopBoom(x.Run())
    case *Variable:
        return t.TopBoom(x.Value())
    case Hash:
        if t.lit != "&" && t.lit != "compact" {
            t.TypeMismatch(a, nil)
        }

        return t.CompactHash(x)
    case Array:
        if t.lit != "&" && t.lit != "compact" {
            t.TypeMismatch(a, nil)
        }

        return t.CompactArray(x)
    case String:
        if t.lit != "&" && t.lit != "compact" {
            t.TypeMismatch(a, nil)
        }

        return t.JoinArray(t.SplitString(x, String(" ")), String(""))
    case Number:
        if t.lit != "&" && t.lit != "popcount" {
            t.TypeMismatch(a, nil)
        }

        return NewNumber(len(t.SearchArray(x.Array(), NewNumber(1))))
    case Boolean:
        return t.TopBoom(x.Number())
    case Null:
        return t.TopBoom(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) CompactHash(x Hash) Hash {
    out := Hash { }

    for key, val := range x {
        switch w := val.(type) {
        case *Variable:
            if w.nom != "null" {
                out[key] = val
            }
        case Null:
        default:
            out[key] = val
        }
    }

    return out
}

func (t *Token) CompactArray(x Array) Array {
    out := Array { }

    for _, val := range x {
        switch w := val.(type) {
        case *Variable:
            if w.nom != "null" {
                out = append(out, val)
            }
        case Null:
        default:
            out = append(out, val)
        }
    }

    return out
}

func (t *Token) IntersectHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key := range x {
        if _, ok := y[key]; ok {
            out[key] = y[key]
        }
    }

    return out
}

func (t *Token) IntersectArray(x Array, y Array) Array {
    out := Array { }

    for _, a := range x {
        for _, b := range y {
            if Equals(a, b) {
                out = append(out, a)
                break
            }
        }
    }

    return t.UniqueArray(out)
}

func (t *Token) IntersectString(x String, y String) String {
    return t.JoinArray(t.IntersectArray(x.Array(), y.Array()), String(""))
}

func (t *Token) BandNumber(x Number, y Number) Number {
    if x.inf == INF || x.inf == -INF || y.inf == INF || y.inf == -INF {
        return NewNumber(0)
    }

    return NewNumber(x.Int() & y.Int())
}
