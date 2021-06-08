package main
import("fmt")

func Equals(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Equals(x.Run(), b)
    case *Variable:
        return Equals(x.Value(), b)
    case Hash:
        return Equals(String(fmt.Sprintf("%v", x)), b)
    case Array:
        return Equals(String(fmt.Sprintf("%v", x)), b)
    case String:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Equals(x, String(fmt.Sprintf("%v", y)))
        case Array:
            return Equals(x, String(fmt.Sprintf("%v", y)))
        case String:
            return Boolean(x == y)
        case Number:
            return Boolean(x.Number().val.Cmp(y.val) == 0)
        case Boolean:
            return Boolean(x.Number().val.Cmp(y.Number().val) == 0)
        case Null:
            return Boolean(x.Number().val.Cmp(NewNumber(0).val) == 0)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == 0)
        case Array:
            return Boolean(x.val.Cmp(NewNumber(len(y)).val) == 0)
        case String:
            return Boolean(x.val.Cmp(y.Number().val) == 0)
        case Number:
            return Boolean(x.val.Cmp(y.val) == 0)
        case Boolean:
            return Boolean(x.val.Cmp(y.Number().val) == 0)
        case Null:
            return Boolean(x.val.Cmp(NewNumber(0).val) == 0)
        }
    case Boolean:
        return Equals(x.Number(), b)
    case Null:
        return Equals(NewNumber(0), b)
    }

    return Boolean(false)
}
