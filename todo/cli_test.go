package todo

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertList(t *testing.T, got, want []string) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %s but wanted %s", got, want)
	}
}

func generateTodoListAsString() string {
	return "1. Iron: Todo\n2. Eat: Complete\n3. Hunker: Complete\n4. Mine: Todo\n5. Shear: Todo\n6. Cut: Todo\n"
}

type SpyList struct {
	List []string
}

func (sl *SpyList) outputTodos(writer io.Writer) {
	fmt.Fprint(writer, "Called")
}

func (sl *SpyList) addTodo(newTodo string) {
	sl.List = append(sl.List, newTodo)
}

func (sl *SpyList) deleteTodo(delTodo string) {
	sl.List = []string{}
}

func (sl *SpyList) completeTodo(delTodo string) {
	sl.List = []string{"Complete"}
}

func (sl *SpyList) list_as_json() ([]byte, error) {
	return []byte{}, fmt.Errorf("Filler")
}

func (sl *SpyList) list_from_json([]byte) {

}

var CliInputTable = map[io.Reader]string{
	strings.NewReader("1"):         "Called",
	strings.NewReader("2\nCalled"): "\"Called\" added\n",
}

func TestCli(t *testing.T) {

	t.Run("That CLI can take input and output correct response", func(t *testing.T) {
		_, cleanFile := createTempFile(t, InitialDataString)
		defer cleanFile()
		for input, want := range CliInputTable {
			todoSpy := &SpyList{}
			output := &bytes.Buffer{}

			in := input

			ReadAndOutput(in, output, todoSpy, test_file_name)

			assertStrings(t, output.String(), want)
		}
	})

	t.Run("That CLI can add todo list", func(t *testing.T) {
		_, cleanFile := createTempFile(t, InitialDataString)
		defer cleanFile()
		todoSpy := &SpyList{}
		output := &bytes.Buffer{}

		in := strings.NewReader("2\nCalled")

		ReadAndOutput(in, output, todoSpy, test_file_name)

		want := []string{"Called"}

		assertList(t, todoSpy.List, want)
	})

	t.Run("That CLI can delete elements from todo list", func(t *testing.T) {
		_, cleanFile := createTempFile(t, InitialDataString)
		defer cleanFile()
		todoSpy := &SpyList{List: []string{"Call"}}
		output := &bytes.Buffer{}

		in := strings.NewReader("3\nCall")

		ReadAndOutput(in, output, todoSpy, test_file_name)

		want := []string{}

		assertList(t, todoSpy.List, want)
	})

	t.Run("Load work", func(t *testing.T) {
		tmpfile, cleanFile := createTempFile(t, InitialDataString)

		todoSpy := &TodoList{List: generateTodoList10()}
		output := &bytes.Buffer{}

		in := strings.NewReader("6")

		ReadAndOutput(in, output, todoSpy, tmpfile.Name())

		assertTodo(t, todoSpy.List, generateTodoList())
		cleanFile()
	})

	t.Run("Save work", func(t *testing.T) {
		tmpfile, cleanFile := createTempFile(t, InitialDataString)

		todo_list_prep := append(generateTodoList(), Todo{"Scale", "Todo"})
		todoSpy := &TodoList{List: todo_list_prep}
		output := &bytes.Buffer{}

		in := strings.NewReader("5")

		ReadAndOutput(in, output, todoSpy, tmpfile.Name())

		in = strings.NewReader("1")

		ReadAndOutput(in, output, todoSpy, tmpfile.Name())

		want := generateTodoListAsString() + "7. Scale: Todo\n"
		assertStrings(t, output.String(), want)
		cleanFile()
	})
}
