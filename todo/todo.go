package todo

import (
	"fmt"
	"io"
)

type baseList interface {
	outputTodos(writer io.Writer)
	addTodo(newTodo string)
}

type TodoList struct {
	List []string
}

func (tl *TodoList) outputTodos(writer io.Writer) {
	for i, todo := range tl.List {
		fmt.Fprintln(writer, todo)
		if i >= 9 {
			break
		}
	}
}

func (tl *TodoList) addTodo(newTodo string) {
	tl.List = append(tl.List, newTodo)
}
