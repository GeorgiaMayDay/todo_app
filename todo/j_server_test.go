package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	t.Run("happy path test", func(t *testing.T) {
		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()
		server, err := NewJsonTodoServer(tempfile.Name())

		request := httptest.NewRequest(http.MethodGet, "/GET", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertNoError(t, err)

		assertStrings(t, response.Body.String(), generateTodoListAsString())
	})
}
