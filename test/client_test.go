package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)

func TestCatalogNoAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Catalog)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Catalog did not return %v when sent a request with an empty BasicAuth but got %v", http.StatusUnauthorized, http.StatusOK)
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
	} else if status != http.StatusOK {
		t.Errorf("Got back a bad error form the tweed catalog: got %v want %v", status, http.StatusOK)
	}
	//if body is totally empty we did something wrong
	if len(rr.Body.String()) <= 0 {
		t.Errorf("The body was empty from the response of the Catalog handler")
	}

}
func TestBindNoAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", util.GetTweedUrl()+"/b/catalog", nil)
	r.Header.Add("Content-Type", "application/json")
	handler := http.HandlerFunc(tweed.Bind)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Catalog did not return %v when sent a request with an empty BasicAuth but got %v", http.StatusUnauthorized, http.StatusOK)
	}
}
func TestBindAuthNoParam(t *testing.T) {
	r, _ := http.NewRequest("GET", util.GetTweedUrl()+"/b/instances", nil)
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	handler := http.HandlerFunc(tweed.Bind)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Passing no parameter for Bind still worked Not Godd check the client.go")
	}
}

func TestBindNoBody(t *testing.T) {
	r, err := http.NewRequest("GET", util.GetTweedUrl()+"/b/instances", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.URL.Query().Add("instance", "i-f5d7a4d99a0f49/tasks")
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Bind)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusBadRequest {
		out := "There was an error with the Bind() in the client---> code" + strconv.FormatInt(int64(status), 10)
		t.Errorf(out)
	}
}

func TestBindWithBody(t *testing.T) {
	cmd := tweed.BindCommand{
		ID:     "",
		NoWait: false,
		Args: struct {
			ID string "positional-arg-name:\"instance\" required:\"true\""
		}{
			ID: "i-f5d7a4d99a0f49/tasks",
		},
	}
	var body bytes.Buffer
	b, _ := json.MarshalIndent(cmd, "", "  ")
	body.Write(b)
	r, err := http.NewRequest("GET", util.GetTweedUrl()+"/b/instances", &body)
	if err != nil {
		t.Fatal(err)
	}
	r.URL.Query().Add("instance", "i-f5d7a4d99a0f49")
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(util.GetUsername(), util.GetPassword())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweed.Bind)
	handler.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		out := "There was an error with the Bind() in the client---> code: " + strconv.FormatInt(int64(status), 10)
		t.Errorf(out)
	}
}
