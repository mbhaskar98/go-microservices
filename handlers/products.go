package handlers

import (
	"log"
	"microservices/data"
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

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		components := strings.Split(r.URL.Path, "/")
		idStr := components[len(components) - 1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Unknown ID",http.StatusBadRequest)
			return
		}
		p.updateProduct(id, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Get products called")
	products := data.GetProducts()

	err := products.ToJSON(w)

	if err != nil {
		http.Error(w, "Some error while getting products", http.StatusInternalServerError)
	}

}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Add new product")

	 prod, err := productFromRequestBody(r)

	if err != nil {
		http.Error(w, "Unable to parse the body", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
	prod.ToJSON(w)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request){
	p.l.Println("Update product", id)
	prod, err := productFromRequestBody(r)
	if err != nil {
		http.Error(w, "Unable to parse the body", http.StatusBadRequest)
		return
	}

	prod, err = data.UpdateProduct(id, prod)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	prod.ToJSON(w)
}


func productFromRequestBody(r *http.Request) (*data.Product, error){
	prod := &data.Product{}

	err := prod.FromJSON(r.Body)

	if err != nil {
		return nil, err
	}
	return prod, nil
}