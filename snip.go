package main
import("fmt"; "regexp")

var patterns = map[string]*regexp.Regexp{ }

func (t *Token) Pattern(a interface{}) *regexp.Regexp {
    pat := ""

    switch x := a.(type) {
    case *Block:
        return t.Pattern(x.Run())
    case *Variable:
        return t.Pattern(x.Value())
    case String:
        pat = string(x)
    case Number, Boolean:
        pat = fmt.Sprintf("%v", x)
    default:
        t.TypeMismatch(a, nil)
    }

    if out, ok := patterns[pat]; ok {
        return out
    }

    out, err := regexp.Compile(pat)

    if err != nil {
        t.InvalidPattern(pat)
    }

    patterns[pat] = out

    return out
}

func (t *Token) Patterns(blk *Block, a interface{}, b interface{}) interface{} {
    args := Array{ }

    if b == nil {
        if arr, ok := blk.Value(a).(Array); ok && len(arr) > 1 {
            a = arr[0]
            args = arr[1:]
        }
    } else {
        for _, arg := range blk.Blockify(b) {
            args = append(args, blk.Value(arg))
        }
    }

    switch t.lit {
    case "match", "scan", "groups":
        if len(args) != 1 {
            t.UnexpectedOperand()
        }
    default:
        if len(args) != 2 {
            t.UnexpectedOperand()
        }
    }

    switch t.lit {
    case "match":
        return t.Match(a, args[0])
    case "scan":
        return t.Scan(a, args[0])
    case "groups":
        return t.Groups(a, args[0])
    }

    return t.Replace(a, args, t.lit == "sub")
}

func (t *Token) Match(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Match(x.Run(), b)
    case *Variable:
        return t.Match(x.Value(), b)
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Match(val, b))
        }

        return out
    case String:
        if loc := t.Pattern(b).FindStringIndex(string(x)); loc != nil {
            return String(string(x)[loc[0]:loc[1]])
        }

        return Null{ }
    case Null:
        return t.Match(String(""), b)
    case Number, Boolean:
        return t.Match(String(fmt.Sprintf("%v", x)), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) Scan(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Scan(x.Run(), b)
    case *Variable:
        return t.Scan(x.Value(), b)
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Scan(val, b))
        }

        return out
    case String:
        out := Array{ }

        for _, found := range t.Pattern(b).FindAllString(string(x), -1) {
            out = append(out, String(found))
        }

        return out
    case Null:
        return t.Scan(String(""), b)
    case Number, Boolean:
        return t.Scan(String(fmt.Sprintf("%v", x)), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) Groups(a interface{}, b interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Groups(x.Run(), b)
    case *Variable:
        return t.Groups(x.Value(), b)
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Groups(val, b))
        }

        return out
    case String:
        out := Array{ }

        for _, found := range t.Pattern(b).FindAllStringSubmatch(string(x), -1) {
            row := Array{ }

            for _, val := range found {
                row = append(row, String(val))
            }

            out = append(out, row)
        }

        return out
    case Null:
        return t.Groups(String(""), b)
    case Number, Boolean:
        return t.Groups(String(fmt.Sprintf("%v", x)), b)
    }

    return t.TypeMismatch(a, b)
}

func (t *Token) Replace(a interface{}, args Array, first bool) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Replace(x.Run(), args, first)
    case *Variable:
        return t.Replace(x.Value(), args, first)
    case Array:
        out := Array{ }

        for _, val := range x {
            out = append(out, t.Replace(val, args, first))
        }

        return out
    case String:
        if len(args) != 2 {
            t.UnexpectedOperand()
        }

        pat := t.Pattern(args[0])
        src := string(x)

        if rep, ok := args[1].(*Block); ok {
            if !first {
                return String(pat.ReplaceAllStringFunc(src, func(found string) string {
                    return string(Stringify(rep.Run(String(found))))
                }))
            }

            if loc := pat.FindStringIndex(src); loc != nil {
                return String(src[:loc[0]] + string(Stringify(rep.Run(String(src[loc[0]:loc[1]])))) + src[loc[1]:])
            }

            return x
        }

        rep := string(Stringify(args[1]))

        if !first {
            return String(pat.ReplaceAllString(src, rep))
        }

        if loc := pat.FindStringSubmatchIndex(src); loc != nil {
            return String(src[:loc[0]] + string(pat.ExpandString(nil, rep, src, loc)) + src[loc[1]:])
        }

        return x
    case Null:
        return t.Replace(String(""), args, first)
    case Number, Boolean:
        return t.Replace(String(fmt.Sprintf("%v", x)), args, first)
    }

    if len(args) > 0 {
        return t.TypeMismatch(a, args[0])
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) Escape(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Escape(x.Run())
    case *Variable:
        return t.Escape(x.Value())
    case String:
        return String(regexp.QuoteMeta(string(x)))
    case Number:
        return String(regexp.QuoteMeta(x.String()))
    case Boolean:
        return t.Escape(x.Number())
    case Null:
        return String("")
    }

    return t.TypeMismatch(a, nil)
}

func (t *Token) Strings(a interface{}) interface{} {
    switch x := a.(type) {
    case *Block:
        return t.Strings(x.Run())
    case *Variable:
        return t.Strings(x.Value())
    case String:
        if len(x) == 0 {
            return t.TypeMismatch(x, nil)
        }

        ord := int(x[0])

        switch t.lit {
        case "digit":
            return Boolean(ord >= 48 && ord <= 57)
        case "letter":
            return Boolean((ord >= 65 && ord <= 90) || (ord >= 97 && ord <= 122))
        case "consonant":
            if (ord < 65 || (ord > 90 && ord < 97) || ord > 122) {
                return Boolean(false)
            }

            switch x[0] {
            case 'A', 'E', 'I', 'O', 'U', 'a', 'e', 'i', 'o', 'u':
                return Boolean(false)
            default:
                return Boolean(true)
            }
        case "vowel":
            switch x[0] {
            case 'A', 'E', 'I', 'O', 'U', 'a', 'e', 'i', 'o', 'u':
                return Boolean(true)
            default:
                return Boolean(false)
            }
        case "upper":
            return Boolean(ord >= 65 && ord <= 90)
        case "lower":
            return Boolean(ord >= 97 && ord <= 122)
        case "space":
            switch x[0] {
            case ' ', '\n', '\t', '\r':
                return Boolean(true)
            default:
                return Boolean(false)
            }
        case "ord":
            if len(x) == 0 {
                t.TypeMismatch(x, nil)
            }

            return NewNumber(int(x[0]))
        }
    case Number:
        switch t.lit {
        case "chr":
            return String([]rune{ rune(x.Int()) })
        }
    case Boolean:
        return t.Strings(x.Number())
    case Null:
        return Boolean(false)
    }

    return t.TypeMismatch(a, nil)
}
