package main

import (
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "nothing to see here\n")
}

func main() {
	log.Printf("Starting server on localhost:4000")
	http.HandleFunc("/", gitserver.MiddlewareFunc(handler))
	http.ListenAndServe(":4000", nil)
}
