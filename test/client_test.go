package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)

func TestCatalogNoAuth(t *testing.T) {
	r := httptest.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Catalog did not return %v when sent a request with an empty BasicAuth but got %v", http.StatusUnauthorized, http.StatusOK)
	}
}

func TestCatalogTestWithAuth(t *testing.T) {
	r := httptest.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf("Catalog did not Correctly handle user that is passing in crednetials. Got code %v", status)
	}
}
