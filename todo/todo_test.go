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

		want := "1. Iron\n2. Eat\n3. Hunker\n4. Mine\n5. Shear\n6. Cut\n7. Griddle\n8. Cook\n9. Host\n10. Grate\n"

		assertStrings(t, output.String(), want)

	})

	t.Run("if more then 10 todo's just print 10", func(t *testing.T) {
		todo_list := &TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut", "Griddle", "Cook", "Host", "Grate", "Scale", "Brush"}}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "1. Iron\n2. Eat\n3. Hunker\n4. Mine\n5. Shear\n6. Cut\n7. Griddle\n8. Cook\n9. Host\n10. Grate\n"
		assertStrings(t, output.String(), want)
	})

	t.Run("if less then 10 todo's just print all Todos", func(t *testing.T) {
		todo_list := &TodoList{[]string{"Iron", "Eat",
			"Hunker", "Mine", "Shear", "Cut"}}

		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "1. Iron\n2. Eat\n3. Hunker\n4. Mine\n5. Shear\n6. Cut\n"
		assertStrings(t, output.String(), want)
	})

}

func TestAddTodo(t *testing.T) {
	todo_list := TodoList{[]string{"Iron", "Eat",
		"Hunker", "Mine", "Shear", "Cut"}}

	todo_list.addTodo("Cook")

	got := todo_list.List
	want := []string{"Iron", "Eat",
		"Hunker", "Mine", "Shear", "Cut", "Cook"}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %s but wanted %s", got, want)
	}

}
