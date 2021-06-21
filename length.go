package main

func Length(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Length(x.Run())
    case *Variable:
        return Length(x.Value())
    case Hash:
        return NewNumber(len(x))
    case Array:
        return NewNumber(len(x))
    case String:
        return NewNumber(len(x))
    case Number:
        return Length(String(x.String()))
    case Boolean:
        return NewNumber(1)
    }

    return NewNumber(0)
}
