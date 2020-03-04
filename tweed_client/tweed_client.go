package tweed_client

import (
	"net/http"
)

type tweed_client struct {
	http     *http.Client
	url      string
	username string
	password string
}

func TweedConnect(url string, user string, pass string) *tweed_client {
	return &tweed_client{
		http:     http.DefaultClient,
		url:      url,
		username: user,
		password: pass,
	}
}

func (c *tweed_client) provision() {
	//provision a service on this tweed
}

func (c *tweed_client) deprovision() {
	//deprovision a tweed on this tweed
}

func (c *tweed_client) catalog() {
	//list catalog for this tweed
}

func (c *tweed_client) purge() {
	//purge a service for this tweed
}
