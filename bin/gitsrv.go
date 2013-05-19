package main

import (
    "net/http"
    "log"
    "gitsrv"
)

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Processing %s!\n", r.URL.Path[0:])

    parsedRoute := gitsrv.MatchRoute(r)
    if parsedRoute != nil {
        parsedRoute.Dispatch(w, r)
    }
}

func main() {
    log.Printf("Starting server on localhost:4000")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":4000", nil)
}

