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
	vars := mux.Vars(r)
	broker := vars["broker"]

	fmt.Fprint(w, "Broker Registration Endpoint Hit: You Registered ", broker)
}

// Provisioning function for services
func provisionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	fmt.Fprint(w, "Provisioning Endpoint Hit: You Provisioned A ", service)
}

// Deprovisioning function for services
func deprovisionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	fmt.Fprint(w, "Deprovisioning Endpoint Hit: You Deprovisioned ", instance)
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
	r.HandleFunc("/register/{broker}", registerHandler)
	r.HandleFunc("/provision/{service}", provisionHandler)
	r.HandleFunc("/deprovision/{instance}", deprovisionHandler)
	r.HandleFunc("/list", listHandler)

	log.Fatal(http.ListenAndServe(":8081", r))

}
