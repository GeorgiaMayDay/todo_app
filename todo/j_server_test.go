package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertSingleTodo(t *testing.T, got, want Todo) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %v but wanted %v", got, want)
	}
}

func TestServer(t *testing.T) {

	t.Run("happy path test", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()
		server, err := NewJsonTodoServer(tempfile.Name())

		request := httptest.NewRequest(http.MethodGet, "/get_todo_list", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertNoError(t, err)

		assertStrings(t, response.Body.String(), generateTodoListAsString())
	})

	t.Run("Posting test", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, "[]")
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name())
		if err != nil {
			fmt.Print(err.Error())
		}

		todo_name, err := json.Marshal("Cheese")
		request := httptest.NewRequest(http.MethodPost, "/add_todo", bytes.NewBuffer(todo_name))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		response_got := check_Todo(server, "Cheese")

		got_Todo := Todo{}

		json.NewDecoder(response_got.Body).Decode(&got_Todo)

		assertNoError(t, err)
		assertSingleTodo(t, got_Todo, Todo{"Cheese", todoStatus})
	})

	t.Run("deleting todo", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name())
		if err != nil {
			fmt.Print(err.Error())
		}

		todo_name, _ := json.Marshal("Cut")
		request := httptest.NewRequest(http.MethodPost, "/delete_todo", bytes.NewBuffer(todo_name))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		response_got := check_Todo(server, "Cut")

		var msg string

		json.NewDecoder(response_got.Body).Decode(&msg)

		assertStrings(t, msg, `{"Message":"Status Not Found"}`)
	})
}

func check_Todo(server *TodoServer, todo_name string) *httptest.ResponseRecorder {
	request_get := httptest.NewRequest(http.MethodGet, "/check_todo/"+todo_name, nil)
	response_got := httptest.NewRecorder()

	server.ServeHTTP(response_got, request_get)
	return response_got
}
