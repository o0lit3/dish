package main
import("unicode")

func (t *Token) Cross(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Cross(x.Run(), b)
    case *Variable:
        return t.Cross(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "+" && t.lit != "+=" && t.lit != "aggregate" {
                    t.TypeMismatch(x, y)
                }

                return t.AggregateHash(x, y)
            }

            return t.Cross(x, y.Run())
        case *Variable:
            return t.Cross(x, y.Value())
        case Hash:
            if t.lit != "+" && t.lit != "+=" && t.lit != "concat" {
                t.TypeMismatch(x, y)
            }

            return t.ConcatHash(x, y)
        case Null:
            return t.Cross(x, Hash{ })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "+" && t.lit != "+=" && t.lit != "aggregate" {
                    t.TypeMismatch(x, y)
                }

                return t.AggregateArray(x, y)
            }

            return t.Cross(x, y.Run())
        case *Variable:
            return t.Cross(x, y.Value())
        case Array:
            if t.lit != "+" && t.lit != "+=" && t.lit != "concat" {
                t.TypeMismatch(x, y)
            }

            return t.ConcatArray(x, y)
        case Number:
            if t.lit != "+" && t.lit != "+=" && t.lit != "pad" {
                t.TypeMismatch(x, y)
            }

            if y.val.Cmp(NewNumber(0).val) < 0 {
                return t.RtruncArray(x, t.NegateNumber(y))
            }

            return t.RpadArray(x, y)
        case Boolean:
            return t.Cross(x, y.Number())
        case Null:
            if t.lit != "+" && t.lit != "+=" && t.lit != "concat" && t.lit != "pad" && t.lit != "trunc" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "+" && t.lit != "+=" && t.lit != "aggregate" {
                    t.TypeMismatch(x, y)
                }

                return t.AggregateString(x, y)
            }

            return t.Cross(x, y.Run())
        case *Variable:
            return t.Cross(x, y.Value())
        case String:
            if t.lit != "+" && t.lit != "+=" && t.lit != "concat" {
                t.TypeMismatch(x, y)
            }

            return t.ConcatString(x, y)
        case Number:
            if t.lit != "+" && t.lit != "+=" && t.lit != "increase" {
                t.TypeMismatch(x, y)
            }

            return t.IncreaseString(x, y)
        case Boolean:
            return t.Cross(x, y.Number())
        case Null:
            if t.lit != "+" && t.lit != "+=" && t.lit != "concat" && t.lit != "increase" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Cross(x, y.Run())
        case *Variable:
            return t.Cross(x, y.Value())
        case Array:
            if t.lit != "+" && t.lit != "+=" && t.lit != "pad" {
                t.TypeMismatch(x, y)
            }

            if x.val.Cmp(NewNumber(0).val) < 0 {
                return t.RtruncArray(y, t.NegateNumber(x))
            }

            return t.RpadArray(y, x)
        case String:
            if t.lit != "+" && t.lit != "+=" && t.lit != "increase" {
                t.TypeMismatch(x, y)
            }

            return t.IncreaseString(y, x)
        case Number:
            if t.lit != "+" && t.lit != "+=" && t.lit != "add" {
                t.TypeMismatch(x, y)
            }

            return t.AddNumber(x, y)
        case Boolean:
            return t.Cross(x, y.Number())
        case Null:
            return t.Cross(x, NewNumber(0))
        }
    case Boolean:
        return t.Cross(x.Number(), b)
    case Null:
        switch y := b.(type) {
        case *Block:
            return t.Cross(x, y.Run())
        case *Variable:
            return t.Cross(x, y.Value())
        case Hash:
            return t.Cross(Hash{ }, y)
        case Array:
            return t.Cross(Array{ }, y)
        case String:
            return t.Cross(String(""), y)
        default:
            return t.Cross(NewNumber(0), y)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopCross(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopCross(x.Run())
    case *Variable:
        return t.TopCross(x.Value())
    case Hash:
        if t.lit == "num" && x["num"] != nil {
            return x["num"]
        }

        return t.TopCross(x.Array())
    case Array:
        if t.lit != "+" && t.lit != "sum" && t.lit != "concat" {
            t.TypeMismatch(x, nil)
        }

        return t.SumArray(x)
    case String:
        return t.TopCross(x.Number())
    case Number:
        if t.lit != "+" && t.lit != "number" && t.lit != "num" {
            t.TypeMismatch(x, nil)
        }

        return x
    case Boolean:
        return t.TopCross(x.Number())
    case Null:
        return t.TopCross(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) TopDoubleCross(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopDoubleCross(x.Run())
    case *Variable:
        return t.TopDoubleCross(x.Value())
    case Array:
        return t.RpadArray(x, NewNumber(1))
    case String:
        return t.IncreaseString(x, NewNumber(1))
    case Number:
        return t.AddNumber(x, NewNumber(1))
    case Boolean:
        return t.TopDoubleCross(x.Number())
    case Null:
        return t.TopDoubleCross(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) AggregateHash(x Hash, y *Block) interface{} {
    var out interface{} = Null{ }

    for key, val := range x {
        out = y.Context(x).Run(out, val, String(key))
    }

    return out
}

func (t *Token) AggregateArray(x Array, y *Block) interface{} {
    var out interface{} = Null{ }

    for i, val := range x {
        out = y.Context(x).Run(out, val, NewNumber(i))
    }

    return out
}

func (t *Token) AggregateString(x String, y *Block) interface{} {
    var out interface{} = Null{ }

    for i, c := range x {
        out = y.Context(x).Run(out, String(string(c)), NewNumber(i))
    }

    return out
}

func (t *Token) ConcatHash(x Hash, y Hash) Hash {
    out := Hash{ }

    for key, val := range x {
        out[key] = val
    }

    for key, val := range y {
        out[key] = val
    }

    return out
}

func (t *Token) ConcatArray(x Array, y Array) Array {
    out := Array{ }

    for _, val := range x {
        out = append(out, val)
    }

    for _, val := range y {
        out = append(out, val)
    }

    return out
}

func (t *Token) ConcatString(x String, y String) String {
    return append(x, y...)
}

func (t *Token) IncreaseString(x String, y Number) String {
    if y.val.Cmp(NewNumber(0).val) < 0 {
        return t.DecreaseString(x, t.NegateNumber(y))
    }

    i := 0
    out := ""
    carry := y.Int()

    for i < len(x) {
        c := x[len(x) - i - 1]

        if carry > 0 {
            switch {
            case unicode.IsLetter(c) && unicode.IsUpper(c):
                if int(c) + carry > int('Z') {
                    if (int(c) + carry - int('Z')) % 26 == 0 {
                        out = "Z" + out
                    } else {
                        out = string(int('A') - 1 + ((int(c) + carry - int('Z')) % 26)) + out
                    }

                    carry = (int(c) + carry - int('A')) / 26

                    if carry > 0 && i == len(x) - 1 {
                        x = append([]rune{'A'}, x...)
                        carry--
                    }
                } else {
                    out = string(int(c) + carry) + out
                    carry = 0
                }
            case unicode.IsLetter(c) && unicode.IsLower(c):
                if int(c) + carry > int('z') {
                    if (int(c) + carry - int('z')) % 26 == 0 {
                        out = "z" + out
                    } else {
                        out = string(int('a') - 1 + ((int(c) + carry - int('z')) % 26)) + out
                    }

                    carry = (int(c) + carry - int('a')) / 26

                    if carry > 0 && i == len(x) - 1 {
                        x = append([]rune{'a'}, x...)
                        carry--
                    }

                } else {
                    out = string(int(c) + carry) + out
                    carry = 0
                }
            case unicode.IsDigit(c):
                if int(c) + carry > int('9') {
                    if (int(c) + carry - int('9') - 1) % 10 == 0 {
                        out = "0" + out
                    } else {
                        out = string(int('0') + (int(c) + carry - int('9') - 1) % 10) + out
                    }

                    carry = (int(c) + carry - int('0')) / 10

                    if carry > 0 && i == len(x) - 1 {
                        x = append([]rune{'0'}, x...)
                    }
                } else {
                    out = string(int(c) + carry) + out
                    carry = 0
                }
            default:
                out = string(c) + out
                carry = 0
            }
        } else {
            out = string(c) + out
        }

        i++
    }

    return String(out)
}

func (t *Token) AddNumber(x Number, b interface{}) interface{} {
    switch y := b.(type) {
    case *Block:
        return t.AddNumber(x, y.Run())
    case *Variable:
        return t.AddNumber(x, y.Value())
    case Number:
        if (x.inf == INF && y.inf == -INF) || (x.inf == -INF && y.inf == INF) {
            return Null{ }
        }

        if x.inf == INF || y.inf == INF {
            return Number{ inf: INF }
        }

        if x.inf == -INF || y.inf == -INF {
            return Number{ inf: -INF }
        }

        return Number{ val: NewNumber(0).val.Add(x.val, y.val) }
    case Boolean:
        return t.AddNumber(x, y.Number())
    case Null:
        return t.AddNumber(x, NewNumber(0))
    }

    return t.TypeMismatch(x, b)
}

func (t *Token) SumArray(x Array) interface{} {
    out := NewNumber(0)

    for _, val := range x {
        switch val := val.(type) {
        case Number:
            out = t.AddNumber(out, val).(Number)
        case Null:
            out = out
        default:
            t.TypeMismatch(out, val)
        }
    }

    return out
}
