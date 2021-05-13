package witchery
func PushStack(item interface{}, stack *Stack) {
	if stack.Cursor == nil {
		stack.Cursor = stack.List.PushFront(item)
	} else {
		stack.Cursor = stack.List.InsertAfter(item, stack.Cursor)
	}
}
