package todo

import (
	"net/http"
)

type TodoServer struct {
	store TodoList
	http.Handler
}

const textType = "text"

func NewJsonTodoServer(file_name string) (*TodoServer, error) {
	List, err := Load_New_Todo_List_From_Json(file_name)
	p := new(TodoServer)

	p.store = List

	router := http.NewServeMux()
	router.Handle("/GET", http.HandlerFunc(p.getBoardHandler))

	p.Handler = router

	return p, err
}

func (p *TodoServer) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	p.store.outputTodos(w)
}
