package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/tweedproject/tweed"
	"github.com/tweedproject/tweed/api"
)

type client struct {
	http     *http.Client
	url      string
	username string
	password string
}

func Connect(url string, user string, pass string) *client {
	return &client{
		http:     http.DefaultClient,
		url:      url,
		username: user,
		password: pass,
	}
}

func (c *client) do(req *http.Request, out interface{}) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	res, err := c.http.Do(req)
	if res != nil {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DEBUG: @W{unable to dump response:} @R{%s}\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n\n", string(b))
		}
	}
	if err != nil || out == nil {
		return res, err
	}
	defer res.Body.Close()
	if err != nil {
		return res, err
	}
	b, err := ioutil.ReadAll(res.Body)

	var e api.ErrorResponse
	if res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 404 || res.StatusCode == 500 {
		if err = json.Unmarshal(b, &e); err != nil {
			return res, err
		}
		return res, e
	}

	return res, json.Unmarshal(b, &out)
}

func (c *client) requst(method string, path string, in interface{}) (*http.Request, error) {
	var body bytes.Buffer
	if in != nil {
		b, err := json.MarshalIndent(in, "", " ")
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
		req.Header.Add("Content-Type", "appliation/json")
	}
	return req, nil
}

//The in is a TweedRequest struct and out is a Tweed Response struct
func (c *client) GET(path string, out interface{}) (*http.Response, error) {
	req, err := c.requst("GET", path, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.do(req, out)
	return res, err
}

func (c *client) POST(path string, in interface{}, out interface{}) error {
	req, err := c.requst("POST", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) PUT(path string, in interface{}, out interface{}) error {
	req, err := c.requst("PUT", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) DELETE(path string, in interface{}, out interface{}) error {
	req, err := c.requst("DELETE", path, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	c := Connect(TweedURL(), GetUserName(), GetPassword())
	var cat tweed.Catalog
	_, err := c.GET("/b/catalog", &cat)
	if err != nil {
		log.Fatal("Error getting response body in catalog from GET reques\n" + err.Error())
	}
	JSON(cat)
	b, err := json.Marshal(cat)
	if err != nil {
		log.Fatal("Failed to convert Catalog to json using the JSON")
	}
	w.WriteHeader(200)
	w.Write(b)
}

func Instances(w http.ResponseWriter, r *http.Request) {
	c := Connect(TweedURL(), GetUserName(), GetPassword())
	var out []api.InstanceResponse
	_, err := c.GET("/b/instances", &out)
	if err != nil {
		log.Fatal("Error getting response body:\t instances from GET request to tweed \n\t" + err.Error())
	}
	JSON(out)
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal("\nError encoding the []api.InstanceResponse from the body using Marshall\n" + err.Error())
	}
	w.WriteHeader(200)
	w.Write(b)
	//	var slice = MakeBody(out)
}

func InstancesId(w http.ResponseWriter, r *http.Request) {
	c := Connect(TweedURL(), GetUserName(), GetPassword())
	var out api.InstanceResponse
	param := r.URL.Query().Get("instance")
	_, err := c.GET("/b/instances/"+param, &out)
	if err != nil {
		log.Fatal("Error getting response body:\t instancesId from GET request to tweed \n\t" + err.Error())
	}
	JSON(out)
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal("\nError encoding the []api.InstanceResponse from the body using Marshall\n" + err.Error())
	}
	w.WriteHeader(200)
	w.Write(b)
}

func InstancesIdTasks(w http.ResponseWriter, r *http.Request) {
	c := Connect(TweedURL(), GetUserName(), GetPassword())
	var out api.TaskResponse
	param := r.URL.Query().Get("instance")
	_, err := c.GET("/b/instances/"+param, &out)
	if err != nil {
		log.Fatal("Error getting response body:\t instancesId from GET request to tweed \n\t" + err.Error())
	}
	JSON(out)
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal("\nError encoding the []api.InstanceResponseTasks from the body using Marshall\n" + err.Error())
	}
	w.WriteHeader(200)
	w.Write(b)
}

//Have not tested this as there are no active tasks currently to test the endpoint
func InstancesIdTaskId(w http.ResponseWriter, r *http.Request) {
	c := Connect(TweedURL(), GetUserName(), GetPassword())
	var out api.TaskResponse
	intanceId := r.URL.Query().Get("instance")
	_, err := c.GET("/b/instances/"+intanceId, &out)
	if err != nil {
		log.Fatal("Error getting response body:\t instancesId from GET request to tweed \n\t" + err.Error())
	}
	JSON(out)
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal("\nError encoding the []api.InstanceResponseTasks from the body using Marshall\n" + err.Error())
	}
	w.WriteHeader(200)
	w.Write(b)
}
