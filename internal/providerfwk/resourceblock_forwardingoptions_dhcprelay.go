package providerfwk

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude struct {
	CircuitType                          types.Bool   `tfsdk:"circuit_type"`
	ClientID                             types.Bool   `tfsdk:"client_id"`
	ClientIDExcludeHeaders               types.Bool   `tfsdk:"client_id_exclude_headers"`
	ClientIDUseAutomaticASCIIHexEncoding types.Bool   `tfsdk:"client_id_use_automatic_ascii_hex_encoding"`
	Delimiter                            types.String `tfsdk:"delimiter"`
	DomainName                           types.String `tfsdk:"domain_name"`
	InterfaceDescription                 types.String `tfsdk:"interface_description"`
	InterfaceName                        types.Bool   `tfsdk:"interface_name"`
	MACAddress                           types.Bool   `tfsdk:"mac_address"`
	Option60                             types.Bool   `tfsdk:"option_60"`
	Option82                             types.Bool   `tfsdk:"option_82"`
	Option82CircuitID                    types.Bool   `tfsdk:"option_82_circuit_id"`
	Option82RemoteID                     types.Bool   `tfsdk:"option_82_remote_id"`
	RelayAgentInterfaceID                types.Bool   `tfsdk:"relay_agent_interface_id"`
	RelayAgentRemoteID                   types.Bool   `tfsdk:"relay_agent_remote_id"`
	RelayAgentSubscriberID               types.Bool   `tfsdk:"relay_agent_subscriber_id"`
	RoutingInstanceName                  types.Bool   `tfsdk:"routing_instance_name"`
	UserPrefix                           types.String `tfsdk:"user_prefix"`
	VlanTags                             types.Bool   `tfsdk:"vlan_tags"`
}

func (block *forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"circuit_type": schema.BoolAttribute{
			Optional:    true,
			Description: "Include circuit type.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"client_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include client ID.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"client_id_exclude_headers": schema.BoolAttribute{
			Optional:    true,
			Description: "Exclude all the headers.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"client_id_use_automatic_ascii_hex_encoding": schema.BoolAttribute{
			Optional:    true,
			Description: "Use automatic ascii hex username encoding.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"delimiter": schema.StringAttribute{
			Optional:    true,
			Description: "Change delimiter/separator character.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"domain_name": schema.StringAttribute{
			Optional:    true,
			Description: "Add domain name.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"interface_description": schema.StringAttribute{
			Optional:    true,
			Description: "Include interface description.",
			Validators: []validator.String{
				stringvalidator.OneOf("device", "logical"),
			},
		},
		"interface_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Include interface name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"mac_address": schema.BoolAttribute{
			Optional:    true,
			Description: "Include MAC address.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"option_60": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option 60.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"option_82": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option 82.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"option_82_circuit_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option 82 circuit-id (sub option 1).",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"option_82_remote_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option 82 remote-id (sub option 2).",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"relay_agent_interface_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include the relay agent interface ID.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"relay_agent_remote_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include the relay agent remote ID.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"relay_agent_subscriber_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Include the relay agent subscriber ID.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"routing_instance_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Include routing instance name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"user_prefix": schema.StringAttribute{
			Optional:    true,
			Description: "Add user defined prefix.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"vlan_tags": schema.BoolAttribute{
			Optional:    true,
			Description: "Include the vlan tag(s).",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude) configSet(
	setPrefix, version string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "authentication username-include "

	if block.CircuitType.ValueBool() {
		configSet = append(configSet, setPrefix+"circuit-type")
	}
	if block.ClientID.ValueBool() {
		configSet = append(configSet, setPrefix+"client-id")
		if block.ClientIDExcludeHeaders.ValueBool() {
			configSet = append(configSet, setPrefix+"client-id exclude-headers")
		}
		if block.ClientIDUseAutomaticASCIIHexEncoding.ValueBool() {
			configSet = append(configSet, setPrefix+"client-id use-automatic-ascii-hex-encoding")
		}
	} else {
		if block.ClientIDExcludeHeaders.ValueBool() {
			return configSet,
				path.Root("authentication_username_include").AtName("client_id_exclude_headers"),
				errors.New("client_id must be specified with client_id_exclude_headers" +
					" in authentication_username_include block")
		}
		if block.ClientIDUseAutomaticASCIIHexEncoding.ValueBool() {
			return configSet,
				path.Root("authentication_username_include").AtName("client_id_use_automatic_ascii_hex_encoding"),
				errors.New("client_id must be specified with client_id_use_automatic_ascii_hex_encoding" +
					" in authentication_username_include block")
		}
	}
	if v := block.Delimiter.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"delimiter \""+v+"\"")
	}
	if v := block.DomainName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"domain-name \""+v+"\"")
	}
	if v := block.InterfaceDescription.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"interface-description "+v)
	}
	if block.InterfaceName.ValueBool() {
		configSet = append(configSet, setPrefix+"interface-name")
	}
	if block.MACAddress.ValueBool() {
		configSet = append(configSet, setPrefix+"mac-address")
	}
	if block.Option60.ValueBool() {
		if version == "v6" {
			return configSet,
				path.Root("authentication_username_include").AtName("option_60"),
				errors.New("option_60 cannot be configured when version = v6" +
					" in authentication_username_include block")
		}

		configSet = append(configSet, setPrefix+"option-60")
	}
	if block.Option82.ValueBool() {
		if version == "v6" {
			return configSet,
				path.Root("authentication_username_include").AtName("option_82"),
				errors.New("option_82 cannot be configured when version = v6" +
					" in authentication_username_include block")
		}

		configSet = append(configSet, setPrefix+"option-82")
		if block.Option82CircuitID.ValueBool() {
			configSet = append(configSet, setPrefix+"option-82 circuit-id")
		}
		if block.Option82RemoteID.ValueBool() {
			configSet = append(configSet, setPrefix+"option-82 remote-id")
		}
	} else {
		if block.Option82CircuitID.ValueBool() {
			return configSet,
				path.Root("authentication_username_include").AtName("option_82_circuit_id"),
				errors.New("option_82 must be specified with option_82_circuit_id" +
					" in authentication_username_include block")
		}
		if block.Option82RemoteID.ValueBool() {
			return configSet,
				path.Root("authentication_username_include").AtName("option_82_remote_id"),
				errors.New("option_82 must be specified with option_82_remote_id" +
					" in authentication_username_include block")
		}
	}
	if block.RelayAgentInterfaceID.ValueBool() {
		if version == "v4" {
			return configSet,
				path.Root("authentication_username_include").AtName("relay_agent_interface_id"),
				errors.New("relay_agent_interface_id cannot be configured when version = v4" +
					" in authentication_username_include block")
		}

		configSet = append(configSet, setPrefix+"relay-agent-interface-id")
	}
	if block.RelayAgentRemoteID.ValueBool() {
		if version == "v4" {
			return configSet,
				path.Root("authentication_username_include").AtName("relay_agent_remote_id"),
				errors.New("relay_agent_remote_id cannot be configured when version = v4" +
					" in authentication_username_include block")
		}

		configSet = append(configSet, setPrefix+"relay-agent-remote-id")
	}
	if block.RelayAgentSubscriberID.ValueBool() {
		if version == "v4" {
			return configSet,
				path.Root("authentication_username_include").AtName("relay_agent_subscriber_id"),
				errors.New("relay_agent_subscriber_id cannot be configured when version = v4" +
					" in authentication_username_include block")
		}

		configSet = append(configSet, setPrefix+"relay-agent-subscriber-id")
	}
	if block.RoutingInstanceName.ValueBool() {
		configSet = append(configSet, setPrefix+"routing-instance-name")
	}
	if v := block.UserPrefix.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-prefix \""+v+"\"")
	}
	if block.VlanTags.ValueBool() {
		configSet = append(configSet, setPrefix+"vlan-tags")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude) read(itemTrim string) {
	switch {
	case itemTrim == "circuit-type":
		block.CircuitType = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "client-id"):
		block.ClientID = types.BoolValue(true)
		switch {
		case itemTrim == " exclude-headers":
			block.ClientIDExcludeHeaders = types.BoolValue(true)
		case itemTrim == " use-automatic-ascii-hex-encoding":
			block.ClientIDUseAutomaticASCIIHexEncoding = types.BoolValue(true)
		}
	case balt.CutPrefixInString(&itemTrim, "delimiter "):
		block.Delimiter = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "domain-name "):
		block.DomainName = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "interface-description "):
		block.InterfaceDescription = types.StringValue(itemTrim)
	case itemTrim == "interface-name":
		block.InterfaceName = types.BoolValue(true)
	case itemTrim == "mac-address":
		block.MACAddress = types.BoolValue(true)
	case itemTrim == "option-60":
		block.Option60 = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "option-82"):
		block.Option82 = types.BoolValue(true)
		switch {
		case itemTrim == " circuit-id":
			block.Option82CircuitID = types.BoolValue(true)
		case itemTrim == " remote-id":
			block.Option82RemoteID = types.BoolValue(true)
		}
	case itemTrim == "relay-agent-interface-id":
		block.RelayAgentInterfaceID = types.BoolValue(true)
	case itemTrim == "relay-agent-remote-id":
		block.RelayAgentRemoteID = types.BoolValue(true)
	case itemTrim == "relay-agent-subscriber-id":
		block.RelayAgentSubscriberID = types.BoolValue(true)
	case itemTrim == "routing-instance-name":
		block.RoutingInstanceName = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user-prefix "):
		block.UserPrefix = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "vlan-tags":
		block.VlanTags = types.BoolValue(true)
	}
}

