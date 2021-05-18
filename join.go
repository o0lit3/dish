package main
import("fmt")

func Join(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return Join(x, y.Run())
        case String:
            return x.Array().Join(string(y))
        default:
            return x.Array().Join(fmt.Sprintf("%v", y))
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return Join(x, y.Run())
        case String:
            return x.Join(string(y))
        default:
            return x.Join(fmt.Sprintf("%v", y))
        }
    case String:
        switch y := b.(type) {
        case Interpreter:
            return Join(x, y.Run())
        case String:
            return String(string(x) + string(y))
        default:
            return String(string(x) + fmt.Sprintf("%v", y))
        }
    default:
        switch y := b.(type) {
        case Interpreter:
            return Join(x, y.Run())
        case String:
            return String(fmt.Sprintf("%v", x) + string(y))
        default:
            return String(fmt.Sprintf("%v", x) + fmt.Sprintf("%v", y))
        }
    }
}

func (a Array) Join(b string) String {
    out := ""

    for i, val := range a {
        addend := ""

        switch x := val.(type) {
        case String:
            addend = string(x)
        default:
            addend = fmt.Sprintf("%v", x)
        }

        if i > 0 {
            out += b + addend
        } else {
            out += addend
        }
    }

    return String(out)
}
