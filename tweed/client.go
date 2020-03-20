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
	req, err := http.NewRequest(method, c.url+path, &body)
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

func Catalog(username, password, url string) tweed.Catalog {
	c := Connect(url, username, password)
	var cat tweed.Catalog
	c.get("/b/catalog", &cat)
	util.JSON(cat)
	return cat
}

func UnBind(username, password, url string, unbindCmd util.UnbindCommand) api.UnbindResponse {
	c := Connect(url, username, password)
	var un api.UnbindResponse
	c.delete("/b/instances/"+unbindCmd.Args.InstanceBinding[0]+"/bindings/"+unbindCmd.Args.InstanceBinding[1], &un)
	return un
}

func Bind(username, password, url string, bindCmd util.BindCommand) api.BindResponse {
	c := Connect(url, username, password)
	var out api.BindResponse
	c.put("/b/instances/"+bindCmd.Args.ID+"/bindings/"+bindCmd.ID, nil, &out)
	return out
}

func Provision(username, password, url string, provCmd util.ProvisionCommand) api.ProvisionResponse {
	c := Connect(url, username, password)
	var out api.ProvisionResponse
	c.put("/b/instances/"+provCmd.ID, provCmd, &out)
	return out
}

func DeProvision(username, password, url string, deprovCmd util.DeprovisionCommand) api.DeprovisionResponse {
	c := Connect(url, username, password)
	var out api.DeprovisionResponse
	c.delete("/b/instances/"+deprovCmd.Args.InstanceIds[0], &out)
	return out
}