type forwardingoptionsDhcprelayBlockLeaseTimeValidation struct {
	LeaseTimeThreshold  types.Int64 `tfsdk:"lease_time_threshold"`
	ViolationActionDrop types.Bool  `tfsdk:"violation_action_drop"`
}

func (forwardingoptionsDhcprelayBlockLeaseTimeValidation) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"lease_time_threshold": schema.Int64Attribute{
			Optional:    true,
			Description: "Threshold for lease time violation seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(60, 2147483647),
			},
		},
		"violation_action_drop": schema.BoolAttribute{
			Optional:    true,
			Description: "Lease time validation violation action is drop.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockLeaseTimeValidation) configSet(setPrefix string) []string {
	setPrefix += "lease-time-validation "

	configSet := []string{
		setPrefix,
	}

	if !block.LeaseTimeThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"lease-time-threshold "+
			utils.ConvI64toa(block.LeaseTimeThreshold.ValueInt64()))
	}
	if block.ViolationActionDrop.ValueBool() {
		configSet = append(configSet, setPrefix+"violation-action drop")
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockLeaseTimeValidation) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "lease-time-threshold "):
		block.LeaseTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "violation-action drop":
		block.ViolationActionDrop = types.BoolValue(true)
	}

	return err
}

type forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd struct {
	DetectionTimeThreshold    types.Int64  `tfsdk:"detection_time_threshold"`
	HolddownInterval          types.Int64  `tfsdk:"holddown_interval"`
	MinimumInterval           types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval    types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                types.Int64  `tfsdk:"multiplier"`
	NoAdaptation              types.Bool   `tfsdk:"no_adaptation"`
	SessionMode               types.String `tfsdk:"session_mode"`
	TransmitIntervalMinimum   types.Int64  `tfsdk:"transmit_interval_minimum"`
	TransmitIntervalThreshold types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                   types.String `tfsdk:"version"`
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"detection_time_threshold": schema.Int64Attribute{
			Optional:    true,
			Description: "High detection-time triggering a trap (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"holddown_interval": schema.Int64Attribute{
			Optional:    true,
			Description: "Time to hold the session-UP notification to the client (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(0, 255000),
			},
		},
		"minimum_interval": schema.Int64Attribute{
			Optional:    true,
			Description: "Minimum transmit and receive interval (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(30000, 255000),
			},
		},
		"minimum_receive_interval": schema.Int64Attribute{
			Optional:    true,
			Description: "Minimum receive interval (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(30000, 255000),
			},
		},
		"multiplier": schema.Int64Attribute{
			Optional:    true,
			Description: "Detection time multiplier.",
			Validators: []validator.Int64{
				int64validator.Between(1, 255),
			},
		},
		"no_adaptation": schema.BoolAttribute{
			Optional:    true,
			Description: "Disable adaptation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"session_mode": schema.StringAttribute{
			Optional:    true,
			Description: "BFD single-hop or multihop session-mode.",
			Validators: []validator.String{
				stringvalidator.OneOf("automatic", "multihop", "single-hop"),
			},
		},
		"transmit_interval_minimum": schema.Int64Attribute{
			Optional:    true,
			Description: "Minimum transmit interval (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(30000, 255000),
			},
		},
		"transmit_interval_threshold": schema.Int64Attribute{
			Optional:    true,
			Description: "High transmit interval triggering a trap (milliseconds).",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"version": schema.StringAttribute{
			Optional:    true,
			Description: "BFD protocol version number.",
			Validators: []validator.String{
				stringvalidator.OneOf("0", "1", "automatic"),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "liveness-detection method bfd "

	if !block.DetectionTimeThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"detection-time threshold "+
			utils.ConvI64toa(block.DetectionTimeThreshold.ValueInt64()))
	}
	if !block.HolddownInterval.IsNull() {
		configSet = append(configSet, setPrefix+"holddown-interval "+
			utils.ConvI64toa(block.HolddownInterval.ValueInt64()))
	}
	if !block.MinimumInterval.IsNull() {
		configSet = append(configSet, setPrefix+"minimum-interval "+
			utils.ConvI64toa(block.MinimumInterval.ValueInt64()))
	}
	if !block.MinimumReceiveInterval.IsNull() {
		configSet = append(configSet, setPrefix+"minimum-receive-interval "+
			utils.ConvI64toa(block.MinimumReceiveInterval.ValueInt64()))
	}
	if !block.Multiplier.IsNull() {
		configSet = append(configSet, setPrefix+"multiplier "+
			utils.ConvI64toa(block.Multiplier.ValueInt64()))
	}
	if block.NoAdaptation.ValueBool() {
		configSet = append(configSet, setPrefix+"no-adaptation")
	}
	if v := block.SessionMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"session-mode "+v)
	}
	if !block.TransmitIntervalMinimum.IsNull() {
		configSet = append(configSet, setPrefix+"transmit-interval minimum-interval "+
			utils.ConvI64toa(block.TransmitIntervalMinimum.ValueInt64()))
	}
	if !block.TransmitIntervalThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"transmit-interval threshold "+
			utils.ConvI64toa(block.TransmitIntervalThreshold.ValueInt64()))
	}
	if v := block.Version.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"version "+v)
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		block.DetectionTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
		block.HolddownInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		block.MinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		block.MinimumReceiveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		block.Multiplier, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "no-adaptation":
		block.NoAdaptation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "session-mode "):
		block.SessionMode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		block.TransmitIntervalMinimum, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		block.TransmitIntervalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "version "):
		block.Version = types.StringValue(itemTrim)
	}

	return err
}

type forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2 struct {
	MaxConsecutiveRetries types.Int64 `tfsdk:"max_consecutive_retries"`
	TransmitInterval      types.Int64 `tfsdk:"transmit_interval"`
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"max_consecutive_retries": schema.Int64Attribute{
			Optional:    true,
			Description: "Retry attempts.",
			Validators: []validator.Int64{
				int64validator.Between(3, 6),
			},
		},
		"transmit_interval": schema.Int64Attribute{
			Optional:    true,
			Description: "Transmit interval for address resolution (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(300, 1800),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "liveness-detection method layer2-liveness-detection "

	if !block.MaxConsecutiveRetries.IsNull() {
		configSet = append(configSet,
			setPrefix+"max-consecutive-retries "+
				utils.ConvI64toa(block.MaxConsecutiveRetries.ValueInt64()))
	}
	if !block.TransmitInterval.IsNull() {
		configSet = append(configSet,
			setPrefix+"transmit-interval "+
				utils.ConvI64toa(block.TransmitInterval.ValueInt64()))
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "max-consecutive-retries "):
		block.MaxConsecutiveRetries, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval "):
		block.TransmitInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

type forwardingoptionsDhcprelayBlockOverridesV4 struct {
	AllowNoEndOption             types.Bool   `tfsdk:"allow_no_end_option"`
	AllowSnoopedClients          types.Bool   `tfsdk:"allow_snooped_clients"`
	NoAllowSnoopedClients        types.Bool   `tfsdk:"no_allow_snooped_clients"`
	AlwaysWriteGiaddr            types.Bool   `tfsdk:"always_write_giaddr"`
	AlwaysWriteOption82          types.Bool   `tfsdk:"always_write_option_82"`
	AsymmetricLeaseTime          types.Int64  `tfsdk:"asymmetric_lease_time"`
	BootpSupport                 types.Bool   `tfsdk:"bootp_support"`
	ClientDiscoverMatch          types.String `tfsdk:"client_discover_match"`
	DelayAuthentication          types.Bool   `tfsdk:"delay_authentication"`
	DeleteBindingOnRenegotiation types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
	DisableRelay                 types.Bool   `tfsdk:"disable_relay"`
	DualStack                    types.String `tfsdk:"dual_stack"`
	InterfaceClientLimit         types.Int64  `tfsdk:"interface_client_limit"`
	Layer2UnicastReplies         types.Bool   `tfsdk:"layer2_unicast_replies"`
	NoBindOnRequest              types.Bool   `tfsdk:"no_bind_on_request"`
	NoUnicastReplies             types.Bool   `tfsdk:"no_unicast_replies"`
	ProxyMode                    types.Bool   `tfsdk:"proxy_mode"`
	RelaySource                  types.String `tfsdk:"relay_source"`
	ReplaceIPSourceWithGiaddr    types.Bool   `tfsdk:"replace_ip_source_with_giaddr"`
	SendReleaseOnDelete          types.Bool   `tfsdk:"send_release_on_delete"`
	TrustOption82                types.Bool   `tfsdk:"trust_option_82"`
	UserDefinedOption82          types.String `tfsdk:"user_defined_option_82"`
}

func (block *forwardingoptionsDhcprelayBlockOverridesV4) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *forwardingoptionsDhcprelayBlockOverridesV4) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockOverridesV4) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"allow_no_end_option": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow packets without end-of-option.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"allow_snooped_clients": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow client creation from snooped PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_allow_snooped_clients": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't allow client creation from snooped PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"always_write_giaddr": schema.BoolAttribute{
			Optional:    true,
			Description: "Overwrite existing 'giaddr' field, when present.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"always_write_option_82": schema.BoolAttribute{
			Optional:    true,
			Description: "Overwrite existing value of option 82, when present.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"asymmetric_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced lease time for the client. In seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"bootp_support": schema.BoolAttribute{
			Optional:    true,
			Description: "Allows relay of bootp req and reply.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"client_discover_match": schema.StringAttribute{
			Optional:    true,
			Description: "Use secondary match criteria for DISCOVER PDU.",
			Validators: []validator.String{
				stringvalidator.OneOf("incoming-interface", "option60-and-option82"),
			},
		},
		"delay_authentication": schema.BoolAttribute{
			Optional:    true,
			Description: "Delay subscriber authentication in DHCP protocol processing until request packet.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"delete_binding_on_renegotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Delete binding on renegotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"disable_relay": schema.BoolAttribute{
			Optional:    true,
			Description: "Disable DHCP relay processing.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"dual_stack": schema.StringAttribute{
			Optional:    true,
			Description: "Dual stack group to use.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"interface_client_limit": schema.Int64Attribute{
			Optional:    true,
			Description: "Limit the number of clients allowed on an interface.",
			Validators: []validator.Int64{
				int64validator.Between(1, 500000),
			},
		},
		"layer2_unicast_replies": schema.BoolAttribute{
			Optional:    true,
			Description: "Do not broadcast client responses.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_bind_on_request": schema.BoolAttribute{
			Optional:    true,
			Description: "Do not bind if stray DHCP request is received.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_unicast_replies": schema.BoolAttribute{
			Optional:    true,
			Description: "Overwrite unicast bit in incoming packet, when present.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"proxy_mode": schema.BoolAttribute{
			Optional:    true,
			Description: "Put the relay in proxy mode.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"relay_source": schema.StringAttribute{
			Optional:    true,
			Description: "Interface for relay source.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
				tfvalidator.String1DotCount(),
			},
		},
		"replace_ip_source_with_giaddr": schema.BoolAttribute{
			Optional:    true,
			Description: "Replace IP source address in request and release packets.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"send_release_on_delete": schema.BoolAttribute{
			Optional:    true,
			Description: "Always send RELEASE to the server when a binding is deleted.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"trust_option_82": schema.BoolAttribute{
			Optional:    true,
			Description: "Trust options-82 option.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"user_defined_option_82": schema.StringAttribute{
			Optional:    true,
			Description: "Set user defined description for option-82.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 252),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockOverridesV4) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if block.AllowNoEndOption.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-no-end-option")
	}
	if block.AllowSnoopedClients.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-snooped-clients")
		if block.NoAllowSnoopedClients.ValueBool() {
			return configSet,
				path.Root("overrides_v4").AtName("no_allow_snooped_clients"),
				errors.New("allow_snooped_clients and no_allow_snooped_clients cannot be configured together" +
					" in overrides_v4 block")
		}
	} else if block.NoAllowSnoopedClients.ValueBool() {
		configSet = append(configSet, setPrefix+"no-allow-snooped-clients")
	}
	if block.AlwaysWriteGiaddr.ValueBool() {
		configSet = append(configSet, setPrefix+"always-write-giaddr")
	}
	if block.AlwaysWriteOption82.ValueBool() {
		configSet = append(configSet, setPrefix+"always-write-option-82")
	}
	if !block.AsymmetricLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+
			utils.ConvI64toa(block.AsymmetricLeaseTime.ValueInt64()))
	}
	if block.BootpSupport.ValueBool() {
		configSet = append(configSet, setPrefix+"bootp-support")
	}
	if v := block.ClientDiscoverMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-discover-match "+v)
	}
	if block.DelayAuthentication.ValueBool() {
		configSet = append(configSet, setPrefix+"delay-authentication")
	}
	if block.DeleteBindingOnRenegotiation.ValueBool() {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if block.DisableRelay.ValueBool() {
		configSet = append(configSet, setPrefix+"disable-relay")
	}
	if v := block.DualStack.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if !block.InterfaceClientLimit.IsNull() {
		configSet = append(configSet, setPrefix+"interface-client-limit "+
			utils.ConvI64toa(block.InterfaceClientLimit.ValueInt64()))
	}
	if block.Layer2UnicastReplies.ValueBool() {
		configSet = append(configSet, setPrefix+"layer2-unicast-replies")
	}
	if block.NoBindOnRequest.ValueBool() {
		configSet = append(configSet, setPrefix+"no-bind-on-request")
	}
	if block.NoUnicastReplies.ValueBool() {
		if block.Layer2UnicastReplies.ValueBool() {
			return configSet,
				path.Root("overrides_v4").AtName("no_unicast_replies"),
				errors.New("layer2_unicast_replies and no_unicast_replies cannot be configured together" +
					" in overrides_v4 block")
		}
		configSet = append(configSet, setPrefix+"no-unicast-replies")
	}
	if block.ProxyMode.ValueBool() {
		configSet = append(configSet, setPrefix+"proxy-mode")
	}
	if v := block.RelaySource.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"relay-source "+v)
	}
	if block.ReplaceIPSourceWithGiaddr.ValueBool() {
		configSet = append(configSet, setPrefix+"replace-ip-source-with giaddr")
	}
	if block.SendReleaseOnDelete.ValueBool() {
		configSet = append(configSet, setPrefix+"send-release-on-delete")
	}
	if block.TrustOption82.ValueBool() {
		configSet = append(configSet, setPrefix+"trust-option-82")
	}
	if v := block.UserDefinedOption82.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-defined-option-82 \""+v+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockOverridesV4) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "allow-no-end-option":
		block.AllowNoEndOption = types.BoolValue(true)
	case itemTrim == "allow-snooped-clients":
		block.AllowSnoopedClients = types.BoolValue(true)
	case itemTrim == "always-write-giaddr":
		block.AlwaysWriteGiaddr = types.BoolValue(true)
	case itemTrim == "always-write-option-82":
		block.AlwaysWriteOption82 = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		block.AsymmetricLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "bootp-support":
		block.BootpSupport = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "client-discover-match "):
		block.ClientDiscoverMatch = types.StringValue(itemTrim)
	case itemTrim == "delay-authentication":
		block.DelayAuthentication = types.BoolValue(true)
	case itemTrim == "delete-binding-on-renegotiation":
		block.DeleteBindingOnRenegotiation = types.BoolValue(true)
	case itemTrim == "disable-relay":
		block.DisableRelay = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		block.DualStack = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		block.InterfaceClientLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "layer2-unicast-replies":
		block.Layer2UnicastReplies = types.BoolValue(true)
	case itemTrim == "no-allow-snooped-clients":
		block.NoAllowSnoopedClients = types.BoolValue(true)
	case itemTrim == "no-bind-on-request":
		block.NoBindOnRequest = types.BoolValue(true)
	case itemTrim == "no-unicast-replies":
		block.NoUnicastReplies = types.BoolValue(true)
	case itemTrim == "proxy-mode":
		block.ProxyMode = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "relay-source "):
		block.RelaySource = types.StringValue(itemTrim)
	case itemTrim == "replace-ip-source-with giaddr":
		block.ReplaceIPSourceWithGiaddr = types.BoolValue(true)
	case itemTrim == "send-release-on-delete":
		block.SendReleaseOnDelete = types.BoolValue(true)
	case itemTrim == "trust-option-82":
		block.TrustOption82 = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user-defined-option-82 "):
		block.UserDefinedOption82 = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return err
}

