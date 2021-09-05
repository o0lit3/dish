package main
import("strconv"; "strings")

func (t *Token) Waka(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Waka(x.Run(), b)
    case *Variable:
        return t.Waka(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.Waka(x, y.Run())
        case *Variable:
            return t.Waka(x, y.Value())
        case Hash:
            return Boolean(len(x) > len(y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.Waka(x, y.Run())
        case *Variable:
            return t.Waka(x, y.Value())
        case Array:
            return Boolean(len(x) > len(y))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.Waka(x, y.Run())
        case *Variable:
            return t.Waka(x, y.Value())
        case String:
            return Boolean(string(x) > string(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Waka(x, y.Run())
        case *Variable:
            return t.Waka(x, y.Value())
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

            return Boolean(x.val.Cmp(y.val) > 0)
        case Boolean:
            return t.Waka(x, y.Number())
        case Null:
            return t.Waka(x, NewNumber(0))
        }
    case Boolean:
        return t.Waka(x.Number(), b)
    case Null:
        return t.Waka(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) WakaBars(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.WakaBars(x.Run(), b)
    case *Variable:
        return t.WakaBars(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.WakaBars(x, y.Run())
        case *Variable:
            return t.WakaBars(x, y.Value())
        case Hash:
            return Boolean(len(x) >= len(y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.WakaBars(x, y.Run())
        case *Variable:
            return t.WakaBars(x, y.Value())
        case Array:
            return Boolean(len(x) >= len(y))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.WakaBars(x, y.Run())
        case *Variable:
            return t.WakaBars(x, y.Value())
        case String:
            return Boolean(string(x) >= string(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.WakaBars(x, y.Run())
        case *Variable:
            return t.WakaBars(x, y.Value())
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

            return Boolean(x.val.Cmp(y.val) >= 0)
        case Boolean:
            return t.WakaBars(x, y.Number())
        case Null:
            return t.WakaBars(x, NewNumber(0))
        }
    case Boolean:
        return t.WakaBars(x.Number(), b)
    case Null:
        return t.WakaBars(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) WakaWaka(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.WakaWaka(x.Run(), b)
    case *Variable:
        return t.WakaWaka(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.WakaWaka(x, y.Run())
        case *Variable:
            return t.WakaWaka(x, y.Value())
        case Hash:
            if t.lit != ">>" {
                t.TypeMismatch(x, y)
            }

            return t.ExtendHash(x, y)
        case Null:
            return t.WakaWaka(x, Hash{ })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.WakaWaka(x, y.Run())
        case *Variable:
            return t.WakaWaka(x, y.Value())
        case Hash:
            return t.WakaWaka(x, y.Array())
        case Array:
            if t.lit != ">>" && t.lit != "unshift" && t.lit != "prepend" {
                t.TypeMismatch(x, y)
            }

            return t.UnshiftArray(x, y)
        case Number:
            if t.lit != ">>" && t.lit != "lpad" && t.lit != "ltrunc" {
                t.TypeMismatch(x, y)
            }

            if t.lit == "ltrunc" {
                return t.LtruncArray(x, y)
            } else if y.val.Cmp(NewNumber(0).val) < 0 {
                return t.LtruncArray(x, t.NegateNumber(y))
            }

            return t.LpadArray(x, y)
        case Null:
            if t.lit != ">>" && t.lit != "prepend" && t.lit != "lpad" && t.lit != "ltrunc" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.WakaWaka(x, y.Run())
        case *Variable:
            return t.WakaWaka(x, y.Value())
        case String:
            if t.lit != ">>" && t.lit != "prepend" {
                t.TypeMismatch(x, y)
            }

            return t.PrependString(x, y)
        case Number:
            if t.lit != ">>" && t.lit != "lpad" && t.lit != "ltrunc" {
                t.TypeMismatch(x, y)
            }

            if t.lit == "ltrunc" {
                return t.LtruncString(x, y)
            } else if y.val.Cmp(NewNumber(0).val) < 0 {
                return t.LtruncString(x, t.NegateNumber(y))
            }

            return t.LpadString(x, y)
        case Null:
            if t.lit != ">>" && t.lit != "prepend" && t.lit != "lpad" && t.lit != "ltrunc" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.WakaWaka(x, y.Run())
        case *Variable:
            return t.WakaWaka(x, y.Value())
        case Number:
            if t.lit != ">>" && t.lit != "rshift" {
                t.TypeMismatch(x, y)
            }

            return t.RshiftNumber(x, y)
        case Boolean:
            return t.WakaWaka(x, y.Number())
        case Null:
            return t.WakaWaka(x, NewNumber(0))
        }
    case Boolean:
        return t.WakaWaka(x.Number(), b)
    case Null:
        switch y := b.(type) {
        case *Block:
            return t.WakaWaka(x, y.Run())
        case *Variable:
            return t.WakaWaka(x, y.Value())
        case Hash:
            return t.WakaWaka(Hash{ }, y)
        case Array:
            return t.WakaWaka(Array{ }, y)
        case String:
            return t.WakaWaka(String(""), y)
        default:
            return t.WakaWaka(NewNumber(0), y)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopWaka(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopWaka(x.Run())
    case *Variable:
        return t.TopWaka(x.Value())
    case Hash:
        return t.TopWaka(x.Array())
    case Array:
        if t.lit != ">" && t.lit != "maximum" && t.lit != "max" {
            t.TypeMismatch(x, nil)
        }

        return t.MaxFromArray(x)
    case String:
        if t.lit != ">" && t.lit != "uppercase" && t.lit != "upcase" && t.lit != "uc" {
            t.TypeMismatch(x, nil)
        }

        return t.UpperString(x)
    case Number:
        if t.lit != ">" && t.lit != "ceiling" && t.lit != "ceil" {
            t.TypeMismatch(x, nil)
        }

        return t.CeilNumber(x)
    case Boolean:
        return t.TopWaka(x.Number())
    case Null:
        return t.TopWaka(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) TopWakaWaka(a interface{}) (interface{}, interface{}) {
    switch x := a.(type) {
    case *Block:
        return t.TopWakaWaka(x.Run())
    case *Variable:
        return t.TopWakaWaka(x.Value())
    case Array:
        if t.lit != ">>" && t.lit != "last" && t.lit != "pop" {
            t.TypeMismatch(x, nil)
        }

        return t.PopArray(x)
    case String:
        if t.lit != ">>" && t.lit != "last" && t.lit != "pop" {
            t.TypeMismatch(x, nil)
        }

        return t.PopString(x)
    case Number:
        if t.lit != ">>" && t.lit != "lsb" && t.lit != "pop" {
            t.TypeMismatch(x, nil)
        }

        return t.PopNumber(x)
    case Boolean:
        return t.PopNumber(x.Number())
    case Null:
        return t.PopNumber(NewNumber(0))
    }

    return a, t.TypeMismatch(a, nil)
}

func (t *Token) UnshiftArray(x Array, y Array) Array {
    out := Array { }

    for _, val := range y {
        out = append(out, val)
    }

    for _, val := range x {
        out = append(out, val)
    }

    return out
}

func (t *Token) LpadArray(x Array, y Number) Array {
    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.LtruncArray(x, t.NegateNumber(y))
    }

    out := x
    i := 0

    for i < y.Int() {
        out = append([]interface{}{ Null{ } }, out...)
        i++
    }

    return out
}

func (t *Token) LpadString(x String, y Number) String {
    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.LtruncString(x, t.NegateNumber(y))
    }

    out := string(x)
    i := 0

    for i < y.Int() {
        out = " " + out
        i++
    }

    return String(out)
}

func (t *Token) LtruncArray(x Array, y Number) Array {
    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.LpadArray(x, t.NegateNumber(y))
    }

    out := x
    i := y.Int()

    for i > 0 {
        out = out[1:]
        i--
    }

    return out
}

func (t *Token) LtruncString(x String, y Number) String {
    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.LpadString(x, t.NegateNumber(y))
    }

    out := x
    i := y.Int()

    for i > 0 {
        out = out[1:]
        i--
    }

    return out
}

func (t *Token) PrependString(x String, y String) String {
    return String(string(y) + string(x))
}

func (t *Token) RshiftNumber(x Number, y Number) Number {
    if x.inf == INF || x.inf == -INF {
        return NewNumber(0)
    }

    if y.inf == INF || y.inf == -INF {
        return x
    }

    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.LshiftNumber(x, t.NegateNumber(y))
    }

    return NewNumber(x.Int() >> uint(y.Int()))
}

func (t *Token) MaxFromArray(x Array) interface{} {
    var out interface{}

    for _, val := range x {
        if out == nil {
            out = val
        } else if b, ok := t.Waka(val, out).(Boolean); bool(b) && ok {
            out = val
        }
    }

    return out
}

func (t *Token) UpperString(x String) String {
    return String(strings.ToUpper(string(x)))
}

func (t *Token) CeilNumber(x Number) Number {
    if x.inf == INF || x.inf == -INF || x.val.IsInt() {
        return x
    }

    if x.val.Cmp(NewNumber(0).val) < 0 {
        return t.NegateNumber(t.FloorNumber(t.NegateNumber(x)))
    }

    return t.AddNumber(t.FloorNumber(x), NewNumber(1)).(Number)
}

func (t *Token) PopArray(x Array) (interface{}, interface{}) {
    if len(x) > 0 {
        val := x[len(x) - 1]
        x = x[:len(x) - 1]
        return val, x
    }

    return Null { }, x
}

func (t *Token) PopString(x String) (interface{}, interface{}) {
    if len(x) > 0 {
        val := String(string(x[len(x) - 1]))
        x = x[:len(x) - 1]
        return val, x
    }

    return Null { }, x
}

func (t *Token) PopNumber(x Number) (interface{}, interface{}) {
    if x.inf == INF || x.inf == -INF {
        return NewNumber(0), x
    }

    bin := strconv.FormatInt(int64(x.Int()), 2)
    first, _ := strconv.Atoi(string(bin[len(bin) - 1]))
    rem, _ := strconv.ParseInt(string(bin[:len(bin) - 1]), 2, 64)
    return NewNumber(first), NewNumber(int(rem))
}
