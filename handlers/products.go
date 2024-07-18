package handlers

import (
	"first/customErrors"
	"first/data"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(rw, r)
		return
	case http.MethodPost:
		p.addProduct(rw, r)
		return
	case http.MethodPut:
		trimmedPath := strings.TrimPrefix(r.URL.Path, "/")
		pathSegments := strings.Split(trimmedPath, "/")

		if pathSegments[0] != "products" || len(pathSegments) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idToPut, err := strconv.Atoi(pathSegments[1])
		if err != nil {
			http.Error(rw, "ID is not a digit", http.StatusBadRequest)
			return
		}

		p.updateProduct(rw, r, idToPut)
		return
	default:
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(rw http.ResponseWriter, _ *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	np := &data.Product{}
	err := np.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// p.l.Printf("New Product: %#v", np)
	data.AddProduct(np)
	pl := data.GetProducts()
	err = pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request, id int) {
	np := &data.Product{}
	err := np.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	pr, err := data.UpdateProduct(id, np)
	if err != nil {
		if customErr, ok := err.(*customErrors.ErrorProductNotFound); ok {
			http.Error(rw, customErr.Error(), http.StatusBadRequest)
			return
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = pr.ToJSON(rw)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	// pl := data.GetProducts()
	// err = pl.ToJSON(rw)
	// if err != nil {
	// 	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }

}
