package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/partite-ai/example-shortlink/server"
	"github.com/partite-ai/example-shortlink/store"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, server.New(store.New())))
}
