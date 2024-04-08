package junos

import (
	"encoding/xml"
	"fmt"

	"github.com/jeremmfr/go-netconf/netconf"
)

const (
	rpcCommandText = "<command format=\"text\">%s</command>"

	rpcLoadConfigSetText = "<load-configuration action=\"set\" format=\"text\">" +
		"<configuration-set>%s</configuration-set>" +
		"</load-configuration>"

	rpcCommitConfig = "<commit-configuration>" +
		"<log>%s</log>" +
		"</commit-configuration>"
	rpcCommitConfigConfirmed = "<commit-configuration>" +
		"<log>%s</log>" +
		"<confirmed/><confirm-timeout>%d</confirm-timeout>" +
		"</commit-configuration>"
	rpcCommitConfigCheck = "<commit-configuration>" +
		"<check/>" +
		"</commit-configuration>"

	rpcLockCandidate   = "<lock><target><candidate/></target></lock>"
	rpcUnlockCandidate = "<unlock><target><candidate/></target></unlock>"

	rpcDeleteConfigCandidate = "<delete-config><target><candidate/></target></delete-config>"

	rpcCloseSession = "<close-session/>"

	rpcGetSystemInformation                 = "<get-system-information/>"
	RPCGetChassisInventory                  = `<get-chassis-inventory></get-chassis-inventory>`
	RPCGetInterfaceInformationInterfaceName = "<get-interface-information><interface-name>%s</interface-name></get-interface-information>" //nolint:lll
	RPCGetInterfacesInformationTerse        = `<get-interface-information><terse/></get-interface-information>`
	RPCGetInterfaceInformationTerse         = `<get-interface-information>%s<terse/></get-interface-information>`
	RPCGetRouteAllInformation               = `<get-route-information><all/></get-route-information>`
	RPCGetRouteAllTableInformation          = `<get-route-information><all/><table>%s</table></get-route-information>`
)

type rpcGetSystemInformationReply struct {
	SystemInformation rpcSystemInformation `xml:"system-information"`
}

type rpcSystemInformation struct {
	HardwareModel string `xml:"hardware-model"`
	OsName        string `xml:"os-name"`
	OsVersion     string `xml:"os-version"`
	SerialNumber  string `xml:"serial-number"`
	HostName      string `xml:"host-name"`
	ClusterNode   *bool  `xml:"cluster-node"`
}

func (i rpcSystemInformation) NotCompatibleMsg() string {
	return fmt.Sprintf(" not compatible with Junos device %q", i.HardwareModel)
}

type commandXMLConfig struct {
	Config string `xml:",innerxml"`
}

type commitResults struct {
	XMLName xml.Name           `xml:"commit-results"`
	Errors  []netconf.RPCError `xml:"rpc-error"`
}

type RPCGetPhysicalInterfaceTerseReply struct {
	XMLName           xml.Name `xml:"interface-information"`
	PhysicalInterface []struct {
		Name             string `xml:"name"`
		AdminStatus      string `xml:"admin-status"`
		OperStatus       string `xml:"oper-status"`
		LogicalInterface []struct {
			Name string `xml:"name"`
		} `xml:"logical-interface"`
	} `xml:"physical-interface"`
}

type RPCGetLogicalInterfaceTerseReply struct {
	XMLName          xml.Name `xml:"interface-information"`
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
}

type RPCGetRouteInformationReply struct {
	XMLName    xml.Name `xml:"route-information"`
	RouteTable []struct {
		TableName string `xml:"table-name"`
		Route     []struct {
			Destination string `xml:"rt-destination"`
			Entry       []struct {
				ASPath          *string   `xml:"as-path"`
				CurrentActive   *struct{} `xml:"current-active"`
				LocalPreference *int      `xml:"local-preference"`
				Metric          *int      `xml:"metric"`
				NextHop         []struct {
					SelectedNextHop *struct{} `xml:"selected-next-hop"`
					LocalInterface  *string   `xml:"nh-local-interface"`
					To              *string   `xml:"to"`
					Via             *string   `xml:"via"`
				} `xml:"nh"`
				NextHopType *string `xml:"nh-type"`
				Preference  *int    `xml:"preference"`
				Protocol    *string `xml:"protocol-name"`
			} `xml:"rt-entry"`
		} `xml:"rt"`
	} `xml:"route-table"`
}

type RPCGetChassisInventoryReply struct {
	XMLName xml.Name `xml:"chassis-inventory"`
	Chassis struct {
		RPCGetChassisInventoryReplyComponent
		Module []struct {
			RPCGetChassisInventoryReplyComponent
			SubModule []struct {
				RPCGetChassisInventoryReplyComponent
				SubSubModule []RPCGetChassisInventoryReplyComponent `xml:"chassis-sub-sub-module"`
			} `xml:"chassis-sub-module"`
		} `xml:"chassis-module"`
	} `xml:"chassis"`
}

type RPCGetChassisInventoryReplyComponent struct {
	Name         *string `xml:"name"`
	Version      *string `xml:"version"`
	PartNumber   *string `xml:"part-number"`
	SerialNumber *string `xml:"serial-number"`
	ModelNumber  *string `xml:"model-number"`
	CleiCode     *string `xml:"clei-code"`
	Description  *string `xml:"description"`
}
