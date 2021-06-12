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
        return Boolean(string(x) == "" || string(x) == "0")
    case Number:
        return Boolean(x.val.Cmp(NewNumber(0).val) == 0)
    case Boolean:
        return Boolean(!x)
    }

    return Boolean(true)
}
