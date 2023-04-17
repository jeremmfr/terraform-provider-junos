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
	sess, err := netconfNewSession(
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
	sess.logFile = func(message string) {
		message = "[" + sess.localAddress + "->" + sess.remoteAddress + "]" + message
		clt.logFile(message)
	}
	sess.sleepLock = clt.sleepLock
	sess.sleepShort = clt.sleepShort
	sess.sleepSSHClosed = clt.sleepSSHClosed
	if clt.fakeCreateSetFile != "" {
		sess.fakeSetFile = clt.appendFakeCreateSetFile
	}
	if sess.SystemInformation.HardwareModel == "" {
		return sess, fmt.Errorf("can't read model of device with <get-system-information/> netconf command")
	}
	sess.logFile("[StartNewSession] session opened")

	return sess, nil
}

func (clt *Client) NewSessionWithoutNetconf(_ context.Context) *Session {
	sess := Session{
		logFile: clt.logFile,
	}
	if clt.fakeCreateSetFile != "" {
		sess.fakeSetFile = clt.appendFakeCreateSetFile
	}

	return &sess
}
