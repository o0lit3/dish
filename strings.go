package main

import(
    "fmt"
    "unicode"
)

func Strings(a interface{}, tok *Token) interface{} {
    switch x := a.(type) {
    case *Block:
        return Strings(x.Run(), tok)
    case *Variable:
        return Strings(x.Value(), tok)
    case Hash:
        return Strings(fmt.Sprintf("%v", x), tok)
    case Array:
        return Strings(fmt.Sprintf("%v", x), tok)
    case String:
        if len(x) == 0 {
            return Boolean(false)
        }

        switch tok.lit {
        case "digit":
            return Boolean(unicode.IsDigit(x[0]))
        case "letter":
            return Boolean(unicode.IsLetter(x[0]))
        case "consonant":
            switch x[0] {
            case 'A', 'a', 'E', 'e', 'I', 'i', 'O', 'o', 'U', 'u':
                return Boolean(false)
            default:
                return Boolean(unicode.IsLetter(x[0]))
            }
        case "vowel":
            switch x[0] {
            case 'A', 'a', 'E', 'e', 'I', 'i', 'O', 'o', 'U', 'u':
                return Boolean(true)
            default:
                return Boolean(false)
            }
        case "upper":
            return Boolean(unicode.IsUpper(x[0]))
        case "lower":
            return Boolean(unicode.IsLower(x[0]))
        case "space":
            return Boolean(unicode.IsSpace(x[0]))
        }
    case Number:
        return Strings(fmt.Sprintf("%v", x), tok)
    case Boolean:
        return Strings(fmt.Sprintf("%v", x), tok)
    case Null:
        return Strings(fmt.Sprintf("%v", x), tok)
    }

    return Null { }
}
