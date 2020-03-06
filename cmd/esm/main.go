package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"
	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
)

type Options struct {
	Help  bool `cli:"-h, --help"`
	Debug bool `cli:"-D, --debug" env:"ESM_DEBUG"`
	Trace bool `cli:"-T, --trace" env:"ESM_TRACE"`

	Catalog struct {
	} `cli:"catalog"`

	Provision struct {
		Service string `cli:"-s, --service" env:"ESM_SERVICE"`
		Plan    string `cli:"-p, --plan" env:"ESM_PLAN"`
	} `cli:"provision"`
}

func main() {
	var options Options
	env.Override(&options)
	options.Provision.Service = "Subcommand --service to be entered here\n"
	options.Provision.Plan = "Subcommand --plan to be entered here \n "
	command, args, err := cli.Parse(&options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
	}

	if command == "catalog" {
		fmt.Printf("List services \n")
		fmt.Printf("running command @G{%s}...\n", command)
		fmt.Printf("with arguments @C{%v}...\n", args)
}
	if command == "provision"{
		fmt.Printf("Provisioning... %s  %s", options.Provision.Service, options.Provision.Plan)
	}
}
