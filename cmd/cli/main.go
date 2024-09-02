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
		keepgoing = todo.ReadAndOutput(os.Stdin, os.Stdout, db, json_file_name, api_address)
	}
	finish <- true
}

const api_address string = "http://localhost:5000"
const json_file_name string = "../../todo/db.json"

func main() {

	db, err := todo.Load_New_Todo_List_From_Json(json_file_name)
	if err != nil {
		fmt.Println("There was an issue accessing the saved todo list and for this session you'll be working from a fresh jotpad!")
		fmt.Println(err.Error())
	}

	fmt.Println("Welcome to GoDo, a application to help you manage you tasks")
	fmt.Println("To end GoDo, either close the application or type Ctrl+C")
	todo.Show_Instructions(os.Stdout)
	finishChannel := make(chan bool, 1)

	GoDo(&db, finishChannel)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quitChannel:
		fmt.Println("Adios!")
	case <-finishChannel:
		fmt.Println("Adios!")
	}
}
