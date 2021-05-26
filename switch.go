package main

func Switch(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Switch(x.Run())
    case Hash:
        return x.Switch()
    case Array:
        return x.Switch()
    case String:
        if len(x) > 0 {
            return Number(0)
        }
    }

    return Number(-1)
}

func (a Hash) Switch() String {
    for key, val := range a {
        if b, ok := Not(val).(Boolean); ok && bool(b) {
            continue
        }

        return String(key)
    }

    return String("-1")
}

func (a Array) Switch() Number {
    for i, val := range a {
        if b, ok := Not(val).(Boolean); ok && bool(b) {
            continue
        }

        return Number(i)
    }

    return Number(-1)
}
