package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Database Schema
type Name struct {
	NConst    string `json:"nconst"`
	Name      string `json:"name"`
	BirthYear string `json:"birthYear"`
	DeathYear string `json:"deathYear"`
}

// Custom Error Message
type Error struct {
	Message string `json:"error"`
}

func renderJSON(rw http.ResponseWriter, val interface{}, statusCode int) {
	rw.WriteHeader(statusCode)
	_ = json.NewEncoder(rw).Encode(val)
}
func main() {
	router := mux.NewRouter()

	//Open The connection to Database using Connect()  SQL
	dbSQL, err := NewPostgreSQLsql()
	if err != nil {
		log.Fatalf("Could not initialize Database Connection using sql %s", err)
		return
	}
	defer dbSQL.Close() // close the connection(mandatory)

	router.HandleFunc("/names/sql/{id}", func(rw http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		name, err := dbSQL.FindNConst(id)
		if err != nil {
			renderJSON(rw, &Error{Message: err.Error()}, http.StatusInternalServerError)
		}
		renderJSON(rw, &name, http.StatusOK)
	})

	//Open The connection to Database using connectContext() SQLX
	dbSQLX, err := NewPostgreSQLsqlx()
	if err != nil {
		log.Fatalf("Could not initialize Database Connection using sqlx %s", err)
		return
	}
	defer dbSQLX.Close() // close the connection(mandatory)

	router.HandleFunc("/names/sqlx/{id}", func(rw http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		name, err := dbSQLX.FindNConst(id)
		if err != nil {
			renderJSON(rw, &Error{Message: err.Error()}, http.StatusInternalServerError)
		}
		renderJSON(rw, &name, http.StatusOK)
	})

	//Open The connection to Database using connectContext()  PGX
	pgxDB, err := NewPostgreSQLpgx()
	if err != nil {
		log.Fatalf("Could not initialize Database Connection using pgx %s", err)
		return
	}
	defer pgxDB.Close() // close the connection(mandatory)

	router.HandleFunc("/names/pgx/{id}", func(rw http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		name, err := pgxDB.FindNConst(id)
		if err != nil {
			renderJSON(rw, &Error{Message: err.Error()}, http.StatusInternalServerError)
		}
		renderJSON(rw, &name, http.StatusOK)
	})

	fmt.Println("Starting server at :9090")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
