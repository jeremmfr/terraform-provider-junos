package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &applicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationsDataSource{}
)

type applicationsDataSource struct {
	client *junos.Client
}

func (dsc *applicationsDataSource) typeName() string {
	return providerName + "_applications"
}

func (dsc *applicationsDataSource) junosName() string {
	return "applications application"
}

func newApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

func (dsc *applicationsDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *applicationsDataSource) Configure(
	ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedDataSourceConfigureType(ctx, req, resp)

		return
	}
	dsc.client = client
}

func (dsc *applicationsDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get configurations from " + dsc.junosName() + " and junos-defaults applications application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source.",
			},
			"match_name": schema.StringAttribute{
				Optional:    true,
				Description: "A regexp to apply a filter on applications name.",
				Validators: []validator.String{
					tfvalidator.StringRegex(),
				},
			},
			"applications": schema.ListAttribute{
				Computed:    true,
				Description: "For each application found.",
				ElementType: types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
					"name":                     types.StringType,
					"application_protocol":     types.StringType,
					"description":              types.StringType,
					"destination_port":         types.StringType,
					"ether_type":               types.StringType,
					"inactivity_timeout":       types.Int64Type,
					"inactivity_timeout_never": types.BoolType,
					"protocol":                 types.StringType,
					"rpc_program_number":       types.StringType,
					"source_port":              types.StringType,
					"term": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"name":                     types.StringType,
						"alg":                      types.StringType,
						"destination_port":         types.StringType,
						"icmp_code":                types.StringType,
						"icmp_type":                types.StringType,
						"icmp6_code":               types.StringType,
						"icmp6_type":               types.StringType,
						"inactivity_timeout":       types.Int64Type,
						"inactivity_timeout_never": types.BoolType,
						"protocol":                 types.StringType,
						"rpc_program_number":       types.StringType,
						"source_port":              types.StringType,
						"uuid":                     types.StringType,
					})),
					"uuid": types.StringType,
				}),
			},
		},
		Blocks: map[string]schema.Block{
			"match_options": schema.SetNestedBlock{
				Description: "List of options to apply a filter on applications.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alg": schema.StringAttribute{
							Optional:    true,
							Description: "Application Layer Gateway.",
						},
						"application_protocol": schema.StringAttribute{
							Optional:    true,
							Description: "Application protocol type.",
						},
						"destination_port": schema.StringAttribute{
							Optional:    true,
							Description: "Match TCP/UDP destination port.",
						},
						"ether_type": schema.StringAttribute{
							Optional:    true,
							Description: "Match ether type.",
						},
						"icmp_code": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP message code.",
						},
						"icmp_type": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP message type.",
						},
						"icmp6_code": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP6 message code.",
						},
						"icmp6_type": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP6 message type.",
						},
						"inactivity_timeout": schema.Int64Attribute{
							Optional:    true,
							Description: "Application-specific inactivity timeout.",
						},
						"inactivity_timeout_never": schema.BoolAttribute{
							Optional:    true,
							Description: "Disables inactivity timeout.",
						},
						"protocol": schema.StringAttribute{
							Optional:    true,
							Description: "Match IP protocol type.",
						},
						"rpc_program_number": schema.StringAttribute{
							Optional:    true,
							Description: "Match range of RPC program numbers.",
						},
						"source_port": schema.StringAttribute{
							Optional:    true,
							Description: "Match TCP/UDP source port.",
						},
						"uuid": schema.StringAttribute{
							Optional:    true,
							Description: "Match universal unique identifier for DCE RPC objects.",
						},
					},
				},
			},
		},
	}
}

type applicationsDataSourceData struct {
	ID           types.String                              `tfsdk:"id"`
	MatchName    types.String                              `tfsdk:"match_name"`
	Applications []applicationsDataSourceBlockApplications `tfsdk:"applications"`
	MatchOptions []applicationsDataSourceBlockMatchOptions `tfsdk:"match_options"`
}

