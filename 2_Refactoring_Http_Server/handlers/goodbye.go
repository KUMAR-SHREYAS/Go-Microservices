// This file will modularise the goodbye handler, it will log the print statements using
// structs and log library for http.Request and http.ResponseWriter object
// passed during curl command.
package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// Goodbye strct which will modularise the log library's Logger function generating logs to io.Writer object
type Goodbye struct {
	l *log.Logger // helps during the execution of multiple concurrent goroutines by providing serialized access to the io.Writer
}

// Creating the Goodbye struct object using a NewGoodBye function, which will return pointer to the struct onject
func NewGoodBye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

// ServeHTTP will implement the http.Handler interface, when handle() function is called from serveMux object
// handler interface has only one method ServeHTTP which takes http.ResponseWriter and http.Request as parameters
//
//	type Handler interface {
//	    ServeHTTP(ResponseWriter, *Request)
//	}
//
// serveHttp contains the logic for handling the request Goodbye
func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("GoodBye Handler")
	fmt.Fprintf(rw, "Goodbye World")
}
