package util

type BindCommand struct {
	ID     string `long:"as" optional:"yes" description:"use given binding id otherwise use random"`
	NoWait bool   `long:"no-wait" description:"don't wait for the binding to be created"`
	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}
type UnbindCommand struct {
	NoWait bool `long:"no-wait" description:"don't wait for the binding to be created"`
	Args   struct {
		InstanceBinding []string `positional-arg-name:"instance/binding" required:"true"`
	} `positional-args:"yes"`
}
type CatalogCommand struct {
}

type ProvisionCommand struct {
	Service string            `json:"service"`
	Plan    string            `json:"plan"`
	Param   map[string]string `json:"params"`
}

type DeprovisionCommand struct {
	NoWait   bool   `long:"no-wait" description:"don't wait for the service be deprovisioned"`
	Instance string `json:"instance"`
}
