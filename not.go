package main

func Not(a interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Not(x.Run())
    case *Variable:
        return Not(x.Value())
    case Hash:
        return Boolean(len(x) == 0)
    case Array:
        return Boolean(len(x) == 0)
    case String:
        return Boolean(x == "" || x == "0")
    case Number:
        return Boolean(x == 0)
    case Boolean:
        return Boolean(!x)
    }

    return Boolean(true)
}
