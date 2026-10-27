package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/partite-ai/example-ticket-pipeline-shortlink/server"
	"github.com/partite-ai/example-ticket-pipeline-shortlink/store"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	ttl := flag.Duration("ttl", store.DefaultTTL, "how long a link lives before it expires")
	flag.Parse()
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, server.New(store.New(store.WithTTL(*ttl)))))
}
