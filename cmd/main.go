package main

import (
	"os"
	"todo_app/todo"
)

// ask about the point and interface interaction

func main() {

	db := &todo.TodoList{
		List: []string{"Grate", "Brush"},
	}

	todo.ReadAndOutput(os.Stdin, os.Stdout, db)
}
