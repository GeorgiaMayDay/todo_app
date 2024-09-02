package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TodoStatus string

const (
	complete   TodoStatus = "Complete"
	todoStatus TodoStatus = "Todo"
)

type Todo struct {
	Name   string
	Status TodoStatus
}

func (t Todo) String() string {
	return t.Name + ": " + string(t.Status)
}

type baseList interface {
	outputTodos(writer io.Writer)
	addTodo(newTodo string)
	deleteTodo(delTodo string)
	completeTodo(compTodo string)
	List_as_json() ([]byte, error)
	List_from_json([]byte)
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
	return json.MarshalIndent(tl.List, "", "    ")
}

func (tl *TodoList) List_from_json(json_list []byte) {
	json.Unmarshal(json_list, &tl.List)
}

func Save_Todo_List_From_Json(list baseList, file_name string) {
	file, _ := os.Create(file_name)
	json_obj, err := list.List_as_json()
	if err != nil {
		fmt.Print(err.Error())
	}
	err = SaveState(*file, json_obj)
	if err != nil {
		fmt.Print(err.Error())
	}
}

func Load_Todo_List_From_Json(list baseList, file_name string) {
	file, _ := os.Open(file_name)
	json_obj, _ := LoadState(*file)
	list.List_from_json(json_obj)
}

func Load_New_Todo_List_From_Json(file_name string) (TodoList, error) {

	current_Todo_list := TodoList{List: []Todo{}}

	file, err := os.Open(file_name)
	if err != nil {
		return current_Todo_list, err
	}

	json_for_todo_list, err := LoadState(*file)

	if err != nil {
		return current_Todo_list, fmt.Errorf("something went wrong reading the file")
	}

	current_Todo_list.List_from_json(json_for_todo_list)

	return current_Todo_list, nil
}

func (tl *TodoList) completeTodo(compTodo string) {
	num, err := strconv.Atoi(compTodo)
	if err != nil {
		for i, todo := range tl.List {
			if todo.Name == compTodo {
				tl.List[i] = flipTodoStatus(todo)
				break
			}
		}
	} else {
		for i, todo := range tl.List {
			if i == num-1 {
				tl.List[i] = flipTodoStatus(todo)
				break
			}
		}
	}
}

func flipTodoStatus(todo Todo) Todo {
	if todo.Status == complete {
		todo.Status = todoStatus
		return todo
	} else {
		todo.Status = complete
		return todo
	}
}
