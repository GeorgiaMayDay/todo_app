package todo

import (
	"fmt"
	"io"
	"strconv"
)

type baseList interface {
	outputTodos(writer io.Writer)
	addTodo(newTodo string)
}

type TodoList struct {
	List []string
}

func (tl *TodoList) outputTodos(writer io.Writer) {
	count := 10
	if len(tl.List) < count {
		count = len(tl.List)
	}
	for i := 0; i < count; i++ {
		todoformat := strconv.Itoa(i+1) + ". " + tl.List[i]
		fmt.Fprintln(writer, todoformat)
	}
}

func (tl *TodoList) addTodo(newTodo string) {
	tl.List = append(tl.List, newTodo)
}

// func (tl *TodoList) deleteTodo(delTodo string) {
// 	num, err := strconv.Atoi(delTodo)
// 	if err != nil {
// 		for i, todo := range tl.List {
// 			if todo == delTodo {
// 				tl.List[i]
// 			}
// 		}
// 	}
// }
