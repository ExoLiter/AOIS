package main

import (
	"fmt"
	"io"
	"os"

	"lab4/internal/hashtable"
)

func main() {
	if err := runFunc(os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		exitFunc(1)
	}
}

func Run(input io.Reader, output io.Writer) error {
	table, err := hashtable.NewTable(hashtable.DefaultTableSize)
	if err != nil {
		return err
	}
	if err := insertInitial(table, initialRecords); err != nil {
		return err
	}
	if err := printReport(output, table); err != nil {
		return err
	}
	return RunCLI(table, input, output)
}

func insertInitial(table *hashtable.Table, records []hashtable.Record) error {
	for _, record := range records {
		if err := table.Insert(record.Key, record.Value); err != nil {
			return fmt.Errorf("insert %s: %w", record.Key, err)
		}
	}
	return nil
}

func printReport(output io.Writer, table *hashtable.Table) error {
	if _, err := fmt.Fprintln(output, "Hash table content:"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, table.Render()); err != nil {
		return err
	}
	_, err := fmt.Fprintf(output, "Load factor: %.2f\n", table.LoadFactor())
	return err
}

var runFunc = Run
var exitFunc = os.Exit
