package junos

import (
	"fmt"
	"strconv"
	"time"
)

type Session struct {
	junosPort        int
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
		return jnpr, fmt.Errorf("can't read Platform Junos with <get-software-information/>")
	}
	return jnpr, nil
}
func (sess *Session) closeSession(jnpr *NetconfObject) {
	jnpr.Close()
}
func (sess *Session) command(cmd string, jnpr *NetconfObject) (string, error) {
	read, err := jnpr.netconfCommand(cmd)
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[command] cmd: %q", cmd), sess.junosLogFile)
		logFile(fmt.Sprintf("[command] read: %q", read), sess.junosLogFile)
	}
	sleepShort()
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
	sleepShort()
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
	sleepShort()
	return nil
}
func (sess *Session) commitConf(jnpr *NetconfObject) error {
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[commitConf] commit check"), sess.junosLogFile)
	}
	err := jnpr.netconfCommitCheck()
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[commitConf] commit check error: %q", err), sess.junosLogFile)
		}
		return err
	}
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[commitConf] commit"), sess.junosLogFile)
	}
	err = jnpr.netconfCommit()
	if err != nil {
		if sess.junosLogFile != "" {
			logFile(fmt.Sprintf("[commitConf] commit error: %q", err), sess.junosLogFile)
		}
		return err
	}
	sleepShort()
	return nil
}

func (sess *Session) configLock(jnpr *NetconfObject) error {
	var lock bool
	for {
		lock = jnpr.netconfConfigLock()
		if lock {
			if sess.junosLogFile != "" {
				logFile(fmt.Sprintf("[configLock] Locked"), sess.junosLogFile)
			}
			break
		} else {
			if sess.junosLogFile != "" {
				logFile(fmt.Sprintf("[configLock] sleep for wait Lock"), sess.junosLogFile)
			}
			sleep()
		}
	}
	return nil
}
func (sess *Session) configClear(jnpr *NetconfObject) {
	err := jnpr.netconfConfigClear()
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[configClear] config clear"), sess.junosLogFile)
	}
	if err != nil {
		jnpr.Close()
		panic(err)
	}
	err = jnpr.netconfConfigUnlock()
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[configClear] config unlock"), sess.junosLogFile)
	}
	if err != nil {
		jnpr.Close()
		panic(err)
	}
}

func sleep() {
	timeSleep := 10
	time.Sleep(time.Duration(timeSleep) * time.Second)
}

func sleepShort() {
	timeSleep := 100
	time.Sleep(time.Duration(timeSleep) * time.Millisecond)
}
