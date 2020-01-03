package junos

// Config : provider config
type Config struct {
	junosPort                int
	junosCmdSleepShort       int
	junosCmdSleepLock        int
	junosIP                  string
	junosUserName            string
	junosSSHKeyFile          string
	junosKeyPass             string
	junosGroupIntDel         string
	junosDebugNetconfLogPath string
}

// Session : read session information for Junos Device
func (c *Config) Session() (*Session, error) {
	sess := &Session{
		junosIP:          c.junosIP,
		junosPort:        c.junosPort,
		junosUserName:    c.junosUserName,
		junosSSHKeyFile:  c.junosSSHKeyFile,
		junosKeyPass:     c.junosKeyPass,
		junosGroupIntDel: c.junosGroupIntDel,
		junosLogFile:     c.junosDebugNetconfLogPath,
		junosSleep:       c.junosCmdSleepLock,
		junosSleepShort:  c.junosCmdSleepShort,
	}
	return sess, nil
}
