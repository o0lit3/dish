package main

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
