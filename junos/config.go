package junos

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config : provider config.
type Config struct {
	junosPort                int
	junosCmdSleepShort       int
	junosCmdSleepLock        int
	junosIP                  string
	junosUserName            string
	junosPassword            string
	junosSSHKeyPEM           string
	junosSSHKeyFile          string
	junosKeyPass             string
	junosGroupIntDel         string
	junosDebugNetconfLogPath string
}

// Session : read session information for Junos Device.
func (c *Config) Session() (*Session, diag.Diagnostics) {
	sess := &Session{
		junosIP:          c.junosIP,
		junosPort:        c.junosPort,
		junosUserName:    c.junosUserName,
		junosPassword:    c.junosPassword,
		junosSSHKeyPEM:   c.junosSSHKeyPEM,
		junosSSHKeyFile:  c.junosSSHKeyFile,
		junosKeyPass:     c.junosKeyPass,
		junosGroupIntDel: c.junosGroupIntDel,
		junosLogFile:     c.junosDebugNetconfLogPath,
		junosSleep:       c.junosCmdSleepLock,
		junosSleepShort:  c.junosCmdSleepShort,
	}

	return sess, nil
}
