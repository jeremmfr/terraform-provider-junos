package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/jeremmfr/go-netconf/netconf"
)

const (
	errorSeverity string = "error"

	rpcCommand         = "<command format=\"text\">%s</command>"
	rpcConfigStringSet = "<load-configuration action=\"set\" format=\"text\">" +
		"<configuration-set>%s</configuration-set></load-configuration>"
	rpcSystemInfo      = "<get-system-information/>"
	rpcCommit          = "<commit-configuration><log>%s</log></commit-configuration>"
	rpcCandidateLock   = "<lock><target><candidate/></target></lock>"
	rpcCandidateUnlock = "<unlock><target><candidate/></target></unlock>"
	rpcClearCandidate  = "<delete-config><target><candidate/></target></delete-config>"
	rpcClose           = "<close-session/>"

	RPCGetInterfaceInformationInterfaceName = "<get-interface-information><interface-name>%s</interface-name></get-interface-information>" //nolint:lll
	RPCGetInterfacesInformationTerse        = `<get-interface-information><terse/></get-interface-information>`
	RPCGetInterfaceInformationTerse         = `<get-interface-information>%s<terse/></get-interface-information>`
	RPCGetRouteAllInformation               = `<get-route-information><all/></get-route-information>`
	RPCGetRouteAllTableInformation          = `<get-route-information><all/><table>%s</table></get-route-information>`

	XMLStartTagConfigOut = "<configuration-output>"
	XMLEndTagConfigOut   = "</configuration-output>"
)

type sysInfoReply struct {
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

type commandXMLConfig struct {
	Config string `xml:",innerxml"`
}

type commitResults struct {
	XMLName xml.Name           `xml:"commit-results"`
	Errors  []netconf.RPCError `xml:"rpc-error"`
}

type GetPhysicalInterfaceTerseReply struct {
	InterfaceInfo struct {
		PhysicalInterface []struct {
			Name        string `xml:"name"`
			AdminStatus string `xml:"admin-status"`
			OperStatus  string `xml:"oper-status"`
		} `xml:"physical-interface"`
	} `xml:"interface-information"`
}

type GetLogicalInterfaceTerseReply struct {
	InterfaceInfo struct {
		LogicalInterface []struct {
			Name          string `xml:"name"`
			AdminStatus   string `xml:"admin-status"`
			OperStatus    string `xml:"oper-status"`
			AddressFamily []struct {
				Name    string `xml:"address-family-name"`
				Address []struct {
					Local string `xml:"ifa-local"`
				} `xml:"interface-address"`
			} `xml:"address-family"`
		} `xml:"logical-interface"`
	} `xml:"interface-information"`
}

type GetRouteInformationReply struct {
	RouteInfo struct {
		RouteTable []struct {
			TableName string `xml:"table-name"`
			Route     []struct {
				Destination string `xml:"rt-destination"`
				Entry       []struct {
					ASPath          string    `xml:"as-path"`
					CurrentActive   *struct{} `xml:"current-active"`
					LocalPreference int       `xml:"local-preference"`
					Metric          int       `xml:"metric"`
					NextHop         []struct {
						SelectedNextHop *struct{} `xml:"selected-next-hop"`
						LocalInterface  string    `xml:"nh-local-interface"`
						To              string    `xml:"to"`
						Via             string    `xml:"via"`
					} `xml:"nh"`
					NextHopType string `xml:"nh-type"`
					Preference  int    `xml:"preference"`
					Protocol    string `xml:"protocol-name"`
				} `xml:"rt-entry"`
			} `xml:"rt"`
		} `xml:"route-table"`
	} `xml:"route-information"`
}

// gatherFacts gathers basic information about the device.
func (sess *Session) gatherFacts() error {
	// Get info for get-system-information and populate SystemInformation Struct
	val, err := sess.netconf.Exec(netconf.RawMethod(rpcSystemInfo))
	if err != nil {
		return fmt.Errorf("executing netconf get-system-information: %w", err)
	}

	if val.Errors != nil {
		var errorsMsg []string
		for _, m := range val.Errors {
			errorsMsg = append(errorsMsg, fmt.Sprintf("%v", m))
		}

		return fmt.Errorf(strings.Join(errorsMsg, "\n"))
	}
	var reply sysInfoReply
	if err := xml.Unmarshal([]byte(val.RawReply), &reply); err != nil {
		return fmt.Errorf("unmarshaling xml reply %q of get-system-information: %w", val.RawReply, err)
	}
	sess.SystemInformation = reply.SystemInformation

	return nil
}

// netconfCommand (show, execute) on Junos device.
func (sess *Session) netconfCommand(cmd string) (string, error) {
	command := fmt.Sprintf(rpcCommand, cmd)
	reply, err := sess.netconf.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("executing netconf command: %w", err)
	}
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Error())
		}
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
	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Error())
		}
	}

	return reply.Data, nil
}

