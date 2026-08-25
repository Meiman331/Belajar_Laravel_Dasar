package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Products")
	}).Methods(http.MethodGet)

	r.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "ID Produk: %s", id)
	}).Methods(http.MethodGet)

	addr := ":8080"
	log.Printf("Server berjalan di %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