type forwardingoptionsDhcprelayBlockOverridesV6 struct {
	AllowSnoopedClients                     types.Bool   `tfsdk:"allow_snooped_clients"`
	NoAllowSnoopedClients                   types.Bool   `tfsdk:"no_allow_snooped_clients"`
	AlwaysProcessOptionRequestOption        types.Bool   `tfsdk:"always_process_option_request_option"`
	AsymmetricLeaseTime                     types.Int64  `tfsdk:"asymmetric_lease_time"`
	AsymmetricPrefixLeaseTime               types.Int64  `tfsdk:"asymmetric_prefix_lease_time"`
	ClientNegotiationMatchIncomingInterface types.Bool   `tfsdk:"client_negotiation_match_incoming_interface"`
	DelayAuthentication                     types.Bool   `tfsdk:"delay_authentication"`
	DeleteBindingOnRenegotiation            types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
	DualStack                               types.String `tfsdk:"dual_stack"`
	InterfaceClientLimit                    types.Int64  `tfsdk:"interface_client_limit"`
	NoBindOnRequest                         types.Bool   `tfsdk:"no_bind_on_request"`
	RelaySource                             types.String `tfsdk:"relay_source"`
	SendReleaseOnDelete                     types.Bool   `tfsdk:"send_release_on_delete"`
}

func (block *forwardingoptionsDhcprelayBlockOverridesV6) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *forwardingoptionsDhcprelayBlockOverridesV6) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockOverridesV6) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"allow_snooped_clients": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow client creation from snooped PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_allow_snooped_clients": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't allow client creation from snooped PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"always_process_option_request_option": schema.BoolAttribute{
			Optional:    true,
			Description: "Always process option even after address allocation failure.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"asymmetric_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced lease time for the client. In seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"asymmetric_prefix_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced prefix lease time for the client. In seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"client_negotiation_match_incoming_interface": schema.BoolAttribute{
			Optional:    true,
			Description: "Use incoming interface match criteria for SOLICIT PDU.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"delay_authentication": schema.BoolAttribute{
			Optional:    true,
			Description: "Delay subscriber authentication in DHCP protocol processing until request packet",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"delete_binding_on_renegotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Delete binding on renegotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"dual_stack": schema.StringAttribute{
			Optional:    true,
			Description: "Dual stack group to use.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"interface_client_limit": schema.Int64Attribute{
			Optional:    true,
			Description: "Limit the number of clients allowed on an interface.",
			Validators: []validator.Int64{
				int64validator.Between(1, 500000),
			},
		},
		"no_bind_on_request": schema.BoolAttribute{
			Optional:    true,
			Description: "Do not bind if stray DHCPv6 RENEW, REBIND is received.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"relay_source": schema.StringAttribute{
			Optional:    true,
			Description: "Interface for relay source.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
				tfvalidator.String1DotCount(),
			},
		},
		"send_release_on_delete": schema.BoolAttribute{
			Optional:    true,
			Description: "Always send RELEASE to the server when a binding is deleted.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockOverridesV6) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if block.AllowSnoopedClients.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-snooped-clients")
		if block.NoAllowSnoopedClients.ValueBool() {
			return configSet,
				path.Root("overrides_v6").AtName("no_allow_snooped_clients"),
				errors.New("allow_snooped_clients and no_allow_snooped_clients cannot be configured together" +
					" in overrides_v6 block")
		}
	} else if block.NoAllowSnoopedClients.ValueBool() {
		configSet = append(configSet, setPrefix+"no-allow-snooped-clients")
	}
	if block.AlwaysProcessOptionRequestOption.ValueBool() {
		configSet = append(configSet, setPrefix+"always-process-option-request-option")
	}
	if !block.AsymmetricLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+
			utils.ConvI64toa(block.AsymmetricLeaseTime.ValueInt64()))
	}
	if !block.AsymmetricPrefixLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-prefix-lease-time "+
			utils.ConvI64toa(block.AsymmetricPrefixLeaseTime.ValueInt64()))
	}
	if block.ClientNegotiationMatchIncomingInterface.ValueBool() {
		configSet = append(configSet, setPrefix+"client-negotiation-match incoming-interface")
	}
	if block.DelayAuthentication.ValueBool() {
		configSet = append(configSet, setPrefix+"delay-authentication")
	}
	if block.DeleteBindingOnRenegotiation.ValueBool() {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if v := block.DualStack.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if !block.InterfaceClientLimit.IsNull() {
		configSet = append(configSet, setPrefix+"interface-client-limit "+
			utils.ConvI64toa(block.InterfaceClientLimit.ValueInt64()))
	}
	if block.NoBindOnRequest.ValueBool() {
		configSet = append(configSet, setPrefix+"no-bind-on-request")
	}
	if v := block.RelaySource.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"relay-source "+v)
	}
	if block.SendReleaseOnDelete.ValueBool() {
		configSet = append(configSet, setPrefix+"send-release-on-delete")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockOverridesV6) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "allow-snooped-clients":
		block.AllowSnoopedClients = types.BoolValue(true)
	case itemTrim == "always-process-option-request-option":
		block.AlwaysProcessOptionRequestOption = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		block.AsymmetricLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-prefix-lease-time "):
		block.AsymmetricPrefixLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "client-negotiation-match incoming-interface":
		block.ClientNegotiationMatchIncomingInterface = types.BoolValue(true)
	case itemTrim == "delay-authentication":
		block.DelayAuthentication = types.BoolValue(true)
	case itemTrim == "delete-binding-on-renegotiation":
		block.DeleteBindingOnRenegotiation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		block.DualStack = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		block.InterfaceClientLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "no-allow-snooped-clients":
		block.NoAllowSnoopedClients = types.BoolValue(true)
	case itemTrim == "no-bind-on-request":
		block.NoBindOnRequest = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "relay-source "):
		block.RelaySource = types.StringValue(itemTrim)
	case itemTrim == "send-release-on-delete":
		block.SendReleaseOnDelete = types.BoolValue(true)
	}

	return err
}

type forwardingoptionsDhcprelayBlockRelayAgentID struct {
	IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
	NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
	PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
	PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
	UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
	UseOption82               types.Bool   `tfsdk:"use_option_82"`
	UseOption82Strict         types.Bool   `tfsdk:"use_option_82_strict"`
	UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
}

func (forwardingoptionsDhcprelayBlockRelayAgentID) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"include_irb_and_l2": schema.BoolAttribute{
			Optional:    true,
			Description: "Include IRB and L2 interface name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_vlan_interface_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Not include vlan or interface name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"prefix_host_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Add router host name to circuit / interface-id or remote-id.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"prefix_routing_instance_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Add routing instance name to circuit / interface-id or remote-id.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"use_interface_description": schema.StringAttribute{
			Optional:    true,
			Description: "Use interface description instead of circuit identifier.",
			Validators: []validator.String{
				stringvalidator.OneOf("device", "logical"),
			},
		},
		"use_option_82": schema.BoolAttribute{
			Optional:    true,
			Description: "Use option-82 circuit-id for interface-id or remote-id.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"use_option_82_strict": schema.BoolAttribute{
			Optional:    true,
			Description: "Drop packet if option-82 circuit-id not present.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"use_vlan_id": schema.BoolAttribute{
			Optional:    true,
			Description: "Use VLAN id instead of name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentID) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)

	if block.IncludeIrbAndL2.ValueBool() {
		configSet = append(configSet, setPrefix+"include-irb-and-l2")
	}
	if block.NoVlanInterfaceName.ValueBool() {
		configSet = append(configSet, setPrefix+"no-vlan-interface-name")
	}
	if block.PrefixHostName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix host-name")
	}
	if block.PrefixRoutingInstanceName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix routing-instance-name")
	}
	if v := block.UseInterfaceDescription.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"use-interface-description "+v)
	}
	if block.UseOption82.ValueBool() {
		configSet = append(configSet, setPrefix+"use-option-82")
		if block.UseOption82Strict.ValueBool() {
			configSet = append(configSet, setPrefix+"use-option-82 strict")
		}
	} else if block.UseOption82Strict.ValueBool() {
		return configSet,
			pathRoot.AtName("use_option_82_strict"),
			errors.New("use_option_82 must be specified with " +
				"use_option_82_strict")
	}
	if block.UseVlanID.ValueBool() {
		configSet = append(configSet, setPrefix+"use-vlan-id")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentID) read(itemTrim string) {
	switch {
	case itemTrim == "include-irb-and-l2":
		block.IncludeIrbAndL2 = types.BoolValue(true)
	case itemTrim == "no-vlan-interface-name":
		block.NoVlanInterfaceName = types.BoolValue(true)
	case itemTrim == "prefix host-name":
		block.PrefixHostName = types.BoolValue(true)
	case itemTrim == "prefix routing-instance-name":
		block.PrefixRoutingInstanceName = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
		block.UseInterfaceDescription = types.StringValue(itemTrim)
	case itemTrim == "use-option-82":
		block.UseOption82 = types.BoolValue(true)
	case itemTrim == "use-option-82 strict":
		block.UseOption82 = types.BoolValue(true)
		block.UseOption82Strict = types.BoolValue(true)
	case itemTrim == "use-vlan-id":
		block.UseVlanID = types.BoolValue(true)
	}
}

