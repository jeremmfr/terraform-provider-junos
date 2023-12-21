package junos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/jeremmfr/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
)

// Session : store Junos device info and session.
type Session struct {
	SystemInformation      sysInfo
	netconf                *netconf.Session
	localAddress           string
	remoteAddress          string
	logFile                func(string)
	fakeSetFile            func([]string) error
	sleepShort             int
	sleepLock              int
	commitConfirmedTimeout int
	commitConfirmedWait    time.Duration
	sleepSSHClosed         int
}

type sshAuthMethod struct {
	Password       string
	Username       string
	PrivateKeyPEM  string
	PrivateKeyFile string
	Passphrase     string
	Ciphers        []string
	Timeout        int
}

type openSSHOptions struct {
	Retry   int
	Timeout int
}

type sshOptions struct {
	*openSSHOptions
	ClientConfig *ssh.ClientConfig
}

// netconfNewSession establishes a new connection to a Junos device that we will use
// to run our commands against.
// Authentication methods are defined using the netconfAuthMethod struct, and are as follows:
// username and password, SSH private key (with or without passphrase).
func netconfNewSession(
	ctx context.Context,
	host string,
	auth *sshAuthMethod,
	openSSH *openSSHOptions,
) (
	*Session, error,
) {
	clientConfig, err := genSSHClientConfig(auth)
	if err != nil {
		return nil, err
	}

	return netconfNewSessionWithConfig(ctx, host, &sshOptions{openSSH, clientConfig})
}

// netconfNewSessionWithConfig establishes a new connection to a Junos device that we will use
// to run our commands against.
func netconfNewSessionWithConfig(
	ctx context.Context,
	host string,
	sshOpts *sshOptions,
) (
	*Session, error,
) {
	netDialer := net.Dialer{
		Timeout: time.Duration(sshOpts.Timeout) * time.Second,
	}
	retry := sshOpts.Retry
	if retry < 1 {
		retry = 1
	}
	if retry > 10 {
		retry = 10
	}
	sleepTime := 0
toretry:
	for retry > 0 {
		retry--
		conn, err := netDialer.DialContext(ctx, "tcp", host)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("error connecting to %s: %w", host, err)
			default:
				if retry != 0 {
					log.Printf("[WARN] connecting to %s: %s, go retry", host, err.Error())
					// sleep with time increasing as things try
					sleepTime++
					utils.Sleep(sleepTime)

					continue toretry
				}

				return nil, fmt.Errorf("error connecting to %s: %w", host, err)
			}
		}
		s, err := netconf.NewSSHSession(conn, sshOpts.ClientConfig)
		if err != nil {
			conn.Close()
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("initializing SSH session to %s: %w", host, err)
			default:
				if retry != 0 {
					log.Printf("[WARN] initializing SSH session to %s: %s, go retry", host, err.Error())
					// sleep with time increasing as things try
					sleepTime++
					utils.Sleep(sleepTime)

					continue toretry
				}

				return nil, fmt.Errorf("initializing SSH session to %s: %w", host, err)
			}
		}

		return newSessionFromNetconf(s, conn.LocalAddr().String(), conn.RemoteAddr().String())
	}
	// this return can't happen
	return nil, fmt.Errorf("connecting to %s: retries exceeded", host)
}

// newSessionFromNetconf uses an existing netconf.Session to run our commands against.
func newSessionFromNetconf(
	netConfSess *netconf.Session,
	localAddress,
	remoteAddress string,
) (
	*Session, error,
) {
	sess := &Session{
		netconf:       netConfSess,
		localAddress:  localAddress,
		remoteAddress: remoteAddress,
	}

	return sess, sess.gatherFacts()
}

// genSSHClientConfig is a wrapper function based around the auth method defined
// (user/password or private key) which returns the SSH client configuration used to
// connect.
func genSSHClientConfig(auth *sshAuthMethod) (*ssh.ClientConfig, error) {
	configs := make([]*ssh.ClientConfig, 0)
	configs = append(configs, &ssh.ClientConfig{})

	// keys method
	switch {
	case len(auth.PrivateKeyPEM) > 0:
		config, err := netconf.SSHConfigPubKeyPem(auth.Username, []byte(auth.PrivateKeyPEM), auth.Passphrase)
		if err != nil {
			return config, fmt.Errorf("creating new SSHConfig with PEM private key: %w", err)
		}
		configs = append(configs, config)
	case len(auth.PrivateKeyFile) > 0:
		config, err := netconf.SSHConfigPubKeyFile(auth.Username, auth.PrivateKeyFile, auth.Passphrase)
		if err != nil {
			return config, fmt.Errorf("creating new SSHConfig with file private key: %w", err)
		}
		configs = append(configs, config)
	case os.Getenv("SSH_AUTH_SOCK") != "":
		config, err := netconf.SSHConfigPubKeyAgent(auth.Username)
		if err != nil {
			log.Printf("[WARN] communicating with SSH agent: %s", err.Error())
		} else {
			configs = append(configs, config)
		}
	}
	if len(auth.Password) > 0 {
		config := netconf.SSHConfigPassword(auth.Username, auth.Password)
		configs = append(configs, config)
	}
	if len(configs) == 1 {
		return configs[0], errors.New("no credentials/keys available")
	}
	configs[0] = configs[1]
	configs[0].Ciphers = auth.Ciphers
	configs[0].HostKeyCallback = ssh.InsecureIgnoreHostKey()
	for _, v := range configs[2:] {
		configs[0].Auth = append(configs[0].Auth, v.Auth...)
	}
	configs[0].Timeout = time.Duration(auth.Timeout) * time.Second

	return configs[0], nil
}

