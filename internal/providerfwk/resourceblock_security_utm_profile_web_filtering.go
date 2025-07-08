package providerfwk

import (
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type securityUtmProfileWebFilteringBlockFallbackSettings struct {
	Default            types.String `tfsdk:"default"`
	ServerConnectivity types.String `tfsdk:"server_connectivity"`
	Timeout            types.String `tfsdk:"timeout"`
	TooManyRequests    types.String `tfsdk:"too_many_requests"`
}

func (securityUtmProfileWebFilteringBlockFallbackSettings) schema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Configure fallback settings.",
		Attributes: map[string]schema.Attribute{
			"default": schema.StringAttribute{
				Optional:    true,
				Description: "Default action.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit"),
				},
			},
			"server_connectivity": schema.StringAttribute{
				Optional:    true,
				Description: "Action when device cannot connect to server.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit"),
				},
			},
			"timeout": schema.StringAttribute{
				Optional:    true,
				Description: "Action when connection to server timeout.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit"),
				},
			},
			"too_many_requests": schema.StringAttribute{
				Optional:    true,
				Description: "Action when requests exceed the limit of engine.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit"),
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			tfplanmodifier.BlockRemoveNull(),
		},
	}
}

func (block *securityUtmProfileWebFilteringBlockFallbackSettings) configSet(setPrefix string) []string {
	setPrefix += "fallback-settings "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if v := block.Default.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"default "+v)
	}
	if v := block.ServerConnectivity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"server-connectivity "+v)
	}
	if v := block.Timeout.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"timeout "+v)
	}
	if v := block.TooManyRequests.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"too-many-requests "+v)
	}

	return configSet
}

func (block *securityUtmProfileWebFilteringBlockFallbackSettings) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "default "):
		block.Default = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "server-connectivity "):
		block.ServerConnectivity = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "too-many-requests "):
		block.TooManyRequests = types.StringValue(itemTrim)
	}
}

type securityUtmProfileWebFilteringBlockCategoryCustom struct {
	Name          types.String `tfsdk:"name"           tfdata:"identifier"`
	Action        types.String `tfsdk:"action"`
	CustomMessage types.String `tfsdk:"custom_message"`
}

func (securityUtmProfileWebFilteringBlockCategoryCustom) schema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "For each name of category, configure custom category actions.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Required:    true,
					Description: "Name of category.",
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 59),
						tfvalidator.StringDoubleQuoteExclusion(),
					},
				},
				"action": schema.StringAttribute{
					Required:    true,
					Description: "Action when web traffic matches category.",
					Validators: []validator.String{
						stringvalidator.OneOf("block", "log-and-permit", "permit", "quarantine"),
					},
				},
				"custom_message": schema.StringAttribute{
					Optional:    true,
					Description: "Custom message.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					},
				},
			},
		},
	}
}

func (block *securityUtmProfileWebFilteringBlockCategoryCustom) configSet(setPrefix string) []string {
	setPrefix += "category \"" + block.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "action " + block.Action.ValueString()

	if v := block.CustomMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-message \""+v+"\"")
	}

	return configSet
}

func (block *securityUtmProfileWebFilteringBlockCategoryCustom) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "action "):
		block.Action = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "custom-message "):
		block.CustomMessage = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}
