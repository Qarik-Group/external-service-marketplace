package api

import (
	"encoding/json"
	"external-service-marketplace/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tweedproject/tweed"
	"github.com/tweedproject/tweed/api"
)

func TweedAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	tweedRouter := r.PathPrefix("/tweed").Subrouter()
	tweedRouter.HandleFunc("/catalog", catalog)
	tweedRouter.HandleFunc("/instances", instances)
	return r
}

func catalog(w http.ResponseWriter, r *http.Request) {
	c := util.Connect(util.TweedURL(), util.GetUserName(), util.GetPassword())
	var cat tweed.Catalog
	_, err := c.GET("/b/catalog", &cat)
	if err != nil {
		log.Fatal("Error getting response body in catalog from GET reques\n" + err.Error())
	}
	util.JSON(cat)
	b, err := json.Marshal(cat)
	if err != nil {
		log.Fatal("Failed to convert Catalog to json using the util.JSON")
	}
	w.WriteHeader(200)
	w.Write(b)
}

func instances(w http.ResponseWriter, r *http.Request) {
	c := util.Connect(util.TweedURL(), util.GetUserName(), util.GetPassword())
	var out []api.InstanceResponse
	_, err := c.GET("/b/instances", &out)
	if err != nil {
		log.Fatal("Error getting response body:\t instances from GET request to tweed \n\t" + err.Error())
	}
	util.JSON(out)
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal("\nError encoding the []api.InstanceResponse from the body using Marshall\n" + err.Error())
	}
	w.WriteHeader(200)
	w.Write(b)
	//	var slice = util.MakeBody(out)
}
