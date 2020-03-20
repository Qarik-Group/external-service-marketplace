package test

import (
	"testing"

	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"
)
func TestUnBind(t *testing.T) {
	var unbindCmd util.UnbindCommand
	ids := []string{"hi", "hello"}
	unbindCmd.Args.InstanceBinding = ids
	res := tweed.UnBind("tweed", "tweed", util.GetTweedUrl(), unbindCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestUnBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestBind(t *testing.T) {
	var bindCmd util.BindCommand
	bindCmd.ID = "hi"
	bindCmd.Args.ID = "bye"
	res := tweed.Bind("tweed", "tweed", "http://10.128.32.138:31666", bindCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestBind()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestDeprovision(t *testing.T) {
	var deprovCmd util.DeprovisionCommand
	ids := []string{"hi", "hello"}
	deprovCmd.Args.InstanceIds = ids
	res := tweed.DeProvision("tweed", "tweed", "http://10.128.32.138:31666", deprovCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestDeprovision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}

func TestProvision(t *testing.T) {
	var provCmd util.ProvisionCommand
	ids := []string{"hi"}
	provCmd.Args.ServicePlan = ids
	provCmd.ID = "hello"
	res := tweed.Provision("tweed", "tweed", "http://10.128.32.138:31666", provCmd)
	if res.Error == "" && res.OK == "" {
		t.Errorf("Error in TestProvision()\n" + res.Error + "\n res: \n" + res.Ref)
	}
}
