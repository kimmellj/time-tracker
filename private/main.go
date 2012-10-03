package main

import (
    "fmt"
    "net/http"
)

type Controller interface {
}

func router (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

    var controller Controller
    
}

func main() {
    http.HandleFunc("/", router)
    http.ListenAndServe(":8080", nil)
}