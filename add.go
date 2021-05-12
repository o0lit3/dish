package main

import ("strconv")

func Add(t Token, a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Map:
        switch y := b.(type) {
        case Map:
            for key, val := range y {
                x[key] = val
            }

            return x
        case Set:
            for idx, val := range y {
                x[strconv.Itoa(idx)] = val
            }

            return x
        }
    case Set:
        switch y := b.(type) {
        case Map:
            for _, val := range y {
                x = append(x, val)
            }

            return x
        case Set:
            for _, val := range y {
                x = append(x, val)
            }

            return x
        }
    case string:
    case float64:
    case int:
        switch y := b.(type) {
        case Map:
        case Set:
        case string:
        case float64:
        case int:
            return x + y
        case bool:
        }
    case bool:
    }

    return ""
}
