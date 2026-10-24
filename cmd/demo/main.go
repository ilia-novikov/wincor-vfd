package main

import (
	"log"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func main() {
	d, err := display.Open(display.CharsetUS)
	if err != nil {
		log.Fatalf("opening display: %v", err)
	}
	defer d.Close()

	if err := d.Initialize(); err != nil {
		log.Fatalf("initializing display: %v", err)
	}
	if err := d.PrintText("Hello, wincor-vfd!"); err != nil {
		log.Fatalf("printing text: %v", err)
	}
}
