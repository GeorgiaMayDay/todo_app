package todo

import (
	"fmt"
	"net/http"
)

type TodoServer struct {
	store TodoList
	http.Handler
}

const textType = "text"

func NewJsonTodoServer(file_name string) (*TodoServer, error) {
	List, err := Load_New_Todo_List_From_Json(file_name)
	if err != nil {
		fmt.Println("There was an issue accessing the saved todo list and for this session you'll be working from a fresh jotpad!")
	}
	p := new(TodoServer)

	p.store = List

	router := http.NewServeMux()
	router.Handle("/get_todo_list", http.HandlerFunc(p.getBoardHandler))

	p.Handler = router

	return p, err
}

func (p *TodoServer) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	p.store.outputTodos(w)
}