type applicationsDataSourceBlockApplications struct {
	Name                   types.String                                       `tfsdk:"name"`
	ApplicationProtocol    types.String                                       `tfsdk:"application_protocol"`
	Description            types.String                                       `tfsdk:"description"`
	DestinationPort        types.String                                       `tfsdk:"destination_port"`
	EtherType              types.String                                       `tfsdk:"ether_type"`
	InactivityTimeout      types.Int64                                        `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool                                         `tfsdk:"inactivity_timeout_never"`
	Protocol               types.String                                       `tfsdk:"protocol"`
	RPCRrogramNumber       types.String                                       `tfsdk:"rpc_program_number"`
	SourcePort             types.String                                       `tfsdk:"source_port"`
	Term                   []applicationsDataSourceBlockApplicationsBlockTerm `tfsdk:"term"`
	UUID                   types.String                                       `tfsdk:"uuid"`
}

type applicationsDataSourceBlockApplicationsBlockTerm struct {
	Name                   types.String `tfsdk:"name"`
	Alg                    types.String `tfsdk:"alg"`
	DestinationPort        types.String `tfsdk:"destination_port"`
	IcmpCode               types.String `tfsdk:"icmp_code"`
	IcmpType               types.String `tfsdk:"icmp_type"`
	Icmp6Code              types.String `tfsdk:"icmp6_code"`
	Icmp6Type              types.String `tfsdk:"icmp6_type"`
	InactivityTimeout      types.Int64  `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool   `tfsdk:"inactivity_timeout_never"`
	Protocol               types.String `tfsdk:"protocol"`
	RPCRrogramNumber       types.String `tfsdk:"rpc_program_number"`
	SourcePort             types.String `tfsdk:"source_port"`
	UUID                   types.String `tfsdk:"uuid"`
}

