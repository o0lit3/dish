package main
import("fmt")

func Stringify(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Stringify(x.Run())
    case String:
        return x
    default:
        return String(fmt.Sprintf("%v", x))
    }
}
