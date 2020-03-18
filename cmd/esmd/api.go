package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
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

var config Config

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}
func catalogFunction(w http.ResponseWriter, r *http.Request) {
	//get config service brokers
	//loop through them
	//add results to response

}
func bindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	binding := vars["binding"]
	//nowait := vars["nowait"]

	var bind bind
	bind.ID = instance
	bind.Binding = binding
	bindCmd := util.BindCommand

	tweed.Bind(config.ServiceBrokers[0].Username, config.ServiceBrokers[0].Password, config.ServiceBrokers[0].URL, bindCmd)
	fmt.Fprint(w, "Bound Service")

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
	unbindCmd := util.UnbindCommand

	tweed.UnBind(config.ServiceBrokers[0].Username, config.ServiceBrokers[0].Password, config.ServiceBrokers[0].URL, unbindCmd)
	fmt.Fprint(w, "Unbound Instance")

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
	provCmd := util.ProvisionCommand

	tweed.Provision(config.ServiceBrokers[0].Username, config.ServiceBrokers[0].Password, config.ServiceBrokers[0].URL, provCmd)
	fmt.Fprint(w, "Provisioned Service")

}

func deprovisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	//nowait := vars["nowait"]

	var deprovision deprovision
	deprovision.ID = instance
	deprovCmd := util.DeprovisionCommand

	tweed.DeProvision(config.ServiceBrokers[0].Username, config.ServiceBrokers[0].Password, config.ServiceBrokers[0].URL, deprovCmd)
	fmt.Fprint(w, "Deprovisioned Service")

}

func (api API) Run() {
	config = *api.Config
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
	r.HandleFunc("/catalog", testResponse)
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
