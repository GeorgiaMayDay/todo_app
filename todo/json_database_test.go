package todo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const test_file_name = "db_testing.json"

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", test_file_name)

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("didn't want an error but got %s", err.Error())
	}
}

const InitialDataString string = `[{"Name":"Iron","Status":"Todo"},{"Name":"Eat","Status":"Complete"},{"Name":"Hunker","Status":"Complete"},{"Name":"Mine","Status":"Todo"},{"Name":"Shear","Status":"Todo"},{"Name":"Cut","Status":"Todo"}]`

func TestRead(t *testing.T) {

	db, cleanDB := createTempFile(t, InitialDataString)
	defer cleanDB()

	json_output, err := LoadState(*db)

	assertNoError(t, err)
	json_string_output := string(json_output)
	require.JSONEq(t, InitialDataString, json_string_output)

}
