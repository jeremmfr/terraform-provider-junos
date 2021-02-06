package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
)

const warningSeverity string = "warning"

var (
	rpcCommand         = "<command format=\"text\">%s</command>"
	rpcConfigStringSet = "<load-configuration action=\"set\" format=\"text\">" +
		"<configuration-set>%s</configuration-set></load-configuration>"
	rpcSystemInfo      = "<get-system-information/>"
	rpcCommit          = "<commit-configuration><log>%s</log></commit-configuration>"
	rpcCandidateLock   = "<lock><target><candidate/></target></lock>"
	rpcCandidateUnlock = "<unlock><target><candidate/></target></unlock>"
	rpcClearCandidate  = "<delete-config><target><candidate/></target></delete-config>"
	rpcClose           = "<close-session/>"
)

// NetconfObject : store Junos device info and session.
type NetconfObject struct {
	Session           *netconf.Session
	SystemInformation sysInfo `xml:"system-information"`
}

type sysInfo struct {
	HardwareModel string `xml:"hardware-model"`
	OsName        string `xml:"os-name"`
	OsVersion     string `xml:"os-version"`
	SerialNumber  string `xml:"serial-number"`
	HostName      string `xml:"host-name"`
	ClusterNode   *bool  `xml:"cluster-node"`
}

// RoutingEngine : store Platform information.
type RoutingEngine struct {
	Model   string
	Version string
}
type commandXMLConfig struct {
	Config string `xml:",innerxml"`
}
type netconfAuthMethod struct {
	Password       string
	Username       string
	PrivateKeyPEM  string
	PrivateKeyFile string
	Passphrase     string
}
type commitError struct {
	Path     string `xml:"error-path"`
	Element  string `xml:"error-info>bad-element"`
	Message  string `xml:"error-message"`
	Severity string `xml:"error-severity"`
}
type commitResults struct {
	XMLName xml.Name      `xml:"commit-results"`
	Errors  []commitError `xml:"rpc-error"`
}

// netconfNewSession establishes a new connection to a NetconfObject device that we will use
// to run our commands against.
// Authentication methods are defined using the netconfAuthMethod struct, and are as follows:
// username and password, SSH private key (with or without passphrase).
func netconfNewSession(host string, auth *netconfAuthMethod) (*NetconfObject, error) {
	clientConfig, err := genSSHClientConfig(auth)
	if err != nil {
		return nil, err
	}

	return netconfNewSessionWithConfig(host, clientConfig)
}

// netconfNewSessionWithConfig establishes a new connection to a NetconfObject device that we will use
// to run our commands against.
func netconfNewSessionWithConfig(host string, clientConfig *ssh.ClientConfig) (*NetconfObject, error) {
	s, err := netconf.DialSSH(host, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s - %w", host, err)
	}

	return newSessionFromNetconf(s)
}

// newSessionFromNetconf uses an existing netconf.Session to run our commands against.
func newSessionFromNetconf(s *netconf.Session) (*NetconfObject, error) {
	n := &NetconfObject{
		Session: s,
	}

	return n, n.gatherFacts()
}

// genSSHClientConfig is a wrapper function based around the auth method defined
// (user/password or private key) which returns the SSH client configuration used to
// connect.
func genSSHClientConfig(auth *netconfAuthMethod) (*ssh.ClientConfig, error) {
	var config *ssh.ClientConfig

	if len(auth.PrivateKeyPEM) > 0 {
		config, err := netconf.SSHConfigPubKeyPem(auth.Username, []byte(auth.PrivateKeyPEM), auth.Passphrase)
		if err != nil {
			return config, fmt.Errorf("failed to create new SSHConfig with PEM private key : %w", err)
		}
		config.Ciphers = append(config.Ciphers,
			"aes128-gcm@openssh.com", "chacha20-poly1305@openssh.com",
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
			"aes128-cbc")
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

		return config, nil
	}
	if len(auth.PrivateKeyFile) > 0 {
		config, err := netconf.SSHConfigPubKeyFile(auth.Username, auth.PrivateKeyFile, auth.Passphrase)
		if err != nil {
			return config, fmt.Errorf("failed to create new SSHConfig with file private key : %w", err)
		}
		config.Ciphers = append(config.Ciphers,
			"aes128-gcm@openssh.com", "chacha20-poly1305@openssh.com",
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
			"aes128-cbc")
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

		return config, nil
	}
	if len(auth.Password) > 0 {
		config = netconf.SSHConfigPassword(auth.Username, auth.Password)
		config.Ciphers = append(config.Ciphers,
			"aes128-gcm@openssh.com", "chacha20-poly1305@openssh.com",
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
			"aes128-cbc")
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

		return config, nil
	}

	return config, errors.New("no credentials/keys available")
}

