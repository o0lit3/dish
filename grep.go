package main
import("math/big"; "strings")

func (t *Token) Grep(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Grep(x.Run(), b)
    case *Variable:
        return t.Grep(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                    t.TypeMismatch(x, y)
                }

                return t.FilterHash(x, y)
            }

            return t.Grep(x, y.Run())
        case *Variable:
            return t.Grep(x, y.Value())
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                    t.TypeMismatch(x, y)
                }

                return t.FilterArray(x, y)
            }

            return t.Grep(x, y.Run())
        case Hash:
            return t.Grep(x, y.Array())
        case Array:
            if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                t.TypeMismatch(x, y)
            }

            return t.GrepArray(x, y)
        case *Variable:
            return t.Grep(x, y.Value())
        case Number:
            if t.lit != "%" && t.lit != "%=" && t.lit != "every" {
                t.TypeMismatch(x, y)
            }

            return t.EveryNthItem(x, y)
        case Boolean:
            return t.Grep(x, y.Number())
        case Null:
            return t.Grep(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                    t.TypeMismatch(x, y)
                }

                return t.FilterString(x, y)
            }

            return t.Grep(x, y.Run())
        case *Variable:
            return t.Grep(x, y.Value())
        case Hash:
            return t.Grep(x, y.Array())
        case Array:
            if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                t.TypeMismatch(x, y)
            }

            return t.GrepString(x, y)
        case String:
            if t.lit != "%" && t.lit != "%=" && t.lit != "filter" && t.lit != "select" && t.lit != "grep" {
                t.TypeMismatch(x, y)
            }

            return t.SelectString(x, y)
        case Number:
            if t.lit != "%" && t.lit != "%=" && t.lit != "every" {
                t.TypeMismatch(x, y)
            }

            return t.EveryNthChar(x, y)
        case Boolean:
            return t.Grep(x, y.Number())
        case Null:
            return t.Grep(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Grep(x, y.Run())
        case *Variable:
            return t.Grep(x, y.Value())
        case Number:
            if t.lit != "%" && t.lit != "%=" && t.lit != "mod" {
                t.TypeMismatch(x, y)
            }

            if y.inf == 0 && y.val.Cmp(NewNumber(0).val) == 0 {
                t.DivideByZero()
            }

            return t.ModNumber(x, y)
        case Boolean:
            return t.Grep(x, y.Number())
        case Null:
            return t.Grep(x, NewNumber(0))
        }
    case Boolean:
        return t.Grep(x.Number(), b)
    case Null:
        return t.Grep(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) FilterHash(x Hash, y *Block) Hash {
    out := Hash{ }

    for key, val := range x {
        if !Boolify(y.Context(x).Run(val, String(key))) {
            continue
        }

        out[key] = val
    }

    return out
}

