package main

func Not(t Token, a interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        return len(x) == 0
    case Set:
        return len(x) == 0
    case string:
        return x == "" || x == "0"
    case float64:
        return x == 0
    case int:
        return x == 0
    case bool:
        return !x
    }

    return ""
}
