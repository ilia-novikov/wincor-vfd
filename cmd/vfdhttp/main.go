package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ilia-novikov/wincor-vfd/httpapi"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "address to listen on")
	flag.Parse()

	srv := &http.Server{
		Addr:    *addr,
		Handler: httpapi.New(),
	}

	log.Printf("vfdhttp listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
