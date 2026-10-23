package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func parseCharset(s string) (display.Charset, error) {
	switch s {
	case "us":
		return display.CharsetUS, nil
	case "ru":
		return display.CharsetRU, nil
	default:
		return 0, fmt.Errorf("unknown charset %q, want %q or %q", s, "us", "ru")
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: vfdctl [--charset us|ru] <command> [args...]")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}

	fmt.Fprintln(os.Stderr, "\nCommands:")
	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %-12s %s\n", name, commands[name].description)
	}
}

func main() {
	flag.Usage = usage
	charsetName := flag.String("charset", "us", `charset to use: "us" or "ru"`)
	flag.Parse()

	if err := run(*charsetName, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(charsetName string, args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}

	charset, err := parseCharset(charsetName)
	if err != nil {
		return err
	}

	cmd, cmdArgs := args[0], args[1:]
	handler, ok := commands[cmd]
	if !ok {
		return fmt.Errorf("unknown command %q", cmd)
	}

	d, err := display.Open(charset)
	if err != nil {
		return fmt.Errorf("opening display: %w", err)
	}
	defer d.Close()

	return handler.run(d, cmdArgs)
}
