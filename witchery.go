package witchery
import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
)
type (
	Dictionary struct {
		Executables []*Executable
		Indices map[string]int
	}
	Executable struct {
		Objects []interface{}
		Varieties []uint8
	}
	Stack struct {
		Cursor *list.Element
		List *list.List
		Parent *Stack
	}
	Thread struct {
		Ceiling *Stack
		CeilingNext *Stack
		Cursor int
		Dictionary *Dictionary
		Executable *Executable
		Mutex *sync.Mutex
		Restrictions map[string]bool
		RestrictionsNext map[string]bool
		Stack *Stack
	}
)
const (
	VarietyDefinition uint8 = iota
	VarietyItem
	VarietyKeyword
)
var Keywords map[string]func(*Thread)
func init() {
	Keywords = map[string]func(*Thread){
		"add": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			PushStack(big.NewFloat(0).Add(DropStack(thread.Stack).(*big.Float), second), thread.Stack)
		},
		"and": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			if DropStack(thread.Stack).(string) != "" && second != "" {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"ascend": func(thread *Thread){
			if thread.CeilingNext == thread.Stack {
				panic("witchery (ascend): Already at ceiling.")
			}
			thread.Stack = thread.Stack.Parent
		},
		"compile-later": func(thread *Thread){
			var temporary *Executable
			temporary, thread.RestrictionsNext = Compile(DropStack(thread.Stack).(string), thread.Dictionary, thread.RestrictionsNext)
			PushStack(temporary, thread.Stack)
		},
		"concatenate-executables": func(thread *Thread){
			second := DropStack(thread.Stack).(*Executable)
			first := DropStack(thread.Stack).(*Executable)
			PushStack(&Executable{
				Objects: append(first.Objects, second.Objects...),
				Varieties: append(first.Varieties, second.Varieties...),
			}, thread.Stack)
		},
		"concatenate-strings": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			thread.Stack.Cursor.Value = thread.Stack.Cursor.Value.(string) + second
		},
		"concatenate-vectors": func(thread *Thread){
			second := DropStack(thread.Stack).([]interface{})
			thread.Stack.Cursor.Value = append(thread.Stack.Cursor.Value.([]interface{}), second...)
		},
		"create-map": func(thread *Thread){
			PushStack(map[string]interface{}{}, thread.Stack)
		},
		"create-stack": func(thread *Thread){
			PushStack(&Stack{
				List: list.New(),
				Parent: thread.Stack,
			}, thread.Stack)
		},
		"create-vector": func(thread *Thread){
			first, _ := DropStack(thread.Stack).(*big.Float).Int64()
			PushStack(make([]interface{}, int(first)), thread.Stack)
		},
		"descend": func(thread *Thread){
			thread.Stack = thread.Stack.Cursor.Value.(*Stack)
		},
		"divide": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			PushStack(big.NewFloat(0).Quo(DropStack(thread.Stack).(*big.Float), second), thread.Stack)
		},
		"drop": func(thread *Thread){
			DropStack(thread.Stack)
		},
		"duplicate": func(thread *Thread){
			PushStack(thread.Stack.Cursor.Value, thread.Stack)
		},
		"echo-line": func(thread *Thread){
			fmt.Println(DropStack(thread.Stack))
		},
		"echo-raw": func(thread *Thread){
			fmt.Print(DropStack(thread.Stack))
		},
		"equal-numbers?": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			if DropStack(thread.Stack).(*big.Float).Cmp(second) == 0 {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"equal-strings?": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			if DropStack(thread.Stack).(string) == second {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"execute": func(thread *Thread){
			thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(DropStack(thread.Stack).(*Executable), thread))
		},
		"get-map": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			PushStack(thread.Stack.Cursor.Value.(map[string]interface{})[second], thread.Stack)
		},
		"get-vector": func(thread *Thread){
			second, _ := DropStack(thread.Stack).(*big.Float).Int64()
			PushStack(thread.Stack.Cursor.Value.([]interface{})[int(second)], thread.Stack)
		},
		"greater?": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			if DropStack(thread.Stack).(*big.Float).Cmp(second) == 1 {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"if": func(thread *Thread){
			second := DropStack(thread.Stack).(*Executable)
			if DropStack(thread.Stack).(string) != "" {
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(second, thread))
			}
		},
		"if-else": func(thread *Thread){
			third := DropStack(thread.Stack).(*Executable)
			second := DropStack(thread.Stack).(*Executable)
			if DropStack(thread.Stack).(string) != "" {
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(second, thread))
			} else {
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(third, thread))
			}
		},
		"lesser?": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			if DropStack(thread.Stack).(*big.Float).Cmp(second) == -1 {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"measure-map": func(thread *Thread){
			PushStack(big.NewFloat(float64(len(thread.Stack.Cursor.Value.(map[string]interface{})))), thread.Stack)
		},
		"measure-stack": func(thread *Thread){
			PushStack(big.NewFloat(float64(thread.Stack.List.Len())), thread.Stack)
		},
		"measure-string": func(thread *Thread){
			PushStack(big.NewFloat(float64(len(thread.Stack.Cursor.Value.(string)))), thread.Stack)
		},
		"measure-vector": func(thread *Thread){
			PushStack(big.NewFloat(float64(len(thread.Stack.Cursor.Value.([]interface{})))), thread.Stack)
		},
		"multiply": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			PushStack(big.NewFloat(0).Mul(DropStack(thread.Stack).(*big.Float), second), thread.Stack)
		},
		"not": func(thread *Thread){
			if thread.Stack.Cursor.Value.(string) == "" {
				thread.Stack.Cursor.Value = "t"
			} else {
				thread.Stack.Cursor.Value = ""
			}
		},
		"numberize": func(thread *Thread){
			number, _, issue := big.NewFloat(0).Parse(DropStack(thread.Stack).(string), 10)
			if issue != nil {
				panic("witchery (numberize): The number is not parseable.")
			}
			PushStack(number, thread.Stack)
		},
		"or": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			if DropStack(thread.Stack).(string) != "" || second != "" {
				PushStack("t", thread.Stack)
			} else {
				PushStack("", thread.Stack)
			}
		},
		"read-file": func(thread *Thread){
			file, issue := ioutil.ReadFile(DropStack(thread.Stack).(string))
			if issue != nil {
				panic("witchery (read-file): The file is not able to be read.")
			}
			PushStack(string(file), thread.Stack)
		},
		"read-line": func(thread *Thread){
			reader := bufio.NewReader(os.Stdin)
			data, issue := reader.ReadString('\n')
			if issue != nil {
				panic("witchery (read-line): There was an issue in collecting the input.")
			}
			PushStack(string(data[:len(data) - 1]), thread.Stack)
		},
		"report": func(thread *Thread){
			fmt.Println("BEGINNING OF REPORT")
			element := thread.Stack.List.Front()
			for element != nil {
				if element == thread.Stack.Cursor {
					fmt.Print("CURSOR: ")
				}
				fmt.Println(element.Value)
				element = element.Next()
			}
			fmt.Println("ENDING OF REPORT")
		},
		"reset-ceiling": func(thread *Thread){
			thread.CeilingNext = thread.Ceiling
		},
		"reset-restrictions": func(thread *Thread){
			thread.RestrictionsNext = map[string]bool{}
			for key := range thread.Restrictions {
				thread.RestrictionsNext[key] = true
			}
		},
		"restrict": func(thread *Thread){
			first := DropStack(thread.Stack).(string)
			if Keywords[first] == nil {
				panic("witchery (restrict): The keyword '" + first + "' does not exist.")
			}
			thread.RestrictionsNext[first] = true
		},
		"set-ceiling": func(thread *Thread){
			thread.CeilingNext = thread.Stack
		},
		"set-map": func(thread *Thread){
			third := DropStack(thread.Stack)
			second := DropStack(thread.Stack).(string)
			thread.Stack.Cursor.Value.(map[string]interface{})[second] = third
		},
		"set-vector": func(thread *Thread){
			third := DropStack(thread.Stack)
			second, _ := DropStack(thread.Stack).(*big.Float).Int64()
			thread.Stack.Cursor.Value.([]interface{})[int(second)] = third
		},
		"stringify": func(thread *Thread){
			PushStack(DropStack(thread.Stack).(*big.Float).String(), thread.Stack)
		},
		"subtract": func(thread *Thread){
			second := DropStack(thread.Stack).(*big.Float)
			PushStack(big.NewFloat(0).Sub(DropStack(thread.Stack).(*big.Float), second), thread.Stack)
		},
		"swap": func(thread *Thread){
			if thread.Stack.Cursor == nil {
				panic("witchery (swap): The cursor is empty.")
			}
			if thread.Stack.Cursor.Prev() == nil {
				panic("witchery (swap): There is nothing before cursor.")
			}
			item := thread.Stack.List.Remove(thread.Stack.Cursor.Prev())
			thread.Stack.Cursor = thread.Stack.List.InsertAfter(item, thread.Stack.Cursor)
		},
		"visit-first": func(thread *Thread){
			if thread.Stack.List.Front() == nil {
				panic("witchery (visit-first): There is nothing in stack.")
			}
			thread.Stack.Cursor = thread.Stack.List.Front()
		},
		"visit-last": func(thread *Thread){
			if thread.Stack.List.Back() == nil {
				panic("witchery (visit-first): There is nothing in stack.")
			}
			thread.Stack.Cursor = thread.Stack.List.Back()
		},
		"visit-next": func(thread *Thread){
			if thread.Stack.Cursor == nil {
				if thread.Stack.List.Len() == 0 {
					panic("witchery (visit-next): There is nothing in stack.")
				}
				thread.Stack.Cursor = thread.Stack.List.Front()
			} else {
				if thread.Stack.Cursor.Next() == nil {
					panic("witchery (visit-next): There is nothing after cursor.")
				}
				thread.Stack.Cursor = thread.Stack.Cursor.Next()
			}
		},
		"visit-previous": func(thread *Thread){
			if thread.Stack.Cursor == nil {
				panic("witchery (visit-previous): The cursor is nil.")
			}
			thread.Stack.Cursor = thread.Stack.Cursor.Prev()
		},
		"while": func(thread *Thread){
			second := DropStack(thread.Stack).(*Executable)
			first := DropStack(thread.Stack).(*Executable)
			thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(first, thread))
			for DropStack(thread.Stack).(string) != "" {
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(second, thread))
				thread.CeilingNext, thread.Dictionary, thread.RestrictionsNext, thread.Stack = Execute(PopulateThread(first, thread))
			}
		},
		"write-file": func(thread *Thread){
			second := DropStack(thread.Stack).(string)
			if ioutil.WriteFile(DropStack(thread.Stack).(string), []byte(second), 0644) != nil {
				panic("witchery (write-file): The file could not be written to.")
			}
		},
	}
}
