//POST
// curl -X POST http://localhost:8080/images/123/test.png \
//   -F "file=@/path/to/test.png"

// GET
// curl -X GET http://localhost:8080/images/123/test.png -o test.png

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"product-images/files"
	"product-images/handlers"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

var bindAddress = ":9090"
var basePath = "D:/Go_Microservices/10_HTTP_Save_And_Retrieve_Files"

func main() {
	// hclog.New(...) → makes a structured logger with level = "debug".
	// l.StandardLogger(...) → wraps it into a regular Go log.Logger.
	// InferLevels: true → lets the wrapper map log.Println(...),
	// log.Printf(...), etc. to the right hclog level (info, warn, error) instead of always treating them as plain “print.”
	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString("debug"),
		},
	)

	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})
	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(basePath, 5*1024*1000)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	//Create handlers
	fh := handlers.NewFiles(stor, l)
	//Create a new serve muc and register handlers
	sm := mux.NewRouter()

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	//get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(basePath))),
	)

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address", bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)

}
