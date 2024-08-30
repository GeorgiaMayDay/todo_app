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
	List map[int]string
}

func (tl *TodoList) outputTodos(writer io.Writer) {
	for i := 1; i < 11; i++ {
		if tl.List[i] == "" {
			break
		}
		fmt.Fprintln(writer, tl.List[i])
	}
}

func (tl *TodoList) addTodo(newTodo string) {
	tl.List[len(tl.List)+1] = newTodo
}
