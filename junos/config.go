package junos

import (
	"os"
)

type Config struct {
	junosPort       int
	junosIP         string
	junosUserName   string
	junosSSHKeyFile string
	junosKeyPass    string
}

func (c *Config) Session() (*Session, error) {
	junosLogFile, _ := os.LookupEnv("TFJUNOS_LOG_PATH")
	sess := &Session{
		junosIP:         c.junosIP,
		junosPort:       c.junosPort,
		junosUserName:   c.junosUserName,
		junosSSHKeyFile: c.junosSSHKeyFile,
		junosKeyPass:    c.junosKeyPass,
		junosLogFile:    junosLogFile,
	}

	return sess, nil
}
