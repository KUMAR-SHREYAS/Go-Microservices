package handlers

import (
	"context"
	"net/http"
	"restful-service/data"
)

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	// A HandlerFunc(f) calls function f as Http Handler , basically conversion
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// p.l.Println("Middleware is called")
		prod := &data.Product{}
		err := data.FromJSON(prod, r.Body) //Read the request body into JSON using FromJSON method of prod , use additonal err handling
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		//validate the product int the middleware only
		errs := p.v.Validate(prod)
		if errs != nil {
			p.l.Println("[ERROR] validating product", err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
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
