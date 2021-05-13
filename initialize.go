package witchery
import (
	"container/list"
	"sync"
)
func Initialize(code string) {
	thread := Thread{
		Dictionary: &Dictionary{
			Indices: map[string]int{},
		},
		Mutex: &sync.Mutex{},
		Restrictions: map[string]bool{},
		Stack: &Stack{
			List: list.New(),
		},
	}
	thread.Ceiling, thread.CeilingNext = thread.Stack, thread.Stack
	thread.Executable, thread.RestrictionsNext = Compile(code, thread.Dictionary, thread.Restrictions)
	Execute(&thread)
}
