package providerfwk

import (
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type ripBlockBfdLivenessDetection struct {
	AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
	AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
	AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
	DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
	MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                      types.Int64  `tfsdk:"multiplier"`
	NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
	TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
	TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                         types.String `tfsdk:"version"`
}

func (block *ripBlockBfdLivenessDetection) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *ripBlockBfdLivenessDetection) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (block *ripBlockBfdLivenessDetection) configSet(setPrefix string) []string {
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

func (block *ripBlockBfdLivenessDetection) read(itemTrim string) (err error) {
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
	case itemTrim == "no-adaptation":
		block.NoAdaptation = types.BoolValue(true)
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
