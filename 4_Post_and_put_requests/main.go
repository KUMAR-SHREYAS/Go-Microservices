// GET - curl request - curl http://localhost:9090/
// POST - curl -X POST "http://localhost:8080/products" -H "Content-Type: application/json" -d "{\"id\":3,\"name\":\"Cappuccino\",\"description\":\"Coffee with steamed milk and foam\",\"price\":3.15,\"sku\":\"capp123\",\"createdOn\":\"2025-09-19T10:00:00Z\",\"updatedOn\":\"2025-09-19T10:00:00Z\"}"
// PUT - curl -X PUT "http://localhost:9090/products/2" -H "Content-Type: application/json" -d "{\"id\":2,\"name\":\"Espresso Double Shot\",\"description\":\"Extra strong coffee without milk\",\"price\":2.49,\"sku\":\"fjd34\",\"createdOn\":\"2025-09-19T09:00:00Z\",\"updatedOn\":\"2025-09-19T12:30:00Z\"}"

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"restful-service/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "RESTFul Service:", log.LstdFlags) // log.LstdFlags contains the date and time
	prodHandler := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", prodHandler)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l, //if nil then logging is done using the standard logger
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
	}
	go func() {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatalf("Error Starting the server: %s\n", err) // Fatalf is equivalent to Printf() followed by a call to os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // catch OS signals
	sig := <-quit
	log.Println("Got Shutdown Signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		l.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
