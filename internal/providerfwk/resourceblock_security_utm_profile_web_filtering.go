package providerfwk

import (
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"

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
