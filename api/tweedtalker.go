package api

import (
	"external-service-marketplace/util"
	"net/http"

	"github.com/gorilla/mux"
)

func TweedAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	tweedRouter := r.PathPrefix("/tweed").Subrouter()
	tweedRouter.HandleFunc("/catalog", catalog)
	return r
}

func catalog(w http.ResponseWriter, r *http.Request) {
	c := util.Connect(util.TweedURL(), util.GetUserName(), util.GetPassword())
	res, _ := c.GET("/b/catalog", nil, nil)
	w.WriteHeader(res.StatusCode)
	w.Write(util.ReadResponse(res))
}
