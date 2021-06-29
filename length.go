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
    default:
        return Null { }
    }
}
