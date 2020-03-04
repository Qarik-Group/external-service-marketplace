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
	r.HandleFunc("/register/{broker}", ism_client.Register)
	r.HandleFunc("/catalog", ism_client.Catalog)
	r.HandleFunc("/provision/{service}/{plan}", ism_client.Provision)
	r.HandleFunc("/deprovision/{instance}", ism_client.Deprovision)
	return r
}
