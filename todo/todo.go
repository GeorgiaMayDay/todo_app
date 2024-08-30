package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Todo struct {
	Name   string
	Status string
}

func (t Todo) String() string {
	return t.Name + ": " + t.Status
}

type baseList interface {
	outputTodos(writer io.Writer)
	addTodo(newTodo string)
	deleteTodo(delTodo string)
	List_as_json() ([]byte, error)
}

type TodoList struct {
	List []Todo
}

func (tl *TodoList) outputTodos(writer io.Writer) {
	count := 10
	if len(tl.List) < count {
		count = len(tl.List)
	}
	for i := 0; i < count; i++ {
		todoformat := strconv.Itoa(i+1) + ". " + tl.List[i].String()
		fmt.Fprintln(writer, todoformat)
	}
}

func (tl *TodoList) addTodo(newTodo string) {
	tl.List = append(tl.List, Todo{newTodo, "Todo"})
}

func (tl *TodoList) deleteTodo(delTodo string) {
	num, err := strconv.Atoi(delTodo)
	var newTodoList []Todo
	if err != nil {
		for _, todo := range tl.List {
			if todo.Name == delTodo {
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

func (tl *TodoList) List_as_json() ([]byte, error) {
	return json.Marshal(tl.List)
}
