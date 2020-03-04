package main

import (
	"external-service-marketplace/api"
	"net/http"
)

func main() {
	r := api.API()
	http.ListenAndServe(":9000", r)
}
