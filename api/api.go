package api

import (
	"external-service-marketplace/ism_client"
	"net/http"

	"github.com/gorilla/mux"
)

func API() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	tweedRouter := r.PathPrefix("/tweed").Subrouter()
	tweedRouter.HandleFunc("/catalog", ism_client.Connect)
	tweedRouter.HandleFunc("/provision/{service}/{plan}", ism_client.Connect)
	tweedRouter.HandleFunc("/deprovision/{service}/{plan}", ism_client.Connect)
	return r
}
