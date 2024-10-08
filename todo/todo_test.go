package todo

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
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
	want := []Todo{{"Iron", "Todo"}, {"Eat", "Complete"},
		{"Hunker", "Complete"}, {"Mine", "Todo"}, {"Shear", "Todo"}, {"Cut", "Todo"}, {"Cook", "Todo"}}

	assertTodo(t, got, want)
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete Todo by name", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}

		todo_list.deleteTodo("Mine")

		got := todo_list.List
		want := []Todo{{"Iron", "Todo"}, {"Eat", "Complete"},
			{"Hunker", "Complete"}, {"Shear", "Todo"}, {"Cut", "Todo"}}

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

func TestTodoJson(t *testing.T) {
	t.Run("Get json from TodoList", func(t *testing.T) {
		todo_list := TodoList{generateTodoList()}
		var output []Todo

		output_json, _ := todo_list.list_as_json()

		json.Unmarshal(output_json, &output)

		want := "[{\"Name\":\"Iron\",\"Status\":\"Todo\"},{\"Name\":\"Eat\",\"Status\":\"Complete\"},{\"Name\":\"Hunker\",\"Status\":\"Complete\"},{\"Name\":\"Mine\",\"Status\":\"Todo\"},{\"Name\":\"Shear\",\"Status\":\"Todo\"},{\"Name\":\"Cut\",\"Status\":\"Todo\"}]"
		require.JSONEq(t, want, string(output_json))
	})

	t.Run("Get TodoList from json", func(t *testing.T) {
		todo_list := TodoList{}

		json_of_list, _ := json.Marshal(generateTodoList())

		todo_list.list_from_json(json_of_list)

		assertTodo(t, todo_list.List, generateTodoList())
	})
}

func TestCompleteTodo(t *testing.T) {

	t.Run("complete Todo by name", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}

		todo_list.completeTodo("Mine")

		got := todo_list.List
		want := []Todo{{"Iron", "Todo"}, {"Eat", "Complete"}, {"Hunker", "Complete"}, {"Mine", "Complete"},
			{"Shear", "Todo"}, {"Cut", "Todo"}}

		assertTodo(t, got, want)
	})

	t.Run("complete Todo by number", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}
		todo_list.completeTodo("5")

		got := todo_list.List
		want := []Todo{{"Iron", "Todo"}, {"Eat", "Complete"}, {"Hunker", "Complete"},
			{"Mine", "Todo"}, {"Shear", "Complete"}, {"Cut", "Todo"}}

		assertTodo(t, got, want)
	})

	t.Run("return Todo by number", func(t *testing.T) {

		todo_list := TodoList{generateTodoList()}
		todo_list.completeTodo("2")

		got := todo_list.List
		want := []Todo{{"Iron", "Todo"}, {"Eat", "Todo"}, {"Hunker", "Complete"},
			{"Mine", "Todo"}, {"Shear", "Todo"}, {"Cut", "Todo"}}

		assertTodo(t, got, want)
	})
}
