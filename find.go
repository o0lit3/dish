package main

import(
    "fmt"
    "strings"
)

func Find(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Find(x.Run(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        default:
            return x.Find(y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        default:
            return x.Find(y)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Find(x, y.Run())
        case String:
            return Number(strings.Index(string(x), string(y)))
        default:
            return Number(strings.Index(string(x), fmt.Sprintf("%v", y)))
        }
    }

    return Null { }
}

func (a Hash) Find(b interface{}) String {
    for key, val := range a {
        if Equals(val, b) {
            return String(key)
        }
    }

    return String("-1")
}

func (a Array) Find(b interface{}) Number {
    for i, val := range a {
        if Equals(val, b) {
            return Number(i)
        }
    }

    return Number(-1)
}
