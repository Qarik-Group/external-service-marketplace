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
		Service string `cli:"-s, --service" env:"ESM_SERVICE"`
		Plan    string `cli:"-p, --plan" env:"ESM_PLAN"`
	} `cli:"provision"`

	Deprovision struct {
		Instance string `cli:"-i, --instance" env:"ESM_INSTANCE"`
	}
}

func main() {

	var options Options

	env.Override(&options)

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

	// if command == "provision" {

	// 	var prov util.ProvisionCommand
	// 	s := make([]string, 1)
	// 	s[0] = "Redis"
	// 	prov.Args.ServicePlan = s
	// 	sp, _ := json.Marshal(prov)
	// 	r := bytes.NewReader(sp)
	// 	username := "tweed"
	// 	passwd := "tweed"
	// 	client := &http.Client{}
	// 	req, err := http.NewRequest("POST", "http://localhost:8081/provision", r)
	// 	req.SetBasicAuth(username, passwd)
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		fmt.Printf("Error: not sent")
	// 	}
	// 	//defer req.Body.Close()
	// 	//body, err := ioutil.ReadAll(req.Body)
	// 	if err != nil {
	// 		fmt.Printf("Error: body not read")
	// 	}
	// 	json.NewDecoder(resp.Body).Decode(&prov)
	// 	//json.Unmarshal(body, &cat)
	// 	util.JSON(prov)

	// 	//fmt.Printf("Provisioning... %s  %s", options.Provision.Service, options.Provision.Plan)
	// }
	// if command == "deprovision" {
	// 	var prov util.ProvisionCommand
	// 	s := make([]string, 1)
	// 	s[0] = "Redis"
	// 	prov.Args.ServicePlan = s
	// 	sp, _ := json.Marshal(prov)
	// 	r := bytes.NewReader(sp)
	// 	username := "tweed"
	// 	passwd := "tweed"
	// 	client := &http.Client{}
	// 	req, err := http.NewRequest("GET", "http://localhost:8081/provision", r)
	// 	req.SetBasicAuth(username, passwd)
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		fmt.Printf("Error: not sent")
	// 	}
	// 	//defer req.Body.Close()
	// 	//body, err := ioutil.ReadAll(req.Body)
	// 	if err != nil {
	// 		fmt.Printf("Error: body not read")
	// 	}
	// 	json.NewDecoder(resp.Body).Decode(&prov)
	// 	//json.Unmarshal(body, &cat)
	// 	util.JSON(prov)

	// }
}
