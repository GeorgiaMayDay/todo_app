package todo

import (
	"fmt"
	"io"
)

type TodoList []string

func (tl TodoList) outputTodos(writer io.Writer) {
	for i, todo := range tl {
		fmt.Fprintln(writer, todo)
		if i >= 9 {
			break
		}
	}
}
