// Create a simple Http server,
// handle two functions distinguished using the path
// / and /goodbye, inside / use error handling as well, and write the response , httpBadRequest
// to the user.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// reqeusts to the path /goodbye with be handled by this function
	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})
	// any other request will be handled by this function
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Running / handler")
		// read the body
		request, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading the body", err)
			http.Error(rw, "Error reading the response Body", http.StatusBadRequest)
			return
		}
		// write the response
		fmt.Fprintf(rw, "Hello %s\n", request)
	})
	// Listen for connections on all ip addresses (0.0.0.0)
	// port 9090
	log.Println("Starting a Server with port 9090")
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}

//curl commands
// curl -v -X POST http://localhost:9090/ -d "Alice" for '/' handler
// curl http://localhost:9090/goodbye for '/goodbye' handler
// first run the 1.go file.
