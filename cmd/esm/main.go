package main

import (
	"encoding/json"
	"net/http"
	"os"

	fmt "github.com/jhunt/go-ansi"
	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/starkandwayne/external-service-marketplace/util"
	"github.com/tweedproject/tweed"
)

type Options struct {
	Help  bool `cli:"-h, --help"`
	Debug bool `cli:"-D, --debug" env:"ESM_DEBUG"`
	Trace bool `cli:"-T, --trace" env:"ESM_TRACE"`

	Catalog struct {
	} `cli:"catalog"`

	Provision struct { // sub commands to be entered here
		Service string `cli:"-s, --service" env:"ESM_SERVICE_PROV"`
		Plan    string `cli:"-p, --plan" env:"ESM_PLAN"`
		Prefix  string `cli:"-r, --prefix" env:"ESM_PREFIX_PROV"`
	} `cli:"provision"`

	Deprovision struct {
		Instance string `cli:"-i, --instance" env:"ESM_INSTANCE_DEPROV"`

		Prefix string `cli:"-r, --prefix" env:"ESM_PREFIX_DEPROV"`
	} `cli:"deprovision"`

	Bind struct {
		Prefix string `cli:"-r, --prefix" env:"ESM_PREFIX_BIND"`

		Instance string `cli:"-i, --instance" env:"ESM_INSTANCE_BIND"`
	} `cli:"bind"`
	Unbind struct {
		Instance string `cli:"-i, --instance" env:"ESM_INSTANCE_UNBIND"`
		Prefix   string `cli:"-r, --prefix" env:"ESM_PREFIX_UNBIND"`
		Binding  string `cli:"-k, --binding" env:"ESM_INST_ID"`
	} `cli:"unbind"`
}

func main() {

	var options Options
	// instance_bind := options.Bind.Instance
	// prefix := options.Prefix

	env.Override(&options)
	// fmt.Printf("this is main %s \n", options.Prefix)

	command, args, err := cli.Parse(&options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
	}

	if command == "catalog" {
		var cat tweed.Catalog
		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://localhost:8081/catalog", nil)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: not sent")
		}
		//defer req.Body.Close()
		//body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Error: body not read")
		}
		json.NewDecoder(resp.Body).Decode(&cat)
		//json.Unmarshal(body, &cat)
		util.JSON(cat)
		fmt.Printf("List services \n")
		fmt.Printf("running command @G{%s}...\n", command)
		fmt.Printf("with arguments @C{%v}...\n", args)
	}

	if command == "provision" {
		serv := options.Provision.Service
		plan := options.Provision.Plan
		prefix := options.Provision.Prefix

		// fmt.Printf(prefix)

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://localhost:8081/provision/"+prefix+"--"+serv+"/"+plan, nil)
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error: not sent")
		}
		if err != nil {
			fmt.Printf("Error: body not read")
		}

	}
	if command == "deprovision" {
		instance := options.Deprovision.Instance
		prefix := options.Deprovision.Prefix

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://localhost:8081/deprovision/"+prefix+"/"+instance, nil)
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error: not sent")
		}
		if err != nil {
			fmt.Printf("Error: body not read")
		}

	}
	if command == "bind" {
		instance := options.Bind.Instance
		prefix := options.Bind.Prefix

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://localhost:8081/bind/"+prefix+"/"+instance, nil)
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error: not sent")
		}
		if err != nil {
			fmt.Printf("Error: body not read")
		}

	}
	if command == "unbind" {
		binding := options.Unbind.Binding
		prefix := options.Unbind.Prefix

		instance := options.Unbind.Instance
		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://localhost:8081/unbind/"+prefix+"/"+instance+"/"+binding, nil)
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error: not sent")
		}
		if err != nil {
			fmt.Printf("Error: body not read")
		}

	}
}
