package main

func Boolify(a interface{}) Boolean {
    switch x := a.(type) {
        case *Block:
            return Boolify(x.Run())
        case *Variable:
            return Boolify(x.Value())
        case Hash:
            return Boolean(len(x) > 0)
        case Array:
            return Boolean(len(x) > 0)
        case String:
            return Boolean(string(x) != "" && string(x) != "0")
        case Number:
            if x.inf == INF || x.inf == -INF {
                return Boolean(true)
            }

            return Boolean(x.val.Cmp(NewNumber(0).val) != 0)
        case Boolean:
            return x
    }

    return Boolean(false)
}
