package handlers

import (
	"net/http"
	"restful-service/data"
)

// // swagger:route POST /products products createProduct
// // Create a new product
// //
// // responses:
// //	200: productResponse
// //  422: errorValidation
// //  501: errorResponse

// // Create handles POST requests to add new products
// func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
// 	// fetch the product from the context
// 	prod := r.Context().Value(KeyProduct{}).(*data.Product)

// 	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
// 	data.AddProduct(*prod)
// }

func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product pointer from the context
	prod, ok := r.Context().Value(KeyProduct{}).(*data.Product)
	if !ok || prod == nil {
		http.Error(rw, "Product missing in context", http.StatusInternalServerError)
		return
	}

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)

	// AddProduct expects a value, not pointer, so dereference
	data.AddProduct(*prod)

	rw.WriteHeader(http.StatusCreated)
}
