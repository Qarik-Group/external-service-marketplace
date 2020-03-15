package tweed

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/starkandwayne/external-service-marketplace/util"
	"github.com/tweedproject/tweed"
	"github.com/tweedproject/tweed/api"
)

type client struct {
	http     *http.Client
	url      string
	username string
	password string
}

func Connect(url, username, password string) *client {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return &client{
		url:      strings.TrimSuffix(url, "/"),
		username: username,
		password: password,
		http:     http.DefaultClient,
	}
}

func (c *client) do(req *http.Request, out interface{}) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	res, err := c.http.Do(req)
	if err != nil || out == nil {
		return res, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}

	var e api.ErrorResponse
	if res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 404 || res.StatusCode == 500 {
		if err = json.Unmarshal(b, &e); err != nil {
			return res, err
		}
		return res, e
	}

	return res, json.Unmarshal(b, &out)
}

func (c *client) request(method, path string, in interface{}) (*http.Request, error) {
	var body bytes.Buffer
	if in != nil {
		b, err := json.MarshalIndent(in, "", "  ")
		if err != nil {
			return nil, err
		}
		body.Write(b)
	}
	req, err := http.NewRequest(method, c.url+"/"+path, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if in != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (c *client) get(path string, out interface{}) error {
	req, err := c.request("GET", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) post(path string, in, out interface{}) error {
	req, err := c.request("POST", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) put(path string, in, out interface{}) error {
	req, err := c.request("PUT", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) delete(path string, out interface{}) error {
	req, err := c.request("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No username or password submited in your request"))
		return
	}
	c := Connect(util.GetTweedUrl(), username, password)
	var cat tweed.Catalog
	c.get("/b/catalog", &cat)
	util.JSON(cat)
	body, err := json.Marshal(cat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func CatalogTweeds(w http.ResponseWriter, r *http.Request) {}

func UnBind(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No username or password submited in your request"))
		return
	}
	c := Connect(util.GetTweedUrl(), username, password)
	var un api.UnbindResponse
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	var in util.UnbindCommand
	json.Unmarshal(body, &in)
	c.delete("/b/instances/:id/bindings/:bid", &un)
	util.JSON(un)
	data, err := json.Marshal(un)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func Bind(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No username or password submited in your request"))
		return
	}
	var bcmd util.BindCommand
	err := json.NewDecoder(r.Body).Decode(&bcmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("There was no struct passed in your bind request please make sure to at least pass an empty BindCommand struct"))
		return
	}
	c := Connect(util.GetTweedUrl(), username, password)
	var out api.BindResponse
	c.put(util.GetTweedUrl()+"/b/instances/"+bcmd.Args.ID+"/bindings/"+bcmd.ID, nil, &out)
	if out.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(out.Error))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out.OK))
}

func Provision(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could not extract username and password from the request"))
		return
	}
	c := Connect(util.GetTweedUrl(), username, password)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read in the body of your request"))
		return
	}
	var in util.ProvisionCommand
	json.Unmarshal(body, &in)
	if len(in.Args.ServicePlan) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("There was no Plan passed in with your Provision Request"))
		return
	}
	var out api.ProvisionResponse
	c.put("/b/instances/"+in.ID, in, &out)
	if out.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(out.OK))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out.Ref))
}

func Deprovision(w http.ResponseWriter, r *http.Request) {

}
func Purge(w http.ResponseWriter, r *http.Request) {

}
