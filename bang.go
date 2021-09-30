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
            if t.lit != "!" && t.lit != "swap" {
                return t.TypeMismatch(a, b)
            }

            return blk.Run()
        }

        if t.lit != "!" && t.lit != "swap" {
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

func (t *Token) Until(a interface{}, y *Block) interface{} {
    var val interface{} = Null{ }
    i := 0

    for !Boolify(a) {
        val = y.Run(NewNumber(i))
        i = i + 1
    }

    return val
}
