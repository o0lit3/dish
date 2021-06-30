package main

func Switch(a Array, b Array) interface{} {
    i := a.Switch()

    if i == -1 && len(b) <= len(a) {
        return Null { }
    }

    if i < 0 && len(b) > 0 && len(b) + i < len(b) {
        if blk, ok := b[len(b) + i].(*Block); ok {
            return blk.Run()
        }

        return b[len(b) + i]
    }

    if len(b) > 0 && i < len(b) {
        if blk, ok := b[i].(*Block); ok {
            return blk.Run()
        }

        return b[i]
    }

    return Null { }
}

func (a Array) Switch() int {
    for i, val := range a {
        if Not(val) {
            continue
        }

        return i
    }

    return -1
}