// gatherFacts gathers basic information about the device.
func (j *NetconfObject) gatherFacts() error {
	if j == nil {
		return errors.New("attempt to call GatherFacts on nil NetconfObject object")
	}
	s := j.Session
	// Get info for get-system-information and populate SystemInformation Struct
	val, err := s.Exec(netconf.RawMethod(rpcSystemInfo))
	if err != nil {
		return fmt.Errorf("failed to netconf get-system-information : %w", err)
	}

	if val.Errors != nil {
		var errorsMsg []string
		for _, m := range val.Errors {
			errorsMsg = append(errorsMsg, fmt.Sprintf("%v", m))
		}

		return fmt.Errorf(strings.Join(errorsMsg, "\n"))
	}
	err = xml.Unmarshal([]byte(val.RawReply), &j)
	if err != nil {
		return fmt.Errorf("failed to xml unmarshal reply : %w", err)
	}

	return nil
}

// netconfCommand (show, execute) on Junos device.
func (j *NetconfObject) netconfCommand(cmd string) (string, error) {
	command := fmt.Sprintf(rpcCommand, cmd)
	reply, err := j.Session.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("failed to netconf command exec : %w", err)
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}
	if reply.Data == "" || strings.Count(reply.Data, "") <= 2 {
		return emptyWord, errors.New("no output available - please check the syntax of your command")
	}
	var output commandXMLConfig
	if err := xml.Unmarshal([]byte(reply.Data), &output); err != nil {
		return "", fmt.Errorf("failed to xml unmarshal reply : %w", err)
	}

	return output.Config, nil
}
func (j *NetconfObject) netconfCommandXML(cmd string) (string, error) {
	reply, err := j.Session.Exec(netconf.RawMethod(cmd))
	if err != nil {
		return "", fmt.Errorf("failed to netconf xml command exec : %w", err)
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	return reply.Data, nil
}

func (j *NetconfObject) netconfConfigSet(cmd []string) (string, error) {
	command := fmt.Sprintf(rpcConfigStringSet, strings.Join(cmd, "\n"))
	reply, err := j.Session.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("failed to netconf set/delete command exec : %w", err)
	}
	// logFile("netconfConfigSet.Reply:" + reply.RawReply)
	message := ""
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			message += m.Message
		}

		return message, nil
	}

	return "", nil
}

// netConfConfigLock locks the candidate configuration.
func (j *NetconfObject) netconfConfigLock() bool {
	reply, err := j.Session.Exec(netconf.RawMethod(rpcCandidateLock))
	if err != nil {
		return false
	}
	if reply.Errors != nil {
		return false
	}

	return true
}

// Unlock unlocks the candidate configuration.
func (j *NetconfObject) netconfConfigUnlock() error {
	reply, err := j.Session.Exec(netconf.RawMethod(rpcCandidateUnlock))
	if err != nil {
		return fmt.Errorf("failed to netconf config unlock : %w", err)
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}
func (j *NetconfObject) netconfConfigClear() error {
	reply, err := j.Session.Exec(netconf.RawMethod(rpcClearCandidate))
	if err != nil {
		return fmt.Errorf("failed to netconf config clear : %w", err)
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// netconfCommit commits the configuration.
func (j *NetconfObject) netconfCommit(logMessage string) (_warn []error, _err error) {
	var errs commitResults
	reply, err := j.Session.Exec(netconf.RawMethod(fmt.Sprintf(rpcCommit, logMessage)))
	if err != nil {
		return []error{}, fmt.Errorf("failed to netconf commit : %w", err)
	}

	if reply.Errors != nil {
		warnings := make([]error, 0)
		for _, m := range reply.Errors {
			if m.Severity != warningSeverity {
				return warnings, errors.New(m.Message)
			}
			warnings = append(warnings, errors.New(m.Message))
		}

		return warnings, nil
	}

	if reply.Data != "\n<ok/>\n" {
		err = xml.Unmarshal([]byte(reply.Data), &errs)
		if err != nil {
			return []error{}, fmt.Errorf("failed to xml unmarshal reply : %w", err)
		}

		if errs.Errors != nil {
			warnings := make([]error, 0)
			for _, m := range errs.Errors {
				if m.Severity != warningSeverity {
					message := fmt.Sprintf("[%s]\n    %s\nError: %s",
						strings.Trim(m.Path, "[\r\n]"),
						strings.Trim(m.Element, "[\r\n]"),
						strings.Trim(m.Message, "[\r\n]"))

					return []error{}, errors.New(message)
				}
				warnings = append(warnings, errors.New(m.Message))
			}

			return warnings, nil
		}
	}

	return []error{}, nil
}

// Close disconnects our session to the device.
func (j *NetconfObject) Close(sleepClosed int) error {
	_, err := j.Session.Exec(netconf.RawMethod(rpcClose))
	j.Session.Transport.Close()
	if err != nil {
		sleep(sleepClosed)

		return fmt.Errorf("failed to netconf close : %w", err)
	}
	sleep(sleepClosed)

	return nil
}
