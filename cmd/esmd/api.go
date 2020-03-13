package main

import (
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

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}

func (api API) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//retrieve clouds
	r.HandleFunc("/clouds", testResponse)
	//retrieve a specific cloud
	r.HandleFunc("/clouds/{cloud}", testResponse)
	//retrieve all catalogs
	r.HandleFunc("/catalog", tweed.CatalogTweeds)
	//provision a service
	r.HandleFunc("/provision/{service}/{plan}", tweed.Provision)
	//get an instance
	r.HandleFunc("/instances/{instance}", testResponse)
	//deprovision an instance
	r.HandleFunc("/deprovision/{instance}", tweed.Deprovision)
	//bind an instance
	r.HandleFunc("/bind/{instance}", tweed.Bind)
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", testResponse)
	//unbind an instance
	r.HandleFunc("/unbind/{instance}", tweed.UnBind)

	//start the server
	log.Fatal(http.ListenAndServe(":8081", r))
}
