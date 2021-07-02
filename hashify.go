package main
import("fmt")

func Hashify(a interface{}) Hash {
    switch x := a.(type) {
    case *Block:
        return Hashify(x.Run())
    case *Variable:
        return Hashify(x.Value())
    case Hash:
        return x
    case Array:
        return x.Hash()
    case String:
        return Hash{ string(x): Boolean(true) }
    default:
        return Hash{ fmt.Sprintf("%v", x): Boolean(true) }
    }
}