func (sess *Session) netconfConfigSet(cmd []string) (string, error) {
	command := fmt.Sprintf(rpcConfigStringSet, strings.Join(cmd, "\n"))
	reply, err := sess.netconf.Exec(netconf.RawMethod(command))
	if err != nil {
		return "", fmt.Errorf("executing netconf apply of set/delete command: %w", err)
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
func (sess *Session) netconfConfigLock() bool {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcCandidateLock))
	if err != nil {
		return false
	}
	if reply.Errors != nil {
		return false
	}

	return true
}

func (sess *Session) netconfConfigClear() []error {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcClearCandidate))
	if err != nil {
		return []error{fmt.Errorf("executing netconf config clear: %w", err)}
	}
	if reply.Errors != nil {
		errs := make([]error, 0)
		for _, m := range reply.Errors {
			errs = append(errs, errors.New("config clear: "+m.Message))
		}

		return errs
	}

	return []error{}
}

// Unlock unlocks the candidate configuration.
func (sess *Session) netconfConfigUnlock() []error {
	reply, err := sess.netconf.Exec(netconf.RawMethod(rpcCandidateUnlock))
	if err != nil {
		return []error{fmt.Errorf("executing netconf config unlock: %w", err)}
	}
	if reply.Errors != nil {
		errs := make([]error, 0)
		for _, m := range reply.Errors {
			errs = append(errs, errors.New("config unlock: "+m.Message))
		}

		return errs
	}

	return []error{}
}

// netconfCommit commits the configuration.
func (sess *Session) netconfCommit(logMessage string) (_warn []error, _err error) {
	reply, err := sess.netconf.Exec(netconf.RawMethod(fmt.Sprintf(rpcCommit, logMessage)))
	if err != nil {
		return []error{}, fmt.Errorf("executing netconf commit: %w", err)
	}

	if reply.Errors != nil {
		warnings := make([]error, 0)
		for _, m := range reply.Errors {
			if m.Severity == errorSeverity {
				return warnings, errors.New(m.Error())
			}
			warnings = append(warnings, errors.New(m.Error()))
		}

		return warnings, nil
	}

	var errs commitResults
	if strings.Contains(reply.Data, "<commit-results>") {
		err = xml.Unmarshal([]byte(reply.Data), &errs)
		if err != nil {
			return []error{}, fmt.Errorf("unmarshaling xml reply %q of commit-configuration: %w", reply.Data, err)
		}

		if errs.Errors != nil {
			warnings := make([]error, 0)
			for _, m := range errs.Errors {
				if m.Severity == errorSeverity {
					return []error{}, errors.New(m.Error())
				}
				warnings = append(warnings, errors.New(m.Error()))
			}

			return warnings, nil
		}
	}

	return []error{}, nil
}

// Close disconnects our session to the device.
func (sess *Session) closeNetconf(sleepClosed int) error {
	_, err := sess.netconf.Exec(netconf.RawMethod(rpcClose))
	sess.netconf.Transport.Close()
	if err != nil {
		utils.Sleep(sleepClosed)

		return fmt.Errorf("closing netconf session: %w", err)
	}
	utils.Sleep(sleepClosed)

	return nil
}
