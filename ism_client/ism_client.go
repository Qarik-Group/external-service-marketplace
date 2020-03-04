package ism_client

import "net/http"

func Provision(w http.ResponseWriter, r *http.Request) {
	//provision a specific tweed
	return
}

func Deprovision(w http.ResponseWriter, r *http.Request) {
	//deprovision a specific tweed
	return
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	//get the master catalog across all tweeds
	return
}

func Purge(w http.ResponseWriter, r *http.Request) {
	//purge a specific tweed instance
}
