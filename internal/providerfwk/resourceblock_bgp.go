package providerfwk

import (
	"errors"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type bgpBlockBfdLivenessDetection struct {
	AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
	AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
	AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
	DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
	HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
	MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                      types.Int64  `tfsdk:"multiplier"`
	SessionMode                     types.String `tfsdk:"session_mode"`
	TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
	TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                         types.String `tfsdk:"version"`
}

func (block *bgpBlockBfdLivenessDetection) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (bgpBlockBfdLivenessDetection) resourceSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Define Bidirectional Forwarding Detection (BFD) options.",
		Attributes: map[string]schema.Attribute{
			"authentication_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "Authentication algorithm name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"authentication_key_chain": schema.StringAttribute{
				Optional:    true,
				Description: "Authentication key chain name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authentication_loose_check": schema.BoolAttribute{
				Optional:    true,
				Description: "Verify authentication only if authentication is negotiated.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"detection_time_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "High detection-time triggering a trap (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
			"holddown_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Time to hold the session-UP notification to the client (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255000),
				},
			},
			"minimum_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum transmit and receive interval (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255000),
				},
			},
			"minimum_receive_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum receive interval (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255000),
				},
			},
			"multiplier": schema.Int64Attribute{
				Optional:    true,
				Description: "Detection time multiplier (1..255).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"session_mode": schema.StringAttribute{
				Optional:    true,
				Description: "BFD single-hop or multihop session-mode.",
				Validators: []validator.String{
					stringvalidator.OneOf("automatic", "multihop", "single-hop"),
				},
			},
			"transmit_interval_minimum_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum transmit interval (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255000),
				},
			},
			"transmit_interval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "High transmit interval triggering a trap (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "BFD protocol version number.",
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "automatic"),
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (block *bgpBlockBfdLivenessDetection) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "bfd-liveness-detection "

	if v := block.AuthenticationAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication algorithm "+v)
	}
	if v := block.AuthenticationKeyChain.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication key-chain \""+v+"\"")
	}
	if block.AuthenticationLooseCheck.ValueBool() {
		configSet = append(configSet, setPrefix+"authentication loose-check")
	}
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
	if v := block.SessionMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"session-mode "+v)
	}
	if !block.TransmitIntervalMinimumInterval.IsNull() {
		configSet = append(configSet, setPrefix+"transmit-interval minimum-interval "+
			utils.ConvI64toa(block.TransmitIntervalMinimumInterval.ValueInt64()))
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

func (block *bgpBlockBfdLivenessDetection) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
		block.AuthenticationAlgorithm = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
		block.AuthenticationKeyChain = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "authentication loose-check":
		block.AuthenticationLooseCheck = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		block.DetectionTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
		block.HolddownInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		block.MinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		block.MinimumReceiveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		block.Multiplier, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "session-mode "):
		block.SessionMode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		block.TransmitIntervalMinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		block.TransmitIntervalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "version "):
		block.Version = types.StringValue(itemTrim)
	}

	return nil
}

type bgpBlockBgpErrorTolerance struct {
	MalformedRouteLimit        types.Int64 `tfsdk:"malformed_route_limit"`
	NoMalformedRouteLimit      types.Bool  `tfsdk:"no_malformed_route_limit"`
	MalformedUpdateLogInterval types.Int64 `tfsdk:"malformed_update_log_interval"`
}

func (bgpBlockBgpErrorTolerance) resourceSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Handle BGP malformed updates softly.",
		Attributes: map[string]schema.Attribute{
			"malformed_route_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of malformed routes from a peer (0..4294967295).",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"malformed_update_log_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Time used when logging malformed update (10..65535 seconds).",
				Validators: []validator.Int64{
					int64validator.Between(10, 65535),
				},
			},
			"no_malformed_route_limit": schema.BoolAttribute{
				Optional:    true,
				Description: "No malformed route limit.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (block *bgpBlockBgpErrorTolerance) configSet(setPrefix string) []string {
	setPrefix += "bgp-error-tolerance"

	configSet := []string{
		setPrefix,
	}

	if !block.MalformedRouteLimit.IsNull() {
		configSet = append(configSet, setPrefix+" malformed-route-limit "+
			utils.ConvI64toa(block.MalformedRouteLimit.ValueInt64()))
	}
	if !block.MalformedUpdateLogInterval.IsNull() {
		configSet = append(configSet, setPrefix+" malformed-update-log-interval "+
			utils.ConvI64toa(block.MalformedUpdateLogInterval.ValueInt64()))
	}
	if block.NoMalformedRouteLimit.ValueBool() {
		configSet = append(configSet, setPrefix+" no-malformed-route-limit")
	}

	return configSet
}

func (block *bgpBlockBgpErrorTolerance) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, " malformed-route-limit "):
		block.MalformedRouteLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, " malformed-update-log-interval "):
		block.MalformedUpdateLogInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == " no-malformed-route-limit":
		block.NoMalformedRouteLimit = types.BoolValue(true)
	}

	return nil
}

