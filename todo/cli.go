package todo

import (
	"bufio"
	"bytes"
	"context"
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

const invalid_opt_msg string = "You've entered an invalid option"

type TodoResult struct {
	Stop bool
	Err  error
}

type RequestError struct {
	StatusCode int

	Err error
}

func (r *RequestError) Error() string {
	error_msg := "There was no specific error"
	if r.Err != nil {
		error_msg = r.Err.Error()
	}
	return fmt.Sprintf("There was an error connecting to the server: status %d: err %s", r.StatusCode, error_msg)
}

func ReadAndOutput(ctx context.Context, in io.Reader, out io.Writer, api_address string, result chan<- TodoResult) {
	ch := make(chan TodoResult, 1)
	reader := bufio.NewScanner(in)
	option := readLine(reader)

	switch option {
	case "1":
		InfoLog("CLI", "Getting TODO list: "+ctx.Value("Trace_id").(string))
		resp, _, shouldExit, keepgoing, server_err := get_Svr(api_address + "/get_todo_list")
		if shouldExit {
			ch <- TodoResult{keepgoing, server_err}
		}
		output := new(bytes.Buffer)
		_, err := output.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resp.Body.Close()
		fmt.Fprint(out, output)
	case "2":
		InfoLog("CLI", "Adding Todo: "+ctx.Value("Trace_id").(string))
		input, todo_name, err := getNameFromScanner(reader, out)
		if err != nil {
			result <- TodoResult{true, err}
		}
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+"/add_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
		}
		out_msg := "\"" + input + "\" added"
		fmt.Fprintln(out, out_msg)
	case "3":
		InfoLog("CLI", "Delete Todo: "+ctx.Value("Trace_id").(string))
		input, todo_name, err := getNameFromScanner(reader, out)
		if err != nil {
			result <- TodoResult{true, err}
		}
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+"/delete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
		}
		out_msg := "\"" + input + "\" deleted"
		fmt.Fprintln(out, out_msg)
	case "4":
		InfoLog("CLI", "Completing Todo: "+ctx.Value("Trace_id").(string))
		input, todo_name, err := getNameFromScanner(reader, out)
		if err != nil {
			result <- TodoResult{true, err}
		}
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+"/complete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
		}
		out_msg := "\"" + input + "\" complete"
		fmt.Fprintln(out, out_msg)
	case "5":
		InfoLog("CLI", "Saving Todo List: "+ctx.Value("Trace_id").(string))
		_, _, shouldExit, keepgoing, server_err := get_Svr(api_address + "/save")
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
		}
		out_msg := "Current Todo List Saved"
		fmt.Println(out, out_msg)
	case "6":
		InfoLog("CLI", "Loading Todo List: "+ctx.Value("Trace_id").(string))
		_, _, shouldExit, keepgoing, server_err := get_Svr(api_address + "/load")
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
		}
		out_msg := "Todo List Loaded"
		fmt.Println(out, out_msg)
	case "7":
		Show_Instructions(out)
	case "Quit", "q", "Q", "":
		result <- TodoResult{false, nil}
	default:
		fmt.Fprintln(out, invalid_opt_msg)
		result <- TodoResult{true, nil}
	}
	if len(result) == 0 {
		result <- TodoResult{true, nil}
	}
}

func get_Svr(url string) (*http.Response, error, bool, bool, error) {
	resp, err := http.Get(url)
	if resp == nil {
		WarnLog("CLI", "Server didn't exist")
		return nil, nil, true, true, &RequestError{0, fmt.Errorf("no response from server")}
	}
	if string(resp.Status[0]) != "2" {
		WarnLog("CLI", "Bad Request")
		return nil, nil, true, true, &RequestError{resp.StatusCode, err}
	}
	if err != nil {
		return nil, nil, true, true, &RequestError{resp.StatusCode, err}
	}
	return resp, err, false, false, nil
}

func post_Svr(url, ct string, reader io.Reader) (*http.Response, error, bool, bool, error) {
	resp, err := http.Post(url, ct, reader)
	if resp == nil {
		WarnLog("CLI", "Server didn't exist")
		return nil, nil, true, true, &RequestError{0, fmt.Errorf("no response from server")}
	}
	if string(resp.Status[0]) != "2" {
		WarnLog("CLI", "Bad Request")
		return nil, nil, true, true, &RequestError{resp.StatusCode, err}
	}
	if err != nil {
		return nil, nil, true, true, &RequestError{resp.StatusCode, err}
	}
	return resp, err, false, false, nil
}

func getNameFromScanner(reader *bufio.Scanner, out io.Writer) (string, []byte, error) {
	input := readLine(reader)
	todo_name, err := json.Marshal(input)
	if err != nil {
		fmt.Fprintln(out, "This is an invalid name")
		return "", nil, fmt.Errorf("%s: This is an invalid name", string(todo_name))
	}
	return input, todo_name, nil
}
