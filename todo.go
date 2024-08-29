package todo

import (
	"fmt"
	"io"
)

type TodoList []string

func (tl TodoList) outputTodos(writer io.Writer) {
	for _, todo := range tl {
		fmt.Fprintln(writer, todo)
	}
}