func (t *Token) FilterArray(x Array, y *Block) Array {
    out := Array{ }

    for i, val := range x {
        if !Boolify(y.Context(x).Run(val, NewNumber(i))) {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (t *Token) GrepArray(x Array, y Array) Array {
    indices := t.SortArray(t.SearchInArray(x, y), nil)
    out := Array{ }

    for _, i := range indices {
        out = append(out, x[i.(Number).Int()])
    }

    return out
}

func (t *Token) FilterString(x String, y *Block) String {
    return t.JoinArray(t.FilterArray(x.Array(), y), String(""))
}

func (t *Token) GrepString(x String, y Array) String {
    indices := t.SortArray(t.SearchInString(x, y), nil)
    out := []rune{ }

    for _, i := range indices {
        out = append(out, x[i.(Number).Int()])
    }

    return String(out)
}

func (t *Token) SelectString(x String, y String) String {
    return t.RepeatString(y, NewNumber(strings.Count(string(x), string(y))))
}

func (t *Token) EveryNthItem(x Array, y Number) Array {
    if y.val.Cmp(NewNumber(0).val) == 0 {
        return Array{ }
    }

    out := Array{ }

    for i, val := range x {
        if i % y.Int() > 0 {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (t *Token) EveryNthChar(x String, y Number) String {
    if y.val.Cmp(NewNumber(0).val) == 0 {
        return String("")
    }

    out := ""

    for i, c := range x {
        if i % y.Int() > 0 {
            continue
        }

        out += string(c)
    }

    return String(out)
}

func (t *Token) ModNumber(x Number, y Number) interface{} {
    if x.inf == INF || x.inf == -INF {
        return Null{ }
    }

    if y.inf == INF || y.inf == -INF {
        return x
    }

    div := new(big.Rat).Quo(x.val, y.val)
    flr := new(big.Rat).SetInt(new(big.Int).Quo(div.Num(), div.Denom()))
    out := Number{ val: new(big.Rat).Sub(x.val, new(big.Rat).Mul(y.val, flr)) }

    if out.val.Cmp(NewNumber(0).val) < 0 && x.val.Cmp(NewNumber(0).val) + y.val.Cmp(NewNumber(0).val) >= 0 {
        out = t.ModNumber(t.AddNumber(out, y).(Number), y).(Number)
    }

    return out
}

func (t *Token) DoubleGrep(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleGrep(x.Run(), b)
    case *Variable:
        return t.DoubleGrep(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%%" && t.lit != "without" {
                    t.TypeMismatch(x, y)
                }

                return t.WithoutHash(x, y)
            }

            return t.DoubleGrep(x, y.Run())
        case *Variable:
            return t.DoubleGrep(x, y.Value())
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%%" && t.lit != "without" {
                    t.TypeMismatch(x, y)
                }

                return t.WithoutArray(x, y)
            }

            return t.DoubleGrep(x, y.Run())
        case *Variable:
            return t.DoubleGrep(x, y.Value())
        case Hash:
            return t.DoubleGrep(x, y.Array())
        case Array:
            if t.lit != "%%" && t.lit != "without" {
                t.TypeMismatch(x, y)
            }

            return t.ExcludeArray(x, y)
        case Number:
            if t.lit != "%%" && t.lit != "xevery" {
                t.TypeMismatch(x, y)
            }

            return t.ExcludingEveryNthItem(x, y)
        case Boolean:
            return t.DoubleGrep(x, y.Number())
        case Null:
            return t.DoubleGrep(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "%%" && t.lit != "without" {
                    t.TypeMismatch(x, y)
                }

                return t.WithoutString(x, y)
            }

            return t.DoubleGrep(x, y.Run())
        case *Variable:
            return t.DoubleGrep(x, y.Value())
        case Hash:
            return t.DoubleGrep(x, y.Array())
        case Array:
            if t.lit != "%%" && t.lit != "without" {
                t.TypeMismatch(x, y)
            }

            return t.ExcludeFromString(x, y)
        case String:
            if t.lit != "%%" && t.lit != "without" {
                t.TypeMismatch(x, y)
            }

            return t.ExcludeString(x, y)
        case Number:
            if t.lit != "%%" && t.lit != "xevery" {
                t.TypeMismatch(x, y)
            }

            return t.ExcludingEveryNthChar(x, y)
        case Boolean:
            return t.DoubleGrep(x, y.Number())
        case Null:
            return t.DoubleGrep(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.DoubleGrep(x, y.Run())
        case *Variable:
            return t.DoubleGrep(x, y.Value())
        case Number:
            if t.lit != "%%" && t.lit != "imod" {
                t.TypeMismatch(x, y)
            }

            return NewNumber(t.ModNumber(x, y).(Number).Int())
        case Boolean:
            return t.DoubleGrep(x, y.Number())
        case Null:
            return t.DoubleGrep(x, NewNumber(0))
        }
    case Boolean:
        return t.DoubleGrep(x.Number(), b)
    case Null:
        return t.DoubleGrep(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopGrep(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopGrep(x.Run())
    case *Variable:
        return t.TopGrep(x.Value())
    case Hash:
        return x
    case Array:
        return Hashify(x)
    case String:
        return Hashify(x)
    case Number:
        return t.RatioNumber(x)
    case Boolean:
        return t.TopGrep(x.Number())
    case Null:
        return t.TopGrep(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) WithoutHash(x Hash, y *Block) Hash {
    out := Hash{ }

    for key, val := range x {
        if Boolify(y.Context(x).Run(val, String(key))) {
            continue
        }

        out[key] = val
    }

    return out
}

func (t *Token) WithoutArray(x Array, y *Block) Array {
    out := Array{ }

    for i, val := range x {
        if Boolify(y.Context(x).Run(val, NewNumber(i))) {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (t *Token) WithoutString(x String, y *Block) String {
    return t.JoinArray(t.WithoutArray(x.Array(), y), String(""))
}

func (t *Token) ExcludingEveryNthItem(x Array, y Number) Array {
    if y.val.Cmp(NewNumber(0).val) == 0 {
        return x
    }

    out := Array{ }

    for i, val := range x {
        if i % y.Int() == 0 {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (t *Token) ExcludingEveryNthChar(x String, y Number) String {
   if y.val.Cmp(NewNumber(0).val) == 0 {
       return x
   }

   out := ""

   for i, c := range x {
        if i % y.Int() == 0 {
            continue;
        }

        out += string(c)
   }

   return String(out)
}

func (t *Token) ExcludeArray(x Array, y Array) Array {
    out := Array{ }

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

    return out
}

func (t *Token) ExcludeFromString(x String, y Array) String {
    out := []rune{ }

    for _, c := range x {
        found := false

        for _, b := range y {
            if Equals(String([]rune{ c }), Stringify(b)) {
                found = true
                break
            }
        }

        if !found {
            out = append(out, c)
        }
    }

    return String(out)
}

func (t *Token) ExcludeString(x String, y String) String {
    return t.JoinArray(t.SplitString(x, y), String(""))
}

func (t *Token) RatioNumber(x Number) Hash {
    if x.inf == INF || x.inf == -INF {
        return Hash{
            "num": Number{ inf: x.inf },
            "denom": NewNumber(1),
        }
    }

    return Hash{
        "num": Number{ val: new(big.Rat).SetInt(x.val.Num()) },
        "denom": Number{ val: new(big.Rat).SetInt(x.val.Denom()) },
    }
}
