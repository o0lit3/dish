package main
import("strings"; "unicode")

func (t *Token) Dash(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Dash(x.Run(), b)
    case *Variable:
        return t.Dash(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "-" && t.lit != "-=" && t.lit != "reduce" {
                    t.TypeMismatch(x, y)
                }

                return t.ReduceHash(x, y)
            }

            return t.Dash(x, y.Run())
        case *Variable:
            return t.Dash(x, y.Value())
        case Hash:
            if t.lit != "-" && t.lit != "-=" && t.lit != "remove" && t.lit != "delete" && t.lit != "del" {
                t.TypeMismatch(x, y)
            }

            return t.RemoveFromHash(x, y)
        case Null:
            return t.Dash(x, Hash{ })
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "-" && t.lit != "-=" && t.lit != "reduce" {
                    t.TypeMismatch(x, y)
                }

                return t.ReduceArray(x, y)
            }

            return t.Dash(x, y.Run())
        case *Variable:
            return t.Dash(x, y.Value())
        case Array:
            if t.lit != "-" && t.lit != "-=" && t.lit != "remove" && t.lit != "delete" && t.lit != "del" {
                t.TypeMismatch(x, y)
            }

            return t.RemoveFromArray(x, y)
        case Number:
            if t.lit != "-" && t.lit != "-=" {
                t.TypeMismatch(x, y)
            }

            if y.val.Cmp(NewNumber(0).val) < 0 {
                return t.RpadArray(x, t.NegateNumber(y))
            }

            return t.RtruncArray(x, y)
        case Null:
            if t.lit != "-" && t.lit != "-=" && t.lit != "remove" && t.lit != "delete" && t.lit != "del" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "-" && t.lit != "-=" && t.lit != "reduce" {
                    t.TypeMismatch(x, y)
                }

                return t.ReduceString(x, y)
            }

            return t.Dash(x, y.Run())
        case *Variable:
            return t.Dash(x, y.Value())
        case String:
            if t.lit != "-" && t.lit != "-=" && t.lit != "remove" && t.lit != "delete" && t.lit != "del" {
                t.TypeMismatch(x, y)
            }

            return t.RemoveFromString(x, y)
        case Number:
            if t.lit != "-" && t.lit != "-=" && t.lit != "decrease" {
                t.TypeMismatch(x, y)
            }

            if y.val.Cmp(NewNumber(0).val) < 0 {
                return t.IncreaseString(x, t.NegateNumber(y))
            }

            return t.DecreaseString(x, y)
        case Boolean:
            return t.Dash(x, y.Number())
        case Null:
            if t.lit != "-" && t.lit != "-=" && t.lit != "remove" && t.lit != "delete" && t.lit != "del" && t.lit != "decrease" {
                t.TypeMismatch(x, y)
            }

            return x
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Dash(x, y.Run())
        case *Variable:
            return t.Dash(x, y.Value())
        case Array:
            if t.lit != "-" && t.lit != "-=" {
                t.TypeMismatch(x, y)
            }

            if x.val.Cmp(NewNumber(0).val) < 0 {
                return t.LpadArray(y, t.NegateNumber(x))
            }

            return t.LtruncArray(y, x)
        case String:
            if t.lit != "-" && t.lit != "-=" {
                t.TypeMismatch(x, y)
            }

            if x.val.Cmp(NewNumber(0).val) < 0 {
                return t.IncreaseString(y, t.NegateNumber(x))
            }

            return t.DecreaseString(y, x)
        case Number:
            if t.lit != "-" && t.lit != "-=" && t.lit != "subtract" && t.lit != "sub" {
                t.TypeMismatch(x, y)
            }

            return t.SubtractNumber(x, y)
        case Boolean:
            return t.Dash(x, y.Number())
        case Null:
            return t.Dash(x, NewNumber(0))
        }
    case Boolean:
        return t.Dash(x.Number(), b)
    case Null:
        switch y := b.(type) {
        case *Block:
            return t.Dash(x, y.Run())
        case *Variable:
            return t.Dash(x, y.Value())
        case Hash:
            return t.Dash(Hash{ }, y)
        case Array:
            return t.Dash(Array{ }, y)
        case String:
            return t.Dash(String(""), y)
        default:
            return t.Dash(NewNumber(0), y)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopDash(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopDash(x.Run())
    case *Variable:
        return t.TopDash(x.Value())
    case Hash:
        return t.TopDash(x.Array())
    case Array:
        if t.lit != "-" && t.lit != "negsum" {
            t.TypeMismatch(a, nil)
        }

        return t.TopDash(t.SumArray(x))
    case String:
        if t.lit != "-" && t.lit != "separate" {
            t.TypeMismatch(a, nil)
        }

        return x.Array()
    case Number:
        if t.lit != "-" && t.lit != "negate" && t.lit != "neg" {
            t.TypeMismatch(a, nil)
        }

        return t.NegateNumber(x)
    case Boolean:
        return t.TopDash(x.Number())
    case Null:
        return t.TopDash(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) TopDoubleDash(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopDoubleDash(x.Run())
    case *Variable:
        return t.TopDoubleDash(x.Value())
    case Array:
        return t.RtruncArray(x, NewNumber(1))
    case String:
        return t.DecreaseString(x, NewNumber(1))
    case Number:
        return t.SubtractNumber(x, NewNumber(1))
    case Boolean:
        return t.TopDoubleDash(x.Number())
    case Null:
        return t.TopDoubleDash(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) ReduceHash(x Hash, y *Block) interface{} {
    var out interface{} = Null{ }

    for key, val := range x {
        out = y.Context(x).Run(out, val, String(key))
    }

    return out
}

func (t *Token) ReduceArray(x Array, y *Block) interface{} {
    var out interface{} = Null{ }

    for i, val := range x {
        out = y.Context(x).Run(out, val, NewNumber(i))
    }

    return out
}

func (t *Token) ReduceString(x String, y *Block) interface{} {
    var out interface{} = Null{ }

    for i, c := range x {
        out = y.Context(x).Run(out, String(string(c)), NewNumber(i))
    }

    return out
}

func (t *Token) RemoveFromHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key, val := range x {
        out[key] = val
    }

    for key, _ := range y {
        if _, ok := out[key]; ok {
            delete(out, key)
        }
    }

    return out
}

func (t *Token) RemoveFromArray(x Array, y Array) Array {
    out := Array { }

    for _, a := range x {
        found := false

        for i, b := range y {
            if Equals(a, b) {
                found = true
                y = append(y[:i], y[i + 1:]...)
                break
            }
        }

        if !found {
            out = append(out, a)
        }
    }

    return out
}

func (t *Token) RemoveFromString(x String, y String) String {
    return String(strings.Replace(string(x), string(y), "", 1))
}

func (t *Token) DecreaseString(x String, y Number) String {
    i := 0
    out := ""
    carry := y.Int()

    for i < len(x) {
        c := x[len(x) - i - 1]

        if carry > 0 {
            switch {
            case unicode.IsLetter(c) && unicode.IsUpper(c):
                if int(c) - carry < int('A') {
                    if (int('A') - int(c) + carry) % 26 == 0 {
                        out = "A" + out
                    } else {
                        out = string(int('Z') + 1 - ((int('A') - int(c) + carry) % 26)) + out
                    }

                    carry = (int('Z') - int(c) + carry) / 26
                } else {
                    out = string(int(c) - carry) + out
                    carry = 0
                }
            case unicode.IsLetter(c) && unicode.IsLower(c):
                if int(c) - carry < int('a') {
                    if (int('a') - int(c) + carry) % 26 == 0 {
                        out = "a" + out
                    } else {
                        out = string(int('z') + 1 - ((int('a') - int(c) + carry) % 26)) + out
                    }

                    carry = (int('z') - int(c) + carry) / 26
                } else {
                    out = string(int(c) - carry) + out
                    carry = 0
                }
            case unicode.IsDigit(c):
                if int(c) - carry < int('0') {
                    if (int('0') - int(c) + carry) % 10 == 0 {
                        out = "0" + out
                    } else {
                        out = string(int('9') + 1 - (int('0') - int(c) + carry) % 10) + out
                    }

                    carry = (int('9') - int(c) + carry) / 10
                } else {
                    out = string(int(c) - carry) + out
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

    for carry > 0 && len(out) > 0 {
        carry = int('z') - int(out[0])
        out = out[1:]
    }

    return String(out)
}

func (t *Token) SubtractNumber(x Number, y Number) interface{} {
    if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
        return Null { }
    }

    if x.inf == INF && y.inf == -INF {
        return Number{ inf: INF }
    }

    if x.inf == -INF && y.inf == INF {
        return Number{ inf: -INF }
    }

    return Number{ val: NewNumber(0).val.Sub(x.val, y.val) }
}

func (t *Token) NegateNumber(x Number) Number {
    if x.inf == INF {
        return Number{ inf: -INF }
    }

    if x.inf == -INF {
        return Number{ inf: INF }
    }

    return Number{ val: NewNumber(0).val.Neg(x.val) }
}