type bgpBlockBgpMultipath struct {
	AllowProtection types.Bool `tfsdk:"allow_protection"`
	Disable         types.Bool `tfsdk:"disable"`
	MultipleAS      types.Bool `tfsdk:"multiple_as"`
}

func (bgpBlockBgpMultipath) resourceSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Allow load sharing among multiple BGP paths.",
		Attributes: map[string]schema.Attribute{
			"allow_protection": schema.BoolAttribute{
				Optional:    true,
				Description: "Allows the BGP multipath and protection to co-exist.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable Multipath.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"multiple_as": schema.BoolAttribute{
				Optional:    true,
				Description: "Use paths received from different ASs.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (block *bgpBlockBgpMultipath) configSet(setPrefix string) []string {
	setPrefix += "multipath"

	configSet := []string{
		setPrefix,
	}

	if block.AllowProtection.ValueBool() {
		configSet = append(configSet, setPrefix+" allow-protection")
	}
	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+" disable")
	}
	if block.MultipleAS.ValueBool() {
		configSet = append(configSet, setPrefix+" multiple-as")
	}

	return configSet
}

func (block *bgpBlockBgpMultipath) read(itemTrim string) {
	switch {
	case itemTrim == " allow-protection":
		block.AllowProtection = types.BoolValue(true)
	case itemTrim == " disable":
		block.Disable = types.BoolValue(true)
	case itemTrim == " multiple-as":
		block.MultipleAS = types.BoolValue(true)
	}
}

type bgpBlockFamily struct {
	NlriType            types.String                    `tfsdk:"nlri_type"`
	AcceptedPrefixLimit *bgpBlockFamilyBlockPrefixLimit `tfsdk:"accepted_prefix_limit"`
	PrefixLimit         *bgpBlockFamilyBlockPrefixLimit `tfsdk:"prefix_limit"`
}

func (bgpBlockFamily) resourceSchema(family string) schema.ListNestedBlock {
	attributes := map[string]schema.Attribute{
		"nlri_type": schema.StringAttribute{
			Required:    true,
			Description: "NLRI type.",
			Validators: []validator.String{
				stringvalidator.OneOf("any", "flow", "labeled-unicast", "unicast", "multicast"),
			},
		},
	}
	if family == "EVPN" {
		attributes = map[string]schema.Attribute{
			"nlri_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("signaling"),
				Description: "NLRI type.",
				Validators: []validator.String{
					stringvalidator.OneOf("signaling"),
				},
			},
		}
	}

	return schema.ListNestedBlock{
		Description: "For each `nlri_type`, configure " + family + " NLRI parameters.",
		NestedObject: schema.NestedBlockObject{
			Attributes: attributes,
			Blocks: map[string]schema.Block{
				"accepted_prefix_limit": bgpBlockFamilyBlockPrefixLimit{}.resourceSchema(true),
				"prefix_limit":          bgpBlockFamilyBlockPrefixLimit{}.resourceSchema(false),
			},
		},
	}
}

type bgpBlockFamilyBlockPrefixLimit struct {
	Maximum                    types.Int64 `tfsdk:"maximum"`
	Teardown                   types.Int64 `tfsdk:"teardown"`
	TeardownIdleTimeout        types.Int64 `tfsdk:"teardown_idle_timeout"`
	TeardownIdleTimeoutForever types.Bool  `tfsdk:"teardown_idle_timeout_forever"`
}

