package todo

import (
	"fmt"
	"io"
	"strconv"
)

type baseList interface {
	outputTodos(writer io.Writer)
	addTodo(newTodo string)
	deleteTodo(delTodo string)
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

func (tl *TodoList) deleteTodo(delTodo string) {
	num, err := strconv.Atoi(delTodo)
	var newTodoList []string
	if err != nil {
		for _, todo := range tl.List {
			if todo == delTodo {
				continue
			}
			newTodoList = append(newTodoList, todo)
		}
	} else {
		for i, todo := range tl.List {
			if i == num-1 {
				continue
			}
			newTodoList = append(newTodoList, todo)
		}
	}
	tl.List = tl.List[:len(tl.List)-1]
	copy(tl.List, newTodoList)
}
