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

	Config string `cli:"-c, --config" env:"ESM_CONFIG"`
	Listen string `cli:"-l, --listen" env:"ESM_LISTEN"`
}

func main() {
	var options Options
	options.Config = "cmd/esm/esmd.yml" //need to change this
	env.Override(&options)
	_, _, err := cli.Parse(&options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
	}

	if options.Listen == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! missing required --listen option}\n")
		os.Exit(1)
	}
	config, err := ReadConfig(options.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s: %s}\n", options.Config, err)
		os.Exit(1)
	}

	api := API{
		Config: config,
		Bind:   options.Listen,
	}

	fmt.Fprintf(os.Stderr, "running api server...\n")
	fmt.Fprintf(os.Stderr, "(listening on %s)\n", api.Bind)
	api.Run()
	fmt.Fprintf(os.Stderr, "api server exited...\n")
	os.Exit(1)
}
