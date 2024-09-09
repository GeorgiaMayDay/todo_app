package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type TodoServer struct {
	store TodoList
	file  string
	http.Handler
}

type apiRequest struct {
	verb     string
	key      string
	response chan<- *bytes.Buffer
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const textType = "text"

var requestBuffer chan<- apiRequest
var loggerDB *slog.Logger

func NewJsonTodoServer(file_name string, ts_filename string) (*TodoServer, error) {
	f, fileCreationErr := os.Create("log_jdb.json")
	check(fileCreationErr)
	loggerDB = slog.New(slog.NewJSONHandler(f, nil))

	List, err := Load_New_Todo_List_From_Json(file_name)
	if err != nil {
		panic(err)
		//fmt.Println("There was an issue accessing the saved todo list and for this session you'll be working from a fresh jotpad!")
	}
	p := new(TodoServer)

	p.store = List
	p.file = file_name

	requests := make(chan apiRequest, 10)
	requestBuffer = requests

	if ts_filename != "" {
		ctx := context.Background()

		go actor(requests, ts_filename, ctx)
		loggerDB.Info("JSON database", "Message", "Set Up")
	}

	router := http.NewServeMux()
	router.Handle("/get_todo_list", http.HandlerFunc(p.getBoardHandler))
	router.Handle("/add_todo", http.HandlerFunc(p.addTodoHandler))
	router.Handle("/check_todo/", http.HandlerFunc(p.checkTodoHandler))
	router.Handle("/delete_todo", http.HandlerFunc(p.deleteTodoHandler))
	router.Handle("/complete_todo", http.HandlerFunc(p.completeTodoHandler))
	router.Handle("/save", http.HandlerFunc(p.saveTodoHandler))
	router.Handle("/load", http.HandlerFunc(p.loadTodoHandler))

	router.Handle("/threadsafe/", http.HandlerFunc(requestHandler))

	p.Handler = router

	return p, err
}

func (p *TodoServer) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", textType)
	p.store.outputTodos(w)
}

func (p *TodoServer) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	loggerDB.Info("JSON database", "Message", "Add Todo called")
	w.Header().Set("content-type", jsonContentType)
	body, _ := io.ReadAll(r.Body)
	p.store.addTodo(string(body[:]))
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
	Save_Todo_List_From_Json(p.store, p.file)
}

func (p *TodoServer) loadTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	Load_Todo_List_From_Json(&p.store, p.file)
}

func (p *TodoServer) checkTodoHandler(w http.ResponseWriter, r *http.Request) {
	todo := strings.TrimPrefix(r.URL.Path, "/check_todo/")

	todo_found, err := checkTodo(p.store.List, todo)

	if err != nil {
		json.NewEncoder(w).Encode("{\"Message\":\"Status Not Found\"}")
		return
	}
	w.Header().Set("Content-Type", jsonContentType)
	json.NewEncoder(w).Encode(todo_found)
}

func actor(requests chan apiRequest, filename string, ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	loggerDB.Info("Actor", "Message", "Set Up")
	//Shutdown gracefully
	go func() {
		<-ctx.Done()
		close(requests)
	}()

	go func() {
		defer close(done)
		processRequests(requests, filename)
	}()

	return done
}

func processRequests(requests <-chan apiRequest, filename string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		List, err := Load_New_Todo_List_From_Json(filename)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer close(done)
		for req := range requests {
			loggerDB.Info("JSON database", "Message", "Threadsafe call: "+req.verb+" Processing")
			switch req.verb {
			case "get_todo_list":
				output := &bytes.Buffer{}
				List.outputTodos(output)
				req.response <- output
			case "add_todo":
				List.addTodo(req.key)
			case "delete_todo":
				List.deleteTodo(req.key)
			case "complete_todo":
				List.completeTodo(req.key)
			}
			close(req.response)
		}
	}()
	return done
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	command := strings.TrimPrefix(r.URL.Path, "/threadsafe/")
	loggerDB.Info("JSON database", "Message", "Threadsafe call: "+command)
	request := apiRequest{}
	responseChan := make(chan *bytes.Buffer)
	switch command {
	case "get_todo_list":
		request = apiRequest{verb: command, key: "", response: responseChan}
	case "add_todo", "delete_todo", "complete_todo":
		body, _ := io.ReadAll(r.Body)
		loggerDB.Info("JSON database", "Message", "Threadsafe Todo PUT with: "+string(body[:]))
		request = apiRequest{verb: command, key: string(body[:]), response: responseChan}
	}

	select {
	case requestBuffer <- request:
	default:
		log.Println("request buffer full or shutdown")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	response, ok := <-responseChan

	loggerDB.Info("JSON database", "Message", "Threadsafe call: "+command+" response received")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}

	if ok {
		w.WriteHeader(http.StatusOK)
		response.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
