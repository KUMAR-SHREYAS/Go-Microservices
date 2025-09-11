// In this file we will refactor the http server to use a separate package for handling routes and requests.
// The handlers goodbye and hello will be moved to a new package called handlers.
// New practices such as spinning up a server using a function will be implemented.
// The code will be organized in a more modular way.
// The main.go file will be responsible for starting the server and setting up the routes.
// The handlers package will contain the logic for handling the requests which will be directly called from the main.go file.

// Step 1: create a bindAddress variable for :9090
// Step 2: create a new logger instance
// Step 3: create the handlers passing the logger instance
// Step 4: create a new http.ServeMux and register the handlers
// Step 5: create a new http.Server instance
// Step 6: start the server using ListenAndServe method using Goroutines
// Step 7: Make a channel to  handle communication for os.Signal to trap SIGTERM or SigInterrupt
// Step 8: gracefully shutdown the server on receiving the signal (after making a new context for timeout)

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"refactoring-http-server/handlers"
	"syscall"
	"time"
)

var bindAddress = ":9090"

func main() {
	// logger instance (where to log, prefix, flags(date/time etc))
	l := log.New(os.Stdout, "Server refactored code: ", log.LstdFlags)
	// create handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodBye(l)
	// create a new serveMux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)
	//spin up a server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	// start the server using goRoutines
	go func() {
		l.Println("Starting the server on port 9090")
		err := s.ListenAndServe() //start listening and serving requests on :9090 port
		if err != nil {
			l.Printf("Error starting the server: %s\n", err)
			os.Exit(1)
		}
	}()
	// trap sigterm or interrupt using channels and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)      // make a channel of type os.Signal having buffer size of 1
	signal.Notify(sigChan, os.Interrupt)    // on ctrl+C
	signal.Notify(sigChan, syscall.SIGTERM) // on termination signal from kubernetes

	actualSignal := <-sigChan //block until a signal is received
	log.Println("Received Signal: ", actualSignal)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30) // new context is created from background and 30 seconds timeout before shutting down
	defer cancel()
	s.Shutdown(ctx)
	log.Println("Shutting down the server")
}
