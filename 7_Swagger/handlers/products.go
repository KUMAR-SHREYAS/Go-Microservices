// this handler file is for GET request for products and customized Logger
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"restful-service/data"
	"strconv"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}
type Products struct {
	l *log.Logger
	v *data.Validation
}

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /product/[id]")

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}
