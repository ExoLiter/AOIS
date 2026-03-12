package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"lab4/internal/hashtable"
)

func TestCLIBasicFlow(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader(
		"1\nAA\nValue\n" +
			"2\nAA\n" +
			"3\nAA\nValue2\n" +
			"2\nAA\n" +
			"4\nAA\n" +
			"6\n" +
			"5\n" +
			"0\n",
	)
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
	text := output.String()
	if !strings.Contains(text, "Record added.") {
		t.Fatalf("expected create confirmation")
	}
	if !strings.Contains(text, "Value: Value") {
		t.Fatalf("expected read output")
	}
	if !strings.Contains(text, "Record updated.") {
		t.Fatalf("expected update confirmation")
	}
	if !strings.Contains(text, "Record deleted.") {
		t.Fatalf("expected delete confirmation")
	}
	if !strings.Contains(text, "Table dump:") {
		t.Fatalf("expected dump output")
	}
	if !strings.Contains(text, "Load factor:") {
		t.Fatalf("expected load factor output")
	}
}

func TestCLIInvalidInput(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("abc\n9\n0\n")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
	text := output.String()
	if !strings.Contains(text, "Invalid option") {
		t.Fatalf("expected invalid option warning")
	}
	if !strings.Contains(text, "Unknown option") {
		t.Fatalf("expected unknown option warning")
	}
}

func TestCLIReadNotFound(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("2\nAA\n0\n")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
	if !strings.Contains(output.String(), "Not found.") {
		t.Fatalf("expected not found message")
	}
}

func TestCLIUpdateDeleteErrors(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("3\nAA\nValue\n4\nAA\n0\n")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
	text := output.String()
	if !strings.Contains(text, "key not found") {
		t.Fatalf("expected not found errors")
	}
}

func TestCLIDuplicateInsert(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("1\nAA\nValue\n1\nAA\nValue2\n0\n")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
	if !strings.Contains(output.String(), "key already exists") {
		t.Fatalf("expected duplicate key error")
	}
}

func TestCLIPromptEOF(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("1\n")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
}

func TestCLIEOF(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	input := strings.NewReader("")
	var output bytes.Buffer
	if err := RunCLI(table, input, &output); err != nil {
		t.Fatalf("cli failed: %v", err)
	}
}

func TestCLIReadError(t *testing.T) {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	var output bytes.Buffer
	if err := RunCLI(table, errorReader{}, &output); err == nil {
		t.Fatalf("expected read error")
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}
