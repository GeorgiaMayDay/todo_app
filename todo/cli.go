package todo

import (
	"bufio"
	"fmt"
	"io"
)

var invalid_opt_msg = "You've entered an invalid option"

func ReadAndOutput(in io.Reader, out io.Writer, list baseList) {
	reader := bufio.NewScanner(in)
	reader.Scan()
	option := reader.Text()

	switch option {
	case "1":
		list.outputTodos(out)
	case "2":
		reader.Scan()
		todo_name := reader.Text()
		out_msg := "\"" + todo_name + "\" added"
		fmt.Fprint(out, out_msg)
	default:
		fmt.Fprint(out, invalid_opt_msg)
	}

}
