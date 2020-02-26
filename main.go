package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	http.ListenAndServe(":80", r)
}
