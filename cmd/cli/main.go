package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo_app/todo"
)

func GoDo(finish chan<- bool) {
	keepgoing := true
	var err error = nil
	for keepgoing {
		keepgoing, err = todo.ReadAndOutput(os.Stdin, os.Stdout, api_address)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	finish <- true
}

const api_address string = "http://localhost:5000"

func main() {
	fmt.Println("Welcome to GoDo, a application to help you manage you tasks")
	fmt.Println("To end GoDo, either close the application or type Ctrl+C")
	todo.Show_Instructions(os.Stdout)
	finishChannel := make(chan bool, 1)

	GoDo(finishChannel)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quitChannel:
		fmt.Println("Adios!")
	case <-finishChannel:
		fmt.Println("Adios!")
	}
}
