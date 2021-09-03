package main
import("strings"; "math/big")

func (t *Token) Split(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Split(x.Run(), b)
    case *Variable:
        return t.Split(x.Value(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "/" && t.lit != "/=" && t.lit != "divide" && t.lit != "div" && t.lit != "split" {
                    t.TypeMismatch(x, y)
                }

                return t.DivideArray(x, y)
            }

            return t.Split(x, y.Context(x).Run())
        case *Variable:
            return t.Split(x, y.Value())
        case Number:
            if t.lit != "/" && t.lit != "/=" && t.lit != "split" {
                t.TypeMismatch(x, y)
            }

            if y.inf == 0 && y.val.Cmp(NewNumber(0).val) == 0 {
                t.DivideByZero()
            }

            return t.SplitArray(x, y)
        case Boolean:
            return t.Split(x, y.Number())
        case Null:
            return t.Split(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "/" && t.lit != "/=" && t.lit != "split" {
                    t.TypeMismatch(x, y)
                }

                return t.DivideString(x, y)
            }

            return t.Split(x, y.Context(x).Run())
        case *Variable:
            return t.Split(x, y.Value())
        case String:
            if t.lit != "/" && t.lit != "/=" && t.lit != "split" {
                t.TypeMismatch(x, y)
            }

            return t.SplitString(x, y)
        case Number:
            if t.lit != "/" && t.lit != "/=" && t.lit != "split" {
                t.TypeMismatch(x, y)
            }

            if y.inf == 0 && y.val.Cmp(NewNumber(0).val) == 0 {
                t.DivideByZero()
            }

            return t.SectionString(x, y)
        case Boolean:
            return t.Split(x, y.Number())
        case Null:
            return t.Split(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.Split(x, y.Run())
        case *Variable:
            return t.Split(x, y.Value())
        case Number:
            if t.lit != "/" && t.lit != "/=" && t.lit != "divide" && t.lit != "div" {
                t.TypeMismatch(x, y)
            }

            if y.inf == 0 && y.val.Cmp(NewNumber(0).val) == 0 {
                t.DivideByZero()
            }

            return t.DivideNumber(x, y)
        case Boolean:
            return t.Split(x, y.Number())
        case Null:
            return t.Split(x, NewNumber(0))
        }
    case Boolean:
        return t.Split(x.Number(), b)
    case Null:
        return t.Split(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) DoubleSplit(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleSplit(x.Run(), b)
    case *Variable:
        return t.DoubleSplit(x.Value(), b)
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "//" && t.lit != "group" {
                    t.TypeMismatch(x, y)
                }

                return t.GroupArray(x, y)
            }

            return t.DoubleSplit(x, y.Context(x).Run())
        case *Variable:
            return t.DoubleSplit(x, y.Value())
        case Number:
            if t.lit != "//" && t.lit != "partition" {
                t.TypeMismatch(x, y)
            }

            if y.inf == 0 && y.val.Cmp(NewNumber(0).val) == 0 {
                t.DivideByZero()
            }

            return t.PartitionArray(x, y)
        case Boolean:
            return t.DoubleSplit(x, y.Number())
        case Null:
            return t.DoubleSplit(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "//" && t.lit != "group" {
                    t.TypeMismatch(x, y)
                }

                return t.GroupString(x, y)
            }

            return t.DoubleSplit(x, y.Context(x).Run())
        case *Variable:
            return t.DoubleSplit(x, y.Value())
        case String:
            if t.lit != "//" && t.lit != "partition" {
                t.TypeMismatch(x, y)
            }

            return t.BreakString(x, y)
        case Number:
            if t.lit != "//" && t.lit != "partition" {
                t.TypeMismatch(x, y)
            }

            return t.PartitionString(x, y)
        case Boolean:
            return t.DoubleSplit(x, y.Number())
        case Null:
            return t.DoubleSplit(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.DoubleSplit(x, y.Run())
        case *Variable:
            return t.DoubleSplit(x, y.Value())
        case Number:
            if t.lit != "//" && t.lit != "idiv" {
                t.TypeMismatch(x, y)
            }

            return NewNumber(t.DivideNumber(x, y).(Number).Int())
        case Boolean:
            return t.DoubleSplit(x, y.Number())
        case Null:
            return t.DoubleSplit(x, NewNumber(0))
        }
    case Boolean:
        return t.DoubleSplit(x.Number(), b)
    case Null:
        return t.DoubleSplit(NewNumber(0), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopSplit(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopSplit(x.Run())
    case *Variable:
        return t.TopSplit(x.Value())
    case Hash:
        if t.lit != "/" && t.lit != "array" && t.lit != "arr" && t.lit != "values" && t.lit != "vals" {
            t.TypeMismatch(a, nil)
        }

        return x.Array()
    case Array:
        if t.lit != "/" && t.lit != "array" && t.lit != "arr" {
            t.TypeMismatch(a, nil)
        }

        return x
    case String:
        if t.lit != "/" && t.lit != "split" && t.lit != "words" {
            t.TypeMismatch(a, nil)
        }

        return t.WordsFromString(x)
    case Number:
        if t.lit != "/" && t.lit != "factors" {
            t.TypeMismatch(a, nil)
        }

        return t.FactorsForNumber(x)
    case Boolean:
        return t.TopSplit(x.Number())
    case Null:
        return t.TopSplit(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) TopDoubleSplit(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopDoubleSplit(x.Run())
    case *Variable:
        return t.TopDoubleSplit(x.Value())
    case Hash:
        if t.lit != "//" && t.lit != "flatten" && t.lit != "flat" {
            t.TypeMismatch(x, nil)
        }

        return t.FlattenHash(x)
    case Array:
        if t.lit != "//" && t.lit != "flatten" && t.lit != "flat" {
            t.TypeMismatch(x, nil)
        }

        return t.FlattenArray(x)
    case String:
        if t.lit != "//" && t.lit != "chars" {
            t.TypeMismatch(x, nil)
        }

        return x.Array()
    case Number:
        if t.lit != "//" && t.lit != "bits" {
            t.TypeMismatch(x, nil)
        }

        return x.Array()
    case Boolean:
        return t.TopDoubleSplit(x.Number())
    case Null:
        return t.TopDoubleSplit(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) SplitString(x String, y String) Array {
    items := strings.Split(string(x), string(y))
    out := Array{ }

    for _, item := range items {
        out = append(out, String(item))
    }

    return out
}

func (t *Token) WordsFromString(x String) Array {
    out := Array{ }

    for _, item := range strings.Fields(string(x)) {
        out = append(out, String(item))
    }

    return out
}

func (t *Token) PartitionArray(x Array, y Number) Array {
    out := Array{ }
    step := y.Int()
    i := 0

    for i < len(x) {
        if i + step < len(x) {
            out = append(out, x[i:i + step])
        } else {
            out = append(out, x[i:len(x)])
        }

        i = i + step
    }

    return out
}

func (t *Token) PartitionString(x String, y Number) Array {
    step := y.Int()
    out := Array{ }
    i := 0

    for i < len(x) {
        if i + step < len(x) {
            out = append(out, x[i:i + step])
        } else {
            out = append(out, x[i:len(x)])
        }

        i = i + step
    }

    return out
}

func (t *Token) DivideArray(x Array, y *Block) Array {
    out := Array{ Array{ } }

    for i, val := range x {
        if bool(Boolify(y.Context(x).Run(val, NewNumber(i))) && i < len(x) - 1) {
            out = append(out, Array{ })
        } else {
            out[len(out) - 1] = append(out[len(out) - 1].(Array), val)
        }
    }

    return out
}

func (t *Token) DivideString(x String, y *Block) Array {
    items := t.DivideArray(x.Array(), y)
    out := Array{ }

    for _, item := range items {
        out = append(out, t.JoinArray(item.(Array), String("")))
    }

    return out
}

func (t *Token) DivideNumber(x Number, y Number) interface{} {
    if (x.inf == INF || x.inf == -INF) && (y.inf == INF || y.inf == -INF) {
        return Null{ }
    }

    if y.inf == INF || y.inf == -INF {
        return NewNumber(0)
    }

    if x.inf == INF || x.inf == -INF {
        if y.val.Cmp(NewNumber(0).val) == -1 {
            return Number{ inf: -x.inf }
        }

        return Number { inf: x.inf }
    }

    if y.val.Cmp(NewNumber(0).val) != 0 {
        return Number{ val: NewNumber(0).val.Quo(x.val, y.val) }
    } else {
        switch x.val.Cmp(NewNumber(0).val) {
        case -1:
            return Number{ inf: -INF }
        case 1:
            return Number{ inf: INF }
        }

        return Null{ }
    }
}

func (t *Token) GroupArray(x Array, y *Block) Hash {
    out := Hash{ }

    for i, val := range x {
        key := string(Stringify(y.Run(val, NewNumber(i))))

        if _, ok := out[key]; !ok {
            out[key] = Array{ }
        }

        out[key] = append(out[key].(Array), val)
    }

    return out
}

func (t *Token) GroupString(x String, y *Block) Hash {
    out := Hash{ }

    for i, c := range x {
        key := string(Stringify(y.Run(String(string(c)), NewNumber(i))))

        if _, ok := out[key]; !ok {
            out[key] = String("")
        }

        out[key] = t.ConcatString(out[key].(String), String(string(c)))
    }

    return out
}

func (t *Token) BreakString(x String, y String) Array {
    items := strings.Split(string(x), string(y))
    out := Array{ }

    for _, item := range items {
        if len(out) < len(items) - 1 {
            out = append(out, String(item + string(y)))
        } else {
            out = append(out, String(item))
        }
    }

    return out
}

func (t *Token) SplitArray(x Array, y Number) Array {
    out := Array{ }
    parts := y.Int()
    size := len(x) / parts + 1
    full := len(x) % parts
    i := 0

    if full == 0 {
        full = parts
        size--
    }

    for len(out) < full {
        out = append(out, Array{ })

        for len(out[len(out) - 1].(Array)) < size {
            out[len(out) - 1] = append(out[len(out) - 1].(Array), x[i])
            i++
        }
    }

    for i < len(x) || len(out) < parts {
        out = append(out, Array{ })

        for len(out[len(out) - 1].(Array)) < size - 1 {
            out[len(out) - 1] = append(out[len(out) - 1].(Array), x[i])
            i++
        }
    }

    return out
}

func (t *Token) SectionString(x String, y Number) Array {
    items := t.SplitArray(x.Array(), y)
    out := Array{ }

    for _, item := range items {
        out = append(out, t.JoinArray(item.(Array), String("")))
    }

    return out
}

func (t *Token) FlattenHash(x Hash) Array {
    out := Array{ }

    for key, val := range x {
        out = append(out, String(key))
        out = append(out, val)
    }

    return out
}

func (t *Token) FlattenArray (a interface{}) Array {
    out := Array{ }

    switch x := a.(type) {
    case *Block:
        return t.FlattenArray(x.Run())
    case *Variable:
        return t.FlattenArray(x.Value())
    case Hash:
        return t.FlattenArray(x.Array())
    case Array:
        for _, item := range x {
            for _, val := range t.FlattenArray(item) {
                out = append(out, val)
            }
        }
    default:
        out = append(out, x)
    }

    return out
}

func (t *Token) FactorsForNumber(x Number) interface{} {
    out := Array{ }
    n := new(big.Int).Quo(x.val.Num(), x.val.Denom())
    mod, div := new(big.Int), new(big.Int)
    i := big.NewInt(2)

    if x.val.Cmp(NewNumber(0).val) < 0 {
        out = append(out, NewNumber(-1))
        n = n.Neg(n)
    }

    for i.Cmp(n) != 1 {
        div.DivMod(n, i, mod)

        for mod.Cmp(big.NewInt(0)) == 0 {
            out = append(out, Number{ val: new(big.Rat).SetInt(i) })
            n.Set(div)
            div.DivMod(n, i, mod)
        }

        i.Add(i, big.NewInt(1))
    }

    return out
}

func (t *Token) DivisorsForNumber(x Number) interface{} {
    i := new(big.Int).Quo(x.val.Num(), x.val.Denom())
    j := big.NewInt(1)

    if i.Cmp(big.NewInt(0)) < 0 {
        i = i.Neg(i)
    }

    if i.Cmp(big.NewInt(2)) < 0 {
        return Array{ }
    }

    out := Array { NewNumber(1) }
    sqrt := new(big.Int).Sqrt(i)

    for j.Cmp(sqrt) <= 0 {
        if j.Cmp(big.NewInt(1)) > 0 && new(big.Int).Rem(i, j).Cmp(big.NewInt(0)) == 0 {
            out = append(out, Number{ val: new(big.Rat).SetInt(j) })
            div := new(big.Int).Div(i, j)

            if j.Cmp(div) != 0 {
                out = append(out, Number{ val: new(big.Rat).SetInt(div) })
            }
        }

        j = new(big.Int).Add(j, big.NewInt(1))
    }

    return t.SortArray(out, nil)
}
