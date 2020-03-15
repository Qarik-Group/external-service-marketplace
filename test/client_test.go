package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)

func TestCatalogNoAuth(t *testing.T) {
	var in util.CatalogCommand
	json.Marshal(in)
	r, _ := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", bytes.NewReader(nil))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Catalog did not return %v when sent a request with an empty BasicAuth but got %v", http.StatusUnauthorized, http.StatusOK)
		return
	}
}

func TestCatalogTestWithAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	//confirms that the status does not get denied when credentials are passed to Catalog.
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf("Catalog did not Correctly handle user that is passing in crednetials. Got code %v", status)
		return
	} else if status != http.StatusOK {
		t.Errorf("Got back a bad error form the tweed catalog: got %v want %v", status, http.StatusOK)
		return
	}
	//if body is totally empty we did something wrong
	if len(rr.Body.String()) <= 0 {
		t.Errorf("The body was empty from the response of the Catalog handler")
		return
	}

}

func TestBindNoAuth(t *testing.T) {
	id := "hi"
	bid := "hello"
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+id+"/bindings/"+bid, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Bind)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Not passing crednetials allowed the Bind Function to fire off not good")
		return
	} else if rr.Code == http.StatusBadRequest {
		t.Errorf("Error constructing your request double check that you made it right")
		return
	}
}
func TestBindAuth(t *testing.T) {
	var in util.BindCommand
	in.ID = "hello"
	in.NoWait = true
	in.ID = "hi"
	in.Args.ID = "what isthis no clue haha"
	body, _ := json.Marshal(&in)
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+in.Args.ID+"/bindings/"+in.ID, bytes.NewReader(body))
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Bind)
	handler.ServeHTTP(rr, r)
	if rr.Code == http.StatusUnauthorized {
		t.Errorf(rr.Body.String())
		return
	} else if rr.Code == http.StatusBadRequest {
		t.Errorf(rr.Body.String())
		return
	}
	if rr.Code != http.StatusOK {
		t.Errorf(rr.Body.String())
	}
}

func TestProvisionNoAuth(t *testing.T) {
	var prov util.ProvisionCommand
	body, _ := json.Marshal(&prov)
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+"", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Provision)
	handler.ServeHTTP(rr, r)
	if rr.Code == http.StatusOK {
		t.Errorf(rr.Body.String())
		return
	}
}
func TestProvisionWithAuth(t *testing.T) {
	var prov util.ProvisionCommand
	body, _ := json.Marshal(&prov)
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+"", bytes.NewReader(body))
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Provision)
	handler.ServeHTTP(rr, r)
	//check no auth should not be OK!
	if rr.Code != http.StatusBadRequest {
		t.Errorf(rr.Body.String())
		return
	}

}

func TestProvisionID(t *testing.T) {
	var prov util.ProvisionCommand
	prov.ID = "hello"
	body, _ := json.Marshal(&prov)
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+"id", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Provision)
	handler.ServeHTTP(rr, r)
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusBadRequest {
		t.Errorf(rr.Body.String())
	}
}

func TestProvisionIDServices(t *testing.T) {
	var prov util.ProvisionCommand
	prov.ID = "hello"
	services := []string{"hello", "goodbye"}
	prov.Args.ServicePlan = services
	body, _ := json.Marshal(&prov)
	r, _ := http.NewRequest("PUT", util.GetTweedUrl()+"/b/instances/"+"id", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Provision)
	handler.ServeHTTP(rr, r)
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		t.Errorf(rr.Body.String())
	}
}
