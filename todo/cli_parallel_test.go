package todo

import (
	"bytes"
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestIntergrationParellelProcessing(t *testing.T) {

	t.Run("That CLI can add Todos", func(t *testing.T) {
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name())

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		in := strings.NewReader("2\nBrush")

		ReadAndOutput(ctx, in, &bytes.Buffer{}, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		in = strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got = <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString()+"7. Brush: Todo\n")
	})

	t.Run("That CLI can complete Todos", func(t *testing.T) {
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name())

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		in := strings.NewReader("4\nCut")

		ReadAndOutput(ctx, in, &bytes.Buffer{}, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		in = strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got = <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString()[:len(generateTodoListAsString())-5]+"Complete\n")
	})
}
