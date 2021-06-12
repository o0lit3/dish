package main

func Mod(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Mod(x.Run(), b)
    case *Variable:
        return Mod(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return x.Filter(y)
            }

            return Mod(x, y.Run())
        case *Variable:
            return Mod(x, y.Value())
        default:
            return Mod(x.Array(), b)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return x.Filter(y)
            }

            return Mod(x, y.Run())
        case *Variable:
            return Mod(x, y.Value())
        case String:
            return x.Mod(y.Number())
        case Number:
            return x.Mod(y)
        case Boolean:
            return x.Mod(y.Number())
        case Null:
            return x.Mod(NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return Join(x.Array().Filter(y), String(""))
            }

            return Mod(x, y.Run())
        case *Variable:
            return Mod(x, y.Value())
        case String:
            return x.Mod(y.Number())
        case Number:
            return x.Mod(y)
        case Boolean:
            return x.Mod(y.Number())
        case Null:
            return x.Mod(NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Mod(x, y.Run())
        case *Variable:
            return Mod(x, y.Value())
        case Hash:
            if len(y) != 0 {
                return NewNumber(x.Int() % len(y))
            }
        case Array:
            if len(y) != 0 {
                return NewNumber(x.Int() % len(y))
            }
        case String:
            if len(y) != 0 {
                return NewNumber(x.Int() % len(y))
            }
        case Number:
            if y.val.Cmp(NewNumber(0).val) != 0 {
                return NewNumber(x.Int() % y.Int())
            }
        case Boolean:
            if y {
                return NewNumber(0)
            }
        }
    case Boolean:
        return Mod(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Filter(b *Block) Hash {
    out := Hash { }

    for key, val := range a {
        if Not(b.Run(val, String(key))) {
            continue
        }

        out[key] = val
    }

    return out
}

func (a Array) Filter(b *Block) Array {
    out := Array { }

    for i, val := range a {
        if Not(b.Run(val, NewNumber(i))) {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (a Array) Mod(b Number) interface{} {
    if b.val.Cmp(NewNumber(0).val) == 0 {
        return Null { }
    }

    out := Array { }

    for i, val := range a{
        if i % b.Int() > 0 {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (a String) Mod(b Number) interface{} {
    if b.val.Cmp(NewNumber(0).val) == 0 {
        return Null { }
    }

    out := ""

    for i, c := range a {
        if i % b.Int() > 0 {
            continue
        }

        out += string(c)
    }

    return String(out)
}
