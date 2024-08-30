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
func assertTodo(t *testing.T, got, want []Todo) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("got todo list print %v but wanted %v", got, want)
	}
}

func generateTodoList() []Todo {
	return []Todo{{"Iron", "Todo"}, {"Eat", "Complete"},
		{"Hunker", "Complete"}, {"Mine", "Todo"}, {"Shear", "Todo"}, {"Cut", "Todo"}}
}

func generateTodoList10() []Todo {
	return []Todo{{"Iron", "Todo"}, {"Eat", "Complete"},
		{"Hunker", "Complete"}, {"Mine", "Todo"}, {"Shear", "Todo"}, {"Cut", "Todo"},
		{"Griddle", "Complete"}, {"Cook", "Todo"}, {"Host", "Complete"}, {"Grate", "Todo"}}
}

func TestTodoOutput(t *testing.T) {
	t.Run("10 todo's are printed", func(t *testing.T) {
		todo_list := &TodoList{generateTodoList10()}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "1. Iron: Todo\n2. Eat: Complete\n3. Hunker: Complete\n4. Mine: Todo\n5. Shear: Todo\n6. Cut: Todo\n7. Griddle: Complete\n8. Cook: Todo\n9. Host: Complete\n10. Grate: Todo\n"

		assertStrings(t, output.String(), want)

	})

	t.Run("if more then 10 todo's just print 10", func(t *testing.T) {
		todo_list_prep := append(generateTodoList10(), Todo{"Scale", "Todo"}, Todo{"Brush", "Complete"})
		todo_list := &TodoList{todo_list_prep}
		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "1. Iron: Todo\n2. Eat: Complete\n3. Hunker: Complete\n4. Mine: Todo\n5. Shear: Todo\n6. Cut: Todo\n7. Griddle: Complete\n8. Cook: Todo\n9. Host: Complete\n10. Grate: Todo\n"
		assertStrings(t, output.String(), want)
	})

	t.Run("if less then 10 todo's just print all Todos", func(t *testing.T) {
		todo_list := &TodoList{generateTodoList()}

		output := &bytes.Buffer{}

		todo_list.outputTodos(output)

		want := "1. Iron: Todo\n2. Eat: Complete\n3. Hunker: Complete\n4. Mine: Todo\n5. Shear: Todo\n6. Cut: Todo\n"
		assertStrings(t, output.String(), want)
	})

}

func TestAddTodo(t *testing.T) {
	todo_list := TodoList{generateTodoList()}

	todo_list.addTodo("Cook")

	got := todo_list.List
	want := []Todo{Todo{"Iron", "Todo"}, Todo{"Eat", "Complete"},
		Todo{"Hunker", "Complete"}, Todo{"Mine", "Todo"}, Todo{"Shear", "Todo"}, Todo{"Cut", "Todo"}, Todo{"Cook", "Todo"}}

	assertTodo(t, got, want)
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete Todo by name", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}

		todo_list.deleteTodo("Mine")

		got := todo_list.List
		want := []Todo{Todo{"Iron", "Todo"}, Todo{"Eat", "Complete"},
			Todo{"Hunker", "Complete"}, Todo{"Shear", "Todo"}, Todo{"Cut", "Todo"}}

		assertTodo(t, got, want)
	})

	t.Run("delete Todo by number", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}
		todo_list.deleteTodo("3")

		got := todo_list.List
		want := []Todo{{"Iron", "Todo"}, {"Eat", "Complete"},
			{"Mine", "Todo"}, {"Shear", "Todo"}, {"Cut", "Todo"}}

		assertTodo(t, got, want)
	})
}

// func TestTodoJson(t *testing.T) {
// 	todo_list := TodoList{[]string{"Iron", "Eat",
// 		"Hunker"}}
// 	var output []Todo

// 	output_json, _ := todo_list.List_as_json()

// 	json.Unmarshal(output_json, &output)

// 	want := `[{"Name": "Iron", "Status": "Complete"}, {"Name": "Eat", "Status": "Todo"}, {"Name": "Hunker", "Status": "Complete"}]`
// 	require.JSONEq(t, want, string(output_json))
// }
