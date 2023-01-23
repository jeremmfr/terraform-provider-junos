package junos

import (
	"context"
	"fmt"
	"net"
	"strconv"
)

func (clt *Client) StartNewSession(ctx context.Context) (*Session, error) {
	var auth sshAuthMethod
	auth.Username = clt.junosUserName
	auth.Ciphers = clt.junosSSHCiphers
	if clt.junosSSHKeyPEM != "" {
		auth.PrivateKeyPEM = clt.junosSSHKeyPEM
		if clt.junosSSHKeyPass != "" {
			auth.Passphrase = clt.junosSSHKeyPass
		}
	}
	if clt.junosSSHKeyFile != "" {
		auth.PrivateKeyFile = clt.junosSSHKeyFile
		if clt.junosSSHKeyPass != "" {
			auth.Passphrase = clt.junosSSHKeyPass
		}
	}
	if clt.junosPassword != "" {
		auth.Password = clt.junosPassword
	}
	auth.Timeout = clt.junosSSHTimeoutToEstab
	junSess, err := netconfNewSession(
		ctx,
		net.JoinHostPort(clt.junosIP, strconv.Itoa(clt.junosPort)),
		&auth,
		&openSSHOptions{
			Retry:   clt.junosSSHRetryToEstab,
			Timeout: clt.junosSSHTimeoutToEstab,
		},
	)
	if err != nil {
		return nil, err
	}
	if junSess.SystemInformation.HardwareModel == "" {
		return junSess, fmt.Errorf("can't read model of device with <get-system-information/> netconf command")
	}
	clt.logFile("[startNewSession] started")

	return junSess, nil
}

func (clt *Client) CloseSession(sess *Session) {
	err := sess.close(clt.sleepSSHClosed)
	if err != nil {
		clt.logFile(fmt.Sprintf("[closeSession] err: %q", err))
	} else {
		clt.logFile("[closeSession] closed")
	}
}
