package main

func Mod(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Mod(x.Run(), b)
    case Hash:
        return Mod(x.Array(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            return Mod(x, y.Run())
        case Hash:
            return x.Mod(Number(len(y)))
        case Array:
            return x.Mod(Number(len(y)))
        case String:
            return x.Mod(y.Number())
        case Number:
            return x.Mod(y)
        case Boolean:
            return x.Mod(y.Number())
        case Null:
            return x.Mod(Number(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Mod(x, y.Run())
        case Hash:
            return x.Mod(Number(len(y)))
        case Array:
            return x.Mod(Number(len(y)))
        case String:
            return x.Mod(y.Number())
        case Number:
            return x.Mod(y)
        case Boolean:
            return x.Mod(y.Number())
        case Null:
            return x.Mod(Number(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Mod(x, y.Run())
        case Hash:
            return Number(int(x) % len(y))
        case Array:
            return Number(int(x) % len(y))
        case String:
            return Number(int(x) % len(y))
        case Number:
            return Number(int(x) % int(y))
        case Boolean:
            if y {
                return Number(0)
            }
        }
    case Boolean:
        return Mod(x.Number(), b)
    }

    return Null { }
}

func (a Array) Mod(b Number) interface{} {
    if b == 0 {
        return Null { }
    }

    out := Array { }

    for i, val := range a{
        if i % int(b) > 0 {
            continue
        }

        out = append(out, val)
    }

    return out
}

func (a String) Mod(b Number) interface{} {
    if b == 0 {
        return Null { }
    }

    out := ""

    for i, c := range a {
        if i % int(b) > 0 {
            continue
        }

        out += string(c)
    }

    return String(out)
}
