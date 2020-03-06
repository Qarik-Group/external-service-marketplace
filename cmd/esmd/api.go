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
	r.HandleFunc("/clouds", test_response)
	r.HandleFunc("/register/{broker}", test_response)
	r.HandleFunc("/catalog", test_response)
	r.HandleFunc("/provision/{service}/{plan}", test_response)
	r.HandleFunc("/deprovision/{instance}", test_response)
	return r
}
