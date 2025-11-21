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
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ action.Action              = &commitFileAction{}
	_ action.ActionWithConfigure = &commitFileAction{}
)

type commitFileAction struct {
	client *junos.Client
}

func newCommitFileAction() action.Action {
	return &commitFileAction{}
}

func (act *commitFileAction) typeName() string {
	return providerName + "_commit_file"
}

func (act *commitFileAction) junosClient() *junos.Client {
	return act.client
}

func (act *commitFileAction) Metadata(
	_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_commit_file"
}

func (act *commitFileAction) Configure(
	ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedActionConfigureType(ctx, req, resp)

		return
	}
	act.client = client
}

func (act *commitFileAction) Schema(
	_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Load a file with set/delete lines on device and commit.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "The path of the file to load.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"append_lines": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of lines to append to the lines in the loaded file.",
			},
			"clear_file_after_commit": schema.BoolAttribute{
				Optional:    true,
				Description: "Truncate file after successful commit.",
			},
		},
	}
}

type commitFileActionData struct {
	Filename             types.String   `tfsdk:"filename"`
	AppendLines          []types.String `tfsdk:"append_lines"`
	ClearFileAfterCommit types.Bool     `tfsdk:"clear_file_after_commit"`
}

func (act *commitFileAction) Invoke(
	ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse,
) {
	var config commitFileActionData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Reading configuration file",
	})
	configSet, err := config.readFile()
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("filename"), tfdiag.ConfigSetErrSummary, err.Error())

		return
	}
	for _, v := range config.AppendLines {
		configSet = append(configSet, v.ValueString())
	}

	clt := act.junosClient()
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Starting session to device",
	})
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Locking candidate configuration",
	})
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
	}()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Loading configuration",
	})
	if err := junSess.ConfigSet(configSet); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Committing configuration",
	})
	warns, err := junSess.CommitConf(ctx, "commit a file with action "+act.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Configuration loaded and committed",
	})

	if config.ClearFileAfterCommit.ValueBool() {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: "Clearing file after commit",
		})
		if err := config.cleanFile(); err != nil {
			resp.Diagnostics.AddWarning("Post Clean Error", err.Error())
		}
	}
}

func (actData *commitFileActionData) readFile() ([]string, error) {
	filename := actData.Filename.ValueString()
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

func (actData *commitFileActionData) cleanFile() error {
	filename := actData.Filename.ValueString()
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
