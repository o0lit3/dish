package main
import("math")

func Max(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Max(x.Run())
    case *Variable:
        return Max(x.Value())
    case Hash:
        return x.Array().Max()
    case Array:
        return x.Max()
    case String:
        return x.Max()
    case Number:
        val, _ := x.val.Float64()
        return NewNumber(int(math.Ceil(val)))
    case Boolean:
        return x.Number()
    }

    return NewNumber(0)
}

func (a Array) Max() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if Above(val, out) {
            out = val
        }
    }

    return out
}

func (a String) Max() interface{} {
    var out interface{}

    for _, c := range a {
        if out == nil {
            out = String(c)
        } else if x, ok := out.(String); ok && String(c) > x {
            out = String(c)
        }
    }

    return out
}
