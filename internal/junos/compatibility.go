package junos

import "strings"

func (sess *Session) CheckCompatibilitySecurity() bool {
	if strings.HasPrefix(strings.ToLower(sess.SystemInformation.HardwareModel), "srx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(sess.SystemInformation.HardwareModel), "vsrx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(sess.SystemInformation.HardwareModel), "j") {
		return true
	}

	return false
}

func (sess *Session) CheckCompatibilityRouter() bool {
	if strings.HasPrefix(strings.ToLower(sess.SystemInformation.HardwareModel), "mx") {
		return true
	}
	if strings.HasPrefix(strings.ToLower(sess.SystemInformation.HardwareModel), "vmx") {
		return true
	}

	return false
}
