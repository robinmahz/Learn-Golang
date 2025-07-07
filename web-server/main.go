package main

import (
    "fmt"
    "log"
    "net/http"
)

// Middleware to log incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log request details
        log.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
        // Call the next handler
        next.ServeHTTP(w, r)
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404 Not Found", http.StatusNotFound)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprintf(w, "Hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    fmt.Fprintf(w, "POST request successful\n")
    name := r.FormValue("name")
    fmt.Fprintf(w, "Name = %s\n", name)
}

func main() {
    // Create a new ServeMux to handle routes
    mux := http.NewServeMux()

    // Set up file server and handlers
    fileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("/", fileServer)
    mux.HandleFunc("/form", formHandler)
    mux.HandleFunc("/hello", helloHandler)

    // Wrap the entire mux with the logging middleware
    fmt.Println("Starting server on port 8000...")
    if err := http.ListenAndServe(":8000", loggingMiddleware(mux)); err != nil {
        log.Fatal(err)
    }
}