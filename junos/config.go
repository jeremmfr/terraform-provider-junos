package junos

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// configProvider.
type configProvider struct {
	junosPort                int
	junosCmdSleepShort       int
	junosCmdSleepLock        int
	junosSSHSleepClosed      int
	junosIP                  string
	junosUserName            string
	junosPassword            string
	junosSSHAgent            string
	junosSSHKeyPEM           string
	junosSSHKeyFile          string
	junosKeyPass             string
	junosGroupIntDel         string
	junosFilePermission      string
	junosDebugNetconfLogPath string
	junosFakeCreateSetFile   string
}

// prepareSession : prepare information to connect to Junos Device and more.
func (c *configProvider) prepareSession() (*Session, diag.Diagnostics) {
	sess := &Session{
		junosIP:             c.junosIP,
		junosPort:           c.junosPort,
		junosUserName:       c.junosUserName,
		junosPassword:       c.junosPassword,
		junosSSHAgent:       c.junosSSHAgent,         
		junosSSHKeyPEM:      c.junosSSHKeyPEM,
		junosKeyPass:        c.junosKeyPass,
		junosGroupIntDel:    c.junosGroupIntDel,
		junosSleepLock:      c.junosCmdSleepLock,
		junosSleepShort:     c.junosCmdSleepShort,
		junosSleepSSHClosed: c.junosSSHSleepClosed,
	}
	// junosSSHKeyFile
	sshKeyFile := c.junosSSHKeyFile
	if err := replaceTildeToHomeDir(&sshKeyFile); err != nil {
		return sess, diag.FromErr(err)
	}
	sess.junosSSHKeyFile = sshKeyFile

	// junosFilePermission
	filePermission, err := strconv.ParseInt(c.junosFilePermission, 8, 64)
	if err != nil {
		return sess, diag.FromErr(fmt.Errorf("failed to convert value from '%s' to int64 : %w",
			c.junosFilePermission, err))
	}
	sess.junosFilePermission = filePermission

	// junosLogFile
	junosLogFile := c.junosDebugNetconfLogPath
	if err := replaceTildeToHomeDir(&junosLogFile); err != nil {
		return sess, diag.FromErr(err)
	}
	sess.junosLogFile = junosLogFile

	// junosFakeCreateSetFile
	junosFakeCreateSetFile := c.junosFakeCreateSetFile
	if err := replaceTildeToHomeDir(&junosFakeCreateSetFile); err != nil {
		return sess, diag.FromErr(err)
	}
	sess.junosFakeCreateSetFile = junosFakeCreateSetFile

	return sess, nil
}
