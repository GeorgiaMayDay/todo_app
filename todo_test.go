package todo

import (
	"bytes"
	"cmp"
	"testing"
)

func TestTodoCLI(t *testing.T) {
	t.Run("10 todo's are printed", func(t *testing.T) {
		todo_list := TodoList{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut", "Griddle", "Cook", "Host", "Grate"}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\nGriddle\nCook\nHost\nGrate\n"

		if cmp.Compare(output.String(), want) != 0 {
			t.Errorf("got todo list print %s but wanted %s", output.String(), want)
		}
	})

}
