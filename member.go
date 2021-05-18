package main

import (
    "fmt"
    "strings"
)

func Member(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(string(y))
        default:
            return x.Member(fmt.Sprintf("%v", y))
        }
    case Array:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(int(y.Number()))
        case Number:
            return x.Member(int(y))
        case Boolean:
            return x.Member(int(y.Number()))
        }
    case String:
        switch y := b.(type) {
        case Interpreter:
            return Member(x, y.Run())
        case Array:
            return x.Members(y)
        case String:
            return x.Substring(y)
        case Number:
            return x.Member(int(y))
        case Boolean:
            return x.Member(int(y.Number()))
        }
    }

    return Null { }
}

func (a Hash) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        var addend interface{}

        switch key := val.(type) {
        case String:
            addend = a[string(key)]
        default:
            addend = a[fmt.Sprintf("%v", key)]
        }

        if addend != nil {
            out = append(out, addend)
        }
    }

    return out
}

func (a Hash) Member(b string) interface{} {
    out := a[b]

    if out != nil {
        return out
    }

    return Null { }
}

func (a Array) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        idx := len(a)

        switch x := val.(type) {
        case String:
            idx = int(x.Number())
        case Number:
            idx = int(x)
        case Boolean:
            idx = int(x.Number())
        }

        if idx < len(a) {
            out = append(out, a[idx])
        }
    }

    return out
}

func (a Array) Member(b int) interface{} {
    if b < len(a) {
        return a[b]
    }

    return Null { }
}

func (a String) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        switch x := val.(type) {
        case String:
            if strings.Contains(string(a), string(x)) {
                out = append(out, x)
            }
        case Number:
            if int(x) < len(a) {
                out = append(out, String(a[int(x)]))
            }
        case Boolean:
            if int(x.Number()) < len(a) {
                out = append(out, String(a[int(x.Number())]))
            }
        }
    }

    return out
}

func (a String) Substring(b String) interface{} {
    if strings.Contains(string(a), string(b)) {
        return b
    }

    return Null { }
}

func (a String) Member(b int) interface{} {
    if b < len(a) {
        return String(string(string(a)[b]))
    }

    return Null { }
}
