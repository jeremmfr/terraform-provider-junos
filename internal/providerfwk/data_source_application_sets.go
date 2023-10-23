package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	_ datasource.DataSource              = &applicationSetsDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationSetsDataSource{}
)

type applicationSetsDataSource struct {
	client *junos.Client
}

func (dsc *applicationSetsDataSource) typeName() string {
	return providerName + "_application_sets"
}

func (dsc *applicationSetsDataSource) junosName() string {
	return "applications application-set"
}

func newApplicationSetsDataSource() datasource.DataSource {
	return &applicationSetsDataSource{}
}

func (dsc *applicationSetsDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *applicationSetsDataSource) Configure(
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

func (dsc *applicationSetsDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get configurations from " + dsc.junosName() + " and junos-defaults applications application-set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source.",
			},
			"match_name": schema.StringAttribute{
				Optional:    true,
				Description: "A regexp to apply a filter on application-sets name.",
				Validators: []validator.String{
					tfvalidator.StringRegex(),
				},
			},
			"match_applications": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of applications to apply a filter on application-sets.",
			},
			"match_application_sets": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of application-sets to apply a filter on application-sets.",
			},
			"application_sets": schema.ListAttribute{
				Computed:    true,
				Description: "For each application-set found.",
				ElementType: types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
					"name":            types.StringType,
					"applications":    types.ListType{}.WithElementType(types.StringType),
					"application_set": types.ListType{}.WithElementType(types.StringType),
					"description":     types.StringType,
				}),
			},
		},
	}
}

type applicationSetsDataSourceData struct {
	ID                   types.String                                    `tfsdk:"id"`
	MatchName            types.String                                    `tfsdk:"match_name"`
	MatchApplications    []types.String                                  `tfsdk:"match_applications"`
	MatchApplicationSets []types.String                                  `tfsdk:"match_application_sets"`
	ApplicationSets      []applicationSetsDataSourceBlockApplicationSets `tfsdk:"application_sets"`
}

type applicationSetsDataSourceBlockApplicationSets struct {
	Name           types.String   `tfsdk:"name"`
	Applications   []types.String `tfsdk:"applications"`
	ApplicationSet []types.String `tfsdk:"application_set"`
	Description    types.String   `tfsdk:"description"`
}

func (dsc *applicationSetsDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data applicationSetsDataSourceData
	var matchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match_name"), &matchName)...)
	var matchApplications, matchApplicationSets []types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match_applications"), &matchApplications)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match_application_sets"), &matchApplicationSets)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MatchName = matchName
	data.MatchApplications = matchApplications
	data.MatchApplicationSets = matchApplicationSets

	junSess, err := dsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	applicationSetMap, err := dsc.search(junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}

	if err := data.filter(applicationSetMap); err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}
	data.fillData(applicationSetMap)
	data.fillID()

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (dsc *applicationSetsDataSource) search(
	junSess *junos.Session,
) (
	map[string]applicationSetsDataSourceBlockApplicationSets, error,
) {
	results := make(map[string]applicationSetsDataSourceBlockApplicationSets, 0)
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
			if !strings.HasPrefix(item, "set application-set ") {
				continue
			}
			itemTrim := strings.TrimPrefix(item, "set application-set ")
			itemTrimFields := strings.Split(itemTrim, " ")
			if _, ok := results[itemTrimFields[0]]; !ok {
				results[itemTrimFields[0]] = applicationSetsDataSourceBlockApplicationSets{
					Name: types.StringValue(itemTrimFields[0]),
				}
			}
			appSet := results[itemTrimFields[0]]
			balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
			switch {
			case balt.CutPrefixInString(&itemTrim, "application "):
				appSet.Applications = append(appSet.Applications, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "application-set "):
				appSet.ApplicationSet = append(appSet.ApplicationSet, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "description "):
				appSet.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			}
			results[itemTrimFields[0]] = appSet
		}
	}

	return results, nil
}

func (dscData *applicationSetsDataSourceData) filter(
	results map[string]applicationSetsDataSourceBlockApplicationSets,
) error {
	if v := dscData.MatchName.ValueString(); v != "" {
		for appSetKey, appSet := range results {
			matched, err := regexp.MatchString(v, appSet.Name.ValueString())
			if err != nil {
				return fmt.Errorf("matching with regexp '%s': %w", v, err)
			}
			if !matched {
				delete(results, appSetKey)
			}
		}
	}
	if len(dscData.MatchApplications) > 0 {
		// for each app-set, check if all applications is matched
		for appSetKey, appSet := range results {
			if len(appSet.Applications) != len(dscData.MatchApplications) {
				delete(results, appSetKey)

				continue
			}
			matchAppsOk := 0
		each_match:
			for _, v := range dscData.MatchApplications {
				for _, app := range appSet.Applications {
					if v.ValueString() == app.ValueString() {
						matchAppsOk++

						continue each_match
					}
				}
			}
			if len(dscData.MatchApplications) != matchAppsOk {
				delete(results, appSetKey)
			}
		}
	}
	if len(dscData.MatchApplicationSets) > 0 {
		// for each app-set, check if all application-sets is matched
		for appSetKey, appSet := range results {
			if len(appSet.ApplicationSet) != len(dscData.MatchApplicationSets) {
				delete(results, appSetKey)

				continue
			}
			matchAppSetsOk := 0
		each_match_sets:
			for _, v := range dscData.MatchApplicationSets {
				for _, appSet := range appSet.ApplicationSet {
					if v.ValueString() == appSet.ValueString() {
						matchAppSetsOk++

						continue each_match_sets
					}
				}
			}
			if len(dscData.MatchApplicationSets) != matchAppSetsOk {
				delete(results, appSetKey)
			}
		}
	}

	return nil
}

func (dscData *applicationSetsDataSourceData) fillID() {
	dscData.ID = types.StringValue(
		"match_name=" + dscData.MatchName.ValueString() +
			junos.IDSeparator +
			"match_applications_n=" + strconv.Itoa(len(dscData.MatchApplications)),
	)
}

func (dscData *applicationSetsDataSourceData) fillData(
	results map[string]applicationSetsDataSourceBlockApplicationSets,
) {
	for _, v := range results {
		dscData.ApplicationSets = append(dscData.ApplicationSets, v)
	}
}
