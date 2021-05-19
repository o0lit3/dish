package main
import("fmt")

func Member(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(string(y))
        default:
            return x.Member(fmt.Sprintf("%v", y))
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(int(y.Number()))
        case Number:
            return x.Member(int(y))
        case Boolean:
            return x.Member(int(y.Number()))
        case Null:
            return x.Member(0)
        }
    case String:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(int(y.Number()))
        case Number:
            return x.Member(int(y))
        case Boolean:
            return x.Member(int(y.Number()))
        case Null:
            return x.Member(0)
        }
    }

    return Null { }
}

func (a Hash) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a Hash) Member(b string) interface{} {
    out := a[b]

    if out != nil {
        return out
    }

    return Null { }
}

func (a Array) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a Array) Member(b int) interface{} {
    if b < len(a) {
        return a[b]
    }

    return Null { }
}

func (a String) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a String) Member(b int) interface{} {
    if b < len(a) {
        return String(string(string(a)[b]))
    }

    return Null { }
}
