package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"lab4/internal/hashtable"
)

const (
	menuExit   = 0
	menuCreate = 1
	menuRead   = 2
	menuUpdate = 3
	menuDelete = 4
	menuLoad   = 5
	menuDump   = 6
)

func RunCLI(table *hashtable.Table, input io.Reader, output io.Writer) error {
	client := newCLI(table, input, output)
	return client.loop()
}

type cli struct {
	table  *hashtable.Table
	reader *bufio.Reader
	output io.Writer
}

func newCLI(table *hashtable.Table, input io.Reader, output io.Writer) *cli {
	return &cli{
		table:  table,
		reader: bufio.NewReader(input),
		output: output,
	}
}

func (c *cli) loop() error {
	for {
		c.printMenu()
		choice, exit, err := c.readChoice()
		if err != nil {
			return err
		}
		if exit {
			return nil
		}
		c.handleChoice(choice)
	}
}

func (c *cli) readChoice() (int, bool, error) {
	line, err := c.readLine("Select option: ")
	if err != nil {
		if err == io.EOF {
			return 0, true, nil
		}
		return 0, false, err
	}
	choice, err := strconv.Atoi(line)
	if err != nil {
		c.println("Invalid option. Enter a number.")
		return 0, false, nil
	}
	if choice == menuExit {
		return choice, true, nil
	}
	return choice, false, nil
}

func (c *cli) handleChoice(choice int) {
	switch choice {
	case menuCreate:
		c.handleCreate()
	case menuRead:
		c.handleRead()
	case menuUpdate:
		c.handleUpdate()
	case menuDelete:
		c.handleDelete()
	case menuLoad:
		c.handleLoad()
	case menuDump:
		c.handleDump()
	default:
		c.println("Unknown option.")
	}
}

func (c *cli) handleCreate() {
	key, ok := c.prompt("Key: ")
	if !ok {
		return
	}
	value, ok := c.prompt("Value: ")
	if !ok {
		return
	}
	if err := c.table.Insert(key, value); err != nil {
		c.println(err.Error())
		return
	}
	c.println("Record added.")
}

func (c *cli) handleRead() {
	key, ok := c.prompt("Key: ")
	if !ok {
		return
	}
	entry, found := c.table.Find(key)
	if !found {
		c.println("Not found.")
		return
	}
	c.println(fmt.Sprintf("Value: %s", entry.Value))
	c.println(fmt.Sprintf("V=%d, h=%d", entry.V, entry.Home))
}

func (c *cli) handleUpdate() {
	key, ok := c.prompt("Key: ")
	if !ok {
		return
	}
	value, ok := c.prompt("New value: ")
	if !ok {
		return
	}
	if err := c.table.Update(key, value); err != nil {
		c.println(err.Error())
		return
	}
	c.println("Record updated.")
}

func (c *cli) handleDelete() {
	key, ok := c.prompt("Key: ")
	if !ok {
		return
	}
	if err := c.table.Delete(key); err != nil {
		c.println(err.Error())
		return
	}
	c.println("Record deleted.")
}

func (c *cli) handleLoad() {
	c.println(fmt.Sprintf("Load factor: %.2f", c.table.LoadFactor()))
}

func (c *cli) handleDump() {
	c.println("Table dump:")
	c.println(c.table.Render())
}

func (c *cli) printMenu() {
	c.println("\nMenu:")
	c.println("1 - Add record")
	c.println("2 - Find record")
	c.println("3 - Update record")
	c.println("4 - Delete record")
	c.println("5 - Show load factor")
	c.println("6 - Show table")
	c.println("0 - Exit")
}

func (c *cli) prompt(label string) (string, bool) {
	value, err := c.readLine(label)
	if err != nil {
		return "", false
	}
	return value, true
}

func (c *cli) readLine(label string) (string, error) {
	if _, err := fmt.Fprint(c.output, label); err != nil {
		return "", err
	}
	line, err := c.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return strings.TrimSpace(line), io.EOF
		}
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (c *cli) println(value string) {
	_, _ = fmt.Fprintln(c.output, value)
}
