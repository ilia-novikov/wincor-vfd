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

func main() {
	charsetName := flag.String("charset", "us", `charset to use: "us" or "ru"`)
	flag.Parse()

	if err := run(*charsetName, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(charsetName string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: vfdctl [--charset us|ru] <command> [args...]")
	}

	charset, err := parseCharset(charsetName)
	if err != nil {
		return err
	}

	fmt.Println(charset)
	return nil
}
