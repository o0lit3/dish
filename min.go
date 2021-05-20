package main

func Min(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Array().Min()
    case Array:
        return x.Min()
    case String:
        return x.Min()
    default:
        return x
    }
}

func (a Array) Min() interface{} {
    var out interface{}

    for _, val := range a {
        if out == nil {
            out = val
        } else if x, ok := Below(val, out).(Boolean); Boolean(ok) && x {
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