type applicationsDataSourceBlockMatchOptions struct {
	Alg                    types.String `tfsdk:"alg"`
	ApplicationProtocol    types.String `tfsdk:"application_protocol"`
	DestinationPort        types.String `tfsdk:"destination_port"`
	EtherType              types.String `tfsdk:"ether_type"`
	IcmpCode               types.String `tfsdk:"icmp_code"`
	IcmpType               types.String `tfsdk:"icmp_type"`
	Icmp6Code              types.String `tfsdk:"icmp6_code"`
	Icmp6Type              types.String `tfsdk:"icmp6_type"`
	InactivityTimeout      types.Int64  `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool   `tfsdk:"inactivity_timeout_never"`
	Protocol               types.String `tfsdk:"protocol"`
	RPCRrogramNumber       types.String `tfsdk:"rpc_program_number"`
	SourcePort             types.String `tfsdk:"source_port"`
	UUID                   types.String `tfsdk:"uuid"`
}

func (dsc *applicationsDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data applicationsDataSourceData
	var matchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match_name"), &matchName)...)
	var matchOptions []applicationsDataSourceBlockMatchOptions
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match_options"), &matchOptions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MatchName = matchName
	data.MatchOptions = matchOptions

	junSess, err := dsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	applicationMap, err := dsc.search(junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}

	if err := data.filter(applicationMap); err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}
	data.fillData(applicationMap)
	data.fillID()

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (dsc *applicationsDataSource) search(
	junSess *junos.Session,
) (
	map[string]applicationsDataSourceBlockApplications, error,
) {
	results := make(map[string]applicationsDataSourceBlockApplications, 0)
	for _, config := range []string{
		"groups junos-defaults applications",
		"applications",
	} {
		showConfig, err := junSess.Command(junos.CmdShowConfig + config + junos.PipeDisplaySetRelative)
		if err != nil {
			return results, err
		}
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			if item == "" {
				continue
			}
			if !strings.HasPrefix(item, "set application ") {
				continue
			}
			itemTrim := strings.TrimPrefix(item, "set application ")
			itemTrimFields := strings.Split(itemTrim, " ")
			if _, ok := results[itemTrimFields[0]]; !ok {
				results[itemTrimFields[0]] = applicationsDataSourceBlockApplications{
					Name: types.StringValue(itemTrimFields[0]),
				}
			}
			app := results[itemTrimFields[0]]
			if err := app.read(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")); err != nil {
				return results, err
			}
			results[itemTrimFields[0]] = app
		}
	}

	return results, nil
}

func (block *applicationsDataSourceBlockApplications) read(itemTrim string) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "application-protocol "):
		block.ApplicationProtocol = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		block.DestinationPort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "ether-type "):
		block.EtherType = types.StringValue(itemTrim)
	case itemTrim == "inactivity-timeout never":
		block.InactivityTimeoutNever = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
		var err error
		block.InactivityTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
		block.RPCRrogramNumber = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		block.SourcePort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "term "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var term applicationsDataSourceBlockApplicationsBlockTerm
		block.Term, term = tfdata.ExtractBlockWithTFTypesString(
			block.Term, "Name", itemTrimFields[0],
		)
		term.Name = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := term.read(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")); err != nil {
			return err
		}
		block.Term = append(block.Term, term)
	case balt.CutPrefixInString(&itemTrim, "uuid "):
		block.UUID = types.StringValue(itemTrim)
	}

	return nil
}

func (block *applicationsDataSourceBlockApplicationsBlockTerm) read(itemTrim string) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "alg "):
		block.Alg = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		block.DestinationPort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "icmp-code "):
		block.IcmpCode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp-type "):
		block.IcmpType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp6-code "):
		block.Icmp6Code = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp6-type "):
		block.Icmp6Type = types.StringValue(itemTrim)
	case itemTrim == "inactivity-timeout never":
		block.InactivityTimeoutNever = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
		var err error
		block.InactivityTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
		block.RPCRrogramNumber = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		block.SourcePort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "uuid "):
		block.UUID = types.StringValue(itemTrim)
	}

	return nil
}

func (dscData *applicationsDataSourceData) filter(
	results map[string]applicationsDataSourceBlockApplications,
) error {
	if v := dscData.MatchName.ValueString(); v != "" {
		for appKey, app := range results {
			matched, err := regexp.MatchString(v, app.Name.ValueString())
			if err != nil {
				return fmt.Errorf("matching with regexp '%s': %w", v, err)
			}
			if !matched {
				delete(results, appKey)
			}
		}
	}
	if len(dscData.MatchOptions) > 0 {
		// for each app, check if all options is matched
		for appKey, app := range results {
			matchOk := true
			// application defined with term or not but not with both
			if len(app.Term) > 0 {
				// check if a term match a options block
				matchOptsNum := 0
			each_opts:
				for _, matchOption := range dscData.MatchOptions {
				each_term:
					for _, term := range app.Term {
						if v := matchOption.Alg.ValueString(); v != "" &&
							v != term.Alg.ValueString() {
							continue each_term
						}
						if v := matchOption.DestinationPort.ValueString(); v != "" &&
							v != term.DestinationPort.ValueString() {
							continue each_term
						}
						if v := matchOption.IcmpCode.ValueString(); v != "" &&
							v != term.IcmpCode.ValueString() {
							continue each_term
						}
						if v := matchOption.IcmpType.ValueString(); v != "" &&
							v != term.IcmpType.ValueString() {
							continue each_term
						}
						if v := matchOption.Icmp6Code.ValueString(); v != "" &&
							v != term.Icmp6Code.ValueString() {
							continue each_term
						}
						if v := matchOption.Icmp6Type.ValueString(); v != "" &&
							v != term.Icmp6Type.ValueString() {
							continue each_term
						}
						if v := matchOption.InactivityTimeout.ValueInt64(); v != 0 &&
							v != term.InactivityTimeout.ValueInt64() {
							continue each_term
						}
						if matchOption.InactivityTimeoutNever.ValueBool() &&
							!term.InactivityTimeoutNever.ValueBool() {
							continue each_term
						}
						if v := matchOption.Protocol.ValueString(); v != "" &&
							v != term.Protocol.ValueString() {
							continue each_term
						}
						if v := matchOption.RPCRrogramNumber.ValueString(); v != "" &&
							v != term.RPCRrogramNumber.ValueString() {
							continue each_term
						}
						if v := matchOption.SourcePort.ValueString(); v != "" &&
							v != term.SourcePort.ValueString() {
							continue each_term
						}
						if v := matchOption.UUID.ValueString(); v != "" &&
							v != term.UUID.ValueString() {
							continue each_term
						}

						// current term match current options block
						matchOptsNum++

						continue each_opts
					}
				}
				// all options block has not been validated
				if matchOptsNum < len(dscData.MatchOptions) {
					matchOk = false
				}
			} else {
				for _, matchOption := range dscData.MatchOptions {
					if v := matchOption.ApplicationProtocol.ValueString(); v != "" &&
						v != app.ApplicationProtocol.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.DestinationPort.ValueString(); v != "" &&
						v != app.DestinationPort.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.EtherType.ValueString(); v != "" &&
						v != app.EtherType.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.InactivityTimeout.ValueInt64(); v != 0 &&
						v != app.InactivityTimeout.ValueInt64() {
						matchOk = false

						break
					}
					if matchOption.InactivityTimeoutNever.ValueBool() &&
						!app.InactivityTimeoutNever.ValueBool() {
						matchOk = false

						break
					}
					if v := matchOption.Protocol.ValueString(); v != "" &&
						v != app.Protocol.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.RPCRrogramNumber.ValueString(); v != "" &&
						v != app.RPCRrogramNumber.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.SourcePort.ValueString(); v != "" &&
						v != app.SourcePort.ValueString() {
						matchOk = false

						break
					}
					if v := matchOption.UUID.ValueString(); v != "" &&
						v != app.UUID.ValueString() {
						matchOk = false

						break
					}
				}
			}
			if !matchOk {
				delete(results, appKey)
			}
		}
	}

	return nil
}

func (dscData *applicationsDataSourceData) fillID() {
	dscData.ID = types.StringValue(
		"match_name=" + dscData.MatchName.ValueString() +
			junos.IDSeparator +
			"match_options_n=" + strconv.Itoa(len(dscData.MatchOptions)),
	)
}

func (dscData *applicationsDataSourceData) fillData(
	results map[string]applicationsDataSourceBlockApplications,
) {
	for _, v := range results {
		dscData.Applications = append(dscData.Applications, v)
	}
}
