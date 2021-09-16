package main

func (t *Token) Bang(a Array, b Array) interface{} {
    i := t.WhichNot(a)

    if i == -1 && len(b) <= len(a) {
        return Null{ }
    }

    if i < 0 && len(b) > 0 && len(b) + i < len(b) {
        i = len(b) + i
    }

    if len(b) > 0 && i < len(b) {
        if blk, ok := b[i].(*Block); ok {
            if len(blk.args) > 0 || t.lit == "until" {
                if t.lit != "!" && t.lit != "until" {
                    t.TypeMismatch(a, b)
                }

                return t.Until(a, blk)
            }

            if t.lit != "!" && t.lit != "swap" && t.lit != "else" {
                return t.TypeMismatch(a, b)
            }

            return blk.Run()
        }

        if t.lit != "!" && t.lit != "swap" && t.lit != "else" {
            return t.TypeMismatch(a, b)
        }

        return b[i]
    }

    return Null{ }
}

func (t *Token) TopBang(a interface{}) interface{} {
    return Boolean(!Boolify(a))
}

func (t *Token) WhichNot(a Array) int {
    for i, val := range a {
        if Boolify(val) {
            continue
        }

        return i
    }

    return -1
}

func (t *Token) Until(x Array, y *Block) interface{} {
    i := 0
    val := y.Run(NewNumber(i))

    for t.WhichNot(x) != -1 {
        val = y.Run(NewNumber(i))
        i = i + 1
    }

    return val
}
