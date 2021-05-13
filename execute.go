package witchery
func Execute(thread *Thread) (*Stack, *Dictionary, map[string]bool, *Stack) {
	for thread.Cursor < len(thread.Executable.Objects) {
		thread.Mutex.Lock()
		switch thread.Executable.Varieties[thread.Cursor] {
			case VarietyDefinition:
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(thread.Dictionary.Executables[thread.Executable.Objects[thread.Cursor].(int)], thread))
			case VarietyItem:
				PushStack(thread.Executable.Objects[thread.Cursor], thread.Stack)
			case VarietyKeyword:
				thread.Executable.Objects[thread.Cursor].(func(*Thread))(thread)
		}
		thread.Mutex.Unlock()
		thread.Cursor++
	}
	return thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack
}
