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
	tweedRouter.HandleFunc("/catalog", util.Catalog)
	tweedRouter.HandleFunc("/instances", util.Instances)
	tweedRouter.HandleFunc("/instances/:id", util.InstancesId)
	tweedRouter.HandleFunc("/instances/:id/tasks", util.InstancesIdTasks)
	tweedRouter.HandleFunc("/instances/:id/tasks/:tid", util.InstancesIdTaskId)
	tweedRouter.HandleFunc("/instances/:id/tasks", util.InstancesIdTasks)
	return r
}
