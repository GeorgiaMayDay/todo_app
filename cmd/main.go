package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo_app/todo"
)

func GoDo(db *todo.TodoList, finish chan<- bool) {
	keepgoing := true
	for keepgoing {
		keepgoing = todo.ReadAndOutput(os.Stdin, os.Stdout, db)
	}
	finish <- true
}

func main() {

	db, err := todo.Create_todo_list_with_json_file("../todo/db.json")
	if err != nil {
		fmt.Println("There was an issue accessing the saved todo list and for this session you'll be working from a fresh jotpad!")
		fmt.Println(err.Error())
	}

	fmt.Println("Welcome to GoDo, a application to help you manage you tasks")
	fmt.Println("To end GoDo, either close the application or type Ctrl+C")
	finishChannel := make(chan bool, 1)

	go GoDo(&db, finishChannel)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quitChannel:
		fmt.Println("Adios!")
	case <-finishChannel:
		fmt.Println("Adios!")
	}
}
