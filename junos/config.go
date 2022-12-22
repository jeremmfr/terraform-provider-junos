package junos

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// configProvider.
type configProvider struct {
	junosFakeUpdateAlso      bool
	junosFakeDeleteAlso      bool
	junosPort                int
	junosCmdSleepShort       int
	junosCmdSleepLock        int
	junosSSHSleepClosed      int
	junosSSHTimeoutToEstab   int
	junosSSHRetryToEstab     int
	junosIP                  string
	junosUserName            string
	junosPassword            string
	junosSSHKeyPEM           string
	junosSSHKeyFile          string
	junosKeyPass             string
	junosGroupIntDel         string
	junosFilePermission      string
	junosDebugNetconfLogPath string
	junosFakeCreateSetFile   string
	junosSSHCiphers          []string
}

// newClient : prepare information to connect to Junos Device and more.
func (c *configProvider) newClient() (*Client, diag.Diagnostics) {
	clt := &Client{
		junosIP:                c.junosIP,
		junosPort:              c.junosPort,
		junosUserName:          c.junosUserName,
		junosPassword:          c.junosPassword,
		junosSSHKeyPEM:         c.junosSSHKeyPEM,
		junosSSHKeyPass:        c.junosKeyPass,
		groupIntDel:            c.junosGroupIntDel,
		sleepLock:              c.junosCmdSleepLock,
		sleepShort:             c.junosCmdSleepShort,
		sleepSSHClosed:         c.junosSSHSleepClosed,
		junosSSHCiphers:        c.junosSSHCiphers,
		junosSSHTimeoutToEstab: c.junosSSHTimeoutToEstab,
		junosSSHRetryToEstab:   c.junosSSHRetryToEstab,
		fakeUpdateAlso:         c.junosFakeUpdateAlso,
		fakeDeleteAlso:         c.junosFakeDeleteAlso,
	}
	// junosSSHKeyFile
	sshKeyFile := c.junosSSHKeyFile
	if err := replaceTildeToHomeDir(&sshKeyFile); err != nil {
		return clt, diag.FromErr(err)
	}
	clt.junosSSHKeyFile = sshKeyFile

	// junosFilePermission
	filePermission, err := strconv.ParseInt(c.junosFilePermission, 8, 64)
	if err != nil {
		return clt, diag.FromErr(fmt.Errorf("failed to convert value from '%s' to int64: %w",
			c.junosFilePermission, err))
	}
	clt.filePermission = filePermission

	// junosLogFile
	junosLogFile := c.junosDebugNetconfLogPath
	if err := replaceTildeToHomeDir(&junosLogFile); err != nil {
		return clt, diag.FromErr(err)
	}
	clt.logFileDst = junosLogFile

	// junosFakeCreateSetFile
	junosFakeCreateSetFile := c.junosFakeCreateSetFile
	if err := replaceTildeToHomeDir(&junosFakeCreateSetFile); err != nil {
		return clt, diag.FromErr(err)
	}
	clt.fakeCreateSetFile = junosFakeCreateSetFile

	return clt, nil
}
