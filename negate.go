package main

func Negate(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Negate(x.Run())
    case *Variable:
        return Negate(x.Value())
    case Hash:
        return x.Array().Negate()
    case Array:
        return x.Negate()
    case String:
        return Negate(x.Number())
    case Number:
        if x.inf == INF {
            return Number{ inf: -INF }
        }

        if x.inf == -INF {
            return Number{ inf: INF }
        }

        return Number{ val: NewNumber(0).val.Neg(x.val) }
    case Boolean:
        return Negate(x.Number())
    case Null:
        return NewNumber(0)
    }

    return Null { }
}

func (a Array) Negate() Number {
    out := NewNumber(0)

    for _, val := range a {
        if x, ok := Negate(val).(Number); ok {
            out = Number{ val: out.val.Add(out.val, x.val) }
        }
    }

    return out
}
