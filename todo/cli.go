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
	"sync"
)

const jsonContentType string = "application/json"

func readLine(reader *bufio.Scanner) string {
	reader.Scan()
	return reader.Text()
}

type TodoErr string

func (e TodoErr) Error() string {
	return string(e)
}

const (
	errNoResponse        = TodoErr("no response from server")
	errWordAlreadyExists = TodoErr("cannot add word because it already exists")
)

func Show_Instructions(printer io.Writer) {
	fmt.Fprintln(printer, "To use:")
	fmt.Fprintln(printer, "type 1 to show the top 10 Todos")
	fmt.Fprintln(printer, "type 2 to add a new Todo")
	fmt.Fprintln(printer, "type 3 to delete a Todo, you can use the name or number")
	fmt.Fprintln(printer, "type 4 to complete or move back to todo a Todo, you can use the name or number")
	fmt.Fprintln(printer, "type 5 and 6 to save or load respectively")
	fmt.Fprintln(printer, "type 7 to see these instructions again")
	fmt.Fprintln(printer, "type 8 to toggle threadsafe mode on and off")
}

const invalid_opt_msg string = "You've entered an invalid option"

var mutex sync.Mutex
var thread_safe_switch string = ""

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
	reader := bufio.NewScanner(in)
	option := readLine(reader)

	switch option {
	case "1":
		InfoLog("CLI", "Getting TODO list: "+ctx.Value("Trace_id").(string))
		resp, _, shouldExit, keepgoing, server_err := get_Svr(api_address + thread_safe_switch + "/get_todo_list")
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
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
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+thread_safe_switch+"/add_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
		}
		out_msg := "\"" + input + "\" added"
		fmt.Fprintln(out, out_msg)
	case "3":
		input, todo_name, err := getNameFromScanner(reader, out)
		if err != nil {
			result <- TodoResult{true, err}
		}
		InfoLog("CLI", "Delete Todo: "+input)
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+thread_safe_switch+"/delete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
		}
		out_msg := "\"" + input + "\" deleted"
		fmt.Fprintln(out, out_msg)
	case "4":
		input, todo_name, err := getNameFromScanner(reader, out)
		if err != nil {
			result <- TodoResult{true, err}
		}
		InfoLog("CLI", "Complete Todo: "+input)
		_, _, shouldExit, keepgoing, server_err := post_Svr(api_address+thread_safe_switch+"/complete_todo", jsonContentType, bytes.NewBuffer(todo_name))
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
		}
		out_msg := "\"" + input + "\" complete"
		fmt.Fprintln(out, out_msg)
	case "5":
		InfoLog("CLI", "Saving Todo List: "+ctx.Value("Trace_id").(string))
		_, _, shouldExit, keepgoing, server_err := get_Svr(api_address + "/save")
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
		}
		out_msg := "Current Todo List Saved"
		fmt.Println(out, out_msg)
	case "6":
		InfoLog("CLI", "Loading Todo List: "+ctx.Value("Trace_id").(string))
		_, _, shouldExit, keepgoing, server_err := get_Svr(api_address + "/load")
		if shouldExit {
			result <- TodoResult{keepgoing, server_err}
			break
		}
		out_msg := "Todo List Loaded"
		fmt.Println(out, out_msg)
	case "7":
		Show_Instructions(out)
	case "8":
		flipThreadSafe()
		if thread_safe_switch == "" {
			fmt.Fprintln(out, "Thread safety off")
		} else {
			fmt.Fprintln(out, "Thread safety on")
		}
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

func flipThreadSafe() {
	InfoLog("CLI", "Switch ThreadSafety")
	mutex.Lock()
	if thread_safe_switch == "" {
		thread_safe_switch = "/threadsafe"
	} else {
		thread_safe_switch = ""
	}
	mutex.Unlock()
}

func get_Svr(url string) (*http.Response, error, bool, bool, error) {
	resp, err := http.Get(url)
	if resp == nil {
		WarnLog("CLI", "Server didn't exist")
		return nil, nil, true, true, &RequestError{500, fmt.Errorf("no response from server")}
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
		return nil, nil, true, true, &RequestError{500, errNoResponse}
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
	return input, todo_name[1 : len(input)+1], nil
}
