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
    return Null { }
}

func (a String) Min() interface{} {
    return Null { }
}
