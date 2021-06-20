package main

func Redo(a *Block, b *Block) interface{} {
    var val interface{}

    val = Boolean(false)
    i := 0

    for Boolify(a.Run()) {
        val = b.Run(a, NewNumber(i))
        i = i + 1
    }

    return val
}
