package witchery
func DropStack(stack *Stack) interface{} {
	if stack.Cursor == nil {
		panic("witchery.Drop: The cursor is empty.")
	}
	var item interface{}
	if stack.Cursor.Prev() == nil {
		item = stack.List.Remove(stack.Cursor)
		stack.Cursor = nil
	} else {
		stack.Cursor = stack.Cursor.Prev()
		item = stack.List.Remove(stack.Cursor.Next())
	}
	return item
}
