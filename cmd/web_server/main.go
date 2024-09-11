package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"todo_app/todo"
)

const json_file_name string = "../../todo/db.json"

func main() {
	ctx, ctxDone := context.WithCancel(context.Background())
	server, err := todo.NewJsonTodoServer(ctx, json_file_name, "../../todo/db_threadsafe.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	go log.Fatal(http.ListenAndServe(":5000", server))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctxDone()
	fmt.Printf("We're done!")
}
