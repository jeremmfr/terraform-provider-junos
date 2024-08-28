package providerfwk

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type vstpVlanAttrData struct {
	BackupBridgePriority types.String `tfsdk:"backup_bridge_priority"`
	BridgePriority       types.String `tfsdk:"bridge_priority"`
	ForwardDelay         types.Int64  `tfsdk:"forward_delay"`
	HelloTime            types.Int64  `tfsdk:"hello_time"`
	MaxAge               types.Int64  `tfsdk:"max_age"`
	SystemIdentifier     types.String `tfsdk:"system_identifier"`
}

func (rscData vstpVlanAttrData) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"backup_bridge_priority": schema.StringAttribute{
			Optional:    true,
			Description: "Priority of the bridge.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^\d\d?k$`),
					"must be a number with increments of 4k - 4k,8k,..60k",
				),
			},
		},
		"bridge_priority": schema.StringAttribute{
			Optional:    true,
			Description: "Priority of the bridge.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^(0|\d\d?k)$`),
					"must be a number with increments of 4k - 0,4k,8k,..60k",
				),
			},
		},
		"forward_delay": schema.Int64Attribute{
			Optional:    true,
			Description: "Time spent in listening or learning state (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(4, 30),
			},
		},
		"hello_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Time interval between configuration BPDUs (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(1, 10),
			},
		},
		"max_age": schema.Int64Attribute{
			Optional:    true,
			Description: "Maximum age of received protocol bpdu (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(6, 40),
			},
		},
		"system_identifier": schema.StringAttribute{
			Optional:    true,
			Description: "System identifier to represent this node.",
			Validators: []validator.String{
				tfvalidator.StringMACAddress().WithMac48ColonHexa(),
			},
		},
	}
}

func (rscData *vstpVlanAttrData) validateConfig(
	_ context.Context, resp *resource.ValidateConfigResponse,
) {
	if !rscData.BackupBridgePriority.IsNull() && !rscData.BackupBridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			rscData.BackupBridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be a multiple of 4k",
				)
			}
			if v < 4 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be between 4k and 60k",
				)
			}
			if !rscData.BridgePriority.IsNull() && !rscData.BridgePriority.IsUnknown() {
				if bridgePriority, err := strconv.Atoi(strings.TrimSuffix(
					rscData.BridgePriority.ValueString(), "k",
				)); err == nil {
					if v <= bridgePriority {
						resp.Diagnostics.AddAttributeError(
							path.Root("backup_bridge_priority"),
							"Bad Value Error",
							"backup_bridge_priority must be worse (higher value) than bridge_priority",
						)
					}
				}
			}
		}
	}
	if !rscData.BridgePriority.IsNull() && !rscData.BridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			rscData.BridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be a multiple of 4k",
				)
			}
			if v < 0 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be between 0 and 60k",
				)
			}
		}
	}
}

func (rscData *vstpVlanAttrData) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 10)

	if v := rscData.BackupBridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if v := rscData.BridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}
	if !rscData.ForwardDelay.IsNull() {
		configSet = append(configSet, setPrefix+"forward-delay "+
			utils.ConvI64toa(rscData.ForwardDelay.ValueInt64()))
	}
	if !rscData.HelloTime.IsNull() {
		configSet = append(configSet, setPrefix+"hello-time "+
			utils.ConvI64toa(rscData.HelloTime.ValueInt64()))
	}
	if !rscData.MaxAge.IsNull() {
		configSet = append(configSet, setPrefix+"max-age "+
			utils.ConvI64toa(rscData.MaxAge.ValueInt64()))
	}
	if v := rscData.SystemIdentifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"system-identifier "+v)
	}

	return configSet
}

func (rscData *vstpVlanAttrData) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "backup-bridge-priority "):
		rscData.BackupBridgePriority = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
		rscData.BridgePriority = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "forward-delay "):
		rscData.ForwardDelay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "hello-time "):
		rscData.HelloTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "max-age "):
		rscData.MaxAge, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "system-identifier "):
		rscData.SystemIdentifier = types.StringValue(itemTrim)
	}

	return nil
}

func (rscData vstpVlanAttrData) configOptsToDel(delPrefix string) []string {
	return []string{
		delPrefix + "backup-bridge-priority",
		delPrefix + "bridge-priority",
		delPrefix + "forward-delay",
		delPrefix + "hello-time",
		delPrefix + "max-age",
		delPrefix + "system-identifier",
	}
}
