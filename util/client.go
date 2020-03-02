package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
)

type client struct {
	http     *http.Client
	url      string
	username string
	password string
}

func Connect(url string, user string, pass string) *client {
	if !strings.HasPrefix(url, "https") {
		url = "https://" + url
	}
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
	if err != nil {
		return res, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read the body into a []byte in do() of the client")
	}
	var e api.ErrorResponse
	if res.StatusCode > 399 {
		if err = json.Unmarshal(b, &e); err != nil {
			return res, err
		}
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
	path = strings.TrimPrefix(path, "/")
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
func (c *client) GET(path string, in interface{}, out interface{}) (*http.Response, error) {
	req, err := c.requst("GET", path, in)
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
