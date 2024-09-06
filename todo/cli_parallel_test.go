package todo

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestIntergrationParellelProcessing(t *testing.T) {

	t.Run("That CLI can add Todos", func(t *testing.T) {
		t.Parallel()
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name(), "")

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

	t.Run("That CLI can get Todo list in a threadsafe manner", func(t *testing.T) {
		t.Parallel()
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		tempfile_useless, cleanUpFile_useless := createTempFile(t, "[]")
		defer cleanUpFile_useless()

		server, err := NewJsonTodoServer(tempfile_useless.Name(), tempfile.Name())

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		toThreadSafeMode(ctx, ts, finishChan, t)

		in := strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString())
	})

	t.Run("That CLI can add Todos in a threadsafe manner", func(t *testing.T) {
		t.Parallel()
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		tempfile_useless, cleanUpFile_useless := createTempFile(t, "[]")
		defer cleanUpFile_useless()

		server, err := NewJsonTodoServer(tempfile_useless.Name(), tempfile.Name())

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		toThreadSafeMode(ctx, ts, finishChan, t)

		in := strings.NewReader("2\nBrush")

		ReadAndOutput(ctx, in, &bytes.Buffer{}, ts.URL, finishChan)

		in = strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString()+"7. Brush: Todo\n")
	})

	t.Run("That CLI can delete Todos", func(t *testing.T) {
		t.Parallel()
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name(), "")

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		in := strings.NewReader("3\n6")
		unTestedReadCmd(ctx, ts, finishChan, t, in)

		in = strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString()[:len(generateTodoListAsString())-13]) // generateTodoListAsString()[:len(generateTodoListAsString())-10]+"Brush: Todo\n")
	})

	t.Run("That CLI can load Todos", func(t *testing.T) {
		t.Parallel()
		var trace_id string = uuid.NewString()
		ctx := context.WithValue(context.Background(), string("Trace_id"), trace_id)
		finishChan := make(chan TodoResult, 1)
		output := &bytes.Buffer{}

		tempfile, cleanUpFile := createTempFile(t, InitialDataString)
		defer cleanUpFile()

		server, err := NewJsonTodoServer(tempfile.Name(), "")

		assertNoError(t, err)

		ts := httptest.NewServer(server.Handler)

		in := strings.NewReader("2\nExample")
		unTestedReadCmd(ctx, ts, finishChan, t, in)

		in = strings.NewReader("6")
		unTestedReadCmd(ctx, ts, finishChan, t, in)

		in = strings.NewReader("1")

		ReadAndOutput(ctx, in, output, ts.URL, finishChan)

		got := <-finishChan

		assertNoError(t, got.Err)

		assertStrings(t, output.String(), generateTodoListAsString())
	})
}

func toThreadSafeMode(ctx context.Context, ts *httptest.Server, finishChan chan TodoResult, t *testing.T) {
	in := strings.NewReader("8")

	ReadAndOutput(ctx, in, &bytes.Buffer{}, ts.URL, finishChan)

	got := <-finishChan

	assertNoError(t, got.Err)
}

func unTestedReadCmd(ctx context.Context, ts *httptest.Server, finishChan chan TodoResult, t *testing.T, in io.Reader) {

	ReadAndOutput(ctx, in, &bytes.Buffer{}, ts.URL, finishChan)

	got := <-finishChan

	assertNoError(t, got.Err)
}
