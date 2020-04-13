package test

import (
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)

/*func TestUnBind(t *testing.T) {
	var unbindCmd util.UnbindCommand
	ids := []string{"hi", "hello"}
	unbindCmd.Args.InstanceBinding = ids
	var conf *util.Config
	client := tweed.Connect(conf)
	res := client.UnBind("http://10.128.32.138:31666", unbindCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestUnBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestBind(t *testing.T) {
	var bindCmd util.BindCommand
	bindCmd.ID = "hi"
	bindCmd.Args.ID = "bye"
	var conf *util.Config
	client := tweed.Connect(conf)
	res := client.Bind("http://10.128.32.138:31666", bindCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestDeprovision(t *testing.T) {
	var deprovCmd util.DeprovisionCommand
	ids := []string{"hi", "hello"}
	deprovCmd.Args.InstanceIds = ids
	var conf *util.Config
	client := tweed.Connect(conf)
	res := client.DeProvision("http://10.128.32.138:31666", deprovCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestDeprovision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
} */

func TestProvision(t *testing.T) {
	config, _ := util.ReadConfig("cmd/esm/esmd.yml")
	var provCmd util.ProvisionCommand
	ids := []string{"redis" + "/" + "shared"}
	provCmd.Args.ServicePlan = ids
	provCmd.ID = "i-37hdg475gdyth2"

	client := tweed.Connect(config)
	res := client.Provision("http://10.128.32.138:32632", provCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestProvision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}
