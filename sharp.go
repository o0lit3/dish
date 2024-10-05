package main
import("fmt"; "strings"; "strconv"; "unicode")

func (t *Token) Sharp(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Sharp(x.Run(), b)
    case *Variable:
        return t.Sharp(x.Value(), b)
    case Hash:
        switch y:= b.(type) {
        case *Block:
            return t.Sharp(x, y.Run())
        case *Variable:
            return t.Sharp(x, y.Value())
        case String:
            return t.FormatArray(x.Array(), y)
        }
    case Array:
        switch y:= b.(type) {
        case *Block:
            return t.Sharp(x, y.Run())
        case *Variable:
            return t.Sharp(x, y.Value())
        case String:
            return t.FormatArray(x, y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
           return t.Sharp(x, y.Run())
        case *Variable:
            return t.Sharp(x, y.Value())
        case String:
            if t.lit != "#" && t.lit != "format" && t.lit != "fmt" {
                t.TypeMismatch(x, y)
            }

            return t.FormatString(x, y)
        case Number:
            if t.lit != "#" && t.lit != "base" && t.lit != "unbase" {
                t.TypeMismatch(x, y)
            }

            return t.BaseString(x, y)
        case Boolean:
            return t.Sharp(x, y.Number())
        case Null:
            return t.Sharp(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Sharp(x, y.Run())
        case *Variable:
            return t.Sharp(x, y.Value())
        case String:
            if t.lit != "#" && t.lit != "format" && t.lit != "fmt" {
                t.TypeMismatch(x, y)
            }

            return t.FormatNumber(x, y)
        case Number:
            if t.lit != "#" && t.lit != "base" {
                t.TypeMismatch(x, y)
            }

            return t.BaseNumber(x, y)
        }
    case Boolean:
        switch b.(type) {
        case Number:
            return t.Sharp(x.Number(), b)
        default:
            return t.Sharp(Stringify(x), b)
        }
    case Null:
        switch b.(type) {
        case Number:
            return t.Sharp(NewNumber(0), b)
        default:
            return t.Sharp(Stringify(x), b)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopSharp(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopSharp(x.Run())
    case *Variable:
        return t.TopSharp(x.Value())
    case Hash:
        if t.lit != "#" && t.lit != "length" && t.lit != "len" && t.lit != "count" {
            t.TypeMismatch(a, nil)
        }

        return t.LengthHash(x)
    case Array:
        if t.lit != "#" && t.lit != "length" && t.lit != "len" && t.lit != "count" {
            t.TypeMismatch(a, nil)
        }

        return t.LengthArray(x)
    case String:
        if t.lit != "#" && t.lit != "length" && t.lit != "len" && t.lit != "count" {
            t.TypeMismatch(a, nil)
        }

        return t.LengthString(x)
    case Number:
        if t.lit != "#" && t.lit != "bitcount" && t.lit != "count" {
            t.TypeMismatch(a, nil)
        }

        if x.inf == INF || x.inf == -INF {
            return Number{ inf: INF }
        }

        if x.val.Cmp(NewNumber(0).val) == 0 {
            return NewNumber(0)
        }

        return t.LengthNumber(x)
    case Boolean:
        return t.TopSharp(x.Number())
    case Null:
        return t.TopSharp(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) FormatArray(x Array, y String) String {
    out := []interface{}{}

    for _, val := range x {
        switch val := val.(type) {
        case String:
            out = append(out, string(val))
        default:
            out = append(out, fmt.Sprintf("%v", val))
        }
    }

    return String(fmt.Sprintf(string(y), out...))
}

func (t *Token) FormatString(x String, y String) String {
    return String(fmt.Sprintf(string(y), string(x)))
}

func (t *Token) FormatNumber(x Number, y String) String {
    if x.inf == INF || x.inf == -INF {
        return String(x.String())
    }

    parts := strings.Split(string(y), ".")

    if len(parts) > 1 {
        out := ""
        dec := ""

        for _, c := range parts[1] {
            if unicode.IsDigit(c) {
                dec += string(c)
            } else {
                break
            }
        }

        if n, err := strconv.Atoi(dec); err == nil {
            out = x.val.FloatString(n)
        }

        return String(fmt.Sprintf(parts[0] + "s", out))
    }

    if val, ok := x.val.Float64(); ok {
        if strings.Contains(parts[0], "f") {
            return String(fmt.Sprintf(parts[0], val))
        }

        return String(fmt.Sprintf(parts[0], int(val)))
    }

    return String(x.val.FloatString(0))
}

func (t *Token) BaseString(x String, y Number) Number {
    if y.val.Cmp(NewNumber(2).val) == -1 || y.val.Cmp(NewNumber(36).val) == 1 {
        panic(fmt.Sprintf("Invalid base \"%s\" used near \"%s\" at %s", y, t.lit, t.pos))
    }

    out, err := strconv.ParseInt(string(x), y.Int(), 64)

    if err != nil {
        panic(fmt.Sprintf("Invalid string %s used near \"%s\" at %s", x, t.lit, t.pos))
    }

    return NewNumber(int(out))
}

func (t *Token) BaseNumber(x Number, y Number) String {
    if y.inf == INF || y.inf == -INF {
        panic(fmt.Sprintf("Invalid base \"%s\" used near \"%s\" at %s", y, t.lit, t.pos))
    }

    if x.inf == INF || x.inf == -INF {
        return String(x.String())
    }

    if y.val.Cmp(NewNumber(2).val) == -1 || y.val.Cmp(NewNumber(36).val) == 1 {
        panic(fmt.Sprintf("Invalid base \"%s\" used near \"%s\" at %s", y, t.lit, t.pos))
    }

    return String(strconv.FormatInt(int64(x.Int()), y.Int()))
}

func (t *Token) LengthHash(x Hash) Number {
    return NewNumber(len(x))
}

func (t *Token) LengthArray(x Array) Number {
    return NewNumber(len(x))
}

func (t *Token) LengthString(x String) Number {
    return NewNumber(len(x))
}

func (t *Token) LengthNumber(x Number) Number {
    return NewNumber(len(x.Array()))
}
