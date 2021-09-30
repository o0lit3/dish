package main

func (t *Token) Wham(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Wham(x.Run(), b)
    case *Variable:
        return t.Wham(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "|" && t.lit != "|=" && t.lit != "any" {
                    t.TypeMismatch(x, y)
                }

                return t.AnyHashItem(x, y)
            }

            return t.Wham(x, y.Run())
        case *Variable:
            return t.Wham(x, y.Value())
        case Hash:
            if t.lit != "|" && t.lit != "|=" && t.lit != "union" {
                t.TypeMismatch(x, y)
            }

            return t.UnionHash(x, y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "|" && t.lit != "|=" && t.lit != "any" {
                    t.TypeMismatch(x, y)
                }

                return t.AnyArrayItem(x, y)
            }

            return t.Wham(x, y.Run())
        case *Variable:
            return t.Wham(x, y.Value())
        case Array:
            if t.lit != "|" && t.lit != "|=" && t.lit != "union" {
                t.TypeMismatch(x, y)
            }

            return t.UnionArray(x, y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "|" && t.lit != "|=" && t.lit != "any" {
                    t.TypeMismatch(x, y)
                }

                for i, val := range(x) {
                    if Boolify(y.Run(String(string(val)), NewNumber(i))) {
                        return Boolean(true)
                    }
                }

                return Boolean(false)
            }

            return t.Wham(x, y.Run())
        case *Variable:
            return t.Wham(x, y.Value())
        case String:
            if t.lit != "|" && t.lit != "|=" && t.lit != "union" {
                t.TypeMismatch(x, y)
            }

            return t.UnionString(x, y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Wham(x, y.Run())
        case *Variable:
            return t.Wham(x, y.Value())
        case Number:
            if t.lit != "|" && t.lit != "|=" && t.lit != "bor" {
                t.TypeMismatch(x, y)
            }

            return t.BorNumber(x, y)
        case Boolean:
            return t.Wham(x, y.Number())
        case Null:
            return t.Wham(x, NewNumber(0))
        }
    case Boolean:
        return t.Wham(x.Number(), b)
    case Null:
        return t.Wham(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) DoubleWham(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleWham(x.Run(), b)
    case *Variable:
        return t.DoubleWham(x.Value(), b)
    case Boolean:
        if x {
            if t.lit == "else" {
                return Null{ }
            }

            return x
        }

        switch y := b.(type) {
        case *Block:
            return t.DoubleWham(Boolean(false), y.Run())
        case *Variable:
            return t.DoubleWham(Boolean(false), y.Value())
        default:
            return y
        }
    default:
        if Boolify(x) {
            return x
        }

        return t.DoubleWham(Boolean(false), b)
    }
}

func (t *Token) TopWham(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopWham(x.Run())
    case *Variable:
        return t.TopWham(x.Value())
    case Hash:
        if t.lit != "|" && t.lit != "unique" && t.lit != "uniq" {
            t.TypeMismatch(a, nil)
        }

        return t.UniqueHash(x)
    case Array:
        if t.lit != "|" && t.lit != "unique" && t.lit != "uniq" {
            t.TypeMismatch(x, nil)
        }

        return t.UniqueArray(x)
    case String:
        if t.lit != "|" && t.lit != "unique" && t.lit != "uniq" {
            t.TypeMismatch(x, nil)
        }

        return t.UniqueString(x)
    case Number:
        if t.lit != "|" && t.lit != "abs" {
            t.TypeMismatch(x, nil)
        }

        return t.AbsNumber(x)
    case Boolean:
        return t.TopWham(x.Number())
    case Null:
        return t.TopWham(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) UniqueHash(x Hash) Hash {
    out := Hash { }
    hash := make(map[string]bool)

    for k, val := range x {
        key := string(Stringify(val))

        if _, ok := hash[key]; !ok {
            out[k] = val
            hash[key] = true
        }
    }

    return out
}

func (t *Token) UniqueArray(x Array) Array {
    out := Array { }
    hash := make(map[string]bool)

    for _, val := range x {
        key := string(Stringify(val))

        if _, ok := hash[key]; !ok {
            out = append(out, val)
            hash[key] = true
        }
    }

    return out
}

func (t *Token) UniqueString(x String) String {
    out := ""
    hash := make(map[rune]bool)

    for _, c := range x {
        if _, ok := hash[c]; !ok {
            out += string(c)
            hash[c] = true
        }
    }

    return String(out)
}

func (t *Token) AbsNumber(x Number) Number {
    if x.inf == INF || x.inf == -INF {
        return Number{ inf: INF }
    }

    return Number{ val: NewNumber(0).val.Abs(x.val) }
}

func (t *Token) AnyHashItem(x Hash, y *Block) Boolean {
    for key, val := range(x) {
        if Boolify(y.Run(val, String(key))) {
            return Boolean(true)
        }
    }

    return Boolean(false)
}

func (t *Token) AnyArrayItem(x Array, y *Block) Boolean {
    for i, val := range(x) {
        if Boolify(y.Run(val, NewNumber(i))) {
            return Boolean(true)
        }
    }

    return Boolean(false)
}

func (t *Token) UnionHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key := range x {
        out[key] = x[key]
    }

    for key := range y {
        out[key] = y[key]
    }

    return out
}

func (t *Token) UnionArray(x Array, y Array) Array {
    out := Array { }

    for i := range x {
        out = append(out, x[i])
    }

    for i := range y {
        out = append(out, y[i])
    }

    return t.UniqueArray(out)
}

func (t *Token) UnionString(x String, y String) String {
    return t.JoinArray(t.UnionArray(x.Array(), y.Array()), String(""))
}

func (t *Token) BorNumber(x Number, y Number) Number {
    if (x.inf == INF || x.inf == -INF) && (y.inf == INF || y.inf == -INF) {
        return NewNumber(0)
    }

    if x.inf == INF || x.inf == -INF {
        return y
    }

    if y.inf == INF || y.inf == -INF {
        return x
    }

    return NewNumber(x.Int() | y.Int())
}
