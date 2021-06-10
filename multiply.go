package main

func Multiply(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Multiply(x.Run(), b)
    case *Variable:
        return Multiply(x.Value(), b)
    case Hash:
        switch y := b.(type) {
        case *Block:
            return x.Map(y)
        case *Variable:
            return Multiply(x, y.Value())
        case Hash:
            return x.Array().Repeat(NewNumber(len(y)))
        case Array:
            return x.Array().Repeat(NewNumber(len(y)))
        case String:
            return x.Array().Repeat(y.Number())
        case Number:
            return x.Array().Repeat(y)
        case Boolean:
            return x.Array().Repeat(y.Number())
        case Null:
            return Hash { }
        }
    case Array:
        switch y := b.(type) {
        case *Block:
            if y.dim == LST {
                if blk, ok := y.Run().(Array); ok {
                    return x.DotProduct(blk)
                }
            }

            return x.Map(y)
        case *Variable:
            return Multiply(x, y.Value())
        case Hash:
            return x.Repeat(NewNumber(len(y)))
        case Array:
            return x.Repeat(NewNumber(len(y)))
        case String:
            return x.Repeat(y.Number())
        case Number:
            return x.Repeat(y)
        case Boolean:
            return x.Repeat(y.Number())
        case Null:
            return Array { }
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Multiply(x, y.Run())
        case *Variable:
            return Multiply(x, y.Value())
        case Hash:
            return x.Repeat(NewNumber(len(y)))
        case Array:
            return x.Repeat(NewNumber(len(y)))
        case String:
            return x.Repeat(y.Number())
        case Number:
            return x.Repeat(y)
        case Boolean:
            return x.Repeat(y.Number())
        case Null:
            return String("")
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Multiply(x, y.Run())
        case *Variable:
            return Multiply(x, y.Value())
        case Hash:
            return y.Multiply(x)
        case Array:
            return y.Multiply(x)
        case String:
            return Number{ val: NewNumber(0).val.Mul(x.val, y.Number().val) }
        case Number:
            return Number{ val: NewNumber(0).val.Mul(x.val, y.val) }
        case Boolean:
            return Number{ val: NewNumber(0).val.Mul(x.val, y.Number().val) }
        case Null:
            return NewNumber(0)
        }
    case Boolean:
        return Multiply(x.Number(), b)
    }

    return Null { }
}

func (a Hash) Map(b *Block) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = b.Run(val, String(key))
    }

    return out
}

func (a Array) Map(b *Block) Array {
    out := make([]interface{}, len(a))

    for i, val := range a {
        out[i] = b.Run(val, NewNumber(i))
    }

    return Array(out)
}

func (a Hash) Multiply(b Number) Hash {
    out := Hash { }

    for key, val := range a {
        out[key] = Multiply(val, b)
    }

    return out
}

func (a Array) Multiply(b Number) Array {
    out := Array { }

    for _, val := range a {
        out = append(out, Multiply(val, b))
    }

    return out
}

func (a Array) Repeat(b Number) Array {
    var out []interface{}
    var blk *Block = &Block{ }

    for n := 0; n < b.Int(); n++ {
        for _, val := range a {
            out = append(out, blk.Eval(val))
        }
    }

    if !b.val.IsInt() {
        val, _ := b.val.Float64()
        rem := int(float64(len(a)) * (val - float64(b.Int())))

        for _, val := range a {
            if len(out) < b.Int() * len(a) + rem {
                out = append(out, blk.Eval(val))
            } else {
                break
            }
        }
    }

    return Array(out)
}

func (a String) Repeat(b Number) String {
    out := ""

    for n := 0; n < b.Int(); n++ {
        out += string(a)
    }

    if !b.val.IsInt() {
        val, _ := b.val.Float64()
        rem := int(float64(len(a)) * (val - float64(b.Int())))

        for _, c := range a {
            if len(out) < b.Int() * len(a) + rem {
                out += string(c)
            } else {
                break
            }
        }
    }

    return String(out)
}

func (a Array) DotProduct(b Array) Array {
    out := Array { }

    for i, arow := range a {
        if acols, ok := arow.(Array); ok {
            row := Array { }

            for k, brow := range b {
                if bcols, ok := brow.(Array); ok {
                    for j, bval := range bcols {
                        if j >= len(row) {
                            row = append(row, NewNumber(0))
                        }

                        if k < len(acols) {
                            row[j] = Add(row[j], Multiply(acols[k], bval))
                        }
                    }
                } else {
                    row = append(row, Multiply(arow, brow))
                }
            }

            out = append(out, row)
        } else {
            if i < len(b) {
                out = append(out, Multiply(a[i], b[i]))
            }
        }
    }

    return out
}
