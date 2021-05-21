package main
import("fmt")

func Unique(a interface{}) interface{} {
    switch x := a.(type) {
    case Hash:
        return x.Array().Unique()
    case Array:
        return x.Unique()
    case String:
        return x.Unique()
    default:
        return x
    }
}

func (a Array) Unique() Array {
    out := Array { }
    hash := make(map[string]bool)

    for _, val := range a {
        key := fmt.Sprintf("%v", val)

        if _, ok := hash[key]; !ok {
            out = append(out, val)
            hash[key] = true
        }
    }

    return out
}

func (a String) Unique() String {
    out := ""
    hash := make(map[rune]bool)

    for _, c := range a {
        if _, ok := hash[c]; !ok {
            out += string(c)
            hash[c] = true
        }
    }

    return String(out)
}
