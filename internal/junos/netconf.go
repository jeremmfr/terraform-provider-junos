package junos

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/jeremmfr/go-netconf/netconf"
)

const (
	errorSeverity string = "error"

	XMLStartTagConfigOut = "<configuration-output>"
	XMLEndTagConfigOut   = "</configuration-output>"
)

// gatherFacts gathers basic information about the device.
func (sess *Session) gatherFacts() error {
	// Get info for get-system-information and populate SystemInformation Struct
	val, err := sess.netconf.Exec(netconf.RawMethod(rpcGetSystemInformation))
	if err != nil {
		return fmt.Errorf("executing netconf get-system-information: %w", err)
	}

	if val.Errors != nil {
		var errorsMsg []string
		for _, m := range val.Errors {
			errorsMsg = append(errorsMsg, fmt.Sprintf("%v", m))
		}

		return errors.New(strings.Join(errorsMsg, "\n"))
	}
	var reply rpcGetSystemInformationReply
	if err := xml.Unmarshal([]byte(val.RawReply), &reply); err != nil {
		return fmt.Errorf("unmarshaling xml reply %q of get-system-information: %w", val.RawReply, err)
	}
	sess.SystemInformation = reply.SystemInformation

	return nil
}

// netconfCommand (show, execute) on Junos device.
func (sess *Session) netconfCommand(cmd string) (string, error) {
	command := fmt.Sprintf(rpcCommandText, cmd)
	reply, err := sess.netconf.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("executing netconf command: %w", err)
	}

	errs := make([]string, len(reply.Errors))
	for i, m := range reply.Errors {
		errs[i] = m.Error()
	}
	if len(errs) > 0 {
		return "", errors.New(strings.Join(errs, "\n"))
	}

	if reply.Data == "" || strings.Count(reply.Data, "") <= 2 {
		return EmptyW, errors.New("no output available - please check the syntax of your command")
	}
	var output commandXMLConfig
	if err := xml.Unmarshal([]byte(reply.Data), &output); err != nil {
		return "", fmt.Errorf("unmarshaling xml reply of command: %w", err)
	}

	return output.Config, nil
}

func (sess *Session) netconfCommandXML(cmd string) (string, error) {
	reply, err := sess.netconf.Exec(netconf.RawMethod(cmd))
	if err != nil {
		return "", fmt.Errorf("executing netconf xml command: %w", err)
	}

	errs := make([]string, len(reply.Errors))
	for i, m := range reply.Errors {
		errs[i] = m.Error()
	}
	if len(errs) > 0 {
		return "", errors.New(strings.Join(errs, "\n"))
	}

	return reply.Data, nil
}

func (sess *Session) netconfConfigSet(cmd []string) (string, error) {
	command := fmt.Sprintf(rpcLoadConfigSetText, strings.Join(cmd, "\n"))
	reply, err := sess.netconf.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("executing netconf apply of set/delete command: %w", err)
	}
	// logFile("netconfConfigSet.Reply:" + reply.RawReply)
	message := ""

	for _, m := range reply.Errors {
		message += m.Message
	}

	return message, nil
}

// netConfConfigLock locks the candidate configuration.
func (sess *Session) netconfConfigLock() bool {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcLockCandidate))
	if err != nil {
		return false
	}
	if len(reply.Errors) > 0 {
		return false
	}

	return true
}

func (sess *Session) netconfConfigClear() (errs []error) {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcDeleteConfigCandidate))
	if err != nil {
		return []error{fmt.Errorf("executing netconf config clear: %w", err)}
	}

	for _, m := range reply.Errors {
		errs = append(errs, errors.New("config clear: "+m.Message))
	}

	return errs
}

// Unlock unlocks the candidate configuration.
func (sess *Session) netconfConfigUnlock() (errs []error) {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcUnlockCandidate))
	if err != nil {
		return []error{fmt.Errorf("executing netconf config unlock: %w", err)}
	}

	for _, m := range reply.Errors {
		errs = append(errs, errors.New("config unlock: "+m.Message))
	}

	return errs
}

// netconfCommit commits the configuration.
//
// return potential warnings and/or error.
func (sess *Session) netconfCommit(logMessage string) (_ []error, _ error) {
	reply, err := sess.netconf.Exec(netconf.RawMethod(fmt.Sprintf(rpcCommitConfig, logMessage)))
	if err != nil {
		return nil, fmt.Errorf("executing netconf commit: %w", err)
	}

	return readNetconfCommitReply(reply, "commit-configuration")
}

// netconfCommitConfirmed commits the configuration with confirmed option and confirmed timeout,
// then wait percentage of timeout and send afterwards the confirmation with commit check.
//
// return potential warnings and/or error.
func (sess *Session) netconfCommitConfirmed(ctx context.Context, logMessage string) (warnings []error, _ error) {
	reply, err := sess.netconf.Exec(
		netconf.RawMethod(fmt.Sprintf(rpcCommitConfigConfirmed, logMessage, sess.commitConfirmedTimeout)),
	)
	if err != nil {
		return warnings, fmt.Errorf("executing netconf commit (confirmed %d): %w", sess.commitConfirmedTimeout, err)
	}

	replyWarns, err := readNetconfCommitReply(reply, "commit-configuration(confirmed)")
	warnings = append(warnings, replyWarns...)
	if err != nil {
		return warnings, err
	}

	select {
	case <-ctx.Done():
		return warnings, errors.New("confirmation of commit with 'confirmed' option aborted before done")
	case <-time.After(sess.commitConfirmedWait):
	}

	replyConfirm, err := sess.netconf.Exec(netconf.RawMethod(rpcCommitConfigCheck))
	if err != nil {
		return warnings, fmt.Errorf("executing netconf commit check (to confirm): %w", err)
	}

	replyWarns, err = readNetconfCommitReply(replyConfirm, "commit-configuration(check)")
	warnings = append(warnings, replyWarns...)
	if err != nil {
		return warnings, err
	}

	return warnings, nil
}

func readNetconfCommitReply(reply *netconf.RPCReply, commitType string) (warnings []error, _ error) {
	errs := make([]string, 0, len(reply.Errors))
	for _, m := range reply.Errors {
		if m.Severity == errorSeverity {
			errs = append(errs, m.Error())
		} else {
			warnings = append(warnings, errors.New(m.Error()))
		}
	}
	if len(errs) > 0 {
		return warnings, errors.New(strings.Join(errs, "\n"))
	}

	var result commitResults
	if strings.Contains(reply.Data, "<commit-results>") {
		if err := xml.Unmarshal([]byte(reply.Data), &result); err != nil {
			return warnings, fmt.Errorf("unmarshaling xml reply %q of %s: %w", reply.Data, commitType, err)
		}

		errs = make([]string, 0, len(result.Errors))
		for _, m := range result.Errors {
			if m.Severity == errorSeverity {
				errs = append(errs, m.Error())
			} else {
				warnings = append(warnings, errors.New(m.Error()))
			}
		}
		if len(errs) > 0 {
			return warnings, errors.New(strings.Join(errs, "\n"))
		}
	}

	return warnings, nil
}

// Close disconnects our session to the device.
func (sess *Session) closeNetconf(sleepClosed int) error {
	_, err := sess.netconf.Exec(netconf.RawMethod(rpcCloseSession))
	sess.netconf.Transport.Close()
	if err != nil {
		utils.Sleep(sleepClosed)

		return fmt.Errorf("closing netconf session: %w", err)
	}
	utils.Sleep(sleepClosed)

	return nil
}
