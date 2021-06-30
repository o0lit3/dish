package main
import("fmt")

func Join(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Join(x.Run(), b)
    case *Variable:
        return Join(x.Value(), b)
    case Hash:
        return Join(x.Array(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            return Join(x, y.Run())
        case *Variable:
            return Join(x, y.Value())
        case String:
            return x.Join(string(y))
        case Null:
            return x.Join("")
        default:
            return x.Join(fmt.Sprintf("%v", y))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Join(x, y.Run())
        case *Variable:
            return Join(x, y.Value())
        case String:
            return String(string(x) + string(y))
        case Null:
            return String(string(x))
        default:
            return String(string(x) + fmt.Sprintf("%v", y))
        }
    case Null:
        return Join(String(""), b)
    default:
        switch y := b.(type) {
        case *Block:
            return Join(x, y.Run())
        case *Variable:
            return Join(x, y.Value())
        case String:
            return String(fmt.Sprintf("%v", x) + string(y))
        case Null:
            return String(fmt.Sprintf("%v", x))
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
