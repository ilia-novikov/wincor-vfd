package main

import (
	"fmt"
	"strings"

	"github.com/ilia-novikov/wincor-vfd/display"
)

var commands = map[string]func(*display.Display, []string) error{
	"print": printCmd,
	"clear": clearCmd,
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
