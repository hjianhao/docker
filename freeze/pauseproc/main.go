package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Function invoked!")
    fmt.Printf("Function invoked!")
}

func main() {
    fmt.Println("Starting server...")
    http.HandleFunc("/invoke", handler)
    http.ListenAndServe(":8088", nil)
}
