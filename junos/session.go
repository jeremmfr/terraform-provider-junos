package junos

import (
	"fmt"
	"strconv"
	"time"
)

// Session information for connect to Junos Device
type Session struct {
	junosPort        int
	junosSleep       int
	junosSleepShort  int
	junosIP          string
	junosUserName    string
	junosSSHKeyFile  string
	junosKeyPass     string
	junosGroupIntDel string
	junosLogFile     string
}

func (sess *Session) startNewSession() (*NetconfObject, error) {
	var auth netconfAuthMethod
	auth.Username = sess.junosUserName
	auth.PrivateKey = sess.junosSSHKeyFile
	if sess.junosKeyPass != "" {
		auth.Passphrase = sess.junosKeyPass
	}
	jnpr, err := netconfNewSession(sess.junosIP+":"+strconv.Itoa(sess.junosPort), &auth)
	if err != nil {
		return nil, err
	}
	if len(jnpr.Platform) == 0 {
		return jnpr, fmt.Errorf("can't read platform junos with <get-software-information/>")
	}
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[startNewSession] started"), sess.junosLogFile)
	}
	return jnpr, nil
}
func (sess *Session) closeSession(jnpr *NetconfObject) {
	err := jnpr.Close()
	if sess.junosLogFile != "" {
		if err != nil {
			logFile(fmt.Sprintf("[closeSession] err: %q", err), sess.junosLogFile)
		} else {
			logFile(fmt.Sprintf("[closeSession] closed"), sess.junosLogFile)
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
func (sess *Session) commitConf(logMessage string, jnpr *NetconfObject) error {
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[commitConf] commit %q", logMessage), sess.junosLogFile)
	}
	err := jnpr.netconfCommit(logMessage)
	sleepShort(sess.junosSleepShort)
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[commitConf] commit error: %q", err), sess.junosLogFile)
		}
		return err
	}
	return nil
}

func (sess *Session) configLock(jnpr *NetconfObject) error {
	var lock bool
	for {
		lock = jnpr.netconfConfigLock()
		if lock {
			if sess.junosLogFile != "" {
				logFile(fmt.Sprintf("[configLock] locked"), sess.junosLogFile)
			}
			sleepShort(sess.junosSleepShort)
			break
		} else {
			if sess.junosLogFile != "" {
				logFile(fmt.Sprintf("[configLock] sleep for wait lock"), sess.junosLogFile)
			}
			sleep(sess.junosSleep)
		}
	}
	return nil
}
func (sess *Session) configClear(jnpr *NetconfObject) {
	err := jnpr.netconfConfigClear()
	sleepShort(sess.junosSleepShort)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[configClear] config clear"), sess.junosLogFile)
	}
	if err != nil {
		err := jnpr.Close()
		if err != nil && sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[configClear] close err: %q", err), sess.junosLogFile)
		}
		panic(err)
	}
	err = jnpr.netconfConfigUnlock()
	sleepShort(sess.junosSleepShort)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[configClear] config unlock"), sess.junosLogFile)
	}
	if err != nil {
		err := jnpr.Close()
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
