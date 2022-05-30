package junos

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

const directoryPermission = 0o755

// Client information to connect on Junos Device and more.
type Client struct {
	fakeUpdateAlso         bool
	fakeDeleteAlso         bool
	junosPort              int
	sleepLock              int
	sleepShort             int
	sleepSSHClosed         int
	junosSSHTimeoutToEstab int
	filePermission         int64
	junosIP                string
	junosUserName          string
	junosPassword          string
	junosSSHKeyPEM         string
	junosSSHKeyFile        string
	junosSSHKeyPass        string
	groupIntDel            string
	logFileDst             string
	fakeCreateSetFile      string
	junosSSHCiphers        []string
}

func (clt *Client) startNewSession(ctx context.Context) (*junosSession, error) {
	var auth netconfAuthMethod
	auth.Username = clt.junosUserName
	auth.Ciphers = clt.junosSSHCiphers
	if clt.junosSSHKeyPEM != "" {
		auth.PrivateKeyPEM = clt.junosSSHKeyPEM
		if clt.junosSSHKeyPass != "" {
			auth.Passphrase = clt.junosSSHKeyPass
		}
	}
	if clt.junosSSHKeyFile != "" {
		auth.PrivateKeyFile = clt.junosSSHKeyFile
		if clt.junosSSHKeyPass != "" {
			auth.Passphrase = clt.junosSSHKeyPass
		}
	}
	if clt.junosPassword != "" {
		auth.Password = clt.junosPassword
	}
	auth.Timeout = clt.junosSSHTimeoutToEstab
	junSess, err := netconfNewSession(ctx, net.JoinHostPort(clt.junosIP, strconv.Itoa(clt.junosPort)), &auth)
	if err != nil {
		return nil, err
	}
	if junSess.SystemInformation.HardwareModel == "" {
		return junSess, fmt.Errorf("can't read model of device with <get-system-information/> netconf command")
	}
	clt.logFile("[startNewSession] started")

	return junSess, nil
}

func (clt *Client) closeSession(junSess *junosSession) {
	err := junSess.close(clt.sleepSSHClosed)
	if err != nil {
		clt.logFile(fmt.Sprintf("[closeSession] err: %q", err))
	} else {
		clt.logFile("[closeSession] closed")
	}
}

func (clt *Client) command(cmd string, junSess *junosSession) (string, error) {
	read, err := junSess.netconfCommand(cmd)
	clt.logFile(fmt.Sprintf("[command] cmd: %q", cmd))
	clt.logFile(fmt.Sprintf("[command] read: %q", read))
	sleepShort(clt.sleepShort)
	if err != nil && read != emptyW {
		clt.logFile(fmt.Sprintf("[command] err: %q", err))

		return "", err
	}

	return read, nil
}

func (clt *Client) commandXML(cmd string, junSess *junosSession) (string, error) {
	read, err := junSess.netconfCommandXML(cmd)
	clt.logFile(fmt.Sprintf("[commandXML] cmd: %q", cmd))
	clt.logFile(fmt.Sprintf("[commandXML] read: %q", read))
	sleepShort(clt.sleepShort)
	if err != nil {
		clt.logFile(fmt.Sprintf("[commandXML] err: %q", err))

		return "", err
	}

	return read, nil
}

func (clt *Client) configSet(cmd []string, junSess *junosSession) error {
	if junSess != nil {
		message, err := junSess.netconfConfigSet(cmd)
		sleepShort(clt.sleepShort)
		clt.logFile(fmt.Sprintf("[configSet] cmd: %q", cmd))
		clt.logFile(fmt.Sprintf("[configSet] message: %q", message))
		if err != nil {
			clt.logFile(fmt.Sprintf("[configSet] err: %q", err))

			return err
		}

		return nil
	} else if clt.fakeCreateSetFile != "" {
		return clt.appendFakeCreateSetFile(cmd)
	}

	return fmt.Errorf("internal error: Client.configSet call without connection on device")
}

func (clt *Client) appendFakeCreateSetFile(lines []string) error {
	dirSetFile := path.Dir(clt.fakeCreateSetFile)
	if _, err := os.Stat(dirSetFile); err != nil {
		if err := os.MkdirAll(dirSetFile, os.FileMode(directoryPermission)); err != nil {
			return fmt.Errorf("failed to create parent directory of '%s': %w", clt.fakeCreateSetFile, err)
		}
	}
	f, err := os.OpenFile(clt.fakeCreateSetFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(clt.filePermission))
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", clt.fakeCreateSetFile, err)
	}
	defer f.Close()
	for _, v := range lines {
		if _, err := f.WriteString(v + "\n"); err != nil {
			return fmt.Errorf("failed to write in file '%s': %w", clt.fakeCreateSetFile, err)
		}
	}

	return nil
}

func (clt *Client) commitConf(logMessage string, junSess *junosSession) (_warnings []error, _err error) {
	clt.logFile(fmt.Sprintf("[commitConf] commit %q", logMessage))
	warns, err := junSess.netconfCommit(logMessage)
	sleepShort(clt.sleepShort)
	if len(warns) > 0 {
		for _, w := range warns {
			clt.logFile(fmt.Sprintf("[commitConf] commit warning: %q", w))
		}
	}
	if err != nil {
		clt.logFile(fmt.Sprintf("[commitConf] commit error: %q", err))

		return warns, err
	}

	return warns, nil
}

func (clt *Client) configLock(ctx context.Context, junSess *junosSession) error {
	for {
		select {
		case <-ctx.Done():
			clt.logFile("[configLock] aborted")

			return fmt.Errorf("candidate configuration lock attempt aborted")
		default:
			if junSess.netconfConfigLock() {
				clt.logFile("[configLock] locked")
				sleepShort(clt.sleepShort)

				return nil
			}
			clt.logFile("[configLock] sleep to wait the lock")
			sleep(clt.sleepLock)
		}
	}
}

func (clt *Client) configClear(junSess *junosSession) (errs []error) {
	errs = append(errs, junSess.netconfConfigClear()...)
	sleepShort(clt.sleepShort)
	clt.logFile("[configClear] config clear")

	errs = append(errs, junSess.netconfConfigUnlock()...)
	sleepShort(clt.sleepShort)
	clt.logFile("[configClear] config unlock")

	return
}

// log message in junosLogFile.
func (clt *Client) logFile(message string) {
	if clt.logFileDst != "" {
		f, err := os.OpenFile(clt.logFileDst,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(clt.filePermission))
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
