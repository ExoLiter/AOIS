package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"lab4/internal/hashtable"
)

func TestRun(t *testing.T) {
	var buffer bytes.Buffer
	input := strings.NewReader("0\n")
	if err := Run(input, &buffer); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	output := buffer.String()
	if !strings.Contains(output, "Hash table content:") {
		t.Fatalf("expected header in output")
	}
	if !strings.Contains(output, "Load factor:") {
		t.Fatalf("expected load factor in output")
	}
}

func TestRunInsertError(t *testing.T) {
	original := initialRecords
	defer func() { initialRecords = original }()
	initialRecords = []hashtable.Record{
		{Key: "AA", Value: "A"},
		{Key: "AA", Value: "B"},
	}
	var buffer bytes.Buffer
	input := strings.NewReader("0\n")
	if err := Run(input, &buffer); err == nil {
		t.Fatalf("expected run error")
	}
}

func TestRunPrintError(t *testing.T) {
	writer := &failingWriter{failAt: 2}
	input := strings.NewReader("0\n")
	if err := Run(input, writer); err == nil {
		t.Fatalf("expected run print error")
	}
}

func TestInsertInitialError(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	records := []hashtable.Record{
		{Key: "Анатомия", Value: "A"},
		{Key: "Анатомия", Value: "B"},
	}
	if err := insertInitial(table, records); err == nil {
		t.Fatalf("expected insert error on duplicate key")
	}
}

func TestPrintReportError(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	if err := table.Insert("Анатомия", "A"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := printReport(&failingWriter{failAt: 1}, table); err == nil {
		t.Fatalf("expected print error on first write")
	}
	if err := printReport(&failingWriter{failAt: 2}, table); err == nil {
		t.Fatalf("expected print error on second write")
	}
	if err := printReport(&failingWriter{failAt: 3}, table); err == nil {
		t.Fatalf("expected print error on third write")
	}
}

func TestMain(t *testing.T) {
	originalRun := runFunc
	originalExit := exitFunc
	defer func() {
		runFunc = originalRun
		exitFunc = originalExit
	}()
	runCalled := false
	runFunc = func(reader io.Reader, writer io.Writer) error {
		runCalled = true
		return nil
	}
	exitFunc = func(code int) {}
	main()
	if !runCalled {
		t.Fatalf("expected main to call runFunc")
	}
}

func TestMainError(t *testing.T) {
	originalRun := runFunc
	originalExit := exitFunc
	defer func() {
		runFunc = originalRun
		exitFunc = originalExit
	}()
	exitCode := 0
	runFunc = func(reader io.Reader, writer io.Writer) error {
		return errors.New("boom")
	}
	exitFunc = func(code int) {
		exitCode = code
	}
	main()
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}

type failingWriter struct {
	called int
	failAt int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	w.called++
	if w.called == w.failAt {
		return 0, errors.New("write failed")
	}
	return len(p), nil
}
