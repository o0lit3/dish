package main

func Sum(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Sum(x.Run())
    case *Variable:
        return Sum(x.Value())
    case Hash:
        return x.Array().Sum()
    case Array:
        return x.Sum()
    case String:
        return x.Number()
    case Number:
        return x
    case Boolean:
        return x.Number()
    case Null:
        return NewNumber(0)
    }

    return Null { }
}

func (a Array) Sum() Number {
    out := NewNumber(0)

    for _, val := range a {
        if x, ok := Sum(val).(Number); ok {
            out = Number{ val: out.val.Add(out.val, x.val) }
        }
    }

    return out
}
