package handlers

import (
	"log"
	"microservices/data"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		break
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
