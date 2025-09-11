package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Hello struct holds a logger instance to log messages
type Hello struct {
	l *log.Logger
}

// NewHello creates a new Hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHttp is the HTTP handler method for Hello
// It reads the request body and responds with "Hello <body content>"
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// log that the handler was called
	h.l.Println("Hello World Handler")

	// read the entire request body
	req, err := io.ReadAll(r.Body)
	if err != nil {
		// log error and send HTTP 400 Bad Request response
		h.l.Println("Error reading request body")
		http.Error(rw, "Error handling the request writing the Bad request", http.StatusBadRequest)
		return
	}

	// write response back to client
	fmt.Fprintf(rw, "Hello %s", req)
}
