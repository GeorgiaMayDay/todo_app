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

func assertTodos(want, got map[int]string, t *testing.T) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %v but wanted %v", got, want)
	}
}

func TestTodoOutput(t *testing.T) {
	t.Run("10 todo's are printed", func(t *testing.T) {
		todo_list := &TodoList{map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut", 7: "Griddle", 8: "Cook", 9: "Host", 10: "Grate"}}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\nGriddle\nCook\nHost\nGrate\n"

		assertStrings(t, output.String(), want)

	})

	t.Run("if their are more then 10 todos, print only 10", func(t *testing.T) {
		todo_list := &TodoList{map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut", 7: "Griddle", 8: "Cook", 9: "Host", 10: "Grate", 11: "Scale", 12: "Brush"}}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\nGriddle\nCook\nHost\nGrate\n"

		assertStrings(t, output.String(), want)
	})

	t.Run("if less then 10 todo's just print 10", func(t *testing.T) {
		todo_list := &TodoList{map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut"}}

		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "Iron\nEat\nHunker\nMine\nShear\nCut\n"
		assertStrings(t, output.String(), want)
	})

}

func TestAddTodo(t *testing.T) {
	t.Run("if less then 10 todo's ju", func(t *testing.T) {
		todo_list := TodoList{map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut"}}

		todo_list.addTodo("Cook")

		got := todo_list.List
		want := map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut", 7: "Cook"}

		assertTodos(want, got, t)

	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete a todo", func(t *testing.T) {
		todo_list := TodoList{map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut"}}

		todo_list.addTodo("Cook")

		got := todo_list.List
		want := map[int]string{1: "Iron", 2: "Eat",
			3: "Hunker", 4: "Mine", 5: "Shear", 6: "Cut"}

		assertTodos(want, got, t)

	})
}
