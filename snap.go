package main

import(
    "time"
    "math"
    "math/big"
    "math/rand"
)

func RootRat(val *big.Rat) (Number, bool) {
    num, ok := Root(val.Num(), 2)

    if !ok {
        return Number{ }, false
    }

    den, ok := Root(val.Denom(), 2)

    if !ok {
        return Number{ }, false
    }

    return NewRat(new(big.Rat).SetFrac(num, den)), true
}

func (t *Token) Approximate(val float64) interface{} {
    switch {
    case math.IsNaN(val):
        t.UnexpectedOperand()
    case math.IsInf(val, 1):
        return Number{ inf: INF }
    case math.IsInf(val, -1):
        return Number{ inf: -INF }
    }

    return NewRat(new(big.Rat).SetFloat64(val))
}

func (t *Token) Numbers(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Numbers(x.Run())
    case *Variable:
        return t.Numbers(x.Value())
    case Hash:
        out := Hash{ }

        for key, val := range x {
            out[key] = t.Numbers(val)
        }

        return out
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Numbers(val))
        }

        return out
    case Number:
        if x.inf == INF || x.inf == -INF {
            switch t.lit {
            case "prime":
                return Boolean(false)
            case "sqrt", "log":
                if x.inf == INF {
                    return Number{ inf: INF }
                }
            }

            t.UnexpectedOperand()
        }

        val, _ := x.Rat().Float64()

        switch t.lit {
        case "rand":
            rand.Seed(time.Now().UnixNano())
            return t.Approximate(rand.Float64() * val)
        case "prime":
            return Boolean(new(big.Int).Quo(x.Rat().Num(), x.Rat().Denom()).ProbablyPrime(0))
        case "sqrt":
            if root, ok := RootRat(x.Rat()); ok {
                return root
            }

            return t.Approximate(math.Sqrt(val))
        case "log":
            return t.Approximate(math.Log(val))
        case "sin":
            return t.Approximate(math.Sin(val))
        case "cos":
            return t.Approximate(math.Cos(val))
        case "tan":
            return t.Approximate(math.Tan(val))
        case "asin":
            return t.Approximate(math.Asin(val))
        case "acos":
            return t.Approximate(math.Acos(val))
        case "atan":
            return t.Approximate(math.Atan(val))
        }
    case Boolean:
        return t.Numbers(x.Number())
    case Null:
        return t.Numbers(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func Tick() Number {
    return Stamped(time.Now())
}

func Stamped(tm time.Time) Number {
    val := new(big.Rat).SetInt64(tm.Unix())

    if nsec := tm.Nanosecond(); nsec != 0 {
        val.Add(val, new(big.Rat).SetFrac64(int64(nsec), 1000000000))
    }

    return NewRat(val)
}

func (t *Token) Ticked(x Number) time.Time {
    if x.inf != 0 {
        t.UnexpectedOperand()
    }

    val := x.Rat()
    sec := new(big.Int).Div(val.Num(), val.Denom())

    if !sec.IsInt64() {
        t.UnexpectedOperand()
    }

    rem := new(big.Rat).Sub(val, new(big.Rat).SetInt(sec))
    rem.Mul(rem, new(big.Rat).SetInt64(1000000000))

    return time.Unix(sec.Int64(), new(big.Int).Div(rem.Num(), rem.Denom()).Int64()).UTC()
}

func Layout(format String) string {
    out := ""
    runes := []rune(string(format))

    for i := 0; i < len(runes); i++ {
        if runes[i] != '%' || i + 1 >= len(runes) {
            out += string(runes[i])
            continue
        }

        i++

        switch runes[i] {
        case 'Y':
            out += "2006"
        case 'y':
            out += "06"
        case 'm':
            out += "01"
        case 'd':
            out += "02"
        case 'e':
            out += "_2"
        case 'H':
            out += "15"
        case 'I':
            out += "03"
        case 'M':
            out += "04"
        case 'S':
            out += "05"
        case 'f':
            out += ".000000"
        case 'p':
            out += "PM"
        case 'a':
            out += "Mon"
        case 'A':
            out += "Monday"
        case 'b':
            out += "Jan"
        case 'B':
            out += "January"
        case 'j':
            out += "002"
        case 'Z':
            out += "MST"
        case 'z':
            out += "-0700"
        default:
            out += string(runes[i])
        }
    }

    return out
}

func (t *Token) Ticks(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Ticks(x.Run(), b)
    case *Variable:
        return t.Ticks(x.Value(), b)
    case Hash:
        out := Hash{ }

        for key, val := range x {
            out[key] = t.Ticks(val, b)
        }

        return out
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Ticks(val, b))
        }

        return out
    case Boolean:
        return t.Ticks(x.Number(), b)
    case Null:
        return t.Ticks(NewNumber(0), b)
    case Number:
        switch t.lit {
        case "seconds":
            return x
        case "minutes":
            return t.MultiplyNumber(x, NewNumber(60))
        case "hours":
            return t.MultiplyNumber(x, NewNumber(3600))
        case "days":
            return t.MultiplyNumber(x, NewNumber(86400))
        case "weeks":
            return t.MultiplyNumber(x, NewNumber(604800))
        }

        tm := t.Ticked(x)

        switch t.lit {
        case "year":
            return NewNumber(tm.Year())
        case "month":
            return NewNumber(int(tm.Month()))
        case "day":
            return NewNumber(tm.Day())
        case "hour":
            return NewNumber(tm.Hour())
        case "minute":
            return NewNumber(tm.Minute())
        case "second":
            return NewNumber(tm.Second())
        case "date":
            if b == nil {
                return String(tm.Format(time.RFC3339))
            }

            return String(tm.Format(Layout(Stringify(b))))
        }
    }

    return t.TypeMismatch(a, b)
}
