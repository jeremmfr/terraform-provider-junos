package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &eventoptionsDestination{}
	_ resource.ResourceWithConfigure      = &eventoptionsDestination{}
	_ resource.ResourceWithValidateConfig = &eventoptionsDestination{}
	_ resource.ResourceWithImportState    = &eventoptionsDestination{}
)

type eventoptionsDestination struct {
	client *junos.Client
}

func newEventoptionsDestinationResource() resource.Resource {
	return &eventoptionsDestination{}
}

func (rsc *eventoptionsDestination) typeName() string {
	return providerName + "_eventoptions_destination"
}

func (rsc *eventoptionsDestination) junosName() string {
	return "event-options destination"
}

func (rsc *eventoptionsDestination) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *eventoptionsDestination) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *eventoptionsDestination) Configure(
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

func (rsc *eventoptionsDestination) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Destination name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"transfer_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Delay before transferring files (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"archive_site": schema.ListNestedBlock{
				Description: "For each archive destination.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Required:    true,
							Description: "URL of destination for file.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
								tfvalidator.StringSpaceExclusion(),
							},
						},
						"password": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "Password for login into the archive site.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
		},
	}
}

type eventoptionsDestinationData struct {
	ID            types.String                              `tfsdk:"id"`
	Name          types.String                              `tfsdk:"name"`
	TransferDelay types.Int64                               `tfsdk:"transfer_delay"`
	ArchiveSite   []eventoptionsDestinationBlockArchiveSite `tfsdk:"archive_site"`
}

type eventoptionsDestinationConfig struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	TransferDelay types.Int64  `tfsdk:"transfer_delay"`
	ArchiveSite   types.List   `tfsdk:"archive_site"`
}

type eventoptionsDestinationBlockArchiveSite struct {
	URL      types.String `tfsdk:"url"`
	Password types.String `tfsdk:"password"`
}

func (rsc *eventoptionsDestination) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config eventoptionsDestinationConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ArchiveSite.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("archive_site"),
			tfdiag.MissingConfigErrSummary,
			"archive_site block must be specified",
		)
	} else if !config.ArchiveSite.IsUnknown() {
		var configArchiveSite []eventoptionsDestinationBlockArchiveSite
		asDiags := config.ArchiveSite.ElementsAs(ctx, &configArchiveSite, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		archiveSiteURL := make(map[string]struct{})
		for i, block := range configArchiveSite {
			if block.URL.IsUnknown() {
				continue
			}
			url := block.URL.ValueString()
			if _, ok := archiveSiteURL[url]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("archive_site").AtListIndex(i).AtName("url"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple archive_site blocks with the same url %q", url),
				)
			}
			archiveSiteURL[url] = struct{}{}
		}
	}
}

func (rsc *eventoptionsDestination) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan eventoptionsDestinationData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "name"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			destinationExists, err := checkEventoptionsDestinationExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if destinationExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			destinationExists, err := checkEventoptionsDestinationExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !destinationExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *eventoptionsDestination) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data eventoptionsDestinationData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *eventoptionsDestination) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state eventoptionsDestinationData
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

func (rsc *eventoptionsDestination) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state eventoptionsDestinationData
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

func (rsc *eventoptionsDestination) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data eventoptionsDestinationData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
	)
}

func checkEventoptionsDestinationExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options destinations \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *eventoptionsDestinationData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *eventoptionsDestinationData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *eventoptionsDestinationData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set event-options destinations \"" + rscData.Name.ValueString() + "\" "

	archiveSiteURL := make(map[string]struct{})
	for i, block := range rscData.ArchiveSite {
		url := block.URL.ValueString()
		if _, ok := archiveSiteURL[url]; ok {
			return path.Root("archive_site").AtListIndex(i).AtName("url"),
				fmt.Errorf("multiple archive_site blocks with the same url %q", url)
		}
		archiveSiteURL[url] = struct{}{}

		configSet = append(configSet, setPrefix+"archive-sites \""+url+"\"")
		if v := block.Password.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"archive-sites \""+url+"\" password \""+v+"\"")
		}
	}
	if !rscData.TransferDelay.IsNull() {
		configSet = append(configSet, setPrefix+"transfer-delay "+
			utils.ConvI64toa(rscData.TransferDelay.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *eventoptionsDestinationData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options destinations \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "archive-sites "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) > 2 { // <url> password <password>
					password, err := tfdata.JunosDecode(strings.Trim(itemTrimFields[2], "\""), "password")
					if err != nil {
						return err
					}
					rscData.ArchiveSite = append(rscData.ArchiveSite, eventoptionsDestinationBlockArchiveSite{
						URL:      types.StringValue(strings.Trim(itemTrimFields[0], "\"")),
						Password: password,
					})
				} else { // <url>
					rscData.ArchiveSite = append(rscData.ArchiveSite, eventoptionsDestinationBlockArchiveSite{
						URL: types.StringValue(strings.Trim(itemTrimFields[0], "\"")),
					})
				}
			case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
				rscData.TransferDelay, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *eventoptionsDestinationData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete event-options destinations \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
