package main
import("fmt")

func Equals(a interface{}, b interface{}) Boolean {
    switch x := a.(type) {
    case *Block:
        return Equals(x.Run(), b)
    case *Variable:
        return Equals(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Hash:
            return Boolean(x.String() == y.String())
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Array:
            return Boolean(x.String() == y.String())
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case String:
            return Boolean(string(x) == string(y))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Number:
            if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
                return Boolean(true)
            }

            if x.inf == INF || x.inf == -INF || y.inf == INF || y.inf == -INF {
                return Boolean(false)
            }

            return Boolean(x.val.Cmp(y.val) == 0)
        }
    case Boolean:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Boolean:
            return Boolean(x == y)
        }
    case Null:
        switch y := b.(type) {
        case *Block:
            return Equals(x, y.Run())
        case *Variable:
            return Equals(x, y.Value())
        case Null:
            return Boolean(true)
        }
    }

    return Boolean(false)
}

func Clone(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return x.Run()
    case *Variable:
        return x.Value()
    case Hash:
        out := Hash{ }

        for key, val := range x {
            out[key] = Clone(val)
        }

        return out
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, Clone(val))
        }

        return out
    case String:
        out := ""

        for _, c := range x {
            out += string(c)
        }

        return String(out)
    case Number:
        return x
    case Boolean:
        return x
    }

    return Null{ }
}

func Boolify(a interface{}) Boolean {
    switch x := a.(type) {
        case *Block:
            return Boolify(x.Run())
        case *Variable:
            return Boolify(x.Value())
        case Hash:
            return Boolean(len(x) > 0)
        case Array:
            return Boolean(len(x) > 0)
        case String:
            return Boolean(string(x) != "")
        case Number:
            if x.inf == INF || x.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(NewNumber(0).val) != 0)
        case Boolean:
            return x
    }

    return Boolean(false)
}

func Numerize (a interface{}) Number {
    switch x := a.(type) {
    case *Block:
        return Numerize(x.Run())
    case *Variable:
        return Numerize(x.Value())
    case Hash:
        return NewNumber(len(x))
    case Array:
        return NewNumber(len(x))
    case String:
        return x.Number()
    case Number:
        return x
    case Boolean:
        return x.Number()
    default:
        return NewNumber(0)
    }
}

func Stringify(a interface{}) String {
    switch x := a.(type) {
    case *Block:
        return Stringify(x.Run())
    case *Variable:
        return Stringify(x.Value())
    case String:
        return x
    default:
        return String(fmt.Sprintf("%v", x))
    }
}

func Itemize (a interface{}) Array {
    switch x:= a.(type) {
    case *Block:
        return Itemize(x.Run())
    case *Variable:
        return Itemize(x.Value())
    case Hash:
        return x.Array()
    case Array:
        return x
    case String:
        return Array{ x }
    case Number:
        return Array{ x }
    case Boolean:
        return Array{ x }
    default:
        return Array { }
    }
}

func Hashify(a interface{}) Hash {
    switch x := a.(type) {
    case *Block:
        return Hashify(x.Run())
    case *Variable:
        return Hashify(x.Value())
    case Hash:
        return x
    case Array:
        return x.Hash()
    case String:
        return Hash{ string(x): Boolean(true) }
    case Null:
        return Hash{ }
    default:
        return Hash{ fmt.Sprintf("%v", x): Boolean(true) }
    }
}
