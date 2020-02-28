package main

import (
	"fmt"
	"log"
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

// Handler function for registering a tweed/broker
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Broker Registration Endpoint Hit")
}

// Provisioning function for services
func provisionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Provisioning Endpoint Hit")
}

// Deprovisioning function for services
func deprovisionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Deprovisioning Endpoint Hit")
}

// Listing function for services
func listHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Listing Endpoint Hit")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/tweed", tweedHandler)
	r.HandleFunc("/svc", svcHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/provision", provisionHandler)
	r.HandleFunc("/deprovision", deprovisionHandler)
	r.HandleFunc("/list", listHandler)

	log.Fatal(http.ListenAndServe(":8081", r))

}
