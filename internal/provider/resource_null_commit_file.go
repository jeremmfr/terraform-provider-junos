package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                 = &nullCommitFile{}
	_ resource.ResourceWithConfigure    = &nullCommitFile{}
	_ resource.ResourceWithUpgradeState = &nullCommitFile{}
)

type nullCommitFile struct {
	client *junos.Client
}

func newNullCommitFileResource() resource.Resource {
	return &nullCommitFile{}
}

func (rsc *nullCommitFile) typeName() string {
	return providerName + "_null_commit_file"
}

func (rsc *nullCommitFile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *nullCommitFile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *nullCommitFile) Configure(
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

func (rsc *nullCommitFile) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Load a file with set/delete lines on device and commit.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<filename>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "The path of the file to load.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"append_lines": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of lines append to lines in the loaded file.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"clear_file_after_commit": schema.BoolAttribute{
				Optional:    true,
				Description: "Truncate file after successful commit.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"triggers": schema.DynamicAttribute{
				Optional:    true,
				Description: "Any value that, when changed, will force the resource to be replaced.",
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type nullCommitFileData struct {
	ID                   types.String   `tfsdk:"id"`
	Filename             types.String   `tfsdk:"filename"`
	AppendLines          []types.String `tfsdk:"append_lines"`
	ClearFileAfterCommit types.Bool     `tfsdk:"clear_file_after_commit"`
	Triggers             types.Dynamic  `tfsdk:"triggers"`
}

func (rsc *nullCommitFile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan nullCommitFileData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clt := rsc.junosClient()
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
	}()

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}

	warns, err := junSess.CommitConf(ctx, "commit a file with resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ClearFileAfterCommit.ValueBool() {
		if err := plan.cleanFile(); err != nil {
			resp.Diagnostics.AddWarning("Post Clean Error", err.Error())
		}
	}
}

func (rsc *nullCommitFile) Read(
	_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse,
) {
	// no-op
}

func (rsc *nullCommitFile) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
	// no-op
}

func (rsc *nullCommitFile) Delete(
	_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse,
) {
	// no-op
}

func (rscData *nullCommitFileData) fillID() {
	rscData.ID = types.StringValue(rscData.Filename.ValueString())
}

func (rscData *nullCommitFileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet, err := rscData.readFile()
	if err != nil {
		return path.Root("filename"), err
	}

	for _, v := range rscData.AppendLines {
		configSet = append(configSet, v.ValueString())
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *nullCommitFileData) readFile() ([]string, error) {
	filename := rscData.Filename.ValueString()
	if err := utils.ReplaceTildeToHomeDir(&filename); err != nil {
		return []string{}, err
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("file %q doesn't exist", filename)
	}
	fileReadByte, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, fmt.Errorf("could not read file %q: %w", filename, err)
	}

	return strings.Split(string(fileReadByte), "\n"), nil
}

func (rscData *nullCommitFileData) cleanFile() error {
	filename := rscData.Filename.ValueString()
	if err := utils.ReplaceTildeToHomeDir(&filename); err != nil {
		return err
	}

	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	f, err := os.OpenFile(filename, os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("could not open file %q to truncate after commit: %w", filename, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file handler for %q after truncation: %w", filename, err)
	}

	return nil
}
