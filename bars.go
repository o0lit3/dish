package main
import("bufio"; "strings")

func (t *Token) TopBars(blk *Block, a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopBars(blk, x.Run())
    case *Variable:
        return t.TopBars(blk, x.Value())
    case Hash:
        return t.TopBars(blk, x.Array())
    case Array:
        if t.lit != "=" && t.lit != "average" && t.lit != "avg" {
            t.TypeMismatch(a, nil)
        }

        return t.AverageArray(x)
    case String:
        if t.lit != "=" && t.lit != "eval" {
            t.TypeMismatch(a, nil)
        }

        return t.EvalString(blk, x)
    case Number:
        return x
    case Boolean:
        return t.TopBars(blk, x.Number())
    case Null:
        return t.TopBars(blk, NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) EvalString(blk *Block, x String) interface{} {
    reader := bufio.NewReader(strings.NewReader(string(x)))
    parser := process(reader, blk.Branch(VAL))
    return parser.blk.Run()
}

func (t *Token) AverageArray(x Array) interface{} {
    return t.DivideNumber(t.SumArray(x).(Number), t.LengthArray(x))
}
