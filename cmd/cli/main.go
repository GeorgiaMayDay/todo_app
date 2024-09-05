package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo_app/todo"

	"github.com/google/uuid"
)

func GoDo(parentCtx context.Context, parentCancel context.CancelFunc) {
	defer parentCancel()
	finishChan := make(chan todo.TodoResult, 1)
	keepgoing := true
	for keepgoing {
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(parentCtx, string("Trace_id"), trace_id)
		go todo.ReadAndOutput(ctx, os.Stdin, os.Stdout, api_address, finishChan)
		select {
		case result := <-finishChan:
			if result.Err != nil {
				fmt.Println(result.Err.Error())
				fmt.Println("There's been an issue with server communication")
			}
			keepgoing = result.Stop
		case <-ctx.Done():
			t_id := ctx.Value("Trace_id").(string)
			todo.InfoLog("CLI", "CLI stopped:"+t_id)
			continue
		}
	}
}

const api_address string = "http://localhost:5000"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create("log.json")
	check(err)
	todo.SetUpLogger(f)

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Welcome to GoDo, a application to help you manage you tasks")
	fmt.Println("To end GoDo, either close the application or type Ctrl+C")
	todo.Show_Instructions(os.Stdout)
	go GoDo(ctx, cancel)
	todo.InfoLog("CLI", "CLI started")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quitChannel:
		fmt.Println("Adios!")
	case <-ctx.Done():
		fmt.Println("Adios!")
	}
}