type forwardingoptionsDhcprelayBlockRelayAgentInterfaceID struct {
	forwardingoptionsDhcprelayBlockRelayAgentID
	KeepIncomingID       types.Bool `tfsdk:"keep_incoming_id"`
	KeepIncomingIDStrict types.Bool `tfsdk:"keep_incoming_id_strict"`
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentInterfaceID) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockRelayAgentInterfaceID) attributesSchema() map[string]schema.Attribute {
	attributes := forwardingoptionsDhcprelayBlockRelayAgentID{}.attributesSchema()

	attributes["keep_incoming_id"] = schema.BoolAttribute{
		Optional:    true,
		Description: "Keep incoming interface identifier.",
		Validators: []validator.Bool{
			tfvalidator.BoolTrue(),
		},
	}
	attributes["keep_incoming_id_strict"] = schema.BoolAttribute{
		Optional:    true,
		Description: "Drop packet if interface identifier not present.",
		Validators: []validator.Bool{
			tfvalidator.BoolTrue(),
		},
	}

	return attributes
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentInterfaceID) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "relay-agent-interface-id "

	configSet := []string{
		setPrefix,
	}

	blockSet, pathErr, err := block.forwardingoptionsDhcprelayBlockRelayAgentID.configSet(setPrefix, pathRoot)
	if err != nil {
		return configSet, pathErr, err
	}
	configSet = append(configSet, blockSet...)

	if block.KeepIncomingID.ValueBool() {
		configSet = append(configSet, setPrefix+"keep-incoming-interface-id")
		if block.KeepIncomingIDStrict.ValueBool() {
			configSet = append(configSet, setPrefix+"keep-incoming-interface-id strict")
		}
	} else if block.KeepIncomingIDStrict.ValueBool() {
		return configSet,
			pathRoot.AtName("keep_incoming_id_strict"),
			errors.New("keep_incoming_id must be specified with keep_incoming_id_strict" +
				" in relay_agent_interface_id block")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentInterfaceID) read(itemTrim string) {
	switch {
	case itemTrim == "keep-incoming-interface-id":
		block.KeepIncomingID = types.BoolValue(true)
	case itemTrim == "keep-incoming-interface-id strict":
		block.KeepIncomingID = types.BoolValue(true)
		block.KeepIncomingIDStrict = types.BoolValue(true)
	default:
		block.forwardingoptionsDhcprelayBlockRelayAgentID.read(itemTrim)
	}
}

type forwardingoptionsDhcprelayBlockRelayAgentRemoteID struct {
	forwardingoptionsDhcprelayBlockRelayAgentID
	KeepIncomingID types.Bool `tfsdk:"keep_incoming_id"`
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentRemoteID) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockRelayAgentRemoteID) attributesSchema() map[string]schema.Attribute {
	attributes := forwardingoptionsDhcprelayBlockRelayAgentID{}.attributesSchema()

	attributes["keep_incoming_id"] = schema.BoolAttribute{
		Optional:    true,
		Description: "Keep incoming interface identifier.",
		Validators: []validator.Bool{
			tfvalidator.BoolTrue(),
		},
	}

	return attributes
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentRemoteID) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "relay-agent-remote-id "

	configSet := []string{
		setPrefix,
	}

	blockSet, pathErr, err := block.forwardingoptionsDhcprelayBlockRelayAgentID.configSet(setPrefix, pathRoot)
	if err != nil {
		return configSet, pathErr, err
	}
	configSet = append(configSet, blockSet...)

	if block.KeepIncomingID.ValueBool() {
		configSet = append(configSet, setPrefix+"keep-incoming-remote-id")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockRelayAgentRemoteID) read(itemTrim string) {
	switch {
	case itemTrim == "keep-incoming-remote-id":
		block.KeepIncomingID = types.BoolValue(true)
	default:
		block.forwardingoptionsDhcprelayBlockRelayAgentID.read(itemTrim)
	}
}

//nolint:lll
type forwardingoptionsDhcprelayBlockRelayOption struct {
	OptionOrder           []types.String                                                       `tfsdk:"option_order"`
	Option15              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN             `tfsdk:"option_15"`
	Option15DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_15_default_action"`
	Option16              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN             `tfsdk:"option_16"`
	Option16DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_16_default_action"`
	Option60              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN             `tfsdk:"option_60"`
	Option60DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_60_default_action"`
	Option77              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN             `tfsdk:"option_77"`
	Option77DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_77_default_action"`
}

func (forwardingoptionsDhcprelayBlockRelayOption) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"option_order": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Options precedence order.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.NoNullValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("15", "16", "60", "77"),
				),
			},
		},
	}
}

