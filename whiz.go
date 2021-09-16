package main

func (t *Token) Whiz(a Array, b Array) interface{} {
    i := t.Which(a)

    if i == -1 && len(b) <= len(a) {
        return Null{ }
    }

    if i < 0 && len(b) > 0 && len(b) + i < len(b) {
        i = len(b) + i
    }

    if len(b) > 0 && i < len(b) {
        if blk, ok := b[i].(*Block); ok {
            if len(blk.args) > 0 || t.lit == "redo" || t.lit == "while" {
                if t.lit != "?" && t.lit != "redo" && t.lit != "while" {
                    t.TypeMismatch(a, b)
                }

                return t.Redo(a, blk)
            }

            if t.lit != "?" && t.lit != "switch" && t.lit != "then" {
                return t.TypeMismatch(a, b)
            }

            return blk.Run()
        }

        if t.lit != "?" && t.lit != "switch" && t.lit != "then" {
            return t.TypeMismatch(a, b)
        }

        return b[i]
    }

    return Null{ }
}

func (t *Token) DoubleWhiz(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.DoubleWhiz(x.Run(), b)
    case *Variable:
        return t.DoubleWhiz(x.Value(), b)
    case Null:
        return b
    default:
        return a
    }
}

func (t *Token) TopWhiz(a interface{}) interface{} {
    return Boolify(a)
}

func (t *Token) Which(a Array) int {
    for i, val := range a {
        if !Boolify(val) {
            continue
        }

        return i
    }

    return -1
}

func (t *Token) Redo(x Array, y *Block) interface{} {
    i := 0
    val := y.Run(NewNumber(i))

    for t.Which(x) != -1 {
        val = y.Run(NewNumber(i))
        i = i + 1
    }

    return val
}
