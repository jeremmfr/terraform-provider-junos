package junos

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

const directoryPermission = 0755

// Session information to connect on Junos Device and more.
type Session struct {
	junosPort              int
	junosSleepLock         int
	junosSleepShort        int
	junosSleepSSHClosed    int
	junosFilePermission    int64
	junosIP                string
	junosUserName          string
	junosPassword          string
	junosSSHKeyPEM         string
	junosSSHKeyFile        string
	junosKeyPass           string
	junosGroupIntDel       string
	junosLogFile           string
	junosFakeCreateSetFile string
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
	sess.logFile("[startNewSession] started")

	return jnpr, nil
}

func (sess *Session) closeSession(jnpr *NetconfObject) {
	err := jnpr.close(sess.junosSleepSSHClosed)
	if err != nil {
		sess.logFile(fmt.Sprintf("[closeSession] err: %q", err))
	} else {
		sess.logFile("[closeSession] closed")
	}
}

func (sess *Session) command(cmd string, jnpr *NetconfObject) (string, error) {
	read, err := jnpr.netconfCommand(cmd)
	sess.logFile(fmt.Sprintf("[command] cmd: %q", cmd))
	sess.logFile(fmt.Sprintf("[command] read: %q", read))
	sleepShort(sess.junosSleepShort)
	if err != nil && read != emptyWord {
		sess.logFile(fmt.Sprintf("[command] err: %q", err))

		return "", err
	}

	return read, nil
}

func (sess *Session) commandXML(cmd string, jnpr *NetconfObject) (string, error) {
	read, err := jnpr.netconfCommandXML(cmd)
	sess.logFile(fmt.Sprintf("[commandXML] cmd: %q", cmd))
	sess.logFile(fmt.Sprintf("[commandXML] read: %q", read))
	sleepShort(sess.junosSleepShort)
	if err != nil {
		sess.logFile(fmt.Sprintf("[commandXML] err: %q", err))

		return "", err
	}

	return read, nil
}

func (sess *Session) configSet(cmd []string, jnpr *NetconfObject) error {
	if jnpr != nil {
		message, err := jnpr.netconfConfigSet(cmd)
		sleepShort(sess.junosSleepShort)
		sess.logFile(fmt.Sprintf("[configSet] cmd: %q", cmd))
		sess.logFile(fmt.Sprintf("[configSet] message: %q", message))
		if err != nil {
			sess.logFile(fmt.Sprintf("[configSet] err: %q", err))

			return err
		}

		return nil
	} else if sess.junosFakeCreateSetFile != "" {
		return sess.appendFakeCreateSetFile(cmd)
	}

	return fmt.Errorf("internal error: sess.configSet call without connection on device")
}

func (sess *Session) appendFakeCreateSetFile(lines []string) error {
	dirSetFile := path.Dir(sess.junosFakeCreateSetFile)
	if _, err := os.Stat(dirSetFile); err != nil {
		if err := os.MkdirAll(dirSetFile, os.FileMode(directoryPermission)); err != nil {
			return fmt.Errorf("failed to create parent directory of `%s` : %w", sess.junosFakeCreateSetFile, err)
		}
	}
	f, err := os.OpenFile(sess.junosFakeCreateSetFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(sess.junosFilePermission))
	if err != nil {
		return fmt.Errorf("failed to openfile `%s` : %w", sess.junosFakeCreateSetFile, err)
	}
	defer f.Close()
	for _, v := range lines {
		if _, err := f.WriteString(v + "\n"); err != nil {
			return fmt.Errorf("failed to write in file `%s` : %w", sess.junosFakeCreateSetFile, err)
		}
	}

	return nil
}

func (sess *Session) commitConf(logMessage string, jnpr *NetconfObject) (_warnings []error, _err error) {
	sess.logFile(fmt.Sprintf("[commitConf] commit %q", logMessage))
	warns, err := jnpr.netconfCommit(logMessage)
	sleepShort(sess.junosSleepShort)
	if len(warns) > 0 {
		for _, w := range warns {
			sess.logFile(fmt.Sprintf("[commitConf] commit warning: %q", w))
		}
	}
	if err != nil {
		sess.logFile(fmt.Sprintf("[commitConf] commit error: %q", err))

		return warns, err
	}

	return warns, nil
}

func (sess *Session) configLock(jnpr *NetconfObject) {
	var lock bool
	for {
		lock = jnpr.netconfConfigLock()
		if lock {
			sess.logFile("[configLock] locked")
			sleepShort(sess.junosSleepShort)

			break
		} else {
			sess.logFile("[configLock] sleep for wait lock")
			sleep(sess.junosSleepLock)
		}
	}
}

func (sess *Session) configClear(jnpr *NetconfObject) {
	err := jnpr.netconfConfigClear()
	sleepShort(sess.junosSleepShort)
	sess.logFile("[configClear] config clear")
	if err != nil {
		err := jnpr.close(sess.junosSleepSSHClosed)
		if err != nil {
			sess.logFile(fmt.Sprintf("[configClear] close err: %q", err))
		}
		panic(err)
	}
	err = jnpr.netconfConfigUnlock()
	sleepShort(sess.junosSleepShort)
	sess.logFile("[configClear] config unlock")
	if err != nil {
		err := jnpr.close(sess.junosSleepSSHClosed)
		if err != nil {
			sess.logFile(fmt.Sprintf("[configClear] close err: %q", err))
		}
		panic(err)
	}
}

// log message in junosLogFile.
func (sess *Session) logFile(message string) {
	if sess.junosLogFile != "" {
		f, err := os.OpenFile(sess.junosLogFile,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(sess.junosFilePermission))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		log.SetOutput(f)
		log.SetPrefix(time.Now().Format("2006-01-02 15:04:05"))

		log.Printf("%s", message)
	}
}

func sleep(timeSleep int) {
	time.Sleep(time.Duration(timeSleep) * time.Second)
}

func sleepShort(timeSleep int) {
	time.Sleep(time.Duration(timeSleep) * time.Millisecond)
}
