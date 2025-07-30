package provider

import (
	"errors"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type dhcpBlockAuthenticationUsernameInclude struct {
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

func (block *dhcpBlockAuthenticationUsernameInclude) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (dhcpBlockAuthenticationUsernameInclude) attributesSchema() map[string]schema.Attribute {
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

func (block *dhcpBlockAuthenticationUsernameInclude) configSet(
	setPrefix, version string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
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

func (block *dhcpBlockAuthenticationUsernameInclude) read(itemTrim string) {
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

type dhcpBlockLivenessDetectionMethodBfd struct {
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

func (block *dhcpBlockLivenessDetectionMethodBfd) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *dhcpBlockLivenessDetectionMethodBfd) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (dhcpBlockLivenessDetectionMethodBfd) attributesSchema() map[string]schema.Attribute {
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

func (block *dhcpBlockLivenessDetectionMethodBfd) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
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

func (block *dhcpBlockLivenessDetectionMethodBfd) read(itemTrim string) (err error) {
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

type dhcpBlockLivenessDetectionMethodLayer2 struct {
	MaxConsecutiveRetries types.Int64 `tfsdk:"max_consecutive_retries"`
	TransmitInterval      types.Int64 `tfsdk:"transmit_interval"`
}

func (block *dhcpBlockLivenessDetectionMethodLayer2) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *dhcpBlockLivenessDetectionMethodLayer2) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (dhcpBlockLivenessDetectionMethodLayer2) attributesSchema() map[string]schema.Attribute {
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

func (block *dhcpBlockLivenessDetectionMethodLayer2) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
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

func (block *dhcpBlockLivenessDetectionMethodLayer2) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "max-consecutive-retries "):
		block.MaxConsecutiveRetries, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval "):
		block.TransmitInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}
