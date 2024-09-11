package todo

import (
	"bytes"
	"context"
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
		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile.Name(), "")
		defer ctxDone()

		request := httptest.NewRequest(http.MethodGet, "/get_todo_list", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertNoError(t, err)

		assertStrings(t, response.Body.String(), generateTodoListAsString())
	})

	t.Run("threadsafe happy path test", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		tempfile_not_used, cleanUpNotUsedFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()
		defer cleanUpNotUsedFile()

		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile_not_used.Name(), tempfile.Name())
		defer ctxDone()

		request := httptest.NewRequest(http.MethodGet, "/threadsafe/get_todo_list", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertNoError(t, err)

		assertStrings(t, response.Body.String(), generateTodoListAsString())
	})

	t.Run("threadsafe add todo", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		tempfile_not_used, cleanUpNotUsedFile := createTempFile(t, "[]")
		defer cleanUpFile()
		defer cleanUpNotUsedFile()

		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile_not_used.Name(), tempfile.Name())
		defer ctxDone()

		var todo_name []byte = []byte("Example")
		request := httptest.NewRequest(http.MethodPut, "/threadsafe/add_todo", bytes.NewBuffer(todo_name))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertNoError(t, err)
		request_get := httptest.NewRequest(http.MethodGet, "/threadsafe/get_todo_list", nil)
		response_get := httptest.NewRecorder()

		server.ServeHTTP(response_get, request_get)

		assertNoError(t, err)
		assertStrings(t, response_get.Body.String(), generateTodoListAsString()+"7. Example: Todo\n")
	})

	t.Run("Posting test", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, "[]")
		defer cleanUpFile()

		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile.Name(), "")
		defer ctxDone()
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

		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile.Name(), "")
		defer ctxDone()
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

	t.Run("completing todo", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		ctx, ctxDone := context.WithCancel(context.Background())
		server, err := NewJsonTodoServer(ctx, tempfile.Name(), "")
		defer ctxDone()
		if err != nil {
			fmt.Print(err.Error())
		}

		todo_name, _ := json.Marshal("Cut")
		request := httptest.NewRequest(http.MethodPost, "/complete_todo", bytes.NewBuffer(todo_name))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		response_got := check_Todo(server, "Cut")

		got_Todo := Todo{}

		json.NewDecoder(response_got.Body).Decode(&got_Todo)

		assertSingleTodo(t, got_Todo, Todo{"Cut", complete})
	})
}

func check_Todo(server *TodoServer, todo_name string) *httptest.ResponseRecorder {
	request_get := httptest.NewRequest(http.MethodGet, "/check_todo/"+todo_name, nil)
	response_got := httptest.NewRecorder()

	server.ServeHTTP(response_got, request_get)
	return response_got
}
