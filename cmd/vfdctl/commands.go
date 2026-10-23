package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ilia-novikov/wincor-vfd/display"
)

var commands = map[string]func(*display.Display, []string) error{
	"print":       printCmd,
	"clear":       clearCmd,
	"init":        initCmd,
	"move":        moveCmd,
	"set-charset": setCharsetCmd,
	"bs":          backspaceCmd,
	"lf":          linefeedCmd,
	"cr":          carriageReturnCmd,
	"delete-eol":  deleteToEndOfLineCmd,
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

func setCharsetCmd(d *display.Display, _ []string) error {
	return d.SetCharset()
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
