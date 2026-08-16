package main
import("math"; "math/big")

func (t *Token) Splat(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Splat(x.Run(), b)
    case *Variable:
        return t.Splat(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "*" && t.lit != "*=" && t.lit != "map" && t.lit != "each" {
                    t.TypeMismatch(x, y)
                }

                return t.MapHash(x, y)
            }

            return t.Splat(x, y.Run())
        case *Variable:
            return t.Splat(x, y.Value())
        case Hash:
            if t.lit != "*" && t.lit != "*=" && t.lit != "dot" {
                t.TypeMismatch(x, y)
            }

            return t.DotHash(x, y)
        case String:
            if t.lit != "*" && t.lit != "*=" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(x.Array(), y)
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "*" && t.lit != "*=" && t.lit != "map" && t.lit != "each" {
                    t.TypeMismatch(x, y)
                }

                return t.MapArray(x, y)
            }

            return t.Splat(x, y.Run())
        case *Variable:
            return t.Splat(x, y.Value())
        case Array:
            if t.lit != "*" && t.lit != "*=" && t.lit != "dot" {
                t.TypeMismatch(x, y)
            }

            return t.DotArray(x, y)
        case String:
            if t.lit != "*" && t.lit != "*=" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(x, y)
        case Number:
            if t.lit != "*" && t.lit != "*=" && t.lit != "repeat" {
                t.TypeMismatch(x, y)
            }

            return t.RepeatArray(x, y)
        case Boolean:
            return t.Splat(x, y.Number())
        case Null:
            if t.lit != "*" && t.lit != "*=" && t.lit != "repeat" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            if t.lit == "join" {
                return t.JoinArray(x, String(""))
            }

            return Array{ }
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "*" && t.lit != "*=" && t.lit != "map" && t.lit != "each" {
                    t.TypeMismatch(x, y)
                }

                return t.JoinArray(t.MapArray(x.Array(), y), String(""))
            }

            return t.Splat(x, y.Run())
        case *Variable:
            return t.Splat(x, y.Value())
        case Hash:
            if t.lit != "*" && t.lit != "*=" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(y.Array(), x)
        case Array:
            if t.lit != "*" && t.lit != "*=" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            return t.JoinArray(y, x)
        case String:
            if t.lit != "*" && t.lit != "*=" && t.lit != "join" {
                t.TypeMismatch(x, y)
            }

            return t.JoinString(x, y)
        case Number:
            if t.lit != "*" && t.lit != "*=" && t.lit != "repeat" {
                t.TypeMismatch(x, y)
            }

            return t.RepeatString(x, y)
        case Boolean:
            return t.Splat(x, y.Number())
        case Null:
            if t.lit != "*" && t.lit != "*=" && t.lit != "repeat" {
                t.TypeMismatch(x, y)
            }

            return String("")
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                if t.lit != "*" && t.lit != "*=" && t.lit != "times" {
                    t.TypeMismatch(x, y)
                }

                return t.MapArray(t.RangeNumber(NewNumber(0), t.SubtractNumber(x, NewNumber(1)).(Number)), y)
            }

            return t.Splat(x, y.Run())
        case *Variable:
            return t.Splat(x, y.Value())
        case Array:
            if t.lit != "*" && t.lit != "*=" {
                t.TypeMismatch(x, y)
            }

            return t.RepeatArray(y, x)
        case String:
            if t.lit != "*" && t.lit != "*=" {
                t.TypeMismatch(x, y)
            }

            return t.RepeatString(y, x)
        case Number:
            if t.lit != "*" && t.lit != "*=" && t.lit != "multiply" && t.lit != "mult" {
                t.TypeMismatch(x, y)
            }

            return t.MultiplyNumber(x, y)
        case Boolean:
            return t.Splat(x, y.Number())
        case Null:
            return t.Splat(x, NewNumber(0))
        }
    case Boolean:
        return t.Splat(x.Number(), b)
    case Null:
        switch y := b.(type) {
        case *Block:
            return t.Splat(x, y.Run())
        case *Variable:
            return t.Splat(x, y.Value())
        case Array:
            return Array{ }
        case String:
            return String("")
        case Number:
            return NewNumber(0)
        case Boolean:
            return NewNumber(0)
        case Null:
            return NewNumber(0)
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) DoubleSplat(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleSplat(x.Run(), b)
    case *Variable:
        return t.DoubleSplat(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return t.DoubleSplat(x, y.Run())
        case *Variable:
            return t.DoubleSplat(x, y.Value())
        case Number:
            if t.lit != "**" && t.lit != "comb" {
                t.TypeMismatch(x, y)
            }

            return t.CombArray(x.Array(), y)
        case Boolean:
            return t.DoubleSplat(x, y.Number())
        case Null:
            return t.DoubleSplat(x, NewNumber(0))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            return t.DoubleSplat(x, y.Run())
        case *Variable:
            return t.DoubleSplat(x, y.Value())
        case Number:
            if t.lit != "**" && t.lit != "comb" {
                t.TypeMismatch(x, y)
            }

            return t.CombArray(x, y)
        case Boolean:
            return t.DoubleSplat(x, y.Number())
        case Null:
            return t.DoubleSplat(x, NewNumber(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return t.DoubleSplat(x, y.Run())
        case *Variable:
            return t.DoubleSplat(x, y.Value())
        case Number:
            if t.lit != "**" && t.lit != "comb" {
                t.TypeMismatch(x, y)
            }

            return t.CombString(x, y)
        case Boolean:
            return t.DoubleSplat(x, y.Number())
        case Null:
            return t.DoubleSplat(x, NewNumber(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return t.DoubleSplat(x, y.Run())
        case *Variable:
            return t.DoubleSplat(x, y.Value())
        case Number:
            if t.lit != "**" && t.lit != "choose" {
                t.TypeMismatch(x, y)
            }

            return t.ChooseNumber(x, y)
        case Boolean:
            return t.ChooseNumber(x, y.Number())
        case Null:
            return t.ChooseNumber(x, NewNumber(0))
        }
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) TopSplat(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopSplat(x.Run())
    case *Variable:
        return t.TopSplat(x.Value())
    case Hash:
        return t.TopSplat(x.Array())
    case Array:
        if t.lit != "*" && t.lit != "*=" && t.lit != "product" && t.lit != "prod" && t.lit != "join" {
            t.TypeMismatch(x, nil)
        }

        switch t.lit {
        case "join":
            return t.JoinArray(x, String("\n"))
        case "product", "prod":
            return t.MultiplyArray(x)
        default:
            if x.Numeric() {
                return t.MultiplyArray(x)
            }

            return t.JoinArray(x, String("\n"))
        }
    case String:
        if t.lit != "*" && t.lit != "*=" && t.lit != "string" && t.lit != "str" {
            t.TypeMismatch(x, nil)
        }

        return append(String{ }, x...)
    default:
        if t.lit != "*" && t.lit != "*=" && t.lit != "string" && t.lit != "str" {
            t.TypeMismatch(x, nil)
        }

        return Stringify(x)
    }
}

func (t *Token) TopDoubleSplat(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopDoubleSplat(x.Run())
    case *Variable:
        return t.TopDoubleSplat(x.Value())
    case Hash:
        if t.lit != "**" && t.lit != "perms" && t.lit != "permutations" {
            t.TypeMismatch(x, nil)
        }

        return t.PermsArray(x.Array())
    case Array:
        if t.lit != "**" && t.lit != "perms" && t.lit != "permutations" {
            t.TypeMismatch(x, nil)
        }

        return t.PermsArray(x)
    case String:
        if t.lit != "**" && t.lit != "perms" && t.lit != "permutations" {
            t.TypeMismatch(x, nil)
        }

        return t.PermsString(x)
    case Number:
        if t.lit != "**" && t.lit != "divisors" && t.lit != "divs" {
            t.TypeMismatch(x, nil)
        }

        return t.DivisorsNumber(x)
    case Boolean:
        return t.TopDoubleSplat(x.Number())
    case Null:
        return t.TopDoubleSplat(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) MapHash(x Hash, y *Block) Hash {
    out := Hash { }

    for key, val := range x {
        out[key] = y.Context(x).Topic(val).Run(val, String(key))
    }

    return out
}

func (t *Token) MapArray(x Array, y *Block) Array {
    out := Array { }

    for i, val := range x {
        out = append(out, y.Context(x).Topic(val).Run(val, NewNumber(i)))
    }

    return out
}

func (t *Token) DotHash(x Hash, y Hash) Hash {
    out := Hash { }

    for key, _ := range x {
        if _, ok := y[key]; ok {
            switch x := x[key].(type) {
            case Number:
                switch y := y[key].(type) {
                case Number:
                    out[key] = t.MultiplyNumber(x, y)
                default:
                    t.TypeMismatch(x, y)
                }
            default:
                t.TypeMismatch(x, y)
            }
        }
    }

    return out
}

func (t *Token) DotArray(x Array, y Array) Array {
    out := Array { }

    for i, arow := range x {
        if acols, ok := arow.(Array); ok {
            row := Array { }

            for k, brow := range y {
                if bcols, ok := brow.(Array); ok {
                    for j, bval := range bcols {
                        if j >= len(row) {
                            row = append(row, NewNumber(0))
                        }

                        if k < len(acols) {
                            switch r := row[j].(type) {
                            case Number:
                                switch c := acols[k].(type) {
                                case Number:
                                    switch b := bval.(type) {
                                    case Number:
                                        row[j] = t.AddNumber(r, t.MultiplyNumber(c, b).(Number))
                                    default:
                                        t.TypeMismatch(r, b)
                                    }
                                default:
                                    t.TypeMismatch(r, c)
                                }
                            default:
                                t.TypeMismatch(r, nil)
                            }
                        }
                    }
                } else {
                    row = append(row, t.Splat(arow, brow))
                }
            }

            out = append(out, row)
        } else {
            if i < len(y) {
                switch r := x[i].(type) {
                case Number:
                    switch c := y[i].(type) {
                    case Number:
                        out = append(out, t.MultiplyNumber(r, c))
                    default:
                        t.TypeMismatch(r, c)
                    }
                default:
                    t.TypeMismatch(r, nil)
                }
            }
        }
    }

    return out
}

func (t *Token) RepeatArray(x Array, y Number) Array {
    out := Array { }

    for n := 0; n < y.Int(); n++ {
        for _, val := range x {
            out = append(out, Clone(val))
        }
    }

    if !y.Fits() && !y.Rat().IsInt() {
        val, _ := y.Rat().Float64()
        rem := int(float64(len(x)) * (val - float64(y.Int())))

        for _, val := range x {
            if len(out) < y.Int() * len(x) + rem {
                out = append(out, Clone(val))
            } else {
                break
            }
        }
    }

    return out
}

func (t *Token) RepeatString(x String, y Number) String {
    out := ""

    for n := 0; n < y.Int(); n++ {
        out += string(x)
    }

    if !y.Fits() && !y.Rat().IsInt() {
        val, _ := y.Rat().Float64()
        rem := int(float64(len(x)) * (val - float64(y.Int())))

        for _, c := range x {
            if len(out) < y.Int() * len(x) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return String(out)
}

func (t *Token) JoinArray(x Array, y String) String {
    out := ""

    for i, val := range x {
        if i > 0 {
            out = out + string(y)
        }

        out = out + string(Stringify(val))
    }

    return String(out)
}

func (t *Token) JoinString(x String, y String) String {
    out := ""

    for i, c := range x {
        if i > 0 {
            out = out + string(y)
        }

        out = out + string(c)
    }

    return String(out)
}

func (t *Token) MultiplyNumber(x Number, y Number) interface{} {
    if (x.inf == INF && y.inf == -INF) || (x.inf == -INF && y.inf == INF) {
        return Number{ inf: -INF }
    }

    if (x.inf == INF && y.inf == INF) || (x.inf == -INF && y.inf == -INF) {
        return Number{ inf: INF }
    }

    if x.inf == INF || x.inf == -INF {
        switch y.Cmp(NewNumber(0)) {
        case -1:
            return Number{ inf: -x.inf }
        case 1:
            return Number{ inf: x.inf }
        }

        return Null { }
    }

    if y.inf == INF || y.inf == -INF {
        switch x.Cmp(NewNumber(0)) {
        case -1:
            return Number{ inf: -y.inf }
        case 1:
            return Number{ inf: y.inf }
        }

        return Null { }
    }

    if x.Fits() && y.Fits() {
        if x.num == 0 || y.num == 0 {
            return Box(Number{ num: 0 })
        }

        if prod := x.num * y.num; prod / y.num == x.num {
            return Box(Number{ num: prod })
        }
    }

    return NewRat(new(big.Rat).Mul(x.Rat(), y.Rat()))
}

func (t *Token) CombArray(x Array, y Number) Array {
    out := Array{ }
    n := y.Int()
    i := 1

    for i < 1 << uint(len(x)) {
        subset := Array{ }
        j := 0

        for j < len(x) {
            if (1 << uint(j)) & i > 0 {
                subset = append(subset, x[j])
            }

            j++
        }

        if len(subset) == n {
            out = append(out, subset)
        }

        i++
    }

    return out
}

func (t *Token) CombString(x String, y Number) Array {
    out := Array{ }
    n := y.Int()
    i := 1

    for i < 1 << uint(len(x)) {
        subset := String("")
        j := 0

        for j < len(x) {
            if (1 << uint(j)) & i > 0 {
                subset = append(subset, x[j])
            }

            j++
        }

        if len(subset) == n {
            out = append(out, subset)
        }

        i++
    }

    return out
}

func (t *Token) PermsArray(x Array) Array {
    var helper func(Array, int)
    out := Array{ }

    helper = func(x Array, n int) {
        if n == 1 {
            tmp := make([]interface{}, len(x))
            copy(tmp, x)
            out = append(out, Array(tmp))
        } else {
            for i := 0; i < n; i++ {
                helper(x, n - 1)

                if n % 2 == 1 {
                    tmp := x[i]
                    x[i] = x[n - 1]
                    x[n - 1] = tmp
                } else {
                    tmp := x[0]
                    x[0] = x[n - 1]
                    x[n - 1] = tmp
                }
            }
        }
    }

    helper(x, len(x))

    return out
}

func (t *Token) PermsString(x String) Array {
    var helper func(String, int)
    out := Array{ }

    helper = func(x String, n int) {
        if n == 1 {
            tmp := make([]rune, len(x))
            copy(tmp, x)
            out = append(out, String(tmp))
        } else {
            for i := 0; i < n; i++ {
                helper(x, n - 1)

                if n % 2 == 1 {
                    tmp := x[i]
                    x[i] = x[n - 1]
                    x[n - 1] = tmp
                } else {
                    tmp := x[0]
                    x[0] = x[n - 1]
                    x[n - 1] = tmp
                }
            }
        }
    }

    helper(x, len(x))

    return out
}

func (t *Token) ChooseNumber(x Number, y Number) Number {
    a := int64(x.Int())
    b := int64(y.Int())

    return NewRat(new(big.Rat).SetInt(new(big.Int).Binomial(a, b)))
}

func (t *Token) MultiplyArray(x Array) interface{} {
    out := NewNumber(1)

    for _, val := range x {
        switch val.(type) {
        case *Block:
            val = val.(*Block).Run()
        case *Variable:
            val = val.(*Variable).Value()
        }

        switch val := val.(type) {
        case Number:
            out = t.MultiplyNumber(out, val).(Number)
        case Boolean:
            out = t.MultiplyNumber(out, val.Number()).(Number)
        case Null:
            out = NewNumber(0)
        default:
            t.UnexpectedItem(val)
        }
    }

    return out
}

func (t *Token) DivisorsNumber(x Number) Array {
    if x.Fits() && x.num != math.MinInt64 {
        n := x.num

        if n < 0 {
            n = -n
        }

        if n < 2 {
            return Array { }
        }

        out := Array { NewNumber(1) }

        for j := int64(2); j * j <= n; j++ {
            if n % j == 0 {
                out = append(out, Number{ num: j })

                if div := n / j; div != j {
                    out = append(out, Number{ num: div })
                }
            }
        }

        return t.SortArray(out, nil)
    }

    i := new(big.Int).Quo(x.Rat().Num(), x.Rat().Denom())
    j := big.NewInt(1)

    if i.Cmp(big.NewInt(0)) < 0 {
        i = i.Neg(i)
    }

    if i.Cmp(big.NewInt(2)) < 0 {
        return Array { }
    }

    out := Array { NewNumber(1) }
    sqrt := new(big.Int).Sqrt(i)

    for j.Cmp(sqrt) <= 0 {
        if j.Cmp(big.NewInt(1)) > 0 && new(big.Int).Rem(i, j).Cmp(big.NewInt(0)) == 0 {
            out = append(out, NewRat(new(big.Rat).SetInt(j)))
            div := new(big.Int).Div(i, j)

            if j.Cmp(div) != 0 {
                out = append(out, NewRat(new(big.Rat).SetInt(div)))
            }
        }

        j = new(big.Int).Add(j, big.NewInt(1))
    }

    return t.SortArray(out, nil)
}
