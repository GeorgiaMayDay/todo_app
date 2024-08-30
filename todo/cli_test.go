package todo

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertInt(t *testing.T, got, want int) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("go %d printed but wanted %d", got, want)
	}
}

type SpyList struct {
}

func (sl *SpyList) outputTodos(writer io.Writer) {
	fmt.Fprint(writer, "Called")
}

func (sl *SpyList) addTodo(newTodo string) {

}

var CliInputTable = map[io.Reader]string{
	strings.NewReader("1"):         "Called",
	strings.NewReader("2\nCalled"): "\"Called\" added\n",
}

func TestCli(t *testing.T) {
	t.Run("That CLI can take input and output correct response", func(t *testing.T) {
		for input, want := range CliInputTable {
			todoSpy := &SpyList{}
			output := &bytes.Buffer{}

			in := input

			ReadAndOutput(in, output, todoSpy)

			assertStrings(t, output.String(), want)
		}
	})

	// t.Run("That CLI can update todo list", func(t *testing.T) {
	// 	todoSpy := &SpyList{}
	// 	output := &bytes.Buffer{}

	// 	in := input

	// 	ReadAndOutput(in, output, todoSpy)

	// 	assertStrings(t, output.String(), want)
	// })
}
