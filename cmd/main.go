package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo_app/todo"
)

// ask about the point and interface interaction

func GoDo(db *todo.TodoList, finish chan<- bool) {
	keepgoing := true
	for keepgoing {
		keepgoing = todo.ReadAndOutput(os.Stdin, os.Stdout, db)
	}
	finish <- true
}

func main() {

	db := &todo.TodoList{
		List: []string{"Grate", "Brush"},
	}

	fmt.Println("Welcome to GoDo, a application to help you manage you tasks")
	fmt.Println("To end GoDo, either close the application or type Ctrl+C")
	finishChannel := make(chan bool, 1)

	go GoDo(db, finishChannel)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quitChannel:
		fmt.Println("Adios!")
	case <-finishChannel:
		fmt.Println("Adios!")

	}
}
