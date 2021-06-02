package main

func Pop(a interface{}) (interface{}, interface{}) {
    switch x := a.(type) {
    case *Block:
        return Pop(x.Run())
    case *Variable:
        return Pop(x.Value())
    case Hash:
        return x.Pop()
    case Array:
        return x.Pop()
    case String:
        return x.Pop()
    default:
        return x, x
    }
}

func (a Hash) Pop() (interface{}, interface{}) {
    var val interface{}
    keys := Keys(a)

    val = Null { }

    if key, ok := keys[len(keys) - 1].(String); ok {
        val = a[string(key)]
        delete(a, string(key))
    }

    return val, a
}

func (a Array) Pop() (interface{}, interface{}) {
    if len(a) > 0 {
        val := a[len(a) - 1]
        a = a[:len(a) - 1]
        return val, a
    }

    return Null { }, a
}

func (a String) Pop() (interface{}, interface{}) {
    if len(a) > 0 {
        val := String(string(a[len(a) - 1]))
        a = a[:len(a) - 1]
        return val, a
    }

    return Null { }, a
}
