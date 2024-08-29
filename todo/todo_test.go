package todo

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got printed %s but wanted %s", got, want)
	}
}

func TestTodoOutput(t *testing.T) {
	t.Run("10 todo's are printed", func(t *testing.T) {
		todo_list := &TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut", "Griddle", "Cook", "Host", "Grate"}}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\nGriddle\nCook\nHost\nGrate\n"

		assertStrings(t, output.String(), want)

	})

	t.Run("if less then 10 todo's ju", func(t *testing.T) {
		todo_list := &TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut", "Griddle", "Cook", "Host", "Grate", "Scale", "Brush"}}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\nGriddle\nCook\nHost\nGrate\n"

		assertStrings(t, output.String(), want)
	})

	t.Run("if less then 10 todo's ju", func(t *testing.T) {
		todo_list := &TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut"}}

		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\n"
		assertStrings(t, output.String(), want)
	})

}

func TestAddTodo(t *testing.T) {
	t.Run("if less then 10 todo's ju", func(t *testing.T) {
		todo_list := TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut"}}

		todo_list.addTodo("Cook")

		got := todo_list.list
		want := []string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut", "Cook"}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got todo list print %s but wanted %s", got, want)
		}

	})
}
