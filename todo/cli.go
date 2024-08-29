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

var invalid_opt_msg = "You've entered an invalid option"

func ReadAndOutput(in io.Reader, out io.Writer, list baseList) bool {
	reader := bufio.NewScanner(in)
	option := readLine(reader)

	switch option {
	case "1":
		list.outputTodos(out)
		return true
	case "2":
		todo_name := readLine(reader)
		out_msg := "\"" + todo_name + "\" added"
		fmt.Fprintln(out, out_msg)
		return true
	case "Quit":
		return false
	case "Q":
		return false
	case "":
		return false
	default:
		fmt.Fprintln(out, invalid_opt_msg)
		return true
	}
}
