package main

func Itemize(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Itemize(x.Run())
    case Hash:
        return x
    case Array:
        return x
    case String:
        return x.Array()
    default:
        return Array { x }
    }
}
