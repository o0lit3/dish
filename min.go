package main
import("math")

func Min(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Min(x.Run())
    case *Variable:
        return Min(x.Value())
    case Hash:
        return x.Array().Min()
    case Array:
        return x.Min()
    case String:
        return x.Min()
    case Number:
        val, _ := x.val.Float64()
        return NewNumber(int(math.Floor(val)))
    case Boolean:
        return x.Number()
    }

    return NewNumber(0)
}

func (a Array) Min() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if Below(val, out) {
            out = val
        }
    }

    return out
}

func (a String) Min() interface{} {
    var out interface{}

    for _, c := range a {
        if out == nil {
            out = String(c)
        } else if x, ok := out.(String); ok && String(c) < x {
            out = String(c)
        }
    }

    return out
}
