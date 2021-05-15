package main

func Sum(a interface{}) interface{} {
    switch a.(type) {
    case Hash:
    case Array:
    case String:
    case Number:
    case Boolean:
    }

    return Null { }
}
