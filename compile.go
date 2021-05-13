package witchery
import (
	"strconv"
	"unicode"
)
func Compile(code string, dictionary *Dictionary, restrictions map[string]bool) (*Executable, map[string]bool) {
	depth := uint(0)
	executable := Executable{}
	item := ""
	restrictionsNext := map[string]bool{}
	for key, value := range restrictions {
		restrictionsNext[key] = value
	}
	word := ""
	for position, symbol := range code + "\n" {
		switch symbol {
			case '[':
				if depth != 0 {
					item += "["
				}
				depth++
			case ']':
				if depth == 0 {
					panic("witchery.Compile: The depth went negative at glyph '" + strconv.Itoa(position) + "'.")
				}
				if depth != 1 {
					item += "]"
				} else {
					PushExecutable(&executable, item, nil, VarietyItem)
					item = ""
				}
				depth--
			default:
				if depth != 0 {
					item += string(symbol)
				} else if unicode.IsSpace(symbol) {
					if word != "" {
						switch word {
							case "comment":
								DropExecutable(&executable)
							case "compile-now":
								first, _ := DropExecutable(&executable)
								var temporary *Executable
								temporary, restrictionsNext = Compile(first.(string), dictionary, restrictionsNext)
								PushExecutable(&executable, temporary, nil, VarietyItem)
							case "define-now":
								second, _ := DropExecutable(&executable)
								first, _ := DropExecutable(&executable)
								_, ok := dictionary.Indices[first.(string)]
								if !ok {
									dictionary.Indices[first.(string)] = len(dictionary.Executables)
									dictionary.Executables = append(dictionary.Executables, second.(*Executable))
								} else {
									dictionary.Executables[dictionary.Indices[first.(string)]] = second.(*Executable)
								}
							case "reset-restrictions":
								restrictionsNext = map[string]bool{}
								for key := range restrictions {
									restrictionsNext[key] = true
								}
								PushExecutable(&executable, "reset-restrictions", restrictionsNext, VarietyKeyword)
							case "restrict":
								object, _ := DropExecutable(&executable)
								if Keywords[object.(string)] == nil {
									panic("witchery.Compile: The keyword '" + object.(string) + "' to be restricted does not exist at glyph '" + strconv.Itoa(position) + "'.")
								}
								restrictionsNext[object.(string)] = true
								PushExecutable(&executable, object.(string), nil, VarietyItem)
								PushExecutable(&executable, "restrict", restrictionsNext, VarietyKeyword)
							default:
								if Keywords[word] != nil {
									PushExecutable(&executable, word, restrictionsNext, VarietyKeyword)
								} else {
									_, ok := dictionary.Indices[word]
									if !ok {
										panic("witchery.Compile: The word '" + word + "' does not exist at glyph '" + strconv.Itoa(position) + "'.")
									}
									PushExecutable(&executable, dictionary.Indices[word], nil, VarietyDefinition)
								}
						}
						word = ""
					}
				} else {
					word += string(symbol)
				}
		}
	}
	return &executable, restrictionsNext
}