func (block bgpBlockFamilyBlockPrefixLimit) resourceSchema(accepted bool) schema.SingleNestedBlock {
	description := "Define maximum number of prefixes from a peer."
	if accepted {
		description = "Define maximum number of prefixes accepted from a peer."
	}

	return schema.SingleNestedBlock{
		Description: description,
		Attributes:  block.resourceSchemaAttributes(),
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (bgpBlockFamilyBlockPrefixLimit) resourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"maximum": schema.Int64Attribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Maximum number of prefixes accepted from a peer.",
			Validators: []validator.Int64{
				int64validator.Between(1, 4294967295),
			},
		},
		"teardown": schema.Int64Attribute{
			Optional:    true,
			Description: "Clear peer connection on reaching limit with this percentage of prefix-limit to start warnings.",
			Validators: []validator.Int64{
				int64validator.Between(1, 100),
			},
		},
		"teardown_idle_timeout": schema.Int64Attribute{
			Optional:    true,
			Description: "Timeout before attempting to restart peer.",
			Validators: []validator.Int64{
				int64validator.Between(1, 2400),
			},
		},
		"teardown_idle_timeout_forever": schema.BoolAttribute{
			Optional:    true,
			Description: "Idle the peer until the user intervenes.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (block *bgpBlockFamily) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += block.NlriType.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if block.AcceptedPrefixLimit != nil {
		if !block.AcceptedPrefixLimit.Maximum.IsNull() {
			configSet = append(configSet, setPrefix+"accepted-prefix-limit maximum "+
				utils.ConvI64toa(block.AcceptedPrefixLimit.Maximum.ValueInt64()))
		}
		if !block.AcceptedPrefixLimit.Teardown.IsNull() {
			configSet = append(configSet, setPrefix+"accepted-prefix-limit teardown "+
				utils.ConvI64toa(block.AcceptedPrefixLimit.Teardown.ValueInt64()))
		}
		if !block.AcceptedPrefixLimit.TeardownIdleTimeout.IsNull() {
			if block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.ValueBool() {
				return configSet,
					pathRoot.AtName("accepted_prefix_limit").AtName("teardown_idle_timeout_forever"),
					errors.New("conflict between teardown_idle_timeout and teardown_idle_timeout_forever")
			}
			configSet = append(configSet, setPrefix+"accepted-prefix-limit teardown idle-timeout "+
				utils.ConvI64toa(block.AcceptedPrefixLimit.TeardownIdleTimeout.ValueInt64()))
		}
		if block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.ValueBool() {
			configSet = append(configSet, setPrefix+"accepted-prefix-limit teardown idle-timeout forever")
		}
	}
	if block.PrefixLimit != nil {
		if !block.PrefixLimit.Maximum.IsNull() {
			configSet = append(configSet, setPrefix+"prefix-limit maximum "+
				utils.ConvI64toa(block.PrefixLimit.Maximum.ValueInt64()))
		}
		if !block.PrefixLimit.Teardown.IsNull() {
			configSet = append(configSet, setPrefix+"prefix-limit teardown "+
				utils.ConvI64toa(block.PrefixLimit.Teardown.ValueInt64()))
		}
		if !block.PrefixLimit.TeardownIdleTimeout.IsNull() {
			if block.PrefixLimit.TeardownIdleTimeoutForever.ValueBool() {
				return configSet,
					pathRoot.AtName("prefix_limit").AtName("teardown_idle_timeout_forever"),
					errors.New("conflict between teardown_idle_timeout and teardown_idle_timeout_forever")
			}
			configSet = append(configSet, setPrefix+"prefix-limit teardown idle-timeout "+
				utils.ConvI64toa(block.PrefixLimit.TeardownIdleTimeout.ValueInt64()))
		}
		if block.PrefixLimit.TeardownIdleTimeoutForever.ValueBool() {
			configSet = append(configSet, setPrefix+"prefix-limit teardown idle-timeout forever")
		}
	}

	return configSet, path.Empty(), nil
}

func (block *bgpBlockFamily) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "accepted-prefix-limit "):
		if block.AcceptedPrefixLimit == nil {
			block.AcceptedPrefixLimit = &bgpBlockFamilyBlockPrefixLimit{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "maximum "):
			block.AcceptedPrefixLimit.Maximum, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "teardown idle-timeout "):
			if itemTrim == "forever" {
				block.AcceptedPrefixLimit.TeardownIdleTimeoutForever = types.BoolValue(true)
			} else {
				block.AcceptedPrefixLimit.TeardownIdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "teardown "):
			block.AcceptedPrefixLimit.Teardown, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "prefix-limit "):
		if block.PrefixLimit == nil {
			block.PrefixLimit = &bgpBlockFamilyBlockPrefixLimit{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "maximum "):
			block.PrefixLimit.Maximum, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "teardown idle-timeout "):
			if itemTrim == "forever" {
				block.PrefixLimit.TeardownIdleTimeoutForever = types.BoolValue(true)
			} else {
				block.PrefixLimit.TeardownIdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "teardown "):
			block.PrefixLimit.Teardown, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type bgpBlockGracefulRestart struct {
	Disable        types.Bool  `tfsdk:"disable"`
	RestartTime    types.Int64 `tfsdk:"restart_time"`
	StaleRouteTime types.Int64 `tfsdk:"stale_route_time"`
}

func (bgpBlockGracefulRestart) resourceSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Define BGP graceful restart options.",
		Attributes: map[string]schema.Attribute{
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable graceful restart.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"restart_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Restart time used when negotiating with a peer (1..600).",
				Validators: []validator.Int64{
					int64validator.Between(1, 600),
				},
			},
			"stale_route_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum time for which stale routes are kept (1..600).",
				Validators: []validator.Int64{
					int64validator.Between(1, 600),
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (block *bgpBlockGracefulRestart) configSet(setPrefix string) []string {
	setPrefix += "graceful-restart"

	configSet := []string{
		setPrefix,
	}

	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+" disable")
	}
	if !block.RestartTime.IsNull() {
		configSet = append(configSet, setPrefix+" restart-time "+
			utils.ConvI64toa(block.RestartTime.ValueInt64()))
	}
	if !block.StaleRouteTime.IsNull() {
		configSet = append(configSet, setPrefix+" stale-routes-time "+
			utils.ConvI64toa(block.StaleRouteTime.ValueInt64()))
	}

	return configSet
}

func (block *bgpBlockGracefulRestart) read(itemTrim string) (err error) {
	switch {
	case itemTrim == " "+junos.DisableW:
		block.Disable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, " restart-time "):
		block.RestartTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, " stale-routes-time "):
		block.StaleRouteTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}