func (forwardingoptionsDhcprelayBlockRelayOption) blocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"option_15": schema.SetNestedBlock{
			Description: "For each combination of `compare`, `value_type` and `value` arguments, add option 15 processing.",
			NestedObject: schema.NestedBlockObject{
				Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN{}.attributesSchema(
					"drop", "forward-only", "relay-server-group",
				),
			},
		},
		"option_15_default_action": schema.SingleNestedBlock{
			Description: "Generic option 15 default action.",
			Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{}.attributesSchema(
				"drop", "forward-only", "relay-server-group",
			),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"option_16": schema.SetNestedBlock{
			Description: "For each combination of `compare`, `value_type` and `value` arguments, add option 16 processing.",
			NestedObject: schema.NestedBlockObject{
				Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN{}.attributesSchema(
					"drop", "forward-only", "relay-server-group",
				),
			},
		},
		"option_16_default_action": schema.SingleNestedBlock{
			Description: "Generic option 16 default action.",
			Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{}.attributesSchema(
				"drop", "forward-only", "relay-server-group",
			),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"option_60": schema.SetNestedBlock{
			Description: "For each combination of `compare`, `value_type` and `value` arguments, add option 60 processing.",
			NestedObject: schema.NestedBlockObject{
				Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN{}.attributesSchema(
					"drop", "forward-only", "local-server-group", "relay-server-group",
				),
			},
		},
		"option_60_default_action": schema.SingleNestedBlock{
			Description: "Generic option 60 default action.",
			Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{}.attributesSchema(
				"drop", "forward-only", "local-server-group", "relay-server-group",
			),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"option_77": schema.SetNestedBlock{
			Description: "For each combination of `compare`, `value_type` and `value` arguments, add option 77 processing.",
			NestedObject: schema.NestedBlockObject{
				Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN{}.attributesSchema(
					"drop", "forward-only", "local-server-group", "relay-server-group",
				),
			},
		},
		"option_77_default_action": schema.SingleNestedBlock{
			Description: "Generic option 77 default action.",
			Attributes: forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{}.attributesSchema(
				"drop", "forward-only", "local-server-group", "relay-server-group",
			),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockRelayOption) configSet(
	setPrefix, version string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "relay-option "

	for _, v := range block.OptionOrder {
		vv := v.ValueString()
		if version == "v4" && (vv == "15" || vv == "16") {
			return configSet,
				path.Root("relay_option").AtName("option_order"),
				errors.New("option_order cannot be configured with 15 or 16 when version = v4" +
					" in relay_option block")
		}
		if version == "v6" && (vv == "60" || vv == "77") {
			return configSet,
				path.Root("relay_option").AtName("option_order"),
				errors.New("option_order cannot be configured with 60 or 77 when version = v6" +
					" in relay_option block")
		}

		configSet = append(configSet, setPrefix+"option-order "+v.ValueString())
	}
	for _, subBlock := range block.Option15 {
		if version == "v4" {
			return configSet,
				path.Root("relay_option").AtName("option_15"),
				errors.New("option_15 cannot be configured when version = v4" +
					" in relay_option block")
		}

		if action := subBlock.
			Action.ValueString(); action == "relay-server-group" {
			if subBlock.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group must be specified when action = relay-server-group" +
						" in option_15 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-15 "))
		} else {
			if subBlock.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group cannot be configured with action != relay-server-group" +
						" in option_15 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-15 "))
		}
	}
	if block.Option15DefaultAction != nil {
		if version == "v4" {
			return configSet,
				path.Root("relay_option").AtName("option_15_default_action"),
				errors.New("option_15_default_action cannot be configured when version = v4" +
					" in relay_option block")
		}

		if action := block.Option15DefaultAction.
			Action.ValueString(); action == "relay-server-group" {
			if block.Option15DefaultAction.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option").AtName("option_15_default_action").AtName("action"),
					errors.New("group must be specified when action = relay-server-group" +
						" in option_15_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+
				"option-15 default-action "+action+
				" \""+block.Option15DefaultAction.Group.ValueString()+"\"",
			)
		} else {
			if block.Option15DefaultAction.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option").AtName("option_15_default_action").AtName("group"),
					errors.New("group cannot be configured with action != relay-server-group" +
						" in option_15_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-15 default-action "+action)
		}
	}
	for _, subBlock := range block.Option16 {
		if version == "v4" {
			return configSet,
				path.Root("relay_option").AtName("option_16"),
				errors.New("option_16 cannot be configured when version = v4" +
					" in relay_option block")
		}

		if action := subBlock.
			Action.ValueString(); action == "relay-server-group" {
			if subBlock.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group must be specified when action = relay-server-group" +
						" in option_16 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-16 "))
		} else {
			if subBlock.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group cannot be configured with action != relay-server-group" +
						" in option_16 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-16 "))
		}
	}
	if block.Option16DefaultAction != nil {
		if version == "v4" {
			return configSet,
				path.Root("relay_option").AtName("option_16_default_action"),
				errors.New("option_16_default_action cannot be configured when version = v4" +
					" in relay_option block")
		}

		if action := block.Option16DefaultAction.
			Action.ValueString(); action == "relay-server-group" {
			if block.Option16DefaultAction.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option").AtName("option_16_default_action").AtName("action"),
					errors.New("group must be specified when action = relay-server-group" +
						" in option_16_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+
				"option-16 default-action "+action+
				" \""+block.Option16DefaultAction.Group.ValueString()+"\"",
			)
		} else {
			if block.Option16DefaultAction.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option").AtName("option_16_default_action").AtName("group"),
					errors.New("group cannot be configured with action != relay-server-group" +
						" in option_16_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-16 default-action "+action)
		}
	}
	for _, subBlock := range block.Option60 {
		if version == "v6" {
			return configSet,
				path.Root("relay_option").AtName("option_60"),
				errors.New("option_60 cannot be configured when version = v6" +
					" in relay_option block")
		}

		if action := subBlock.
			Action.ValueString(); action == "local-server-group" || action == "relay-server-group" {
			if subBlock.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group must be specified when action = local-server-group or relay-server-group" +
						" in option_60 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-60 "))
		} else {
			if subBlock.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group cannot be configured with action != local-server-group && relay-server-group " +
						" in option_60 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-60 "))
		}
	}
	if block.Option60DefaultAction != nil {
		if version == "v6" {
			return configSet,
				path.Root("relay_option").AtName("option_60_default_action"),
				errors.New("option_60_default_action cannot be configured when version = v6" +
					" in relay_option block")
		}

		if action := block.Option60DefaultAction.
			Action.ValueString(); action == "local-server-group" || action == "relay-server-group" {
			if block.Option60DefaultAction.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option").AtName("option_60_default_action").AtName("action"),
					errors.New("group must be specified when action = local-server-group or relay-server-group" +
						" in option_60_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+
				"option-60 default-action "+action+
				" \""+block.Option60DefaultAction.Group.ValueString()+"\"",
			)
		} else {
			if block.Option60DefaultAction.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option").AtName("option_60_default_action").AtName("group"),
					errors.New("group cannot be configured with action != local-server-group && relay-server-group" +
						" in option_60_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-60 default-action "+action)
		}
	}
	for _, subBlock := range block.Option77 {
		if version == "v6" {
			return configSet,
				path.Root("relay_option").AtName("option_77"),
				errors.New("option_77 cannot be configured when version = v6" +
					" in relay_option block")
		}

		if action := subBlock.
			Action.ValueString(); action == "local-server-group" || action == "relay-server-group" {
			if subBlock.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group must be specified when action = local-server-group or relay-server-group" +
						" in option_77 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-77 "))
		} else {
			if subBlock.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option"),
					errors.New("group cannot be configured with action != local-server-group && relay-server-group " +
						" in option_77 block in relay_option block")
			}
			configSet = append(configSet, subBlock.configSet(setPrefix+"option-77 "))
		}
	}
	if block.Option77DefaultAction != nil {
		if version == "v6" {
			return configSet,
				path.Root("relay_option").AtName("option_77_default_action"),
				errors.New("option_77_default_action cannot be configured when version = v6" +
					" in relay_option block")
		}

		if action := block.Option77DefaultAction.
			Action.ValueString(); action == "local-server-group" || action == "relay-server-group" {
			if block.Option77DefaultAction.Group.ValueString() == "" {
				return configSet,
					path.Root("relay_option").AtName("option_77_default_action").AtName("action"),
					errors.New("group must be specified when action = local-server-group or relay-server-group" +
						" in option_77_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+
				"option-77 default-action "+action+
				" \""+block.Option77DefaultAction.Group.ValueString()+"\"",
			)
		} else {
			if block.Option77DefaultAction.Group.ValueString() != "" {
				return configSet,
					path.Root("relay_option").AtName("option_77_default_action").AtName("group"),
					errors.New("group cannot be configured with action != local-server-group && relay-server-group" +
						" in option_77_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-77 default-action "+action)
		}
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockRelayOption) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "option-15 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")

		block.Option15DefaultAction = &forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{
			Action: types.StringValue(itemTrimFields[0]),
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			block.Option15DefaultAction.Group = types.StringValue(strings.Trim(strings.Join(itemTrimFields[1:], " "), "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "option-15 "):
		var subBlock forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN
		if err := subBlock.read(itemTrim, "option-15"); err != nil {
			return err
		}
		block.Option15 = append(block.Option15, subBlock)
	case balt.CutPrefixInString(&itemTrim, "option-16 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")

		block.Option16DefaultAction = &forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{
			Action: types.StringValue(itemTrimFields[0]),
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			block.Option16DefaultAction.Group = types.StringValue(strings.Trim(strings.Join(itemTrimFields[1:], " "), "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "option-16 "):
		var subBlock forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN
		if err := subBlock.read(itemTrim, "option-16"); err != nil {
			return err
		}
		block.Option16 = append(block.Option16, subBlock)
	case balt.CutPrefixInString(&itemTrim, "option-60 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")

		block.Option60DefaultAction = &forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{
			Action: types.StringValue(itemTrimFields[0]),
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			block.Option60DefaultAction.Group = types.StringValue(strings.Trim(strings.Join(itemTrimFields[1:], " "), "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "option-60 "):
		var subBlock forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN
		if err := subBlock.read(itemTrim, "option-60"); err != nil {
			return err
		}
		block.Option60 = append(block.Option60, subBlock)
	case balt.CutPrefixInString(&itemTrim, "option-77 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")

		block.Option77DefaultAction = &forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction{
			Action: types.StringValue(itemTrimFields[0]),
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			block.Option77DefaultAction.Group = types.StringValue(strings.Trim(strings.Join(itemTrimFields[1:], " "), "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "option-77 "):
		var subBlock forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN
		if err := subBlock.read(itemTrim, "option-77"); err != nil {
			return err
		}
		block.Option77 = append(block.Option77, subBlock)
	case balt.CutPrefixInString(&itemTrim, "option-order "):
		block.OptionOrder = append(block.OptionOrder, types.StringValue(itemTrim))
	}

	return nil
}

//nolint:lll
type forwardingoptionsDhcprelayBlockRelayOptionConfig struct {
	OptionOrder           types.List                                                           `tfsdk:"option_order"`
	Option15              types.Set                                                            `tfsdk:"option_15"`
	Option15DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_15_default_action"`
	Option16              types.Set                                                            `tfsdk:"option_16"`
	Option16DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_16_default_action"`
	Option60              types.Set                                                            `tfsdk:"option_60"`
	Option60DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_60_default_action"`
	Option77              types.Set                                                            `tfsdk:"option_77"`
	Option77DefaultAction *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_77_default_action"`
}

type forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN struct {
	Compare   types.String `tfsdk:"compare"`
	ValueType types.String `tfsdk:"value_type"`
	Value     types.String `tfsdk:"value"`
	Action    types.String `tfsdk:"action"`
	Group     types.String `tfsdk:"group"`
}

func (forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN) attributesSchema(
	actions ...string,
) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"compare": schema.StringAttribute{
			Required:    true,
			Description: "How to compare.",
			Validators: []validator.String{
				stringvalidator.OneOf("equals", "starts-with"),
			},
		},
		"value_type": schema.StringAttribute{
			Required:    true,
			Description: "Type of string.",
			Validators: []validator.String{
				stringvalidator.OneOf("ascii", "hexadecimal"),
			},
		},
		"value": schema.StringAttribute{
			Required:    true,
			Description: "String to compare.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 256),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"action": schema.StringAttribute{
			Required:    true,
			Description: "Action on match.",
			Validators: []validator.String{
				stringvalidator.OneOf(actions...),
			},
		},
		"group": schema.StringAttribute{
			Optional:    true,
			Description: "Group for action.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN) configSet(setPrefix string) string {
	if block.Group.ValueString() != "" {
		return setPrefix +
			block.Compare.ValueString() + " " +
			block.ValueType.ValueString() + " " +
			"\"" + block.Value.ValueString() + "\" " +
			block.Action.ValueString() + " " +
			"\"" + block.Group.ValueString() + "\""
	}

	return setPrefix +
		block.Compare.ValueString() + " " +
		block.ValueType.ValueString() + " " +
		"\"" + block.Value.ValueString() + "\" " +
		block.Action.ValueString()
}

func (block *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN) read(itemTrim, option string) error {
	itemTrimFields := strings.Split(itemTrim, " ")
	if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action> <group>?
		return fmt.Errorf(junos.CantReadValuesNotEnoughFields, option, itemTrim)
	}
	value := itemTrimFields[2]
	actionIndex := 3
	if (strings.HasPrefix(itemTrimFields[2], "\"") && !strings.HasSuffix(itemTrimFields[2], "\"")) ||
		itemTrimFields[2] == "\"" {
		for k, v := range itemTrimFields[3:] {
			value += " " + v
			if strings.Contains(v, "\"") {
				actionIndex = 3 + k + 1

				break
			}
		}
	}

	block.Compare = types.StringValue(itemTrimFields[0])
	block.ValueType = types.StringValue(itemTrimFields[1])
	block.Value = types.StringValue(html.UnescapeString(strings.Trim(value, "\"")))
	block.Action = types.StringValue(itemTrimFields[actionIndex])
	if len(itemTrimFields) > actionIndex+1 {
		block.Group = types.StringValue(strings.Trim(strings.Join(itemTrimFields[actionIndex+1:], " "), "\""))
	}

	return nil
}

type forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction struct {
	Action types.String `tfsdk:"action"`
	Group  types.String `tfsdk:"group"`
}

func (block *forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction) attributesSchema(
	actions ...string,
) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Action.",
			Validators: []validator.String{
				stringvalidator.OneOf(actions...),
			},
		},
		"group": schema.StringAttribute{
			Optional:    true,
			Description: "Group for action.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

//nolint:lll
type forwardingoptionsDhcprelayBlockRelayOption82 struct {
	ExcludeRelayAgentIdentifier types.Bool                                                  `tfsdk:"exclude_relay_agent_identifier"`
	LinkSelection               types.Bool                                                  `tfsdk:"link_selection"`
	ServerIDOverride            types.Bool                                                  `tfsdk:"server_id_override"`
	VendorSpecificHostName      types.Bool                                                  `tfsdk:"vendor_specific_host_name"`
	VendorSpecificLocation      types.Bool                                                  `tfsdk:"vendor_specific_location"`
	CircuitID                   *forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID `tfsdk:"circuit_id"`
	RemoteID                    *forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID  `tfsdk:"remote_id"`
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (forwardingoptionsDhcprelayBlockRelayOption82) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"exclude_relay_agent_identifier": schema.BoolAttribute{
			Optional:    true,
			Description: "Exclude relay agent identifier from packets to server.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"link_selection": schema.BoolAttribute{
			Optional:    true,
			Description: "Add link-selection sub-option on packets to server.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"server_id_override": schema.BoolAttribute{
			Optional:    true,
			Description: "Add link-selection and server-id sub-options on packets to server.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"vendor_specific_host_name": schema.BoolAttribute{
			Optional:    true,
			Description: "Add vendor-specific information, add router host name.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"vendor_specific_location": schema.BoolAttribute{
			Optional:    true,
			Description: "Add vendor-specific information, add location information expressed as interface name format.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (forwardingoptionsDhcprelayBlockRelayOption82) blocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"circuit_id": schema.SingleNestedBlock{
			Description: "Add circuit identifier.",
			Attributes: map[string]schema.Attribute{
				"include_irb_and_l2": schema.BoolAttribute{
					Optional:    true,
					Description: "Include IRB and L2 interface name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"keep_incoming_circuit_id": schema.BoolAttribute{
					Optional:    true,
					Description: "Keep incoming circuit identifier.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"no_vlan_interface_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Not include vlan or interface name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"prefix_host_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Add router host name to circuit / interface-id or remote-id.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"prefix_routing_instance_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Add routing instance name to circuit / interface-id or remote-id.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"use_interface_description": schema.StringAttribute{
					Optional:    true,
					Description: "Use interface description instead of circuit identifier.",
					Validators: []validator.String{
						stringvalidator.OneOf("device", "logical"),
					},
				},
				"use_vlan_id": schema.BoolAttribute{
					Optional:    true,
					Description: "Use VLAN id instead of name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"user_defined": schema.BoolAttribute{
					Optional:    true,
					Description: "Include user defined string.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"vlan_id_only": schema.BoolAttribute{
					Optional:    true,
					Description: "Use only VLAN id.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"remote_id": schema.SingleNestedBlock{
			Description: "Add remote identifier.",
			Attributes: map[string]schema.Attribute{
				"hostname_only": schema.BoolAttribute{
					Optional:    true,
					Description: "Include hostname only.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"include_irb_and_l2": schema.BoolAttribute{
					Optional:    true,
					Description: "Include IRB and L2 interface name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"keep_incoming_remote_id": schema.BoolAttribute{
					Optional:    true,
					Description: "Keep incoming remote identifier.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"no_vlan_interface_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Not include vlan or interface name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"prefix_host_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Add router host name to circuit / interface-id or remote-id.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"prefix_routing_instance_name": schema.BoolAttribute{
					Optional:    true,
					Description: "Add routing instance name to circuit / interface-id or remote-id.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
				"use_interface_description": schema.StringAttribute{
					Optional:    true,
					Description: "Use interface description instead of circuit identifier.",
					Validators: []validator.String{
						stringvalidator.OneOf("device", "logical"),
					},
				},
				"use_string": schema.StringAttribute{
					Optional:    true,
					Description: "Use raw string instead of the default remote id.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					},
				},
				"use_vlan_id": schema.BoolAttribute{
					Optional:    true,
					Description: "Use VLAN id instead of name.",
					Validators: []validator.Bool{
						tfvalidator.BoolTrue(),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82) configSet(setPrefix string) []string {
	setPrefix += "relay-option-82 "

	configSet := []string{
		setPrefix,
	}

	if block.ExcludeRelayAgentIdentifier.ValueBool() {
		configSet = append(configSet, setPrefix+"exclude-relay-agent-identifier")
	}
	if block.LinkSelection.ValueBool() {
		configSet = append(configSet, setPrefix+"link-selection")
	}
	if block.ServerIDOverride.ValueBool() {
		configSet = append(configSet, setPrefix+"server-id-override")
	}
	if block.VendorSpecificHostName.ValueBool() {
		configSet = append(configSet, setPrefix+"vendor-specific host-name")
	}
	if block.VendorSpecificLocation.ValueBool() {
		configSet = append(configSet, setPrefix+"vendor-specific location")
	}

	if block.CircuitID != nil {
		configSet = append(configSet, block.CircuitID.configSet(setPrefix)...)
	}
	if block.RemoteID != nil {
		configSet = append(configSet, block.RemoteID.configSet(setPrefix)...)
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82) read(itemTrim string) {
	switch {
	case itemTrim == "exclude-relay-agent-identifier":
		block.ExcludeRelayAgentIdentifier = types.BoolValue(true)
	case itemTrim == "link-selection":
		block.LinkSelection = types.BoolValue(true)
	case itemTrim == "server-id-override":
		block.ServerIDOverride = types.BoolValue(true)
	case itemTrim == "vendor-specific host-name":
		block.VendorSpecificHostName = types.BoolValue(true)
	case itemTrim == "vendor-specific location":
		block.VendorSpecificLocation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "circuit-id"):
		if block.CircuitID == nil {
			block.CircuitID = &forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.CircuitID.read(itemTrim)
		}

	case balt.CutPrefixInString(&itemTrim, "remote-id"):
		if block.RemoteID == nil {
			block.RemoteID = &forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.RemoteID.read(itemTrim)
		}
	}
}

type forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID struct {
	IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
	KeepIncomingCircuitID     types.Bool   `tfsdk:"keep_incoming_circuit_id"`
	NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
	PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
	PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
	UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
	UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
	UserDefined               types.Bool   `tfsdk:"user_defined"`
	VlanIDOnly                types.Bool   `tfsdk:"vlan_id_only"`
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID) configSet(setPrefix string) []string {
	setPrefix += "circuit-id "

	configSet := []string{
		setPrefix,
	}

	if block.IncludeIrbAndL2.ValueBool() {
		configSet = append(configSet, setPrefix+"include-irb-and-l2")
	}
	if block.KeepIncomingCircuitID.ValueBool() {
		configSet = append(configSet, setPrefix+"keep-incoming-circuit-id")
	}
	if block.NoVlanInterfaceName.ValueBool() {
		configSet = append(configSet, setPrefix+"no-vlan-interface-name")
	}
	if block.PrefixHostName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix host-name")
	}
	if block.PrefixRoutingInstanceName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix routing-instance-name")
	}
	if v := block.UseInterfaceDescription.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"use-interface-description "+v)
	}
	if block.UseVlanID.ValueBool() {
		configSet = append(configSet, setPrefix+"use-vlan-id")
	}
	if block.UserDefined.ValueBool() {
		configSet = append(configSet, setPrefix+"user-defined")
	}
	if block.VlanIDOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"vlan-id-only")
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID) read(itemTrim string) {
	switch {
	case itemTrim == "include-irb-and-l2":
		block.IncludeIrbAndL2 = types.BoolValue(true)
	case itemTrim == "keep-incoming-circuit-id":
		block.KeepIncomingCircuitID = types.BoolValue(true)
	case itemTrim == "no-vlan-interface-name":
		block.NoVlanInterfaceName = types.BoolValue(true)
	case itemTrim == "prefix host-name":
		block.PrefixHostName = types.BoolValue(true)
	case itemTrim == "prefix routing-instance-name":
		block.PrefixRoutingInstanceName = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
		block.UseInterfaceDescription = types.StringValue(itemTrim)
	case itemTrim == "use-vlan-id":
		block.UseVlanID = types.BoolValue(true)
	case itemTrim == "user-defined":
		block.UserDefined = types.BoolValue(true)
	case itemTrim == "vlan-id-only":
		block.VlanIDOnly = types.BoolValue(true)
	}
}

type forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID struct {
	HostnameOnly              types.Bool   `tfsdk:"hostname_only"`
	IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
	KeepIncomingRemoteID      types.Bool   `tfsdk:"keep_incoming_remote_id"`
	NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
	PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
	PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
	UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
	UseString                 types.String `tfsdk:"use_string"`
	UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID) configSet(setPrefix string) []string {
	setPrefix += "remote-id "

	configSet := []string{
		setPrefix,
	}

	if block.HostnameOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"hostname-only")
	}
	if block.IncludeIrbAndL2.ValueBool() {
		configSet = append(configSet, setPrefix+"include-irb-and-l2")
	}
	if block.KeepIncomingRemoteID.ValueBool() {
		configSet = append(configSet, setPrefix+"keep-incoming-remote-id")
	}
	if block.NoVlanInterfaceName.ValueBool() {
		configSet = append(configSet, setPrefix+"no-vlan-interface-name")
	}
	if block.PrefixHostName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix host-name")
	}
	if block.PrefixRoutingInstanceName.ValueBool() {
		configSet = append(configSet, setPrefix+"prefix routing-instance-name")
	}
	if v := block.UseInterfaceDescription.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"use-interface-description "+v)
	}
	if v := block.UseString.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"use-string \""+v+"\"")
	}
	if block.UseVlanID.ValueBool() {
		configSet = append(configSet, setPrefix+"use-vlan-id")
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID) read(itemTrim string) {
	switch {
	case itemTrim == "hostname-only":
		block.HostnameOnly = types.BoolValue(true)
	case itemTrim == "include-irb-and-l2":
		block.IncludeIrbAndL2 = types.BoolValue(true)
	case itemTrim == "keep-incoming-remote-id":
		block.KeepIncomingRemoteID = types.BoolValue(true)
	case itemTrim == "no-vlan-interface-name":
		block.NoVlanInterfaceName = types.BoolValue(true)
	case itemTrim == "prefix host-name":
		block.PrefixHostName = types.BoolValue(true)
	case itemTrim == "prefix routing-instance-name":
		block.PrefixRoutingInstanceName = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
		block.UseInterfaceDescription = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "use-string "):
		block.UseString = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "use-vlan-id":
		block.UseVlanID = types.BoolValue(true)
	}
}

type forwardingoptionsDhcprelayBlockServerMatchAddress struct {
	Address types.String `tfsdk:"address"`
	Action  types.String `tfsdk:"action"`
}

func (forwardingoptionsDhcprelayBlockServerMatchAddress) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Required:    true,
			Description: "Server address.",
			Validators: []validator.String{
				tfvalidator.StringCIDR(),
			},
		},
		"action": schema.StringAttribute{
			Required:    true,
			Description: "Action on address.",
			Validators: []validator.String{
				stringvalidator.OneOf("create-relay-entry", "forward-only"),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockServerMatchAddress) configSet(setPrefix string) []string {
	return []string{
		setPrefix + "server-match address " + block.Address.ValueString() + " " + block.Action.ValueString(),
	}
}

func (block *forwardingoptionsDhcprelayBlockServerMatchAddress) read(itemTrim string) error {
	itemTrimFields := strings.Split(itemTrim, " ")
	if len(itemTrimFields) < 2 { // <address> <action>
		return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "server-match address", itemTrim)
	}

	block.Address = types.StringValue(itemTrimFields[0])
	block.Action = types.StringValue(itemTrimFields[1])

	return nil
}

type forwardingoptionsDhcprelayBlockServerMatchDuid struct {
	Compare   types.String `tfsdk:"compare"`
	ValueType types.String `tfsdk:"value_type"`
	Value     types.String `tfsdk:"value"`
	Action    types.String `tfsdk:"action"`
}

func (forwardingoptionsDhcprelayBlockServerMatchDuid) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"compare": schema.StringAttribute{
			Required:    true,
			Description: "How to compare.",
			Validators: []validator.String{
				stringvalidator.OneOf("equals", "starts-with"),
			},
		},
		"value_type": schema.StringAttribute{
			Required:    true,
			Description: "Type of string.",
			Validators: []validator.String{
				stringvalidator.OneOf("ascii", "hexadecimal"),
			},
		},
		"value": schema.StringAttribute{
			Required:    true,
			Description: "String to compare.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 256),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"action": schema.StringAttribute{
			Required:    true,
			Description: "Action on match.",
			Validators: []validator.String{
				stringvalidator.OneOf("create-relay-entry", "forward-only"),
			},
		},
	}
}

func (block *forwardingoptionsDhcprelayBlockServerMatchDuid) configSet(setPrefix string) []string {
	return []string{
		setPrefix + "server-match duid " +
			block.Compare.ValueString() + " " +
			block.ValueType.ValueString() + " " +
			"\"" + block.Value.ValueString() + "\" " +
			block.Action.ValueString(),
	}
}

func (block *forwardingoptionsDhcprelayBlockServerMatchDuid) read(itemTrim string) error {
	itemTrimFields := strings.Split(itemTrim, " ")
	if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action>
		return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "server-match duid", itemTrim)
	}
	if strings.Contains(itemTrimFields[2], "\"") {
		action := itemTrimFields[len(itemTrimFields)-1]
		value := strings.Trim(strings.Join(itemTrimFields[2:len(itemTrimFields)-1], " "), "\"")
		itemTrimFields[2] = value
		itemTrimFields[3] = action
	}
	block.Compare = types.StringValue(itemTrimFields[0])
	block.ValueType = types.StringValue(itemTrimFields[1])
	block.Value = types.StringValue(html.UnescapeString(itemTrimFields[2]))
	block.Action = types.StringValue(itemTrimFields[3])

	return nil
}
