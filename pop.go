package main

func Pop(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Pop(x.Run())
    case Hash:
        return x.Array().Pop()
    case Array:
        return x.Pop()
    case String:
        return x.Pop()
    default:
        return x
    }
}

func (a Array) Pop() interface{} {
    if len(a) > 0 {
        return a[len(a) - 1]
    }

    return Null { }
}

func (a String) Pop() interface{} {
    if len(a) > 0 {
        return String(string(a[len(a) - 1]))
    }

    return Null { }
}
