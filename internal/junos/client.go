package junos

import (
	"errors"
)

const directoryPermission = 0o755

// Client information to connect on Junos Device and more.
type Client struct {
	junosIP                         string
	junosPort                       int
	junosUserName                   string
	junosPassword                   string
	junosSSHKeyPEM                  string
	junosSSHKeyFile                 string
	junosSSHKeyPass                 string
	groupIntDel                     string
	sleepShort                      int
	sleepLock                       int
	junosCommitConfirmed            int
	junosCommitConfirmedWaitPercent int
	sleepSSHClosed                  int
	junosSSHCiphers                 []string
	junosSSHTimeoutToEstab          int
	junosSSHRetryToEstab            int
	filePermission                  int64
	logFileDst                      string
	fakeCreateSetFile               string
	fakeUpdateAlso                  bool
	fakeDeleteAlso                  bool
}

func NewClient(ip string) *Client {
	return &Client{
		junosIP:                         ip,
		junosPort:                       830,
		junosUserName:                   "netconf",
		junosPassword:                   "",
		junosSSHKeyPEM:                  "",
		junosSSHKeyFile:                 "",
		junosSSHKeyPass:                 "",
		groupIntDel:                     "",
		sleepShort:                      100,
		sleepLock:                       10,
		junosCommitConfirmed:            0,
		junosCommitConfirmedWaitPercent: 90,
		sleepSSHClosed:                  0,
		junosSSHCiphers:                 DefaultSSHCiphers(),
		junosSSHTimeoutToEstab:          0,
		junosSSHRetryToEstab:            1,
		filePermission:                  0o644,
		logFileDst:                      "",
		fakeCreateSetFile:               "",
		fakeUpdateAlso:                  false,
		fakeDeleteAlso:                  false,
	}
}

func (clt *Client) WithPort(port int) *Client {
	clt.junosPort = port

	return clt
}

func (clt *Client) WithUserName(userName string) *Client {
	clt.junosUserName = userName

	return clt
}

func (clt *Client) WithPassword(password string) *Client {
	clt.junosPassword = password

	return clt
}

func (clt *Client) WithSSHKeyPEM(sshKeyPEM string) *Client {
	clt.junosSSHKeyPEM = sshKeyPEM

	return clt
}

func (clt *Client) WithSSHKeyFile(sshKeyFile string) *Client {
	clt.junosSSHKeyFile = sshKeyFile

	return clt
}

func (clt *Client) WithSSHKeyPassphrase(sshKeyPass string) *Client {
	clt.junosSSHKeyPass = sshKeyPass

	return clt
}

func (clt *Client) WithGroupInterfaceDelete(groupIntDel string) *Client {
	clt.groupIntDel = groupIntDel

	return clt
}

func (clt *Client) WithSleepShort(sleep int) *Client {
	clt.sleepShort = sleep

	return clt
}

func (clt *Client) WithSleepLock(sleep int) *Client {
	clt.sleepLock = sleep

	return clt
}

func (clt *Client) WithCommitConfirmed(timeout int) (*Client, error) {
	if timeout < 1 || timeout > 65535 {
		return clt, errors.New("bad value for timeout of commit confirmed")
	}
	clt.junosCommitConfirmed = timeout

	return clt, nil
}

func (clt *Client) WithCommitConfirmedWaitPercent(percent int) (*Client, error) {
	if percent < 0 || percent > 99 {
		return clt, errors.New("bad value for wait percent of timeout before the commit confirm")
	}
	clt.junosCommitConfirmedWaitPercent = percent

	return clt, nil
}

func (clt *Client) WithSleepSSHClosed(sleep int) *Client {
	clt.sleepSSHClosed = sleep

	return clt
}

func (clt *Client) WithSSHCiphers(ciphers []string) *Client {
	clt.junosSSHCiphers = ciphers

	return clt
}

func (clt *Client) WithSSHTimeoutToEstablish(timeout int) *Client {
	clt.junosSSHTimeoutToEstab = timeout

	return clt
}

func (clt *Client) WithSSHRetryToEstablish(retry int) (*Client, error) {
	if retry < 1 || retry > 10 {
		return clt, errors.New("bad value for number of retry to establishing SSH connection")
	}
	clt.junosSSHRetryToEstab = retry

	return clt, nil
}

func (clt *Client) WithFilePermission(perm int64) (*Client, error) {
	if perm > 0o777 || perm < 0 {
		return clt, errors.New("bad value for file permision, must be three octal digits")
	}
	clt.filePermission = perm

	return clt, nil
}

func (clt *Client) WithDebugLogFile(file string) *Client {
	clt.logFileDst = file

	return clt
}

func (clt *Client) WithFakeCreateSetFile(file string) *Client {
	clt.fakeCreateSetFile = file

	return clt
}

func (clt *Client) WithFakeUpdateAlso() *Client {
	clt.fakeUpdateAlso = true

	return clt
}

func (clt *Client) WithFakeDeleteAlso() *Client {
	clt.fakeDeleteAlso = true

	return clt
}

func (clt *Client) FakeCreateSetFile() bool {
	return clt.fakeCreateSetFile != ""
}

func (clt *Client) FakeUpdateAlso() bool {
	return clt.fakeUpdateAlso
}

func (clt *Client) FakeDeleteAlso() bool {
	return clt.fakeDeleteAlso
}

func (clt *Client) GroupInterfaceDelete() string {
	return clt.groupIntDel
}

func DefaultSSHCiphers() []string {
	return []string{
		"aes128-gcm@openssh.com", "aes256-gcm@openssh.com",
		"chacha20-poly1305@openssh.com",
		"aes128-ctr", "aes192-ctr", "aes256-ctr",
	}
}
