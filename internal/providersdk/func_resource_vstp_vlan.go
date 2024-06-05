package providersdk

import "github.com/jeremmfr/terraform-provider-junos/internal/junos"

func checkVstpVlanExists(vlanID, routingInstance string, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	if routingInstance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols vstp vlan " + vlanID + junos.PipeDisplaySet)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols vstp vlan " + vlanID + junos.PipeDisplaySet)
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}
