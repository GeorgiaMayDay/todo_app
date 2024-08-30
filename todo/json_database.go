package todo

import (
	"os"
)

func GetAll(db os.File) ([]byte, error) {
	var json_slice []byte
	db.Seek(0, 0)
	json_slice, err := os.ReadFile(db.Name())
	return json_slice, err
}
