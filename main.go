package main

import (
    "net/http"
)

func main() {
    // Define a handler function to serve the index.html file
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    // Start the server on port 8080
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
