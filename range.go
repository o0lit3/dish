package main

func Range(a interface{}, b interface{}) interface{} {
    str := false
    start := 0
    end := 0

    switch x := a.(type) {
    case *Block:
        return Range(x.Run(), b)
    case *Variable:
        return Range(x.Value(), b)
    case Hash:
        start = len(x)
    case Array:
        start = len(x)
    case String:
        start = int([]rune(string(x))[len(x) - 1])
        str = true
    case Number:
        start = x.Int()
    case Boolean:
        start = map[bool]int{true: 1, false: 0}[bool(x)]
    default:
        start = 0
    }

    switch y := b.(type) {
    case *Block:
        return Range(a, y.Run())
    case *Variable:
        return Range(a, y.Value())
    case Hash:
        end = len(y)
    case Array:
        end = len(y)
    case String:
        end = int([]rune(string(y))[len(y) - 1])
    case Number:
        end = y.Int()
    case Boolean:
        end = map[bool]int{true: 1, false: 0}[bool(y)]
    default:
        end = 0
    }

    if end > start {
        out := Array { }
        n := start

        for n <= end {
            if str {
                out = append(out, String(string(n)))
            } else {
                out = append(out, NewNumber(n))
            }

            n = n + 1
        }

        return out
    } else {
        out := Array { }
        n := start

        for n >= end {
            if str {
                out = append(out, String(string(n)))
            } else {
                out = append(out, NewNumber(n))
            }

            n = n - 1
        }

        return out
    }

    return Array { }
}
