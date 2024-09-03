package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TodoServer struct {
	store TodoList
	file  string
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
	p.file = file_name

	router := http.NewServeMux()
	router.Handle("/get_todo_list", http.HandlerFunc(p.getBoardHandler))
	router.Handle("/add_todo", http.HandlerFunc(p.addTodoHandler))
	router.Handle("/check_todo/", http.HandlerFunc(p.checkTodoHandler))
	router.Handle("/delete_todo", http.HandlerFunc(p.deleteTodoHandler))
	router.Handle("/complete_todo", http.HandlerFunc(p.completeTodoHandler))
	router.Handle("/save", http.HandlerFunc(p.saveTodoHandler))
	router.Handle("/load", http.HandlerFunc(p.loadTodoHandler))

	p.Handler = router

	return p, err
}

func (p *TodoServer) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	p.store.outputTodos(w)
}

func (p *TodoServer) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	var output string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&output)
	p.store.addTodo(string(output[:]))
}

func (p *TodoServer) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	var output string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&output)
	p.store.deleteTodo(string(output[:]))
}

func (p *TodoServer) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	var output string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&output)
	p.store.completeTodo(string(output[:]))
}

func (p *TodoServer) saveTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	Save_Todo_List_From_Json(&p.store, p.file)
}

func (p *TodoServer) loadTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	Load_Todo_List_From_Json(&p.store, p.file)
}

func (p *TodoServer) checkTodoHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/check_todo/")

	todo_found, err := checkTodo(p.store.List, player)

	if err != nil {
		json.NewEncoder(w).Encode("{\"Message\":\"Status Not Found\"}")
		return
	}
	w.Header().Set("Content-Type", jsonContentType)
	json.NewEncoder(w).Encode(todo_found)
}
