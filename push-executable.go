package witchery
func PushExecutable(executable *Executable, object interface{}, restrictions map[string]bool, variety uint8) {
	if variety == VarietyKeyword {
		_, ok := Keywords[object.(string)]
		if !ok {
			panic("witchery.PushExecutable: The keyword '" + object.(string) + "' does not exist.")
		}
		if restrictions[object.(string)] == true {
			panic("witchery.PushExecutable: The keyword '" + object.(string) + "' is restricted.")
		}
		executable.Objects = append(executable.Objects, Keywords[object.(string)])
	} else {
		executable.Objects = append(executable.Objects, object)
	}
	executable.Varieties = append(executable.Varieties, variety)
}
