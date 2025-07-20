package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &applications{}
	_ resource.ResourceWithConfigure      = &applications{}
	_ resource.ResourceWithValidateConfig = &applications{}
	_ resource.ResourceWithImportState    = &applications{}
)

type applications struct {
	client *junos.Client
}

func newApplicationsResource() resource.Resource {
	return &applications{}
}

func (rsc *applications) typeName() string {
	return providerName + "_applications"
}

func (rsc *applications) junosName() string {
	return "applications"
}

func (rsc *applications) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *applications) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *applications) Configure(
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

func (rsc *applications) Schema(
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
			"application": schema.SetNestedBlock{
				Description: "For each name, define an application.",
				NestedObject: schema.NestedBlockObject{
					Attributes: applicationAttrData{}.attributesSchema(),
					Blocks:     applicationAttrData{}.blocksSchema(),
				},
			},
			"application_set": schema.SetNestedBlock{
				Description: "For each name, define an application set.",
				NestedObject: schema.NestedBlockObject{
					Attributes: applicationSetAttrData{}.attributesSchema(),
				},
			},
		},
	}
}

type applicationsData struct {
	ID             types.String             `tfsdk:"id"`
	Application    []applicationAttrData    `tfsdk:"application"`
	ApplicationSet []applicationSetAttrData `tfsdk:"application_set"`
}

type applicationsConfig struct {
	ID             types.String `tfsdk:"id"`
	Application    types.Set    `tfsdk:"application"`
	ApplicationSet types.Set    `tfsdk:"application_set"`
}

func (rsc *applications) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config applicationsConfig
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

		for _, block := range configApplication {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := applicationName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple application blocks with the same name %q", name),
					)
				}
				applicationName[name] = struct{}{}
			}

			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("application"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of arguments need to be set (in addition to `name`)"+
						" in application block %q", block.Name.ValueString()),
				)
			}
			rootPath := path.Root("application")
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
		for _, block := range configApplicationSet {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := applicationSetName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application_set"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple application_set blocks with the same name %q", name),
					)
				}
				applicationSetName[name] = struct{}{}
				if _, ok := applicationName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("application"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("application and application_set blocks with the same name %q", name),
					)
				}
			}

			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("application_set"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of applications, application_set or description must be specified"+
						" in application_set block %q", block.Name.ValueString()),
				)
			}
			rootPath := path.Root("application_set")
			block.validateConfig(
				ctx,
				&rootPath,
				fmt.Sprintf(" in application_set block %q", block.Name.ValueString()),
				resp,
			)
		}
	}
}

func (rsc *applications) Create(
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

func (rsc *applications) Read(
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

func (rsc *applications) Update(
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

func (rsc *applications) Delete(
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

func (rsc *applications) ImportState(
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

func (rscData *applicationsData) fillID() {
	rscData.ID = types.StringValue("applications")
}

func (rscData *applicationsData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *applicationsData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)

	applicationName := make(map[string]struct{})
	for i, block := range rscData.Application {
		name := block.Name.ValueString()
		if name == "" {
			return path.Root("application").AtListIndex(i).AtName("name"),
				errors.New("name argument in application block is empty")
		}
		if _, ok := applicationName[name]; ok {
			return path.Root("application").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple application blocks with the same name %q", name)
		}
		applicationName[name] = struct{}{}

		blockErrorSuffix := fmt.Sprintf(" in application block %q", name)
		if block.isEmpty() {
			return path.Root("application").AtListIndex(i).AtName("name"),
				errors.New("at least one of arguments need to be set (in addition to `name`)" +
					blockErrorSuffix)
		}

		dataConfigSet, _, err := block.configSet(blockErrorSuffix)
		if err != nil {
			return path.Root("application").AtListIndex(i).AtName("name"), err
		}
		configSet = append(configSet, dataConfigSet...)
	}
	applicationSetName := make(map[string]struct{})
	for i, block := range rscData.ApplicationSet {
		name := block.Name.ValueString()
		if name == "" {
			return path.Root("application_set").AtListIndex(i).AtName("name"),
				errors.New("name argument in application_set block is empty")
		}
		if _, ok := applicationSetName[name]; ok {
			return path.Root("application_set").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple application_set blocks with the same name %q", name)
		}
		if _, ok := applicationName[name]; ok {
			return path.Root("application_set").AtListIndex(i).AtName("name"),
				fmt.Errorf("application and application_set blocks with the same name %q", name)
		}
		applicationSetName[name] = struct{}{}

		if block.isEmpty() {
			return path.Root("application_set").AtListIndex(i).AtName("name"),
				fmt.Errorf("at least one of applications, application_set or description must be specified"+
					" in application_set block %q", name)
		}

		configSet = append(configSet, block.configSet()...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *applicationsData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"applications " + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "application "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var application applicationAttrData
				rscData.Application, application = tfdata.ExtractBlock(rscData.Application, types.StringValue(name))
				balt.CutPrefixInString(&itemTrim, name+" ")

				if err := application.read(itemTrim); err != nil {
					return err
				}
				rscData.Application = append(rscData.Application, application)
			case balt.CutPrefixInString(&itemTrim, "application-set "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var applicationSet applicationSetAttrData
				rscData.ApplicationSet, applicationSet = tfdata.ExtractBlock(rscData.ApplicationSet, types.StringValue(name))
				balt.CutPrefixInString(&itemTrim, name+" ")

				applicationSet.read(itemTrim)
				rscData.ApplicationSet = append(rscData.ApplicationSet, applicationSet)
			}
		}
	}

	return nil
}

func (rscData *applicationsData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete applications",
	}

	return junSess.ConfigSet(configSet)
}
