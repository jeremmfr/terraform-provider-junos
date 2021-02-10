package junos

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Session information to connect on Junos Device.
type Session struct {
	junosPort           int
	junosSleepLock      int
	junosSleepShort     int
	junosSleepSSHClosed int
	junosIP             string
	junosUserName       string
	junosPassword       string
	junosSSHKeyPEM      string
	junosSSHKeyFile     string
	junosKeyPass        string
	junosGroupIntDel    string
	junosLogFile        string
}

func (sess *Session) startNewSession() (*NetconfObject, error) {
	var auth netconfAuthMethod
	auth.Username = sess.junosUserName
	if sess.junosSSHKeyPEM != "" {
		auth.PrivateKeyPEM = sess.junosSSHKeyPEM
		if sess.junosKeyPass != "" {
			auth.Passphrase = sess.junosKeyPass
		}
	}
	if sess.junosSSHKeyFile != "" {
		auth.PrivateKeyFile = sess.junosSSHKeyFile
		if strings.HasPrefix(sess.junosSSHKeyFile, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to read user home directory : %w", err)
			}
			auth.PrivateKeyFile = homeDir + sess.junosSSHKeyFile[1:]
		}
		if sess.junosKeyPass != "" {
			auth.Passphrase = sess.junosKeyPass
		}
	}
	if sess.junosPassword != "" {
		auth.Password = sess.junosPassword
	}
	jnpr, err := netconfNewSession(sess.junosIP+":"+strconv.Itoa(sess.junosPort), &auth)
	if err != nil {
		return nil, err
	}
	if jnpr.SystemInformation.HardwareModel == "" {
		return jnpr, fmt.Errorf("can't read model of device with <get-system-information/> netconf command")
	}
	if sess.junosLogFile != "" {
		logFile("[startNewSession] started", sess.junosLogFile)
	}

	return jnpr, nil
}
func (sess *Session) closeSession(jnpr *NetconfObject) {
	err := jnpr.Close(sess.junosSleepSSHClosed)
	if sess.junosLogFile != "" {
		if err != nil {
			logFile(fmt.Sprintf("[closeSession] err: %q", err), sess.junosLogFile)
		} else {
			logFile("[closeSession] closed", sess.junosLogFile)
		}
	}
}
func (sess *Session) command(cmd string, jnpr *NetconfObject) (string, error) {
	read, err := jnpr.netconfCommand(cmd)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[command] cmd: %q", cmd), sess.junosLogFile)
		logFile(fmt.Sprintf("[command] read: %q", read), sess.junosLogFile)
	}
	sleepShort(sess.junosSleepShort)
	if err != nil && read != emptyWord {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[command] err: %q", err), sess.junosLogFile)
		}

		return "", err
	}

	return read, nil
}
func (sess *Session) commandXML(cmd string, jnpr *NetconfObject) (string, error) {
	read, err := jnpr.netconfCommandXML(cmd)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[commandXML] cmd: %q", cmd), sess.junosLogFile)
		logFile(fmt.Sprintf("[commandXML] read: %q", read), sess.junosLogFile)
	}
	sleepShort(sess.junosSleepShort)
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[commandXML] err: %q", err), sess.junosLogFile)
		}

		return "", err
	}

	return read, nil
}
func (sess *Session) configSet(cmd []string, jnpr *NetconfObject) error {
	message, err := jnpr.netconfConfigSet(cmd)
	sleepShort(sess.junosSleepShort)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[configSet] cmd: %q", cmd), sess.junosLogFile)
		logFile(fmt.Sprintf("[configSet] message: %q", message), sess.junosLogFile)
	}
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[configSet] err: %q", err), sess.junosLogFile)
		}

		return err
	}

	return nil
}
func (sess *Session) commitConf(logMessage string, jnpr *NetconfObject) (_warnings []error, _err error) {
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[commitConf] commit %q", logMessage), sess.junosLogFile)
	}
	warns, err := jnpr.netconfCommit(logMessage)
	sleepShort(sess.junosSleepShort)
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[commitConf] commit error: %q", err), sess.junosLogFile)
			if len(warns) > 0 {
				for _, w := range warns {
					logFile(fmt.Sprintf("[commitConf] commit warning: %q", w), sess.junosLogFile)
				}
			}
		}

		return warns, err
	}
	if len(warns) > 0 && sess.junosLogFile != "" {
		for _, w := range warns {
			logFile(fmt.Sprintf("[commitConf] commit warning: %q", w), sess.junosLogFile)
		}
	}

	return warns, nil
}

func (sess *Session) configLock(jnpr *NetconfObject) {
	var lock bool
	for {
		lock = jnpr.netconfConfigLock()
		if lock {
			if sess.junosLogFile != "" {
				logFile("[configLock] locked", sess.junosLogFile)
			}
			sleepShort(sess.junosSleepShort)

			break
		} else {
			if sess.junosLogFile != "" {
				logFile("[configLock] sleep for wait lock", sess.junosLogFile)
			}
			sleep(sess.junosSleepLock)
		}
	}
}
func (sess *Session) configClear(jnpr *NetconfObject) {
	err := jnpr.netconfConfigClear()
	sleepShort(sess.junosSleepShort)
	if sess.junosLogFile != "" {
		logFile("[configClear] config clear", sess.junosLogFile)
	}
	if err != nil {
		err := jnpr.Close(sess.junosSleepSSHClosed)
		if err != nil && sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[configClear] close err: %q", err), sess.junosLogFile)
		}
		panic(err)
	}
	err = jnpr.netconfConfigUnlock()
	sleepShort(sess.junosSleepShort)
	if sess.junosLogFile != "" {
		logFile("[configClear] config unlock", sess.junosLogFile)
	}
	if err != nil {
		err := jnpr.Close(sess.junosSleepSSHClosed)
		if err != nil && sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[configClear] close err: %q", err), sess.junosLogFile)
		}
		panic(err)
	}
}

func sleep(timeSleep int) {
	time.Sleep(time.Duration(timeSleep) * time.Second)
}

func sleepShort(timeSleep int) {
	time.Sleep(time.Duration(timeSleep) * time.Millisecond)
}
