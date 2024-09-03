package main

import (
	"fmt"
	"log"
	"net/http"
	"todo_app/todo"
)

const json_file_name string = "../../todo/db.json"

func main() {
	server, err := todo.NewJsonTodoServer(json_file_name)
	if err != nil {
		fmt.Println(err.Error())
	}
	go log.Fatal(http.ListenAndServe(":5000", server))
	fmt.Printf("Set Up!")
}
