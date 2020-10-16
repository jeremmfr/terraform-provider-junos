package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jeremmfr/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
)

var (
	rpcCommand         = "<command format=\"text\">%s</command>"
	rpcConfigStringSet = "<load-configuration action=\"set\" format=\"text\">" +
		"<configuration-set>%s</configuration-set></load-configuration>"
	rpcVersion         = "<get-software-information/>"
	rpcCommit          = "<commit-configuration><log>%s</log></commit-configuration>"
	rpcCandidateLock   = "<lock><target><candidate/></target></lock>"
	rpcCandidateUnlock = "<unlock><target><candidate/></target></unlock>"
	rpcClearCandidate  = "<delete-config><target><candidate/></target></delete-config>"
	rpcClose           = "<close-session/>"
)

// NetconfObject : store Junos device info and session.
type NetconfObject struct {
	Session        *netconf.Session
	Hostname       string
	RoutingEngines int
	Platform       []RoutingEngine
	CommitTimeout  time.Duration
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
	Password   string
	Username   string
	PrivateKey string
	Passphrase string
}
type versionRouteEngines struct {
	XMLName xml.Name             `xml:"multi-routing-engine-results"`
	RE      []versionRouteEngine `xml:"multi-routing-engine-item>software-information"`
}

type versionRouteEngine struct {
	XMLName     xml.Name             `xml:"software-information"`
	Hostname    string               `xml:"host-name"`
	Platform    string               `xml:"product-model"`
	PackageInfo []versionPackageInfo `xml:"package-information"`
}
type versionPackageInfo struct {
	XMLName         xml.Name `xml:"package-information"`
	PackageName     []string `xml:"name"`
	SoftwareVersion []string `xml:"comment"`
}
type commitError struct {
	Path    string `xml:"error-path"`
	Element string `xml:"error-info>bad-element"`
	Message string `xml:"error-message"`
}
type commitResults struct {
	XMLName xml.Name      `xml:"commit-results"`
	Errors  []commitError `xml:"rpc-error"`
}

// netconfNewSession establishes a new connection to a NetconfObject device that we will use
// to run our commands against.
// Authentication methods are defined using the netconfAuthMethod struct, and are as follows:
//
// username and password, SSH private key (with or without passphrase)
//
// Please view the package documentation for netconfAuthMethod on how to use these methods.
//
// NOTE: most users should use this function, instead of the other NewSession* functions.
func netconfNewSession(host string, auth *netconfAuthMethod) (*NetconfObject, error) {
	clientConfig, err := genSSHClientConfig(auth)
	if err != nil {
		return nil, err
	}

	return netconfNewSessionWithConfig(host, clientConfig)
}

// netconfNewSessionWithConfig establishes a new connection to a NetconfObject device that we will use
// to run our commands against.
//
// This is especially useful if you need to customize the SSH connection beyond
// what's supported in NewSession().
func netconfNewSessionWithConfig(host string, clientConfig *ssh.ClientConfig) (*NetconfObject, error) {
	s, err := netconf.DialSSH(host, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s - %s", host, err)
	}

	return newSessionFromNetconf(s)
}

// newSessionFromNetconf uses an existing netconf.Session to run our commands against
//
// This is especially useful if you need to customize the SSH connection beyond
// what's supported in NewSession().
func newSessionFromNetconf(s *netconf.Session) (*NetconfObject, error) {
	n := &NetconfObject{
		Session: s,
	}

	return n, n.GatherFacts()
}

// genSSHClientConfig is a wrapper function based around the auth method defined
// (user/password or private key) which returns the SSH client configuration used to
// connect.
func genSSHClientConfig(auth *netconfAuthMethod) (*ssh.ClientConfig, error) {
	var config *ssh.ClientConfig

	if len(auth.PrivateKey) > 0 {
		config, err := netconf.SSHConfigPubKeyFile(auth.Username, auth.PrivateKey, auth.Passphrase)
		if err != nil {
			return config, err
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

// GatherFacts gathers basic information about the device.
//
// It's automatically called when using the provided NewSession* functions, but can be
// used if you create your own NetconfObject sessions.
func (j *NetconfObject) GatherFacts() error {
	if j == nil {
		return errors.New("attempt to call GatherFacts on nil NetconfObject object")
	}
	s := j.Session
	rex := regexp.MustCompile(`^.*\[(.*)\]`)

	reply, err := s.Exec(netconf.RawMethod(rpcVersion))
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		var errorsMsg []string
		for _, m := range reply.Errors {
			errorsMsg = append(errorsMsg, fmt.Sprintf("%v", m))
		}

		return fmt.Errorf(strings.Join(errorsMsg, "\n"))
	}

	formatted := strings.ReplaceAll(reply.Data, "\n", "")
	if strings.Contains(reply.Data, "multi-routing-engine-results") {
		var facts versionRouteEngines
		err = xml.Unmarshal([]byte(formatted), &facts)
		if err != nil {
			return err
		}

		numRE := len(facts.RE)
		hostname := facts.RE[0].Hostname
		res := make([]RoutingEngine, 0, numRE)

		for i := 0; i < numRE; i++ {
			version := rex.FindStringSubmatch(facts.RE[i].PackageInfo[0].SoftwareVersion[0])
			model := strings.ToUpper(facts.RE[i].Platform)
			res = append(res, RoutingEngine{Model: model, Version: version[1]})
		}

		j.Hostname = hostname
		j.RoutingEngines = numRE
		j.Platform = res
		j.CommitTimeout = 0

		return nil
	}

	var facts versionRouteEngine
	if err := xml.Unmarshal([]byte(formatted), &facts); err != nil {
		return err
	}

	// res := make([]RoutingEngine, 0)
	var res []RoutingEngine
	hostname := facts.Hostname
	version := rex.FindStringSubmatch(facts.PackageInfo[0].SoftwareVersion[0])
	model := strings.ToUpper(facts.Platform)
	res = append(res, RoutingEngine{Model: model, Version: version[1]})

	j.Hostname = hostname
	j.RoutingEngines = 1
	j.Platform = res
	j.CommitTimeout = 0

	return nil
}

// Command (show, execute) on Junos device.
func (j *NetconfObject) netconfCommand(cmd string) (string, error) {
	command := fmt.Sprintf(rpcCommand, cmd)
	reply, err := j.Session.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", err
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
		return "", err
	}

	return output.Config, nil
}
func (j *NetconfObject) netconfCommandXML(cmd string) (string, error) {
	reply, err := j.Session.Exec(netconf.RawMethod(cmd))
	if err != nil {
		return "", err
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
		return "", err
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
		return err
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
		return err
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// netconfCommit commits the configuration.
func (j *NetconfObject) netconfCommit(logMessage string) error {
	var errs commitResults
	reply, err := j.Session.Exec(netconf.RawMethod(fmt.Sprintf(rpcCommit, logMessage)))
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	if reply.Data != "\n<ok/>\n" {
		err = xml.Unmarshal([]byte(reply.Data), &errs)
		if err != nil {
			return err
		}

		if errs.Errors != nil {
			for _, m := range errs.Errors {
				message := fmt.Sprintf("[%s]\n    %s\nError: %s",
					strings.Trim(m.Path, "[\r\n]"),
					strings.Trim(m.Element, "[\r\n]"),
					strings.Trim(m.Message, "[\r\n]"))

				return errors.New(message)
			}
		}
	}

	return nil
}

// Close disconnects our session to the device.
func (j *NetconfObject) Close() error {
	_, err := j.Session.Exec(netconf.RawMethod(rpcClose))
	j.Session.Transport.Close()
	if err != nil {
		return err
	}

	return nil
}
