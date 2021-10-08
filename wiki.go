package main
import("strconv"; "strings"; "math/big")

func (t *Token) Wiki(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Wiki(x.Run(), b)
    case *Variable:
        return t.Wiki(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.Wiki(x, y.Run())
        case *Variable:
            return t.Wiki(x, y.Value())
        case Hash:
            return Boolean(len(x) < len(y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.Wiki(x, y.Run())
        case *Variable:
            return t.Wiki(x, y.Value())
        case Array:
            return Boolean(len(x) < len(y))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.Wiki(x, y.Run())
        case *Variable:
            return t.Wiki(x, y.Value())
        case String:
            return Boolean(string(x) < string(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Wiki(x, y.Run())
        case *Variable:
            return t.Wiki(x, y.Value())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(false)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(false)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(y.val) < 0)
        case Boolean:
            return t.Wiki(x, y.Number())
        case Null:
            return t.Wiki(x, NewNumber(0))
        }
    case Boolean:
        return t.Wiki(x.Number(), b)
    case Null:
        return t.Wiki(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) WikiBars(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.WikiBars(x.Run(), b)
    case *Variable:
        return t.WikiBars(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.WikiBars(x, y.Run())
        case *Variable:
            return t.WikiBars(x, y.Value())
        case Hash:
            return Boolean(len(x) <= len(y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.WikiBars(x, y.Run())
        case *Variable:
            return t.WikiBars(x, y.Value())
        case Array:
            return Boolean(len(x) <= len(y))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.WikiBars(x, y.Run())
        case *Variable:
            return t.WikiBars(x, y.Value())
        case String:
            return Boolean(string(x) <= string(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.WikiBars(x, y.Run())
        case *Variable:
            return t.WikiBars(x, y.Value())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(false)
            }

            if x.inf == -INF || y.inf == INF {
                return Boolean(false)
            }

            if x.inf == INF || y.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(y.val) <= 0)
        case Boolean:
            return t.WikiBars(x, y.Number())
        case Null:
            return t.WikiBars(x, NewNumber(0))
        }
    case Boolean:
        return t.WikiBars(x.Number(), b)
    case Null:
        return t.WikiBars(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) WikiWiki(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.WikiWiki(x.Run(), b)
    case *Variable:
        return t.WikiWiki(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.WikiWiki(x, y.Run())
        case *Variable:
            return t.WikiWiki(x, y.Value())
        case Hash:
            if t.lit != "<<" && t.lit != "extend" {
                t.TypeMismatch(x, y)
            }

            return t.ExtendHash(x, y)
        case Null:
            return t.WikiWiki(x, Hash{ })
        default:
            return t.WikiWiki(x, Hashify(y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.WikiWiki(x, y.Run())
        case *Variable:
            return t.WikiWiki(x, y.Value())
        case Array:
            if t.lit != "<<" && t.lit != "push" && t.lit != "append" {
                t.TypeMismatch(x, y)
            }

            return t.PushArray(x, t.FlattenArray(y))
        case Null:
            if t.lit != "<<" && t.lit != "push" && t.lit != "append" {
                t.TypeMismatch(x, y)
            }

            return x
        default:
            return t.WikiWiki(x, Array{ y })
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.WikiWiki(x, y.Run())
        case *Variable:
            return t.WikiWiki(x, y.Value())
        case String:
            if t.lit != "<<" && t.lit != "append" {
                t.TypeMismatch(x, y)
            }

            return t.AppendString(x, y)
        case Null:
            if t.lit != "<<" && t.lit != "append" {
                t.TypeMismatch(x, y)
            }

            return x
        default:
            return t.WikiWiki(x, Stringify(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.WikiWiki(x, y.Run())
        case *Variable:
            return t.WikiWiki(x, y.Value())
        case Hash:
            return t.ExtendHash(Hashify(x), y)
        case Array:
            return t.PushArray(Array{ x }, y)
        case String:
            return t.AppendString(Stringify(x), y)
        case Number:
            if t.lit != "<<" && t.lit != "lshift" {
                t.TypeMismatch(x, y)
            }

            return t.LshiftNumber(x, y)
        case Boolean:
            return t.WikiWiki(x, y.Number())
        case Null:
            return t.WikiWiki(x, NewNumber(0))
        }
    case Boolean:
        return t.WikiWiki(x.Number(), b)
    case Null:
        return t.WikiWiki(Array{ }, b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopWiki(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopWiki(x.Run())
    case *Variable:
        return t.TopWiki(x.Value())
    case Hash:
        return t.TopWiki(x.Array())
    case Array:
        if t.lit != "<" && t.lit != "minimum" && t.lit != "min" {
            t.TypeMismatch(x, nil)
        }

        return t.MinFromArray(x)
    case String:
        if t.lit != "<" && t.lit != "lowercase" && t.lit != "downcase" && t.lit != "lc" {
            t.TypeMismatch(x, nil)
        }

        return t.LowerString(x)
    case Number:
        if t.lit != "<" && t.lit != "floor" && t.lit != "int" {
            t.TypeMismatch(x, nil)
        }

        return t.FloorNumber(x)
    case Boolean:
        return t.TopWiki(x.Number())
    case Null:
        return t.TopWiki(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) TopWikiWiki(a interface{}) (interface{}, interface{}) {
    switch x := a.(type) {
    case *Block:
        return t.TopWikiWiki(x.Run())
    case *Variable:
        return t.TopWikiWiki(x.Value())
    case Array:
        if t.lit != "<<" && t.lit != "first" && t.lit != "shift" {
            t.TypeMismatch(x, nil)
        }

        return t.ShiftArray(x)
    case String:
        if t.lit != "<<" && t.lit != "first" && t.lit != "shift" {
            t.TypeMismatch(x, nil)
        }

        return t.ShiftString(x)
    case Number:
        if t.lit != "<<" && t.lit != "msb" && t.lit != "shift" {
            t.TypeMismatch(x, nil)
        }

        return t.ShiftNumber(x)
    case Boolean:
        return t.ShiftNumber(x.Number())
    case Null:
        return t.ShiftNumber(NewNumber(0))
    }

    return a, t.TypeMismatch(a, nil)
}

func (t *Token) ExtendHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key, val := range x {
        out[key] = val
    }

    for key, val := range y {
        out[key] = val
    }

    return out
}

func (t *Token) PushArray(x Array, y Array) Array {
    out := Array { }

    for _, val := range x {
        out = append(out, val)
    }

    for _, val := range y {
        out = append(out, val)
    }

    return out
}

func (t *Token) AppendString(x String, y String) String {
    return String(string(x) + string(y))
}

func (t *Token) LshiftNumber(x Number, y Number) Number {
    if x.inf == INF || x.inf == -INF {
        return NewNumber(0)
    }

    if y.inf == INF || y.inf == -INF {
        return x
    }

    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.RshiftNumber(x, t.NegateNumber(y))
    }

    return NewNumber(x.Int() << uint(y.Int()))
}

func (t *Token) ShiftArray(x Array) (interface{}, interface{}) {
    if len(x) > 0 {
        val := x[0]
        x = x[1:]
        return val, x
    }

    return Null { }, x
}

func (t *Token) ShiftString(x String) (interface{}, interface{}) {
    if len(x) > 0 {
        val := String(string(x[0]))
        x = x[1:]
        return val, x
    }

    return Null { }, x
}

func (t *Token) ShiftNumber(x Number) (interface{}, interface{}) {
    if x.inf == INF || x.inf == -INF {
        return NewNumber(1), x
    }

    bin := strconv.FormatInt(int64(x.Int()), 2)
    last, _ := strconv.Atoi(string(bin[0]))
    rem, _ := strconv.ParseInt(string(bin[1:]), 2, 64)
    return NewNumber(last), NewNumber(int(rem))
}

func (t *Token) MinFromArray(x Array) interface{} {
    var out interface{}

    for _, val := range x {
        if out == nil {
            out = val
        } else if b, ok := t.Wiki(val, out).(Boolean); bool(b) && ok {
            out = val
        }
    }

    return out
}

func (t *Token) LowerString(x String) String {
    return String(strings.ToLower(string(x)))
}

func (t *Token) FloorNumber(x Number) Number {
    if x.inf == INF || x.inf == -INF || x.val.IsInt() {
        return x
    }

    if x.val.Cmp(NewNumber(0).val) < -1 {
        return t.NegateNumber(t.CeilNumber(t.NegateNumber(x)))
    }

    return Number{ val: new(big.Rat).SetInt(new(big.Int).Quo(x.val.Num(), x.val.Denom())) }
}
