package main

import (
	"external-service-marketplace/api"
	"net/http"
)

func main() {
	tr := api.TweedAPI()
	http.ListenAndServe(":9000", tr)
}
