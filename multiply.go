package main

func Multiply(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return x.Array().Map(y)
        case Array:
            return x.Array().Multiply(Number(len(y)))
        case String:
            return x.Array().Multiply(y.Number())
        case Number:
            return x.Array().Multiply(y)
        case Boolean:
            return x.Array().Multiply(y.Number())
        case Null:
            return x.Array().Multiply(Number(0))
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return x.Map(y)
        case Array:
            return x.Multiply(Number(len(y)))
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
        case Interpreter:
            return x.Array().Map(y)
        case Array:
            return x.Multiply(Number(len(y)))
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
        case Interpreter:
            return NewArray(int(x)).Map(y)
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
        switch y := b.(type) {
        case Interpreter:
            return NewArray(int(x.Number())).Map(y)
        case Array:
            return y.Multiply(x.Number())
        case String:
            return y.Multiply(x.Number())
        case Number:
            return y * x.Number()
        case Boolean:
            return y.Number() * x.Number()
        case Null:
            return x.Number() * Number(0)
        }
    }

    return Null { }
}

func (a Array) Map(b Interpreter) Array {
    out := Array { }

    for i, val := range a {
        out = append(out, b.IndexRun(Number(i), val))
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
