package main

func Multiply(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Multiply(x.Run(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return x.Map(y)
        default:
            return Multiply(x.Array(), y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return x.Map(y)
        case String:
            return x.Multiply(y.Number())
        case Number:
            return x.Multiply(y)
        case Boolean:
            return x.Multiply(y.Number())
        case Null:
            return x.Multiply(Number(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return x.Array().Map(y)
        case String:
            return x.Multiply(y.Number())
        case Number:
            return x.Multiply(y)
        case Boolean:
            return x.Multiply(y.Number())
        case Null:
            return x.Multiply(Number(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Multiply(x, y.Run())
        case Hash:
            return y.Array().Multiply(x)
        case Array:
            return y.Multiply(x)
        case String:
            return y.Multiply(x)
        case Number:
            return x * y
        case Boolean:
            return x * y.Number()
        case Null:
            return x * Number(0)
        }
    case Boolean:
        return Multiply(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Map(b *Block) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = b.Run(val, String(key))
    }

    return out
}

func (a Array) Map(b *Block) Array {
    out := Array { }

    for i, val := range a {
        out = append(out, b.Run(val, Number(i)))
    }

    return out
}

func (a Array) Multiply(b Number) Array {
    out := Array { }

    for n := 0; n < int(b); n++ {
        for _, val := range a {
            out = append(out, val)
        }
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, val := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out = append(out, val)
            } else {
                break
            }
        }
    }

    return out
}

func (a String) Multiply(b Number) String {
    out := ""

    for n := 0; n < int(b); n++ {
        out += string(a)
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, c := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return String(out)
}
