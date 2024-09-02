package todo

import (
	"os"
)

func LoadState(db os.File) ([]byte, error) {
	var json_slice []byte
	db.Seek(0, 0)
	json_slice, err := os.ReadFile(db.Name())
	return json_slice, err
}

func SaveState(db os.File, json_obj []byte) error {
	db.Seek(0, 0)
	_, err := db.Write(json_obj)
	if err != nil {
		return err
	}
	return nil
}
