package main
import("strings"; "strconv")

func (t *Token) Dot(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Dot(x.Run(), b)
    case *Variable:
        switch val := x.Value().(type) {
        case Null:
            switch y := b.(type) {
            case *Block:
                if len(y.args) > 0 {
                    if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                        t.TypeMismatch(x, y)
                    }

                    return y.Run(val)
                }

                return t.Dot(x, y.Run())
            case Array:
                if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                    t.TypeMismatch(x, y)
                }

                x.blk.cur.vars[x.nom] = Array { }

                return &Variable{ par: x, obj: x.blk.cur.vars[x.nom], arr: t.ArrayMembers(Array{ }, y).arr }
            case String:
                if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                    t.TypeMismatch(x, y)
                }

                x.blk.cur.vars[x.nom] = Hash { }

                return &Variable{ par: x, obj: x.blk.cur.vars[x.nom], nom: string(y) }
            case Number:
                if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                    t.TypeMismatch(x, y)
                }

                x.blk.cur.vars[x.nom] = Array { }

                return &Variable{ par: x, obj: x.blk.cur.vars[x.nom], idx: y.Int() }
            }
        default:
            switch out := t.Dot(val, b).(type) {
            case *Variable:
                out.par = x
                return out
            default:
                return out
            }
        }
    case Hash:
        switch y := b.(type) {
        case *Block:
            switch len(y.args) {
            case 0:
                return t.Dot(x, y.Run())
            case 1:
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x)
            default:
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x.Array()...)
            }
        case *Variable:
            return t.Dot(x, y.Value())
        case Hash:
            if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.HashMembers(x, t.HashKeys(y))
        case Array:
            if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.HashMembers(x, y)
        case String:
            if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.HashMember(x, y)
        case Number:
            if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.HashMember(x, String(y.String()))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            switch len(y.args) {
            case 0:
                return t.Dot(x, y.Run())
            case 1:
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x)
            default:
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x...)
            }
        case *Variable:
            return t.Dot(x, y.Value())
        case Array:
            if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.ArrayMembers(x, y)
        case Number:
            if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.ArrayMember(x, y)
        case Boolean:
            return t.Dot(x, y.Number())
        case Null:
            return t.Dot(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x)
            }

            return t.Dot(x, y.Run())
        case *Variable:
            return t.Dot(x, y.Value())
        case Array:
            if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.StringMembers(x, y)
        case Number:
            if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.StringMember(x, y)
        case Boolean:
            return t.Dot(x, y.Number())
        case Null:
            return t.Dot(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit == "at" || t.lit == "item" || t.lit == "items" || t.lit == "subset" {
                    t.TypeMismatch(x, y)
                }

                return y.Run(x)
            }

            return t.Dot(x, y.Run())
        case *Variable:
            return t.Dot(x, y.Value())
        case Array:
            if t.lit == "at" || t.lit == "item" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.NumberMembers(x, y)
        case Number:
            if t.lit == "items" || t.lit == "subset" || t.lit == "call" {
                t.TypeMismatch(x, y)
            }

            return t.NumberMember(x, y)
        case Boolean:
            return t.Dot(x, y.Number())
        case Null:
            return t.Dot(x, NewNumber(0))
        }
    case Boolean:
        return t.Dot(x.Number(), b)
    case Null:
        return t.Dot(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) DoubleDot(a interface{}, b interface{}) Array {
    switch x := a.(type) {
    case *Block:
        return t.DoubleDot(x.Run(), b)
    case *Variable:
        return t.DoubleDot(x.Value(), b)
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.DoubleDot(x, y.Run())
        case *Variable:
            return t.DoubleDot(x, y.Value())
        case Number:
            return t.RangeNumber(x, y)
        case Boolean:
            return t.DoubleDot(x, y.Number())
        case Null:
            return t.DoubleDot(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.DoubleDot(x, y.Run())
        case *Variable:
            return t.DoubleDot(x, y.Value())
        case String:
            return t.RangeString(x, y)
        }
    case Boolean:
        return t.DoubleDot(x.Number(), b)
    case Null:
        return t.DoubleDot(NewNumber(0), b)
    }

    t.TypeMismatch(a, b)

    return Array { }
}

func (t *Token) HashMember(x Hash, y String) *Variable {
    return &Variable{ obj: x, nom: string(y) }
}

