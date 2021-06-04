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
        case String:
            return x.Multiply(y.Number())
        case Number:
            return x.Multiply(y)
        case Boolean:
            return x.Multiply(y.Number())
        case Null:
            return x.Multiply(Number(0))
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
        case String:
            return x.Multiply(y.Number())
        case Number:
            return x.Multiply(y)
        case Boolean:
            return x.Multiply(y.Number())
        case Null:
            return x.Multiply(Number(0))
        }
    case String:
        switch y := b.(type) {
        case *Block:
            return Multiply(x, y.Run())
        case *Variable:
            return Multiply(x, y.Value())
        case String:
            return x.Multiply(y.Number())
        case Number:
            return x.Multiply(y)
        case Boolean:
            return x.Multiply(y.Number())
        case Null:
            return x.Multiply(Number(0))
        }
    case Number:
        switch y := b.(type) {
        case *Block:
            return Multiply(x, y.Run())
        case *Variable:
            return Multiply(x, y.Value())
        case Hash:
            return y.Array().Multiply(x)
        case Array:
            return y.Multiply(x)
        case String:
            return y.Multiply(x)
        case Number:
            return x * y
        case Boolean:
            return x * y.Number()
        case Null:
            return x * Number(0)
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
    out := Array { }

    for i, val := range a {
        out = append(out, b.Run(val, Number(i)))
    }

    return out
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

func (a String) Multiply(b Number) String {
    out := ""

    for n := 0; n < int(b); n++ {
        out += string(a)
    }

    if b != Number(int(b)) {
        rem := Number(len(a)) * (b - Number(int(b)))

        for _, c := range a {
            if Number(len(out)) < Number(int(b) * len(a)) + rem {
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
                            row = append(row, Number(0))
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
