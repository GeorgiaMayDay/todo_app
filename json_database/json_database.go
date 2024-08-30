package todo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func GetAll(db os.File) ([]byte, error) {
	total_todo_list := ""
	db.Seek(0, 0)
	s := bufio.NewScanner(&db)
	for s.Scan() {
		fmt.Println(s.Text())
		total_todo_list += s.Text()
	}
	json_slice, err := json.MarshalIndent(total_todo_list, "\n", "	")
	return json_slice, err
}
