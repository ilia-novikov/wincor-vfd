package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ilia-novikov/wincor-vfd/display"
)

type command struct {
	description string
	run         func(*display.Display, []string) error
}

var commands = map[string]command{
	"print":      {"print <text> — print text at the current cursor position", printCmd},
	"clear":      {"clear the whole display", clearCmd},
	"init":       {"clear, select charset, and reset the cursor to (0, 0)", initCmd},
	"move":       {"move <x> <y> — move the cursor to the given 0-based column/row", moveCmd},
	"bs":         {"move the cursor one space to the left", backspaceCmd},
	"lf":         {"move the cursor down a line", linefeedCmd},
	"cr":         {"move the cursor to the start of the current line", carriageReturnCmd},
	"delete-eol": {"erase text to the end of the current line", deleteToEndOfLineCmd},
}

func printCmd(d *display.Display, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: vfdctl print <text>")
	}
	return d.PrintText(strings.Join(args, " "))
}

func clearCmd(d *display.Display, _ []string) error {
	return d.ClearDisplay()
}

func initCmd(d *display.Display, _ []string) error {
	return d.Initialize()
}

func moveCmd(d *display.Display, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: vfdctl move <x> <y>")
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid x %q: %w", args[0], err)
	}

	y, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid y %q: %w", args[1], err)
	}

	return d.MoveCursor(x, y)
}

func backspaceCmd(d *display.Display, _ []string) error {
	return d.SendBackspace()
}

func linefeedCmd(d *display.Display, _ []string) error {
	return d.SendLineFeed()
}

func carriageReturnCmd(d *display.Display, _ []string) error {
	return d.SendCarriageReturn()
}

func deleteToEndOfLineCmd(d *display.Display, _ []string) error {
	return d.DeleteToEndOfLine()
}
