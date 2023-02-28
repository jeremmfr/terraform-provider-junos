package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	_ resource.Resource                = &securityIpsecProposal{}
	_ resource.ResourceWithConfigure   = &securityIpsecProposal{}
	_ resource.ResourceWithImportState = &securityIpsecProposal{}
)

type securityIpsecProposal struct {
	client *junos.Client
}

func newSecurityIpsecProposalResource() resource.Resource {
	return &securityIpsecProposal{}
}

func (rsc *securityIpsecProposal) typeName() string {
	return providerName + "_security_ipsec_proposal"
}

func (rsc *securityIpsecProposal) junosName() string {
	return "security ipsec proposal"
}

func (rsc *securityIpsecProposal) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIpsecProposal) Configure(
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

func (rsc *securityIpsecProposal) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Provides a " + rsc.junosName() + ".",
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
				Description: "The name of IPSec proposal.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					newStringFormatValidator(defaultFormat),
				},
			},
			"authentication_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "Authentication algorithm.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					newStringSpaceExclusionValidator(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of IPSec proposal.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					newStringDoubleQuoteExclusionValidator(),
				},
			},
			"encryption_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "Encryption algorithm.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					newStringSpaceExclusionValidator(),
				},
			},
			"lifetime_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "Lifetime, in seconds.",
				Validators: []validator.Int64{
					int64validator.Between(180, 86400),
				},
			},
			"lifetime_kilobytes": schema.Int64Attribute{
				Optional:    true,
				Description: "Lifetime, in kilobytes.",
				Validators: []validator.Int64{
					int64validator.Between(64, 4294967294),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "IPSec protocol. ",
				Validators: []validator.String{
					stringvalidator.OneOf("ah", "esp"),
				},
			},
		},
	}
}

type securityIpsecProposalData struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	Description             types.String `tfsdk:"description"`
	LifetimeSeconds         types.Int64  `tfsdk:"lifetime_seconds"`
	LifetimeKilobytes       types.Int64  `tfsdk:"lifetime_kilobytes"`
	Protocol                types.String `tfsdk:"protocol"`
}

func (rsc *securityIpsecProposal) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIpsecProposalData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			"could not create "+rsc.junosName()+" with empty name",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		resp.Diagnostics.AddError(
			"Compatibility Error",
			fmt.Sprintf(rsc.junosName()+" not compatible "+
				"with Junos device %q", junSess.SystemInformation.HardwareModel),
		)

		return
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	proposalExists, err := checkSecurityIpsecProposalExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if proposalExists {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists", plan.Name.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(diagWarns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	proposalExists, err = checkSecurityIpsecProposalExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !proposalExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q not exists after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityIpsecProposal) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIpsecProposalData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(ctx, state.Name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *securityIpsecProposal) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIpsecProposalData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(diagWarns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityIpsecProposal) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIpsecProposalData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(diagWarns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(diagWarns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *securityIpsecProposal) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	proposalExists, err := checkSecurityIpsecProposalExists(ctx, req.ID, junSess)
	if err != nil {
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if !proposalExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <name>)", req.ID),
		)

		return
	}

	var data securityIpsecProposalData
	err = data.read(ctx, req.ID, junSess)
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if !data.ID.IsNull() {
		resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	}
}

func checkSecurityIpsecProposalExists(_ context.Context, name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec proposal " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIpsecProposalData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIpsecProposalData) set(_ context.Context, junSess *junos.Session) (path.Path, error) {
	configSet := make([]string, 0)

	setPrefix := "set security ipsec proposal " + rscData.Name.ValueString() + " "
	if v := rscData.AuthenticationAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-algorithm "+v)
	}
	if v := rscData.EncryptionAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"encryption-algorithm "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.LifetimeSeconds.IsNull() {
		configSet = append(configSet, setPrefix+"lifetime-seconds "+
			utils.ConvI64toa(rscData.LifetimeSeconds.ValueInt64()))
	}
	if !rscData.LifetimeKilobytes.IsNull() {
		configSet = append(configSet, setPrefix+"lifetime-kilobytes "+
			utils.ConvI64toa(rscData.LifetimeKilobytes.ValueInt64()))
	}
	if v := rscData.Protocol.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityIpsecProposalData) read(_ context.Context, name string, junSess *junos.Session,
) (err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec proposal " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "authentication-algorithm "):
				rscData.AuthenticationAlgorithm = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "encryption-algorithm "):
				rscData.EncryptionAlgorithm = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "lifetime-kilobytes "):
				rscData.LifetimeKilobytes, err = convAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "lifetime-seconds "):
				rscData.LifetimeSeconds, err = convAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "protocol "):
				rscData.Protocol = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityIpsecProposalData) del(_ context.Context, junSess *junos.Session) error {
	configSet := []string{
		"delete security ipsec proposal " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}