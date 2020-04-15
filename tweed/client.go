package tweed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/starkandwayne/external-service-marketplace/util"
	"github.com/tweedproject/tweed"
	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
)

type Config struct {
	ServiceBrokers []struct {
		Prefix     string `yaml:"prefix"`
		URL        string `yaml:"url"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		SkipVerify bool   `yaml:"skip-verify"`
	} `yaml:"service-brokers"`

	Clouds []struct {
		ID   string `yaml:"id"`
		Name string `yaml:"name"`
		Type string `yaml:"type"`

		// figure out: how to specify creds for CF / K8s
	} `yaml:"clouds"`
}

type client struct {
	http      *http.Client
	config    *util.Config
	connected bool
}

func Connect(config *util.Config) *client {
	return &client{
		http:      http.DefaultClient,
		config:    config,
		connected: true,
	}
}

func (c *client) do(req *http.Request, out interface{}) (*http.Response, error) {
	req.SetBasicAuth("tweed", "tweed")
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

func (c *client) request(method, url string, in interface{}) (*http.Request, error) {
	var body bytes.Buffer
	if in != nil {
		b, err := json.MarshalIndent(in, "", "  ")
		if err != nil {
			return nil, err
		}
		body.Write(b)
	}
	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if in != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (c *client) get(url, route string, out interface{}) error {
	req, err := c.request("GET", url+route, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) post(url, route string, in, out interface{}) error {
	req, err := c.request("POST", url+route, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) put(url, route string, in, out interface{}) error {
	req, err := c.request("PUT", url+route, in)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *client) delete(url, route string, out interface{}) error {
	req, err := c.request("DELETE", url+route, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}
func (c *client) status() bool {
	if c.connected {
		return true
	}
	return false
}
func (c *client) SingleCatalog(url string) tweed.Catalog {
	var cat tweed.Catalog
	c.get(url, "/b/catalog", &cat)
	util.JSON(cat)
	return cat
}
func (c *client) Catalog() []tweed.Catalog {
	var cat tweed.Catalog
	broker := c.config.ServiceBrokers[0]
	var cats []tweed.Catalog
	for i := 0; i < len(c.config.ServiceBrokers); i++ {
		broker = c.config.ServiceBrokers[i]
		c.get(broker.URL, "/b/catalog", &cat)
		cats = append(cats, cat)
		util.JSON(cat)
	}
	return cats
}

func (c *client) UnBind(url string, unbindCmd util.UnbindCommand) api.UnbindResponse {
	var un api.UnbindResponse
	c.delete(url, "/b/instances/"+unbindCmd.Args.InstanceBinding[0]+"/bindings/"+unbindCmd.Args.InstanceBinding[1], &un)
	return un
}

func (c *client) Bind(url string, bindCmd util.BindCommand) api.BindResponse {
	var out api.BindResponse
	c.put(url, "/b/instances/"+bindCmd.Args.ID+"/bindings/"+bindCmd.ID, nil, &out)
	return out
}

func (c *client) Provision(url string, provCmd util.ProvisionCommand) api.ProvisionResponse {
	var out api.ProvisionResponse
	id := random.ID("i")
	fmt.Println("From Tweed Provision")
	util.JSON(provCmd)
	c.put(url, "/b/instances/"+id, provCmd, &out)
	return out
}

func (c *client) DeProvision(url string, deprovCmd util.DeprovisionCommand) api.DeprovisionResponse {
	var out api.DeprovisionResponse
	c.delete(url, "/b/instances/"+deprovCmd.Args.InstanceIds[0], &out)
	return out
}
