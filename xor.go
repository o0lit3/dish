package main

func Xor(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Xor(x.Run(), b)
    case *Variable:
        return Xor(x.Value(), b)
    case Hash:
        if len(x) == 0 {
            switch y := b.(type) {
            case *Block:
                return y.Run()
            case *Variable:
                return y.Value()
            default:
                return y
            }
        }

        if Not(b) {
            return x
        }

        return Boolean(false)
    case Array:
        if len(x) == 0 {
            switch y := b.(type) {
            case *Block:
                return y.Run()
            case *Variable:
                return y.Value()
            default:
                return y
            }
        }

        if Not(b) {
            return x
        }

        return Boolean(false)
    case String:
        if x == "" || x == "0" {
            switch y := b.(type) {
            case *Block:
                return y.Run()
            case *Variable:
                return y.Value()
            default:
                return y
            }
        }

        if Not(b) {
            return x
        }

        return Boolean(false)
    case Number:
        if x == 0 {
            switch y := b.(type) {
            case *Block:
                return y.Run()
            case *Variable:
                return y.Value()
            default:
                return y
            }
        }

        if Not(b) {
            return x
        }

        return Boolean(false)
    case Boolean:
        if !x {
            switch y := b.(type) {
            case *Block:
                return y.Run()
            case *Variable:
                return y.Value()
            default:
                return y
            }
        }

        if Not(b) {
            return x
        }

        return Boolean(false)
    case Null:
        switch y := b.(type) {
        case *Block:
            return y.Run()
        case *Variable:
            return y.Value()
        default:
            return y
        }
    }

    return Boolean(false)
}