func (sess *Session) HasNetconf() bool {
	return sess.netconf != nil
}

// Command (show, execute) on Junos device via netconf.
func (sess *Session) Command(cmd string) (string, error) {
	read, err := sess.netconfCommand(cmd)
	sess.logFile(fmt.Sprintf("[Command] cmd: %q", cmd))
	sess.logFile(fmt.Sprintf("[Command] read: %q", read))
	utils.SleepShort(sess.sleepShort)
	if err != nil && read != EmptyW {
		sess.logFile(fmt.Sprintf("[Command] err: %q", err))

		return "", err
	}

	return read, nil
}

// CommandXML send XML cmd on Junos device via netconf.
func (sess *Session) CommandXML(cmd string) (string, error) {
	read, err := sess.netconfCommandXML(cmd)
	sess.logFile(fmt.Sprintf("[CommandXML] cmd: %q", cmd))
	sess.logFile(fmt.Sprintf("[CommandXML] read: %q", read))
	utils.SleepShort(sess.sleepShort)
	if err != nil {
		sess.logFile(fmt.Sprintf("[CommandXML] err: %q", err))

		return "", err
	}

	return read, nil
}

// ConfigSet append candidate configuration with set/delete lines
// on Junos device via netconf or in fake file if set.
func (sess *Session) ConfigSet(cmd []string) error {
	if sess.netconf != nil {
		message, err := sess.netconfConfigSet(cmd)
		utils.SleepShort(sess.sleepShort)
		sess.logFile(fmt.Sprintf("[ConfigSet] cmd: %q", cmd))
		sess.logFile(fmt.Sprintf("[ConfigSet] message: %q", message))
		if err != nil {
			sess.logFile(fmt.Sprintf("[ConfigSet] err: %q", err))

			return err
		}

		return nil
	} else if sess.fakeSetFile != nil {
		return sess.fakeSetFile(cmd)
	}

	return fmt.Errorf("internal error: call Session.ConfigSet without netconf session or fake set file")
}

// ConfigLock lock candidate configuration and retry with sleep between when fail.
func (sess *Session) ConfigLock(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			sess.logFile("[ConfigLock] lock aborted")

			return fmt.Errorf("candidate configuration lock attempt aborted")
		default:
			if sess.netconfConfigLock() {
				sess.logFile("[ConfigLock] config locked")
				utils.SleepShort(sess.sleepShort)

				return nil
			}
			sess.logFile("[ConfigLock] sleep to wait the lock")
			utils.Sleep(sess.sleepLock)
		}
	}
}

// ConfigClear clear potential candidate configuration and unlock it.
func (sess *Session) ConfigClear() (errs []error) {
	errs = append(errs, sess.netconfConfigClear()...)
	errs = append(errs, sess.netconfConfigUnlock()...)

	sess.logFile("[ConfigClear] config cleared/unlocked")
	utils.SleepShort(sess.sleepShort)

	return
}

// CommitConf commit the configuration with message via netconf.
func (sess *Session) CommitConf(ctx context.Context, logMessage string) (warnings []error, err error) {
	if sess.commitConfirmedTimeout > 0 {
		sess.logFile(fmt.Sprintf(
			"[CommitConf] commit confirmed %d (wait %s) %q",
			sess.commitConfirmedTimeout, sess.commitConfirmedWait, logMessage,
		))
		warnings, err = sess.netconfCommitConfirmed(ctx, logMessage)
	} else {
		sess.logFile(fmt.Sprintf("[CommitConf] commit %q", logMessage))
		warnings, err = sess.netconfCommit(logMessage)
	}
	utils.SleepShort(sess.sleepShort)
	if len(warnings) > 0 {
		for _, w := range warnings {
			sess.logFile(fmt.Sprintf("[CommitConf] commit warning: %q", w))
		}
	}
	if err != nil {
		sess.logFile(fmt.Sprintf("[CommitConf] commit error: %q", err))

		return warnings, err
	}

	return warnings, nil
}

func (sess *Session) Close() {
	if sess.HasNetconf() {
		err := sess.closeNetconf(sess.sleepSSHClosed)
		if err != nil {
			sess.logFile(fmt.Sprintf("[Close] err: %q", err))
		} else {
			sess.logFile("[Close] session closed")
		}
	}
}
