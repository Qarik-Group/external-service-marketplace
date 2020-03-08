package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)

func TestCatalogNoAuth(t *testing.T) {
	r, err := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Catalog did not return %v when sent a request with an empty BasicAuth but got %v", http.StatusUnauthorized, http.StatusOK)
	}
}

func TestCatalogTestWithAuth(t *testing.T) {
	r, err := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	//confirms that the status does not get denied when credentials are passed to Catalog.
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf("Catalog did not Correctly handle user that is passing in crednetials. Got code %v", status)
	} else if status != http.StatusOK {
		t.Errorf("Got back a bad error form the tweed catalog: got %v want %v", status, http.StatusOK)
	}
	//if body is totally empty we did something wrong
	if len(rr.Body.String()) <= 0 {
		t.Errorf("The body was empty from the response of the Catalog handler")
	}
}
