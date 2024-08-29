package todo

import (
	"fmt"
	"io"
)

type TodoList struct {
	list []string
}

func (tl *TodoList) outputTodos(writer io.Writer) {
	for i, todo := range tl.list {
		fmt.Fprintln(writer, todo)
		if i >= 9 {
			break
		}
	}
}

func (tl *TodoList) addTodo(newTodo string) {
	tl.list = append(tl.list, newTodo)
}
