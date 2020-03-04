package tweed

type Client struct {
	url string
	username string
	password string
}

type Service struct {
	ID string
	Name string
	// ...

	Plans []Plan
}

type Plan struct {
	ID string
	Name string
	// ...
}

func Connect(url, username, password string) (*Client, error) {
	return nil, nil
}

func (c *Client) Catalog() ([]Service, error) {
	return nil, nil
}

func (c *Client) Provision(service, plan string, params map[string]interface{}) error {
	return nil
}

func (c *Client) Deprovision(instance string) error {
	return nil
}

func (c *Client) Bind(instance string) (string, error) {
	return "", nil
}

func (c *Client) Unbind(instance, binding string) error {
	return nil
}

func (c *Client) Purge(instance string) error {
	return nil
}
