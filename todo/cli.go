package todo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const jsonContentType string = "application/json"

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

func ReadAndOutput(in io.Reader, out io.Writer, list baseList, api_address string) bool {
	reader := bufio.NewScanner(in)
	option := readLine(reader)

	switch option {
	case "1":
		resp, err := http.Get(api_address + "/get_todo_list")
		if err != nil {
			fmt.Println(err)
		}
		output := new(bytes.Buffer)
		_, err = output.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		fmt.Fprint(out, output)
	case "2":
		input, todo_name := getNameFromScanner(reader, out)
		if input == "" {
			break
		}
		_, err := http.Post(api_address+"/add_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if err != nil {
			fmt.Println(err)
		}
		out_msg := "\"" + input + "\" added"
		fmt.Fprintln(out, out_msg)
	case "3":
		input, todo_name := getNameFromScanner(reader, out)
		if input == "" {
			break
		}
		_, err := http.Post(api_address+"/delete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if err != nil {
			fmt.Println(err)
		}
		out_msg := "\"" + input + "\" deleted"
		fmt.Fprintln(out, out_msg)
	case "4":
		input, todo_name := getNameFromScanner(reader, out)
		if input == "" {
			break
		}
		_, err := http.Post(api_address+"/complete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if err != nil {
			fmt.Println(err)
		}
		out_msg := "\"" + input + "\" complete"
		fmt.Fprintln(out, out_msg)
	case "5":
		_, err := http.Get(api_address + "/save")
		if err != nil {
			fmt.Println(err)
			return false
		}
		out_msg := "Current Todo List Saved"
		fmt.Println(out, out_msg)
	case "6":
		_, err := http.Get(api_address + "/load")
		if err != nil {
			fmt.Println(err)
			return false
		}
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

func getNameFromScanner(reader *bufio.Scanner, out io.Writer) (string, []byte) {
	input := readLine(reader)
	todo_name, err := json.Marshal(input)
	if err != nil {
		fmt.Fprintln(out, "This is an invalid name")
		return "", nil
	}
	return input, todo_name
}
