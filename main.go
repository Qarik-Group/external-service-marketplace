package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func tweedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Tweed Endpoint Hit")
}

func svcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Service Catalog Endpoint Hit")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/tweed", tweedHandler)
	r.HandleFunc("/svc", svcHandler)

	http.ListenAndServe(":80", r)
}
