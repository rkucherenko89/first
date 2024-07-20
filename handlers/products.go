package handlers

import (
	"context"
	"errors"
	"first/customErrors"
	"first/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		p.getProducts(rw, r)
// 		return
// 	case http.MethodPost:
// 		p.addProduct(rw, r)
// 		return
// 	case http.MethodPut:
// 		trimmedPath := strings.TrimPrefix(r.URL.Path, "/")
// 		pathSegments := strings.Split(trimmedPath, "/")

// 		if pathSegments[0] != "products" || len(pathSegments) != 2 {
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		idToPut, err := strconv.Atoi(pathSegments[1])
// 		if err != nil {
// 			http.Error(rw, "ID is not a digit", http.StatusBadRequest)
// 			return
// 		}

// 		p.updateProduct(rw, r, idToPut)
// 		return
// 	default:
// 		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 	}
// }

func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	np, ok := r.Context().Value(KeyProduct{}).(*data.Product)
	if !ok {
		http.Error(rw, "Oops! The Product is in invalid format", http.StatusBadRequest)
		return
	}
	data.AddProduct(np)
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	id, err := strconv.Atoi(pathVars["id"])
	if err != nil {
		http.Error(rw, "The ID should be an integer", http.StatusBadRequest)
		return
	}
	np, ok := r.Context().Value(KeyProduct{}).(*data.Product)
	if !ok {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	np.ID = id
	pr, err := data.UpdateProduct(np)
	if err != nil {
		var customErr *customErrors.ErrorProductNotFound
		if errors.As(err, &customErr) {
			http.Error(rw, customErr.Error(), http.StatusBadRequest)
			return
		}
	}
	err = pr.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing Product", err)
			http.Error(w, "Error reading Product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, &prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
