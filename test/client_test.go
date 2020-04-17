package test

import (
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)
<<<<<<< HEAD
=======

<<<<<<< HEAD
>>>>>>> multiple-tweed-routes
func TestUnBind(t *testing.T) {
	var unbindCmd util.UnbindCommand
	ids := []string{"hi", "hello"}
	unbindCmd.Args.InstanceBinding = ids
	var conf tweed.Config
	client := tweed.Connect(conf)
	res := client.UnBind("http://10.128.32.138:31666", unbindCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestUnBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

=======
>>>>>>> multiple-tweed-routes
func TestBind(t *testing.T) {
	var conf *util.Config
	client := tweed.Connect(conf)
	instance := "i-82421ffd2c6522" //change this to a different instance
	res := client.Bind("http://10.128.32.138:32632", instance)
	if res.Error != "" {
		t.Errorf("Error in TestBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}
func TestUnBind(t *testing.T) {
	//var unbindCmd util.UnbindCommand
	//ids := []string{"hi", "hello"}
	//unbindCmd.Args.InstanceBinding = ids
	var conf *util.Config
	client := tweed.Connect(conf)
	instance := "i-82421ffd2c6522" //replace this instance if not provisioned
	binding := "b-bb7a954fdc0680"  //replace this with a new binding
	res := client.UnBind("http://10.128.32.138:32632", instance, binding)
	if res.Error != "" {
		t.Errorf("Error in TestUnBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestProvision(t *testing.T) {
	config, _ := util.ReadConfig("cmd/esm/esmd.yml")
	var provCmd util.ProvisionCommand
	//ids := []string{"redis" + "/" + "shared"}
	provCmd.Service = "redis"
	provCmd.Plan = "shared"
	client := tweed.Connect(config)
	res := client.Provision("http://10.128.32.138:32632", provCmd)
	if res.Error != "" {
		t.Errorf("Error in TestProvision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}
func TestDeprovision(t *testing.T) {
	config, _ := util.ReadConfig("cmd/esm/esmd.yml")
	instance := "i-9d62e7231649bf" //replace this with something relevant
	client := tweed.Connect(config)
	res := client.DeProvision("http://10.128.32.138:32632", instance)
	if res.Error != "" {
		t.Errorf("Error in TestDeprovision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}
