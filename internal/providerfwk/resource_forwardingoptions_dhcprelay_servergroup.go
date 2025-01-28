package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &forwardingoptionsDhcprelayServergroup{}
	_ resource.ResourceWithConfigure   = &forwardingoptionsDhcprelayServergroup{}
	_ resource.ResourceWithImportState = &forwardingoptionsDhcprelayServergroup{}
)

type forwardingoptionsDhcprelayServergroup struct {
	client *junos.Client
}

func newForwardingoptionsDhcprelayServergroupResource() resource.Resource {
	return &forwardingoptionsDhcprelayServergroup{}
}

func (rsc *forwardingoptionsDhcprelayServergroup) typeName() string {
	return providerName + "_forwardingoptions_dhcprelay_servergroup"
}

func (rsc *forwardingoptionsDhcprelayServergroup) junosName() string {
	return "forwarding-options dhcp-relay server-group"
}

func (rsc *forwardingoptionsDhcprelayServergroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *forwardingoptionsDhcprelayServergroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsDhcprelayServergroup) Configure(
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

func (rsc *forwardingoptionsDhcprelayServergroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format" +
					" `<name>" + junos.IDSeparator + "<routing_instance>" + junos.IDSeparator + "<version>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Server group name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Version for DHCP or DHCPv6.",
				Default:     stringdefault.StaticString("v4"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("v4", "v6"),
				},
			},
			"ip_address": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "IP Addresses of DHCP servers.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress(),
					),
				},
			},
		},
	}
}

type forwardingoptionsDhcprelayServergroupData struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	RoutingInstance types.String   `tfsdk:"routing_instance"`
	Version         types.String   `tfsdk:"version"`
	IPAddress       []types.String `tfsdk:"ip_address"`
}

func (rsc *forwardingoptionsDhcprelayServergroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsDhcprelayServergroupData
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
			if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
				instanceExists, err := checkRoutingInstanceExists(fnCtx, v, junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !instanceExists {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("routing instance %q doesn't exist", v),
					)

					return false
				}
			}
			groupExists, err := checkForwardingoptionsDhcprelayServergroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Version.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				if plan.Version.ValueString() == "v6" {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"forwarding-options dhcp-relay dhcpv6 server-group %q already exists in routing-instance %q",
								plan.Name.ValueString(), v,
							),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"forwarding-options dhcp-relay dhcpv6 server-group %q already exists",
								plan.Name.ValueString(),
							),
						)
					}
				} else {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Name, v),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							defaultResourceAlreadyExistsMessage(rsc, plan.Name),
						)
					}
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkForwardingoptionsDhcprelayServergroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Version.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
				if plan.Version.ValueString() == "v6" {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"forwarding-options dhcp-relay dhcpv6 server-group %q does not exists in routing-instance %q after commit "+
									"=> check your config", plan.Name.ValueString(), v,
							),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("forwarding-options dhcp-relay dhcpv6 server-group %q does not exists after commit "+
								"=> check your config", plan.Name.ValueString()),
						)
					}
				} else {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Name, v),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
						)
					}
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *forwardingoptionsDhcprelayServergroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsDhcprelayServergroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom3String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
			state.Version.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *forwardingoptionsDhcprelayServergroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsDhcprelayServergroupData
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

func (rsc *forwardingoptionsDhcprelayServergroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsDhcprelayServergroupData
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

func (rsc *forwardingoptionsDhcprelayServergroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data forwardingoptionsDhcprelayServergroupData

	var _ resourceDataReadFrom3String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>"+junos.IDSeparator+"<version>)",
	)
}

func checkForwardingoptionsDhcprelayServergroupExists(
	_ context.Context, name, routingInstance, version string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "forwarding-options dhcp-relay "
	if version == "v6" {
		showPrefix += "dhcpv6 "
	}
	showConfig, err := junSess.Command(showPrefix +
		"server-group " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *forwardingoptionsDhcprelayServergroupData) fillID() {
	routingInstance := rscData.RoutingInstance.ValueString()
	version := rscData.Version.ValueString()
	if routingInstance == "" {
		routingInstance = junos.DefaultW
	}
	if version == "" {
		version = "v4"
	}

	rscData.ID = types.StringValue(rscData.Name.ValueString() +
		junos.IDSeparator + routingInstance +
		junos.IDSeparator + version)
}

func (rscData *forwardingoptionsDhcprelayServergroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *forwardingoptionsDhcprelayServergroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, len(rscData.IPAddress)+1)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "forwarding-options dhcp-relay "
	if rscData.Version.ValueString() == "v6" {
		setPrefix += "dhcpv6 "
	}
	setPrefix += "server-group " + rscData.Name.ValueString() + " "

	configSet = append(configSet, setPrefix)
	for _, v := range rscData.IPAddress {
		configSet = append(configSet, setPrefix+v.ValueString())
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *forwardingoptionsDhcprelayServergroupData) read(
	_ context.Context, name, routingInstance, version string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "forwarding-options dhcp-relay "
	if version == "v6" {
		showPrefix += "dhcpv6 "
	}
	showConfig, err := junSess.Command(showPrefix +
		"server-group " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
		if version == "v6" {
			rscData.Version = types.StringValue(version)
		} else {
			rscData.Version = types.StringValue("v4")
		}
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			if itemTrim != "" {
				rscData.IPAddress = append(rscData.IPAddress, types.StringValue(itemTrim))
			}
		}
	}

	return nil
}

func (rscData *forwardingoptionsDhcprelayServergroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "forwarding-options dhcp-relay "
	if rscData.Version.ValueString() == "v6" {
		delPrefix += "dhcpv6 "
	}

	configSet := []string{
		delPrefix + "server-group " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
