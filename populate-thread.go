package witchery
import "sync"
func PopulateThread(executable *Executable, old *Thread) *Thread {
	neww := Thread{
		Ceiling: old.CeilingNext,
		CeilingNext: old.CeilingNext,
		Dictionary: old.Dictionary,
		Executable: executable,
		Mutex: &sync.Mutex{},
		Restrictions: old.RestrictionsNext,
		RestrictionsNext: map[string]bool{},
		Stack: old.Stack,
	}
	for key := range old.RestrictionsNext {
		neww.RestrictionsNext[key] = true
	}
	return &neww
}
