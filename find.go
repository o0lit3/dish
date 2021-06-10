package main

import(
    "fmt"
    "strings"
    "strconv"
    "unicode"
    "math/big"
)

func Find(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Find(x.Run(), b)
    case *Variable:
        return Find(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        default:
            return x.Find(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        default:
            return x.Find(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        case String:
            return x.Format(y)
        default:
            return NewNumber(strings.Index(string(x), fmt.Sprintf("%v", y)))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case *Variable:
            return Find(x, y.Value())
        case String:
            return x.Format(y)
        case Number:
            return x.Round(y)
        case Boolean:
            return x.Round(y.Number())
        case Null:
            return x.Round(NewNumber(0))
        }
    case Boolean:
       return Find(x.Number(), b)
    }

    return NewNumber(0)
}

func (a Hash) Find(b interface{}) String {
    for key, val := range a {
        if Equals(val, b) {
            return String(key)
        }
    }

    return String("-1")
}

func (a Array) Find(b interface{}) Number {
    for i, val := range a {
        if Equals(val, b) {
            return NewNumber(i)
        }
    }

    return NewNumber(-1)
}

func (a String) Format(b String) String {
    return String(fmt.Sprintf(string(b), string(a)))
}

func (a Number) Format(b String) String {
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

func (a Number) Round(b Number) Number {
    if pow, ok := Power(NewNumber(10), b).(Number); ok {
        o := new(big.Rat).Mul(a.val, pow.val)
        i := new(big.Int).Quo(o.Num(), o.Denom())

        o = o.SetInt(i)
        return Number{ val: o.Quo(o, pow.val) }
    }

    return NewNumber(0)
}
