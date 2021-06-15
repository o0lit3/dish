package main
import("fmt")

func Stringify(a interface{}) String {
    switch x := a.(type) {
    case *Block:
        return Stringify(x.Run())
    case *Variable:
        return Stringify(x.Value())
    case Hash:
        return Stringify(x.Array())
    case Array:
        return Join(x, "").(String)
    case String:
        return x
    default:
        return String(fmt.Sprintf("%v", x))
    }
}
