package junos

import (
	"fmt"
	"strings"
)

func (sess *Session) CheckInterfaceExists(interFace string) (bool, error) {
	reply, err := sess.CommandXML(fmt.Sprintf(RPCGetInterfaceInformationInterfaceName, interFace))
	if err != nil {
		if strings.Contains(err.Error(), " not found\n") ||
			strings.HasSuffix(err.Error(), " not found") {
			return false, nil
		}

		return false, err
	}
	if strings.Contains(reply, " not found\n") {
		return false, nil
	}

	return true, nil
}
