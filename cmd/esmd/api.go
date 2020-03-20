package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var config Config
var url string

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}
func catalogFunction(w http.ResponseWriter, r *http.Request) {
	//get config service brokers
	//loop through them
	//add results to response
	username, password, _ := r.BasicAuth()

	res := tweed.Catalog(username, password, url)
	body, _ := json.Marshal(res)

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}
func bindFunction(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//instance := vars["instance"]
	//binding := vars["binding"]
	//nowait := vars["nowait"]
	/*username := config.ServiceBrokers[0].Username
	password := config.ServiceBrokers[0].Password
	url := config.ServiceBrokers[0].URL */
	username, password, _ := r.BasicAuth()

	var bindCmd util.BindCommand
	err := json.NewDecoder(r.Body).Decode(&bindCmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}

	res := tweed.Bind(username, password, url, bindCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res.Error))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func unbindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	binding := vars["binding"]
	//nowait := vars["nowait"]

	/*username := config.ServiceBrokers[0].Username
	password := config.ServiceBrokers[0].Password
	url := config.ServiceBrokers[0].URL*/
	username, password, _ := r.BasicAuth()

	instancebinding := make([]string, 2)
	instancebinding[0] = instance
	instancebinding[1] = binding

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}

	var unbindCmd util.UnbindCommand
	err = json.Unmarshal(body, &unbindCmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Did you use the UnbindCommand struct in util directory when you created your request. Please use that format"))
		return
	}

	//unbindCmd.Args.InstanceBinding = instancebinding

	res := tweed.UnBind(username, password, url, unbindCmd)
	util.JSON(res)
	data, _ := json.Marshal(res)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func provisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//instance := vars["instance"]
	service := vars["service"]
	plan := vars["plan"]
	//nowait := vars["nowait"]
	/*username := config.ServiceBrokers[0].Username
	password := config.ServiceBrokers[0].Password
	url := config.ServiceBrokers[0].URL*/
	username, password, _ := r.BasicAuth()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}
	var provCmd util.ProvisionCommand
	s := make([]string, 2)
	s[0] = service
	s[1] = plan
	json.Unmarshal(body, &provCmd)

	//provCmd.Args.ServicePlan = s

	res := tweed.Provision(username, password, url, provCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error In Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func deprovisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	//nowait := vars["nowait"]
	/*username := config.ServiceBrokers[0].Username
	password := config.ServiceBrokers[0].Password
	url := config.ServiceBrokers[0].URL*/
	username, password, _ := r.BasicAuth()

	s := make([]string, 1)
	s[0] = instance

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}
	var deprovCmd util.DeprovisionCommand
	//deprovCmd.Args.InstanceIds = instance
	json.Unmarshal(body, &deprovCmd)

	res := tweed.DeProvision(username, password, url, deprovCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error In Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func (api API) Run() {
	//config = *api.Config
	url = "http://10.128.32.138:31666"
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
