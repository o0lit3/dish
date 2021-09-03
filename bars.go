package main
import("bufio"; "strings")

func (t *Token) TopBars(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.TopBars(x.Run())
    case *Variable:
        return t.TopBars(x.Value())
    case Hash:
        return t.TopBars(x.Array())
    case Array:
        if t.lit != "=" && t.lit != "average" && t.lit != "avg" {
            t.TypeMismatch(a, nil)
        }

        return t.AverageArray(x)
    case String:
        if t.lit != "=" && t.lit != "eval" {
            t.TypeMismatch(a, nil)
        }

        return t.EvalString(x)
    case Boolean:
        return t.TopBars(x.Number())
    case Null:
        return t.TopBars(NewNumber(0))
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) EvalString(x String) interface{} {
    reader := bufio.NewReader(strings.NewReader(string(x)))
    parser := process(reader, program.Branch(VAL))
    return parser.blk.Run()
}

func (t *Token) AverageArray(x Array) interface{} {
    return t.DivideNumber(t.SumArray(x).(Number), t.LengthArray(x))
}
