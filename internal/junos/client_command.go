package junos

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"
)

func (clt *Client) Command(cmd string, junSess *Session) (string, error) {
	read, err := junSess.netconfCommand(cmd)
	clt.logFile(fmt.Sprintf("[command] cmd: %q", cmd))
	clt.logFile(fmt.Sprintf("[command] read: %q", read))
	utils.SleepShort(clt.sleepShort)
	if err != nil && read != EmptyW {
		clt.logFile(fmt.Sprintf("[command] err: %q", err))

		return "", err
	}

	return read, nil
}

func (clt *Client) CommandXML(cmd string, junSess *Session) (string, error) {
	read, err := junSess.netconfCommandXML(cmd)
	clt.logFile(fmt.Sprintf("[commandXML] cmd: %q", cmd))
	clt.logFile(fmt.Sprintf("[commandXML] read: %q", read))
	utils.SleepShort(clt.sleepShort)
	if err != nil {
		clt.logFile(fmt.Sprintf("[commandXML] err: %q", err))

		return "", err
	}

	return read, nil
}

func (clt *Client) ConfigSet(cmd []string, junSess *Session) error {
	if junSess != nil {
		message, err := junSess.netconfConfigSet(cmd)
		utils.SleepShort(clt.sleepShort)
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

func (clt *Client) CommitConf(logMessage string, junSess *Session) (_warnings []error, _err error) {
	clt.logFile(fmt.Sprintf("[commitConf] commit %q", logMessage))
	warns, err := junSess.netconfCommit(logMessage)
	utils.SleepShort(clt.sleepShort)
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

func (clt *Client) ConfigLock(ctx context.Context, junSess *Session) error {
	for {
		select {
		case <-ctx.Done():
			clt.logFile("[configLock] aborted")

			return fmt.Errorf("candidate configuration lock attempt aborted")
		default:
			if junSess.netconfConfigLock() {
				clt.logFile("[configLock] locked")
				utils.SleepShort(clt.sleepShort)

				return nil
			}
			clt.logFile("[configLock] sleep to wait the lock")
			utils.Sleep(clt.sleepLock)
		}
	}
}

func (clt *Client) ConfigClear(junSess *Session) (errs []error) {
	errs = append(errs, junSess.netconfConfigClear()...)
	utils.SleepShort(clt.sleepShort)
	clt.logFile("[configClear] config clear")

	errs = append(errs, junSess.netconfConfigUnlock()...)
	utils.SleepShort(clt.sleepShort)
	clt.logFile("[configClear] config unlock")

	return
}
