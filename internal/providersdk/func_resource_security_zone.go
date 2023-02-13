package providersdk

import (
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

func checkSecurityZonesExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security zones security-zone " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}
