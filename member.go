package main

import(
    "fmt"
    "strconv"
)

func Member(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Member(x.Run(), b)
    case *Variable:
        switch val := x.Value().(type) {
        case Null:
            switch y := b.(type) {
            case *Block:
                if len(y.args) > 0 {
                    return y.Run(val)
                }

                return Member(x, y.Run())
            case String:
                x.blk.cur.vars[x.nom] = Hash { }
                return &Variable{ par: x, obj: x.blk.cur.vars[x.nom], nom: string(y) }
            case Number:
                x.blk.cur.vars[x.nom] = Array { }
                return &Variable{ par: x, obj: x.blk.cur.vars[x.nom], idx: y.Int() }
            }
        default:
            switch out := Member(val, b).(type) {
            case *Variable:
                out.par = x
                return out
            default:
                return out
            }
        }
    case Hash:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x.Array()...)
            }

            return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        case Hash:
            return x.Members(y.Array())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(string(y))
        default:
            return x.Member(fmt.Sprintf("%v", y))
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x...)
             }

             return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        case Hash:
            return x.Members(y.Array())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(y.Number().Int())
        case Number:
            return x.Member(y.Int())
        case Boolean:
            return x.Member(y.Number().Int())
        case Null:
            return x.Member(0)
        }
    case String:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x)
            }

            return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        case Hash:
            return x.Members(y.Array())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(y.Number().Int())
        case Number:
            return x.Member(y.Int())
        case Boolean:
            return x.Member(y.Number().Int())
        case Null:
            return x.Member(0)
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x)
            }

            return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        case Hash:
            return x.Members(y.Array())
        case Array:
            return x.Members(y)
        case String:
            return x.Member(y.Number().Int())
        case Number:
            return x.Member(y.Int())
        case Boolean:
            return x.Member(y.Number().Int())
        case Null:
            return x.Member(0)
        }
    case Boolean:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x)
            }

            return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        default:
            return Member(x.Number(), y)
        }
    case Null:
        switch y := b.(type) {
        case *Block:
            if len(y.args) > 0 {
                return y.Run(x)
            }

            return Member(x, y.Run())
        case *Variable:
            return Member(x, y.Value())
        }
    }

    return Null { }
}

func (a Hash) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a Hash) Member(b string) *Variable {
    return &Variable{ obj: a, nom: b }
}

func (a Array) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a Array) Member(b int) *Variable {
    if b < 0 && len(a) > 0 && len(a) + b < len(a) {
        return &Variable{ obj: a, idx: len(a) + b }
    }

    if len(a) > 0 && b < len(a) {
        return &Variable{ obj: a, idx: b }
    }

    if b < 0 {
        return &Variable{ obj: a, idx: -b }
    }

    return &Variable{ obj: a, idx: b }
}

func (a String) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a String) Member(b int) String {
    if b < 0 && len(a) > 0 && len(a) + b < len(a) {
        return String(a[len(a) + b])
    }

    if len(a) > 0 && b < len(a) {
        return String(a[b])
    }

    if b < 0 {
        return String(a[-b])
    }

    return String("")
}

func (a Number) Members(b Array) Array {
    out := Array { }

    for _, val := range b {
        x := Member(a, val)

        if _, ok := x.(Null); !ok {
            out = append(out, x)
        }
    }

    return out
}

func (a Number) Member(b int) Number {
    bin := strconv.FormatInt(int64(a.Int()), 2)
    bit := false

    if b < 0 && len(bin) > 0 && -b - 1 < len(bin) {
        bit = bin[-b - 1] == '1'
    }

    if len(bin) > 0 && len(bin) - b - 1 >= 0 && len(bin) - b - 1 < len(bin) {
        bit = bin[len(bin) - b - 1] == '1'
    }

    if bit {
        return NewNumber(1)
    }

    return NewNumber(0)
}