func (t *Token) HashMembers(x Hash, y Array) *Variable {
    out := []*Variable{ }

    for _, val := range y {
        switch val := val.(type) {
        case String:
            v := t.HashMember(x, val)
            v.sub = true
            out = append(out, v)
        default:
            t.TypeMismatch(x, y)
        }
    }

    return &Variable{ obj: x, arr: out }
}

func (t *Token) ArrayMember(x Array, y Number) *Variable {
    b := y.Int()

    if b < 0 && len(x) > 0 && len(x) + b < len(x) {
        return &Variable{ obj: x, idx: len(x) + b }
    }

    if len(x) > 0 && b < len(x) {
        return &Variable{ obj: x, idx: b }
    }

    if b < 0 {
        return &Variable{ obj: x, idx: -b }
    }

    return &Variable{ obj: x, idx: b }
}

func (t *Token) ArrayMembers(x Array, y Array) *Variable {
    out := []*Variable{ }

    for _, val := range y {
        switch val := val.(type) {
        case Number:
            v := t.ArrayMember(x, val)
            v.sub = true
            out = append(out, v)
        default:
            t.TypeMismatch(x, val)
        }
    }

    return &Variable{ obj: x, arr: out }
}

func (t *Token) StringMember(x String, y Number) *Variable {
    b := y.Int()

    if b < 0 && len(x) > 0 && len(x) + b < len(x) {
        return &Variable{ obj: x, idx: len(x) + b }
    }

    if len(x) > 0 && b < len(x) {
        return &Variable{ obj: x, idx: b }
    }

    if b < 0 {
        return &Variable{ obj: x, idx: -b }
    }

    return &Variable{ obj: x, idx: b }
}

func (t *Token) StringMembers(x String, y Array) *Variable {
    out := []*Variable{ }

    for _, val := range y {
        switch val := val.(type) {
        case Number:
            v := t.StringMember(x, val)
            v.sub = true
            out = append(out, v)
        default:
            t.TypeMismatch(x, val)
        }
    }

    return &Variable{ obj: x, arr: out }
}

func (t *Token) NumberMember(x Number, y Number) *Variable {
    b := y.Int()
    bin := strconv.FormatInt(int64(x.Int()), 2)

    if b < 0 && len(bin) > 0 && len(bin) + b < len(bin) {
        return &Variable{ obj: x, idx: len(bin) + b }
    }

    if len(bin) > 0 && b < len(bin) {
        return &Variable{ obj: x, idx: b }
    }

    if b < 0 {
        return &Variable{ obj: x, idx: -b }
    }

    return &Variable{ obj: x, idx: b }
}

func (t *Token) NumberMembers(x Number, y Array) *Variable {
    out := []*Variable{ }

    for _, val := range y {
        switch val := val.(type) {
        case Number:
            v := t.NumberMember(x, val)
            v.sub = true
            out = append(out, v)
        default:
            t.TypeMismatch(x, val)
        }
    }

    return &Variable{ obj: x, arr: out }
}

func (t *Token) RangeNumber(x Number, y Number) Array {
    out := Array { }
    a := x.Int()
    b := y.Int()

    if x.inf == INF || x.inf == -INF || y.inf == INF || y.inf == -INF {
        return out
    }

    if a > b {
        for a >= b {
            out = append(out, NewNumber(a))
            a--
        }
    } else {
        for a <= b {
            out = append(out, NewNumber(a))
            a++
        }
    }

    return out
}

func (t *Token) RangeString(x String, y String) Array {
    out := Array { }
    a := string(x)
    b := string(y)

    if len(x) > len(y) || (len(x) == len(y) && strings.ToLower(string(x)) > strings.ToLower(b)) {
        for len(x) >= len(y) && string(x) != b {
            out = append(out, x)
            x = t.DecreaseString(x, NewNumber(1))

            if (string(x) == a) {
                break
            }
        }

        if string(x) == b {
            out = append(out, x)
        }
    } else {
        for len(x) <= len(y) && string(x) != b {
            out = append(out, x)
            x = t.IncreaseString(x, NewNumber(1))

            if (string(x) == a) {
                break
            }
        }

        if string(x) == b {
            out = append(out, x)
        }
    }

    return out
}
