// this handler file is for GET request for products and customized Logger
package handlers

import (
	"log"
	"net/http"
	"regexp"
	"restful-service/data"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// serveHttp is the main entry point for all requests through the products handler
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
	}
	if r.Method == http.MethodPut {
		p.l.Println("PUT ", r.URL.Path)
		//Expect the id in the URL
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid Url more than one id")
			http.Error(rw, "Invalid Url", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid Url more than one capture group")
			http.Error(rw, "Invalid Url", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid Url unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProducts(id, rw, r)
		return
	}
	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter) {
	p.l.Println("Handle GET Products")
	prod := data.GetProducts()
	err := prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}
func (p Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT request")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}
