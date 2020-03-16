package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/external-service-marketplace/tweed"
)

type API struct {
	Config *Config
	Bind   string
}

type bind struct {
	Binding string `json:"as"`
	ID      string `json:"instance"`
	NoWait  bool   `json:"no-wait"`
}

type unbind struct {
	InstanceBinding []string `json:"instance/binding"`
	NoWait          bool     `json:"no-wait"`
}

type provision struct {
	ID          string   `json:"as"`
	Params      []string `json:"P"`
	NoWait      bool     `json:"no-wait"`
	ServicePlan []string `json:"service/plan"`
}

type deprovision struct {
	ID     string `json:"instance"`
	NoWait bool   `json:"no-wait"`
}

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}

func bindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	binding := vars["binding"]
	//nowait := vars["nowait"]

	var bind bind
	bind.ID = instance
	bind.Binding = binding

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&bind)

	//tweed.Connect(...,...)?
	tweed.Bind(w, r)

}

func unbindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	binding := vars["binding"]
	//nowait := vars["nowait"]

	var unbind unbind

	instancebinding := make([]string, 2)
	instancebinding[0] = instance
	instancebinding[1] = binding
	unbind.InstanceBinding = instancebinding

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&unbind)

	//tweed.Connect(...,...)?
	tweed.Bind(w, r)

}

func provisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	service := vars["service"]
	plan := vars["plan"]
	//nowait := vars["nowait"]

	var provision provision
	provision.ID = instance
	s := make([]string, 2)
	s[0] = service
	s[1] = plan
	provision.ServicePlan = s

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&provision)

	//tweed.Connect(...,...)?
	tweed.Provision(w, r)

}

func deprovisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	//nowait := vars["nowait"]

	var deprovision deprovision
	deprovision.ID = instance

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&deprovision)

	//tweed.Connect(...,...)?
	tweed.Deprovision(w, r)

}

func (api API) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//register broker
	r.HandleFunc("/reigster/{broker}", testResponse)
	//retrieve clouds
	r.HandleFunc("/clouds", testResponse)
	//retrieve a specific cloud
	r.HandleFunc("/clouds/{cloud}", testResponse)
	//retrieve all catalogs
	r.HandleFunc("/catalog", tweed.CatalogTweeds)
	//provision a service
	r.HandleFunc("/provision/{service}/{plan}", provisionFunction)
	//get an instance
	r.HandleFunc("/instances/{instance}", testResponse)
	//deprovision an instance
	r.HandleFunc("/deprovision/{instance}", deprovisionFunction)
	//bind an instance
	r.HandleFunc("/bind/{instance}/{binding}/{nowait}", bindFunction)
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", testResponse)
	//unbind an instance
	r.HandleFunc("/unbind/{instance}", unbindFunction)

	//start the server
	log.Fatal(http.ListenAndServe(":8081", r))
}
