package todo

import (
	"bufio"
	"fmt"
	"io"
)

func readLine(reader *bufio.Scanner) string {
	reader.Scan()
	return reader.Text()
}

func Show_Instructions(printer io.Writer) {
	fmt.Fprintln(printer, "To use:")
	fmt.Fprintln(printer, "type 1 to show the top 10 Todos")
	fmt.Fprintln(printer, "type 2 to add a new Todo")
	fmt.Fprintln(printer, "type 3 to delete a Todo, you can use the name or number")
	fmt.Fprintln(printer, "type 4 to complete or move back to todo a Todo, you can use the name or number")
	fmt.Fprintln(printer, "type 5 and 6 to save or load respectively")
	fmt.Fprintln(printer, "type 7 to see these instructions again")
}

var invalid_opt_msg = "You've entered an invalid option"

func ReadAndOutput(in io.Reader, out io.Writer, list baseList, storage_name string) bool {
	reader := bufio.NewScanner(in)
	option := readLine(reader)

	switch option {
	case "1":
		list.outputTodos(out)
	case "2":
		todo_name := readLine(reader)
		list.addTodo(todo_name)
		out_msg := "\"" + todo_name + "\" added"
		fmt.Fprintln(out, out_msg)
	case "3":
		todo_name := readLine(reader)
		list.deleteTodo(todo_name)
		out_msg := "\"" + todo_name + "\" deleted"
		fmt.Fprintln(out, out_msg)
	case "4":
		todo_name := readLine(reader)
		list.completeTodo(todo_name)
		out_msg := "\"" + todo_name + "\" marked complete"
		fmt.Fprintln(out, out_msg)
	case "5":
		Save_Todo_List_From_Json(list, storage_name)
		out_msg := "Current Todo List Saved"
		fmt.Println(out, out_msg)
	case "6":
		Load_Todo_List_From_Json(list, storage_name)
		out_msg := "Todo List Loaded"
		fmt.Println(out, out_msg)
	case "7":
		Show_Instructions(out)
	case "Quit":
		return false
	case "Q":
		return false
	case "":
		return false
	default:
		fmt.Fprintln(out, invalid_opt_msg)
	}
	return true
}
