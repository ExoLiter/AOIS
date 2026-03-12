package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainProgram(t *testing.T) {
	// Prepare stdin with a simple valid expression.
	tmpIn, err := os.CreateTemp("", "aois-lab2-stdin-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp stdin: %v", err)
	}
	defer os.Remove(tmpIn.Name())
	if _, err := tmpIn.WriteString("a\n"); err != nil {
		t.Fatalf("failed to write stdin: %v", err)
	}
	if _, err := tmpIn.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("failed to seek stdin: %v", err)
	}

	// Capture stdout.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	oldStdin, oldStdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = tmpIn, w
	defer func() {
		os.Stdin, os.Stdout = oldStdin, oldStdout
	}()

	main()

	_ = w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "T0:") || !strings.Contains(out, "T1:") {
		t.Fatalf("unexpected output, main may not have completed: %q", out)
	}
}
