package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	Config *Config
	Bind   string
}

func test_response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}
func (api API) Run() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//retrieve clouds
	r.HandleFunc("/clouds", test_response)
	//retrieve a specific cloud
	r.HandleFunc("/clouds/{cloud}", test_response)
	//retrieve all catalogs
	r.HandleFunc("/catalog", test_response)
	//provision a service
	r.HandleFunc("/provision/{service}/{plan}", test_response)
	//get an instance
	r.HandleFunc("/instances/{instance}", test_response)
	//deprovision an instance
	r.HandleFunc("/deprovision/{instance}", test_response)
	//bind an instance
	r.HandleFunc("/bind/{instance}", test_response)
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", test_response)
	//unbind an instance
	r.HandleFunc("/unbind/{instance}", test_response)
	return r
}
