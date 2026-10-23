package main

import (
	"log"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func main() {
	d, err := display.Open(display.US)
	if err != nil {
		log.Fatalf("opening display: %v", err)
	}
	defer d.Close()

	d.Initialize()
}
