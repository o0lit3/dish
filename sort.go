package main
import("sort")

func Sort(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return Sort(x.Run())
    case *Variable:
        return Sort(x.Value())
    case Hash:
        return x.Array().Sort()
    case Array:
        return x.Sort()
    case String:
        return Join(x.Array().Sort(), String(""))
    default:
        return x
    }
}

func (a Array) Sort() Array {
    sort.Slice(a, func(i, j int) bool {
        return bool(Below(a[i], a[j]))
    })

    return a
}

func (a Array) UserSort(b *Block) Array {
    sort.Slice(a, func(i, j int) bool {
        if b, ok := b.Run(a[i], a[j]).(Boolean); ok {
            return bool(b)
        }

        return bool(Below(a[i], a[j]))
    })

    return a
}