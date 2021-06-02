package main
import("strings")

func Divide(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Divide(x.Run(), b)
    case *Variable:
        return Divide(x.Value(), b)
    case Hash:
        return Divide(x.Array(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            return x.Split(y)
        case *Variable:
            return Divide(x, y.Value())
        case String:
            if y.Number() != 0 {
                return x.Divide(y.Number())
            }
        case Number:
            if y != 0 {
                return x.Divide(y)
            }
        case Boolean:
            if y {
                return x.Divide(Number(1))
            }
        }
    case String:
        switch y := b.(type) {
        case *Block:
            items := x.Array().Split(y)
            out := Array { }

            for _, item := range items {
                out = append(out, Join(item, ""))
            }

            return out
        case *Variable:
            return Divide(x, y.Value())
        case String:
            return x.Split(y)
        case Number:
            if y != 0 {
                return x.Divide(y)
            }
        case Boolean:
            if y {
                return x.Divide(Number(1))
            }
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Divide(x, y.Run())
        case *Variable:
            return Divide(x, y.Value())
        case Hash:
            if len(y) != 0 {
                return x / Number(len(y))
            }
        case Array:
            if len(y) != 0 {
                return x / Number(len(y))
            }
        case String:
            if len(y) != 0 {
                return x / Number(len(y))
            }
        case Number:
            if y != 0 {
                return x / y
            }
        case Boolean:
            if y {
                return x
            }
        }
    case Boolean:
        return Divide(x.Number(), b)
    }

    return Null { }
}

func (a Array) Split(b *Block) Array {
    items := []Array { Array { } }
    out := Array { }

    for _, val := range a {
        switch y := b.Run(val).(type) {
        case Hash:
            return a.Divide(Number(len(y)))
        case Array:
            return a.Divide(Number(len(y)))
        case Boolean:
            if y {
                items = append(items, Array { })
            } else {
                items[len(items) - 1] = append(items[len(items) - 1], val)
            }
        default:
            if Equals(val, y) {
                items = append(items, Array { })
            } else {
                items[len(items) - 1] = append(items[len(items) - 1], val)
            }
        }
    }

    for _, item := range items {
        out = append(out, Array(item))
    }

    return out
}

func (a String) Split(b String) Array {
    items := strings.Split(string(a), string(b))
    out := Array { }

    for _, item := range items {
        out = append(out, String(item))
    }

    return out
}

func (a Array) Divide(b Number) Array {
    out := Array { }
    x := int(len(a) / int(b))
    i := 0

    for len(out) < len(a) % int(b) {
        set := Array { }

        for len(set) < x + 1 {
            set = append(set, a[i])
            i = i + 1
        }

        out = append(out, set)
    }

    for len(out) < int(b) {
        set := Array { }

        for len(set) < x {
            set = append(set, a[i])
            i = i + 1
        }

        out = append(out, set)
    }

    return out
}

func (a String) Divide(b Number) Array {
    out := Array { }
    x := int(len(a) / int(b))
    i := 0

    for len(out) < len(a) % int(b) {
        set := String("")

        for len(set) < x + 1 {
            set += String(a[i])
            i = i + 1
        }

        out = append(out, set)
    }

    for len(out) < int(b) {
        set := String("")

        for len(set) < x {
            set += String(a[i])
            i = i + 1
        }

        out = append(out, set)
    }

    return out
}
