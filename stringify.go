package main
import("fmt")

func Stringify(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Stringify(x.Run())
    case Hash:
        return Stringify(x.Array())
    case Array:
        return Join(x, "")
    case String:
        return x
    default:
        return String(fmt.Sprintf("%v", x))
    }
}
