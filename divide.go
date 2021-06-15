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
            if len(y.args) > 0 {
                return x.Split(y)
            }

            return Divide(x, y.Run())
        case *Variable:
            return Divide(x, y.Value())
        case String:
            if y.Number().val.Cmp(NewNumber(0).val) != 0 {
                return x.Divide(y.Number())
            }
        case Number:
            if y.val.Cmp(NewNumber(0).val) != 0 {
                return x.Divide(y)
            }
        case Boolean:
            if y {
                return x.Divide(NewNumber(1))
            }
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                items := x.Array().Split(y)
                out := Array { }

                for _, item := range items {
                    out = append(out, Join(item, ""))
                }

                return out
            }

            return Divide(x, y.Run())
        case *Variable:
            return Divide(x, y.Value())
        case String:
            return x.Split(y)
        case Number:
            if y.val.Cmp(NewNumber(0).val) != 0 {
                return x.Divide(y)
            }
        case Boolean:
            if y {
                return x.Divide(NewNumber(1))
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
                return Number{ val: NewNumber(0).val.Quo(x.val, NewNumber(len(y)).val) }
            }
        case Array:
            if len(y) != 0 {
                return Number{ val: NewNumber(0).val.Quo(x.val, NewNumber(len(y)).val) }
            }
        case String:
            if len(y) != 0 {
                return Number{ val: NewNumber(0).val.Quo(x.val, NewNumber(len(y)).val) }
            }
        case Number:
            if y.val.Cmp(NewNumber(0).val) != 0 {
                return Number{ val: NewNumber(0).val.Quo(x.val, y.val) }
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
            return a.Divide(NewNumber(len(y)))
        case Array:
            return a.Divide(NewNumber(len(y)))
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

func (a Array) Divide(b Number) Array {
    out := Array { }
    x := int(len(a) / b.Int())
    i := 0

    for len(out) < len(a) % b.Int() {
        set := Array { }

        for len(set) < x + 1 {
            set = append(set, a[i])
            i = i + 1
        }

        out = append(out, set)
    }

    for len(out) < b.Int() {
        set := Array { }

        for len(set) < x {
            set = append(set, a[i])
            i = i + 1
        }

        out = append(out, set)
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

func (a String) Divide(b Number) Array {
    out := Array { }
    x := int(len(a) / b.Int())
    i := 0

    for len(out) < len(a) % b.Int() {
        set := ""

        for len(set) < x + 1 {
            set += string(a[i])
            i = i + 1
        }

        out = append(out, String(set))
    }

    for len(out) < b.Int() {
        set := ""

        for len(set) < x {
            set += string(a[i])
            i = i + 1
        }

        out = append(out, String(set))
    }

    return out
}
