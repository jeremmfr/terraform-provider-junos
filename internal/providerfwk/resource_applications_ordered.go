package providerfwk

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &applicationsOrdered{}
	_ resource.ResourceWithConfigure      = &applicationsOrdered{}
	_ resource.ResourceWithValidateConfig = &applicationsOrdered{}
	_ resource.ResourceWithImportState    = &applicationsOrdered{}
)

type applicationsOrdered struct {
	client *junos.Client
}

func newApplicationsOrderedResource() resource.Resource {
	return &applicationsOrdered{}
}

func (rsc *applicationsOrdered) typeName() string {
	return providerName + "_applications_ordered"
}

func (rsc *applicationsOrdered) junosName() string {
	return "applications"
}

func (rsc *applicationsOrdered) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *applicationsOrdered) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *applicationsOrdered) Configure(
	ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedResourceConfigureType(ctx, req, resp)

		return
	}
	rsc.client = client
}

func (rsc *applicationsOrdered) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure entirely `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `applications`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"application": schema.ListNestedBlock{
				Description: "For each name, define an application.",
				NestedObject: schema.NestedBlockObject{
					Attributes: applicationAttrData{}.attributesSchema(),
					Blocks:     applicationAttrData{}.blocksSchema(),
				},
			},
			"application_set": schema.ListNestedBlock{
				Description: "For each name, define an application set.",
				NestedObject: schema.NestedBlockObject{
					Attributes: applicationSetAttrData{}.attributesSchema(),
				},
			},
		},
	}
}

type applicationsOrderedConfig struct {
	ID             types.String `tfsdk:"id"`
	Application    types.List   `tfsdk:"application"`
	ApplicationSet types.List   `tfsdk:"application_set"`
}

func (rsc *applicationsOrdered) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config applicationsOrderedConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationName := make(map[string]struct{})
	if !config.Application.IsNull() &&
		!config.Application.IsUnknown() {
		var configApplication []applicationAttrConfig
		asDiags := config.Application.ElementsAs(ctx, &configApplication, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		for i, block := range configApplication {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := applicationName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple application blocks with the same name %q", name),
					)
				}
				applicationName[name] = struct{}{}
			}

			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("application").AtListIndex(i).AtName("*"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of arguments need to be set (in addition to `name`)"+
						" in application block %q", block.Name.ValueString()),
				)
			}
			rootPath := path.Root("application").AtListIndex(i)
			block.validateConfig(
				ctx,
				&rootPath,
				fmt.Sprintf(" in application block %q", block.Name.ValueString()),
				resp,
			)
		}
	}

	if !config.ApplicationSet.IsNull() &&
		!config.ApplicationSet.IsUnknown() {
		var configApplicationSet []applicationSetAttrConfig
		asDiags := config.ApplicationSet.ElementsAs(ctx, &configApplicationSet, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		applicationSetName := make(map[string]struct{})
		for i, block := range configApplicationSet {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := applicationSetName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application_set").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple application_set blocks with the same name %q", name),
					)
				}
				applicationSetName[name] = struct{}{}
				if _, ok := applicationName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("application and application_set blocks with the same name %q", name),
					)
				}
			}

			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("application_set").AtListIndex(i).AtName("*"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of applications, application_set or description must be specified"+
						" in application_set block %q", block.Name.ValueString()),
				)
			}
			rootPath := path.Root("application_set").AtListIndex(i)
			block.validateConfig(
				ctx,
				&rootPath,
				fmt.Sprintf(" in application_set block %q", block.Name.ValueString()),
				resp,
			)
		}
	}
}

func (rsc *applicationsOrdered) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan applicationsData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *applicationsOrdered) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data applicationsData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		nil,
		resp,
	)
}

func (rsc *applicationsOrdered) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state applicationsData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *applicationsOrdered) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state applicationsData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceDelete(
		ctx,
		rsc,
		&state,
		resp,
	)
}

func (rsc *applicationsOrdered) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data applicationsData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}
