package main
import("strconv")

func Shift(a interface{}) (interface{}, interface{}) {
    switch x := a.(type) {
    case *Block:
        return Shift(x.Run())
    case *Variable:
        return Shift(x.Value())
    case Hash:
        return x.Shift()
    case Array:
        return x.Shift()
    case String:
        return x.Shift()
    case Number:
        return x.Shift()
    case Boolean:
        return x, Boolean(false)
    default:
        return x, Null{ }
    }
}

func (a Hash) Shift() (interface{}, interface{}) {
    var val interface{}
    keys := Keys(a).(Array)

    val = Null { }

    if key, ok := keys[0].(String); ok {
        val = a[string(key)]
        delete(a, string(key))
    }

    return val, a
}

func (a Array) Shift() (interface{}, interface{}) {
    if len(a) > 0 {
        val := a[0]
        a = a[1:]
        return val, a
    }

    return Null { }, a
}

func (a String) Shift() (interface{}, interface{}) {
    if len(a) > 0 {
        val := String(string(a[0]))
        a = a[1:]
        return val, a
    }

    return Null { }, a
}

func (a Number) Shift() (interface{}, interface{}) {
    if a.inf == INF || a.inf == -INF {
        return NewNumber(0), a
    }

    bin := strconv.FormatInt(int64(a.Int()), 2)
    first, _ := strconv.Atoi(string(bin[len(bin) - 1]))
    rem, _ := strconv.ParseInt(string(bin[:len(bin) - 1]), 2, 64)
    return NewNumber(first), NewNumber(int(rem))
}
