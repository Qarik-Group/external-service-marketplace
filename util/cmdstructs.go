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
	ID     string   `long:"as" optional:"yes" description:"use given service id otherwise use random"`
	NoWait bool     `long:"no-wait" description:"don't wait for the service to be created"`
	Params []string `short:"P" optional:"yes" long:"params" description:"params passed to the service"`
	Args   struct {
		ServicePlan []string `positional-arg-name:"service/plan" required:"true"`
	} `positional-args:"yes"`
}

type DeprovisionCommand struct {
	NoWait bool `long:"no-wait" description:"don't wait for the service be deprovisioned"`
	Args   struct {
		InstanceIds []string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}
