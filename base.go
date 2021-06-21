package main

import(
    "fmt"
    "strings"
    "strconv"
    "unicode"
)

func Base(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Base(x.Run(), b)
    case *Variable:
        return Base(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Format(y)
        case Number:
            return x.Base(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Format(y)
        case Number:
            return x.Base(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
           return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Format(y)
        case Number:
            return x.Base(y)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Base(x, y.Run())
        case *Variable:
            return Base(x, y.Value())
        case Hash:
            return x.Base(NewNumber(len(y)))
        case Array:
            return x.Base(NewNumber(len(y)))
        case String:
            return x.Format(y)
        case Number:
            return x.Base(y)
        }
    case Boolean:
        return Base(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Base(b Number) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Base(val, b)
    }

    return out
}

func (a Array) Base(b Number) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Base(val, b))
    }

    return out
}

func (a Hash) Format(b String) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Base(val, b)
    }

    return out
}

func (a Array) Format(b String) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Base(val, b))
    }

    return out
}

func (a String) Format(b String) String {
    return String(fmt.Sprintf(string(b), string(a)))
}

func (a String) Base(b Number) Number {
    if b.val.Cmp(NewNumber(2).val) == -1 || b.val.Cmp(NewNumber(36).val) == 1 {
        return NewNumber(0)
    }

    out, _ := strconv.ParseInt(string(a), b.Int(), 64)

    return NewNumber(int(out))
}

func (a Number) Base(b Number) String {
    if b.inf == INF || b.inf == -INF {
        return String("")
    }

    if a.inf == INF || a.inf == -INF {
        return String(a.String())
    }

    if b.val.Cmp(NewNumber(2).val) == -1 || b.val.Cmp(NewNumber(36).val) == 1 {
        return String("")
    }

    return String(strconv.FormatInt(int64(a.Int()), b.Int()))
}

func (a Number) Format(b String) String {
    if a.inf == INF || a.inf == -INF {
        return String(a.String())
    }

    parts := strings.Split(string(b), ".")

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
            out = a.val.FloatString(n)
        }

        return String(fmt.Sprintf(parts[0] + "s", out))
    }

    if val, ok := a.val.Float64(); ok {
        if strings.Contains(parts[0], "f") {
            return String(fmt.Sprintf(parts[0], val))
        }

        return String(fmt.Sprintf(parts[0], int(val)))
    }

    return String(a.val.FloatString(0))
}
