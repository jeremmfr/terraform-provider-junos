package junos

import (
	"os"
	"strconv"
)

// Config : provider config
type Config struct {
	junosPort        int
	junosIP          string
	junosUserName    string
	junosSSHKeyFile  string
	junosKeyPass     string
	junosGroupIntDel string
}

// Session : read session information for Junos Device
func (c *Config) Session() (*Session, error) {
	junosLogFile, _ := os.LookupEnv("TFJUNOS_LOG_PATH")
	var junosSleep, junosSleepShort int
	var err error
	junosSleepEnv, _ := os.LookupEnv("TFJUNOS_SLEEP")
	if junosSleepEnv == "" {
		junosSleep = 10
	} else {
		junosSleep, err = strconv.Atoi(junosSleepEnv)
		if err != nil {
			return nil, err
		}
	}
	junosSleepShortEnv, _ := os.LookupEnv("TFJUNOS_SLEEP_SHORT")
	if junosSleepShortEnv == "" {
		junosSleepShort = 100
	} else {
		junosSleepShort, err = strconv.Atoi(junosSleepShortEnv)
		if err != nil {
			return nil, err
		}
	}
	sess := &Session{
		junosIP:          c.junosIP,
		junosPort:        c.junosPort,
		junosUserName:    c.junosUserName,
		junosSSHKeyFile:  c.junosSSHKeyFile,
		junosKeyPass:     c.junosKeyPass,
		junosGroupIntDel: c.junosGroupIntDel,
		junosLogFile:     junosLogFile,
		junosSleep:       junosSleep,
		junosSleepShort:  junosSleepShort,
	}

	return sess, nil
}
