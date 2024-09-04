package todo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func assertList(t *testing.T, got, want []string) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %s but wanted %s", got, want)
	}
}

func generateTodoListAsString() string {
	return "1. Iron: Todo\n2. Eat: Complete\n3. Hunker: Complete\n4. Mine: Todo\n5. Shear: Todo\n6. Cut: Todo\n"
}

func TestCli(t *testing.T) {

	t.Run("That CLI can print todos", func(t *testing.T) {
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}
		testSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text")
			fmt.Fprint(w, generateTodoListAsString())
		}))
		defer testSvr.Close()

		svrUrl := testSvr.URL

		in := strings.NewReader("1")

		ReadAndOutput(ctx, in, output, svrUrl, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)
		assertStrings(t, output.String(), generateTodoListAsString())
	})

	t.Run("That CLI can handle multiple inputs", func(t *testing.T) {

		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}
		testSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text")
			fmt.Fprint(w, generateTodoListAsString())
		}))
		defer testSvr.Close()

		svrUrl := testSvr.URL

		in := strings.NewReader("2\nBrush")

		ReadAndOutput(ctx, in, output, svrUrl, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		in = strings.NewReader("1\n")

		ReadAndOutput(ctx, in, output, svrUrl, finishChan)

		got = <-finishChan

		assertNoError(t, got.Err)
		assertStrings(t, output.String(), "\"Brush\" added\n"+generateTodoListAsString())
	})

	t.Run("That CLI can graceful handle server sending a bad status back", func(t *testing.T) {
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}
		testSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer testSvr.Close()

		svrUrl := testSvr.URL

		// test_Int := rand.Intn(6) + 1
		test_Int := 1
		in := strings.NewReader(strconv.Itoa(test_Int))

		ReadAndOutput(ctx, in, output, svrUrl, finishChan)

		got := <-finishChan

		want := RequestError{502, nil}

		if cmp.Equal(got.Err, want) {
			t.Errorf("got an error but got %s, when wants %s", got.Err.Error(), want.Error())
		}

		assertStrings(t, output.String(), "")
	})

	t.Run("That CLI can graceful handle no server", func(t *testing.T) {
		ctx := context.Background()
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}
		testSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		testSvr.Close()

		svrUrl := testSvr.URL

		// test_Int := rand.Intn(6) + 1
		in := strings.NewReader(strconv.Itoa(5))

		ReadAndOutput(ctx, in, output, svrUrl, finishChan)

		got := <-finishChan

		want := RequestError{0, nil}

		if cmp.Equal(got.Err, want) {
			t.Errorf("got an error but got %s, when wants %s", got.Err.Error(), want.Error())
		}

		assertStrings(t, output.String(), "")
	})
}

// 	t.Run("That CLI can add todo list", func(t *testing.T) {
// 		_, cleanFile := createTempFile(t, InitialDataString)
// 		defer cleanFile()
// 		todoSpy := &SpyList{}
// 		output := &bytes.Buffer{}

// 		in := strings.NewReader("2\nCalled")

// 		ReadAndOutput(in, output, todoSpy, test_file_name, "")

// 		want := []string{"Called"}

// 		assertList(t, todoSpy.List, want)
// 	})

// 	t.Run("That CLI can delete elements from todo list", func(t *testing.T) {
// 		_, cleanFile := createTempFile(t, InitialDataString)
// 		defer cleanFile()
// 		todoSpy := &SpyList{List: []string{"Call"}}
// 		output := &bytes.Buffer{}

// 		in := strings.NewReader("3\nCall")

// 		ReadAndOutput(in, output, todoSpy, test_file_name, "")

// 		want := []string{}

// 		assertList(t, todoSpy.List, want)
// 	})

// 	t.Run("Load work", func(t *testing.T) {
// 		tmpfile, cleanFile := createTempFile(t, InitialDataString)

// 		todoSpy := &TodoList{List: generateTodoList10()}
// 		output := &bytes.Buffer{}

// 		in := strings.NewReader("6")

// 		ReadAndOutput(in, output, todoSpy, tmpfile.Name(), "")

// 		assertTodo(t, todoSpy.List, generateTodoList())
// 		cleanFile()
// 	})

// 	t.Run("Save work", func(t *testing.T) {
// 		tmpfile, cleanFile := createTempFile(t, InitialDataString)

// 		todo_list_prep := append(generateTodoList(), Todo{"Scale", "Todo"})
// 		todoSpy := &TodoList{List: todo_list_prep}
// 		output := &bytes.Buffer{}

// 		in := strings.NewReader("5")

// 		ReadAndOutput(in, output, todoSpy, tmpfile.Name(), "")

// 		in = strings.NewReader("1")

// 		ReadAndOutput(in, output, todoSpy, tmpfile.Name(), "")

// 		want := generateTodoListAsString() + "7. Scale: Todo\n"
// 		assertStrings(t, output.String(), want)
// 		cleanFile()
// 	})
// }
