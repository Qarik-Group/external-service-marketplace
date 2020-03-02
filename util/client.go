package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

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
