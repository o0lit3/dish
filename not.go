package main

func Not(t Token, a interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        return Boolean(len(x) == 0)
    case Set:
        return Boolean(len(x) == 0)
    case String:
        return Boolean(x == "" || x == "0")
    case Number:
        return Boolean(x == 0)
    case Boolean:
        return Boolean(!x)
    }

    return String("")
}
