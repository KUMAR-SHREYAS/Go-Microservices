/*
Curl commands to test the API:
# create

	curl -s -X POST -H "Content-Type: application/json" \
	  -d '{"name":"Alice","email":"alice@example.com"}' \
	  http://localhost:8080/users

# get by id
curl http://localhost:8080/users/1

# list
curl http://localhost:8080/users
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// define global vars
var (
	users     = make(map[int64]*User) // id-User Struct
	usersMu   sync.RWMutex            // Mutex for users map
	idCounter int64                   //user ids counter
)

func writeJSON(rw http.ResponseWriter, v interface{}, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_ = json.NewEncoder(rw).Encode(v)
}

// handler functions
// GET POST for /users
func getAllUsers(rw http.ResponseWriter) {
	//lock-unlock and append users to list
	usersMu.RLock()
	list := make([]*User, 0, len(users))
	for _, u := range users {
		list = append(list, u)
	}
	usersMu.RUnlock()
	//return list as json using http.ResponseWriter
	writeJSON(rw, list, http.StatusOK)
}
func createUser(rw http.ResponseWriter, r *http.Request) {
	//some initial checks
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {

	}
	var payload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Email = strings.TrimSpace(payload.Email)
	if payload.Name == "" || payload.Email == "" {
		http.Error(rw, "name and email are required", http.StatusBadRequest)
		return
	}

	// increment id counter and create user struct
	id := atomic.AddInt64(&idCounter, 1)
	u := &User{
		ID:        id,
		Name:      payload.Name,
		Email:     payload.Email,
		CreatedAt: time.Now().UTC(),
	}

	usersMu.Lock()
	users[id] = u
	usersMu.Unlock()

	rw.Header().Set("Location", fmt.Sprintf("/users/%d", id))
	writeJSON(rw, u, http.StatusCreated)
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllUsers(rw)
	case http.MethodPost:
		createUser(rw, r)
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET user by id
func userHandler(rw http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		http.Error(rw, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64) // base and bitSize (int64 on decimal)
	if err != nil || id < 1 {
		http.Error(rw, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		usersMu.RLock()
		u, found := users[id]
		usersMu.RUnlock()
		if !found {
			http.Error(rw, "user not found", http.StatusNotFound)
			return
		}
		writeJSON(rw, u, http.StatusOK)
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// loggingMiddlware handler for server logs
func loggingMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(rw, r)
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
}

// main function
func main() {
	// create handlers
	mux := http.NewServeMux()
	//handlers
	// /users /GET or /POST method
	mux.HandleFunc("/users", usersHandler)
	// /users/id
	mux.HandleFunc("/users/", userHandler)

	// spin up server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddlware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// run server using goroutines
	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err) //equivalent to printf followed by os.Exit(1)
		}
	}()
	// shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // similar to reading os.Interrupt using quit channel
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Println("server stopped")
}
