package witchery
func DropExecutable(executable *Executable) (interface{}, uint8) {
	if len(executable.Objects) == 0 {
		panic("witchery.DropExecutable: The executable has no contents.")
	}
	object, variety := executable.Objects[len(executable.Objects) - 1], executable.Varieties[len(executable.Varieties) - 1]
	executable.Objects = executable.Objects[:len(executable.Objects) - 1]
	executable.Varieties = executable.Varieties[:len(executable.Varieties) - 1]
	return object, variety
}
