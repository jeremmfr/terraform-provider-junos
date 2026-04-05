package junos

import (
	"context"
	"fmt"
	"strings"
)

func (sess *Session) CheckInterfaceExists(ctx context.Context, interFace string) (bool, error) {
	reply, err := sess.CommandXML(ctx, fmt.Sprintf(RPCGetInterfaceInformationInterfaceName, interFace))
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
