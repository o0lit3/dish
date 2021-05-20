package main

func Max(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Array().Max()
    case Array:
        return x.Max()
    case String:
        return x.Max()
    default:
        return x
    }
}

func (a Array) Max() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if x, ok := Above(val, out).(Boolean); Boolean(ok) && x {
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
