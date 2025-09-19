// this handler file is for GET request for products and customized Logger
package handlers

import (
	"context"
	"log"
	"net/http"
	"restful-service/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// No need to ServeHTTP
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	prod := data.GetProducts()
	err := prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	// prod := &data.Product{}
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	// }
	// p.l.Println("AddProduct is called")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT request")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	// A HandlerFunc(f) calls function f as Http Handler , basically conversion
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// p.l.Println("Middleware is called")
		prod := data.Product{}
		err := prod.FromJSON(r.Body) //Read the request body into JSON using FromJSON method of prod , use additonal err handling
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}
		//define the context prod using KeyProduct{} as a key to make a
		// derived context out of parent context in this case r.Context()
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		// Change the parent context  of r to derived context of ctx
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
