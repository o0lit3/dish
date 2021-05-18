package main

func Length(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return Number(len(x))
    case Array:
        return Number(len(x))
    case String:
        return Number(len(x))
    case Number:
        return Number(1)
    case Boolean:
        return Number(1)
    }

    return Number(0)
}
