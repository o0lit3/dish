package main

func Redo(a *Block, b *Block) interface{} {
    var val interface{}

    val = Null { }
    i := 0

    for !Not(a.Run()) {
        val = b.Run(a, Number(i))
        i = i + 1
    }

    return val
}
