package providersdk

import "github.com/jeremmfr/terraform-provider-junos/internal/junos"

func checkRoutingInstanceExists(instance string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}
