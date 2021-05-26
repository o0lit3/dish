package main

func Divide(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Divide(x.Run(), b)
    case Hash:
        return Divide(x.Array(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            return Divide(x, y.Run())
        case Hash:
            return Divide(x, y.Array())
        case Array:
            if len(y) != 0 {
                return x.Divide(Number(len(y)))
            }
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
            return Divide(x, y.Run())
        case Hash:
            return Divide(x, y.Array())
        case Array:
            if len(y) != 0 {
                return x.Divide(Number(len(y)))
            }
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
    case Number:
        switch y := b.(type) {
        case *Block:
            return Divide(x, y.Run())
        case Hash:
            return Divide(x, y.Array())
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

func (a Array) Divide(b Number) interface{} {
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

func (a String) Divide(b Number) interface{} {
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
