package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &services{}
	_ resource.ResourceWithConfigure      = &services{}
	_ resource.ResourceWithValidateConfig = &services{}
	_ resource.ResourceWithImportState    = &services{}
	_ resource.ResourceWithUpgradeState   = &services{}
)

type services struct {
	client *junos.Client
}

func newServicesResource() resource.Resource {
	return &services{}
}

func (rsc *services) typeName() string {
	return providerName + "_services"
}

func (rsc *services) junosName() string {
	return "services"
}

func (rsc *services) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *services) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *services) Configure(
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

func (rsc *services) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `services`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean supported lines when destroy this resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"advanced_anti_malware": schema.SingleNestedBlock{
				Description: "Declare `advanced-anti-malware` static configuration.",
				Blocks: map[string]schema.Block{
					"connection": schema.SingleNestedBlock{
						Description: "Declare `connection` configuration.",
						Attributes: map[string]schema.Attribute{
							"auth_tls_profile": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Authentication TLS profile.",
								PlanModifiers: []planmodifier.String{
									tfplanmodifier.StringUseStateNullForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"proxy_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Proxy profile.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"source_address": schema.StringAttribute{
								Optional:    true,
								Description: "The source ip for connecting to the cloud server.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress(),
								},
							},
							"source_interface": schema.StringAttribute{
								Optional:    true,
								Description: "The source interface for connecting to the cloud server.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
									tfvalidator.String1DotCount(),
								},
							},
							"url": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "The url of the cloud server.",
								PlanModifiers: []planmodifier.String{
									tfplanmodifier.StringUseStateNullForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"default_policy": schema.SingleNestedBlock{
						Description: "Declare `default-policy` configuration.",
						Attributes: map[string]schema.Attribute{
							"blacklist_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware blacklist hit.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"default_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware action.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"fallback_options_action": schema.StringAttribute{
								Optional:    true,
								Description: "Notification action taken for fallback action.",
								Validators: []validator.String{
									stringvalidator.OneOf("block", "permit"),
								},
							},
							"fallback_options_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware fallback action.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"http_action": schema.StringAttribute{
								Optional:    true,
								Description: "Action taken for contents with verdict meet threshold for HTTP.",
								Validators: []validator.String{
									stringvalidator.OneOf("block", "permit"),
								},
							},
							"http_client_notify_file": schema.StringAttribute{
								Optional: true,
								Description: "File name for http response to client notification action taken" +
									" for contents with verdict meet threshold.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 255),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"http_client_notify_message": schema.StringAttribute{
								Optional:    true,
								Description: "Block message to client notification action taken for contents with verdict meet threshold.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"http_client_notify_redirect_url": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect url to client notification action taken for contents with verdict meet threshold.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1023),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"http_file_verdict_unknown": schema.StringAttribute{
								Optional:    true,
								Description: "Action taken for contents with verdict unknown.",
								Validators: []validator.String{
									stringvalidator.OneOf("block", "permit"),
								},
							},
							"http_inspection_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Advanced Anti-malware inspection-profile name for HTTP.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"http_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware actions for HTTP.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"imap_inspection_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Advanced Anti-malware inspection-profile name for IMAP.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"imap_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware actions for IMAP.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"smtp_inspection_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Advanced Anti-malware inspection-profile name for SMTP.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"smtp_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware actions for SMTP.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"verdict_threshold": schema.StringAttribute{
								Optional:    true,
								Description: "Verdict threshold.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"recommended", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
									),
								},
							},
							"whitelist_notification_log": schema.BoolAttribute{
								Optional:    true,
								Description: "Logging option for Advanced Anti-malware whitelist hit.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"application_identification": schema.SingleNestedBlock{
				Description: "Enable `application-identification`.",
				Attributes: map[string]schema.Attribute{
					"application_system_cache_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Application system cache entry lifetime.",
						Validators: []validator.Int64{
							int64validator.Between(0, 1000000),
						},
					},
					"global_offload_byte_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Global byte limit to offload AppID inspection.",
						Validators: []validator.Int64{
							int64validator.Between(0, 4294967295),
						},
					},
					"imap_cache_size": schema.Int64Attribute{
						Optional:    true,
						Description: "IMAP cache size, it will be effective only after next appid sigpack install.",
						Validators: []validator.Int64{
							int64validator.Between(60, 512000),
						},
					},
					"imap_cache_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "IMAP cache entry timeout in seconds.",
						Validators: []validator.Int64{
							int64validator.Between(1, 86400),
						},
					},
					"max_memory": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum amount of object cache memory JDPI can use (in MB).",
						Validators: []validator.Int64{
							int64validator.Between(1, 200000),
						},
					},
					"max_transactions": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of transaction finals to terminate application classification.",
						Validators: []validator.Int64{
							int64validator.Between(0, 25),
						},
					},
					"micro_apps": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable Micro Apps identifcation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_application_system_cache": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable storing AI result in application system cache.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"statistics_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Configure application statistics information with collection interval (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(1, 1440),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"application_system_cache": schema.SingleNestedBlock{
						Description: "Enable application system cache.",
						Attributes: map[string]schema.Attribute{
							"no_miscellaneous_services": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable ASC for miscellaneous services APBR,...",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"security_services": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable ASC for security services (appfw, appqos, idp, skyatp..).",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"download": schema.SingleNestedBlock{
						Description: "Declare `download` configuration.",
						Attributes: map[string]schema.Attribute{
							"automatic_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Attempt to download new application package (hours).",
								Validators: []validator.Int64{
									int64validator.Between(6, 720),
								},
							},
							"automatic_start_time": schema.StringAttribute{
								Optional:    true,
								Description: "Start time to scheduled download and update.",
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile(
										`^([0-9]{4}-)?(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]).(2[0-3]|[01][0-9]):[0-5][0-9](:[0-5][0-9])?$`),
										"must be in the format MM-DD.hh:mm / YYYY-MM-DD.hh:mm:ss",
									),
								},
							},
							"ignore_server_validation": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable server authentication for Application Signature download.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"proxy_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Configure web proxy for Application signature download.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 128),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"url": schema.StringAttribute{
								Optional:    true,
								Description: "URL for application package download.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									stringvalidator.RegexMatches(regexp.MustCompile(
										`^(?i)(https?|file):`),
										"URL starts with http, https or file",
									),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"enable_performance_mode": schema.SingleNestedBlock{
						Description: "Enable performance mode knobs for best DPI performance.",
						Attributes: map[string]schema.Attribute{
							"max_packet_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Set the maximum packet threshold for DPI performance mode.",
								Validators: []validator.Int64{
									int64validator.Between(1, 100),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"inspection_limit_tcp": schema.SingleNestedBlock{
						Description: "Enable TCP byte/packet inspection limit.",
						Attributes: map[string]schema.Attribute{
							"byte_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "TCP byte inspection limit.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"packet_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "TCP packet inspection limit.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"inspection_limit_udp": schema.SingleNestedBlock{
						Description: "Enable UDP byte/packet inspection limit.",
						Attributes: map[string]schema.Attribute{
							"byte_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "UDP byte inspection limit.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"packet_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "UDP packet inspection limit.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"security_intelligence": schema.SingleNestedBlock{
				Description: "Declare `security-intelligence` configuration.",
				Attributes: map[string]schema.Attribute{
					"authentication_tls_profile": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "TLS profile for authentication to use feed update services.",
						PlanModifiers: []planmodifier.String{
							tfplanmodifier.StringUseStateNullForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"authentication_token": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Token string for authentication to use feed update services.",
						PlanModifiers: []planmodifier.String{
							tfplanmodifier.StringUseStateNullForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^[a-zA-Z0-9]{32}$`),
								"must be consisted of 32 alphanumeric characters",
							),
						},
					},
					"category_disable": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Categories to be disabled.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							),
						},
					},
					"proxy_profile": schema.StringAttribute{
						Optional:    true,
						Description: "The proxy profile name.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"url": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Configure the url of feed server.",
						PlanModifiers: []planmodifier.String{
							tfplanmodifier.StringUseStateNullForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"url_parameter": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Configure the parameter of url.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"default_policy": schema.ListNestedBlock{
						Description: "For each name of category, configure default-policy for a category.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"category_name": schema.StringAttribute{
									Required:    true,
									Description: "Name of security intelligence category.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"profile_name": schema.StringAttribute{
									Required:    true,
									Description: "Name of profile.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"user_identification": schema.SingleNestedBlock{
				Description: "Declare `user-identification` configuration.",
				Attributes: map[string]schema.Attribute{
					"device_info_auth_source": schema.StringAttribute{
						Optional:    true,
						Description: "Configure authentication-source on device information configuration.",
						Validators: []validator.String{
							stringvalidator.OneOf("active-directory", "network-access-controller"),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"ad_access": schema.SingleNestedBlock{
						Description: "Enable `active-directory-access`.",
						Attributes: map[string]schema.Attribute{
							"auth_entry_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Authentication entry timeout number (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(0, 1440),
								},
							},
							"filter_exclude": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Exclude addresses.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.NoNullValues(),
									setvalidator.ValueStringsAre(
										tfvalidator.StringCIDR(),
									),
								},
							},
							"filter_include": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Include addresses.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.NoNullValues(),
									setvalidator.ValueStringsAre(
										tfvalidator.StringCIDR(),
									),
								},
							},
							"firewall_auth_forced_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Firewall auth fallback authentication entry forced timeout number (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(10, 1440),
								},
							},
							"invalid_auth_entry_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Invalid authentication entry timeout number (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(0, 1440),
								},
							},
							"no_on_demand_probe": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable on-demand probe.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"wmi_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Wmi timeout number (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(3, 120),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"identity_management": schema.SingleNestedBlock{
						Description: "Declare `identity-management` configuration.",
						Attributes: map[string]schema.Attribute{
							"authentication_entry_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Authentication entry timeout number (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(0, 1440),
								},
							},
							"batch_query_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Query interval for batch query (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 60),
								},
							},
							"batch_query_items_per_batch": schema.Int64Attribute{
								Optional:    true,
								Description: "Items number per batch query.",
								Validators: []validator.Int64{
									int64validator.Between(100, 1000),
								},
							},
							"filter_domain": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Domain filter.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.NoNullValues(),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthBetween(1, 64),
										tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
									),
								},
							},
							"filter_exclude_ip_address_book": schema.StringAttribute{
								Optional:    true,
								Description: "Referenced address book to exclude IP filter.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"filter_exclude_ip_address_set": schema.StringAttribute{
								Optional:    true,
								Description: "Referenced address set to exclude IP filter.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"filter_include_ip_address_book": schema.StringAttribute{
								Optional:    true,
								Description: "Referenced address book to include IP filter.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"filter_include_ip_address_set": schema.StringAttribute{
								Optional:    true,
								Description: "Referenced address set to include IP filter.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"invalid_authentication_entry_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Invalid authentication entry timeout number (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(0, 1440),
								},
							},
							"ip_query_disable": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable IP query.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"ip_query_delay_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Delay time to send IP query (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 60),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"connection": schema.SingleNestedBlock{
								Description: "Declare connection configuration.",
								Attributes: map[string]schema.Attribute{
									"primary_address": schema.StringAttribute{
										Required:    false, // true when SingleNestedBlock is specified
										Optional:    true,
										Description: "IP address of Primary server.",
										Validators: []validator.String{
											tfvalidator.StringIPAddress(),
										},
									},
									"primary_client_id": schema.StringAttribute{
										Required:    false, // true when SingleNestedBlock is specified
										Optional:    true,
										Description: "Client ID of Primary server for OAuth2 grant.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"primary_client_secret": schema.StringAttribute{
										Required:    false, // true when SingleNestedBlock is specified
										Optional:    true,
										Sensitive:   true,
										Description: "Client secret of Primary server for OAuth2 grant.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 128),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"connect_method": schema.StringAttribute{
										Optional:    true,
										Description: "Method of connection.",
										Validators: []validator.String{
											stringvalidator.OneOf("http", "https"),
										},
									},
									"port": schema.Int64Attribute{
										Optional:    true,
										Description: "Server port.",
										Validators: []validator.Int64{
											int64validator.Between(1, 65535),
										},
									},
									"primary_ca_certificate": schema.StringAttribute{
										Optional:    true,
										Description: "Ca-certificate file name of Primary server.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 256),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"query_api": schema.StringAttribute{
										Optional:    true,
										Description: "Query API.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(4, 128),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"secondary_address": schema.StringAttribute{
										Optional:    true,
										Description: "IP address of Secondary server.",
										Validators: []validator.String{
											tfvalidator.StringIPAddress(),
										},
									},
									"secondary_ca_certificate": schema.StringAttribute{
										Optional:    true,
										Description: "Ca-certificate file name of Secondary server.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 256),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"secondary_client_id": schema.StringAttribute{
										Optional:    true,
										Description: "Client ID of Secondary server for OAuth2 grant.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"secondary_client_secret": schema.StringAttribute{
										Optional:    true,
										Sensitive:   true,
										Description: "Client secret of Secondary server for OAuth2 grant.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 128),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"token_api": schema.StringAttribute{
										Optional:    true,
										Description: "API of acquiring token for OAuth2 authentication.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 128),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type servicesData struct {
	ID                        types.String                            `tfsdk:"id"`
	CleanOnDestroy            types.Bool                              `tfsdk:"clean_on_destroy"`
	AdvancedAntiMalware       *servicesBlockAdvancedAntiMalware       `tfsdk:"advanced_anti_malware"`
	ApplicationIdentification *servicesBlockApplicationIdentification `tfsdk:"application_identification"`
	SecurityIntelligence      *servicesBlockSecurityIntelligence      `tfsdk:"security_intelligence"`
	UserIdentification        *servicesBlockUserIdentification        `tfsdk:"user_identification"`
}

type servicesConfig struct {
	ID                        types.String                             `tfsdk:"id"`
	CleanOnDestroy            types.Bool                               `tfsdk:"clean_on_destroy"`
	AdvancedAntiMalware       *servicesBlockAdvancedAntiMalware        `tfsdk:"advanced_anti_malware"`
	ApplicationIdentification *servicesBlockApplicationIdentification  `tfsdk:"application_identification"`
	SecurityIntelligence      *servicesBlockSecurityIntelligenceConfig `tfsdk:"security_intelligence"`
	UserIdentification        *servicesBlockUserIdentificationConfig   `tfsdk:"user_identification"`
}

type servicesBlockAdvancedAntiMalware struct {
	Connection    *servicesBlockAdvancedAntiMalwareBlockConnection    `tfsdk:"connection"`
	DefaultPolicy *servicesBlockAdvancedAntiMalwareBlockDefaultPolicy `tfsdk:"default_policy"`
}

func (block *servicesBlockAdvancedAntiMalware) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type servicesBlockAdvancedAntiMalwareBlockConnection struct {
	AuthTLSProfile  types.String `tfsdk:"auth_tls_profile"`
	ProxyProfile    types.String `tfsdk:"proxy_profile"`
	SourceAddress   types.String `tfsdk:"source_address"`
	SourceInterface types.String `tfsdk:"source_interface"`
	URL             types.String `tfsdk:"url"`
}

type servicesBlockAdvancedAntiMalwareBlockDefaultPolicy struct {
	BlacklistNotificationLog       types.Bool   `tfsdk:"blacklist_notification_log"`
	DefaultNotificationLog         types.Bool   `tfsdk:"default_notification_log"`
	FallbackOptionsAction          types.String `tfsdk:"fallback_options_action"`
	FallbackOptionsNotificationLog types.Bool   `tfsdk:"fallback_options_notification_log"`
	HTTPAction                     types.String `tfsdk:"http_action"`
	HTTPClientNotifyFile           types.String `tfsdk:"http_client_notify_file"`
	HTTPClientNotifyMessage        types.String `tfsdk:"http_client_notify_message"`
	HTTPClientNotifyRedirectURL    types.String `tfsdk:"http_client_notify_redirect_url"`
	HTTPFileVerdictUnknown         types.String `tfsdk:"http_file_verdict_unknown"`
	HTTPInspectionProfile          types.String `tfsdk:"http_inspection_profile"`
	HTTPNotificationLog            types.Bool   `tfsdk:"http_notification_log"`
	IMAPInspectionProfile          types.String `tfsdk:"imap_inspection_profile"`
	IMAPNotificationLog            types.Bool   `tfsdk:"imap_notification_log"`
	SMTPInspectionProfile          types.String `tfsdk:"smtp_inspection_profile"`
	SMTPNotificationLog            types.Bool   `tfsdk:"smtp_notification_log"`
	VerdictThreshold               types.String `tfsdk:"verdict_threshold"`
	WhitelistNotificationLog       types.Bool   `tfsdk:"whitelist_notification_log"`
}

//nolint:lll
type servicesBlockApplicationIdentification struct {
	ApplicationSystemCacheTimeout types.Int64                                                        `tfsdk:"application_system_cache_timeout"`
	GlobalOffloadByteLimit        types.Int64                                                        `tfsdk:"global_offload_byte_limit"`
	IMAPCacheSize                 types.Int64                                                        `tfsdk:"imap_cache_size"`
	IMAPCacheTimeout              types.Int64                                                        `tfsdk:"imap_cache_timeout"`
	MaxMemory                     types.Int64                                                        `tfsdk:"max_memory"`
	MaxTransactions               types.Int64                                                        `tfsdk:"max_transactions"`
	MicroApps                     types.Bool                                                         `tfsdk:"micro_apps"`
	NoApplicationSystemCache      types.Bool                                                         `tfsdk:"no_application_system_cache"`
	StatisticsInterval            types.Int64                                                        `tfsdk:"statistics_interval"`
	ApplicationSystemCache        *servicesBlockApplicationIdentificationBlockApplicationSystemCache `tfsdk:"application_system_cache"`
	Download                      *servicesBlockApplicationIdentificationBlockDownload               `tfsdk:"download"`
	EnablePerformanceMode         *servicesBlockApplicationIdentificationBlockEnablePerformanceMode  `tfsdk:"enable_performance_mode"`
	InspectionLimitTCP            *servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP  `tfsdk:"inspection_limit_tcp"`
	InspectionLimitUDP            *servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP  `tfsdk:"inspection_limit_udp"`
}

type servicesBlockApplicationIdentificationBlockApplicationSystemCache struct {
	NoMiscellaneousServices types.Bool `tfsdk:"no_miscellaneous_services"`
	SecurityServices        types.Bool `tfsdk:"security_services"`
}

func (block *servicesBlockApplicationIdentificationBlockApplicationSystemCache) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type servicesBlockApplicationIdentificationBlockDownload struct {
	AutomaticInterval      types.Int64  `tfsdk:"automatic_interval"`
	AutomaticStartTime     types.String `tfsdk:"automatic_start_time"`
	IgnoreServerValidation types.Bool   `tfsdk:"ignore_server_validation"`
	ProxyProfile           types.String `tfsdk:"proxy_profile"`
	URL                    types.String `tfsdk:"url"`
}

func (block *servicesBlockApplicationIdentificationBlockDownload) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type servicesBlockApplicationIdentificationBlockEnablePerformanceMode struct {
	MaxPacketThreshold types.Int64 `tfsdk:"max_packet_threshold"`
}

type servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP struct {
	ByteLimit   types.Int64 `tfsdk:"byte_limit"`
	PacketLimit types.Int64 `tfsdk:"packet_limit"`
}

type servicesBlockSecurityIntelligence struct {
	AuthenticationTLSProfile types.String                                          `tfsdk:"authentication_tls_profile"`
	AuthenticationToken      types.String                                          `tfsdk:"authentication_token"`
	CategoryDisable          []types.String                                        `tfsdk:"category_disable"`
	ProxyProfile             types.String                                          `tfsdk:"proxy_profile"`
	URL                      types.String                                          `tfsdk:"url"`
	URLParameter             types.String                                          `tfsdk:"url_parameter"`
	DefaultPolicy            []servicesBlockSecurityIntelligenceBlockDefaultPolicy `tfsdk:"default_policy"`
}

func (block *servicesBlockSecurityIntelligence) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type servicesBlockSecurityIntelligenceConfig struct {
	AuthenticationTLSProfile types.String `tfsdk:"authentication_tls_profile"`
	AuthenticationToken      types.String `tfsdk:"authentication_token"`
	CategoryDisable          types.Set    `tfsdk:"category_disable"`
	ProxyProfile             types.String `tfsdk:"proxy_profile"`
	URL                      types.String `tfsdk:"url"`
	URLParameter             types.String `tfsdk:"url_parameter"`
	DefaultPolicy            types.List   `tfsdk:"default_policy"`
}

type servicesBlockSecurityIntelligenceBlockDefaultPolicy struct {
	CategoryName types.String `tfsdk:"category_name" tfdata:"identifier"`
	ProfileName  types.String `tfsdk:"profile_name"`
}

type servicesBlockUserIdentification struct {
	DeviceInfoAuthSource types.String                                            `tfsdk:"device_info_auth_source"`
	ADAccess             *servicesBlockUserIdentificationBlockADAccess           `tfsdk:"ad_access"`
	IdentityManagement   *servicesBlockUserIdentificationBlockIdentityManagement `tfsdk:"identity_management"`
}

func (block *servicesBlockUserIdentification) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type servicesBlockUserIdentificationConfig struct {
	DeviceInfoAuthSource types.String                                                  `tfsdk:"device_info_auth_source"`
	ADAccess             *servicesBlockUserIdentificationBlockADAccessConfig           `tfsdk:"ad_access"`
	IdentityManagement   *servicesBlockUserIdentificationBlockIdentityManagementConfig `tfsdk:"identity_management"`
}

func (block *servicesBlockUserIdentificationConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type servicesBlockUserIdentificationBlockADAccess struct {
	AuthEntryTimeout          types.Int64    `tfsdk:"auth_entry_timeout"`
	FilterExclude             []types.String `tfsdk:"filter_exclude"`
	FilterInclude             []types.String `tfsdk:"filter_include"`
	FirewallAuthForcedTimeout types.Int64    `tfsdk:"firewall_auth_forced_timeout"`
	InvalidAuthEntryTimeout   types.Int64    `tfsdk:"invalid_auth_entry_timeout"`
	NoOnDemandProbe           types.Bool     `tfsdk:"no_on_demand_probe"`
	WmiTimeout                types.Int64    `tfsdk:"wmi_timeout"`
}

type servicesBlockUserIdentificationBlockADAccessConfig struct {
	AuthEntryTimeout          types.Int64 `tfsdk:"auth_entry_timeout"`
	FilterExclude             types.Set   `tfsdk:"filter_exclude"`
	FilterInclude             types.Set   `tfsdk:"filter_include"`
	FirewallAuthForcedTimeout types.Int64 `tfsdk:"firewall_auth_forced_timeout"`
	InvalidAuthEntryTimeout   types.Int64 `tfsdk:"invalid_auth_entry_timeout"`
	NoOnDemandProbe           types.Bool  `tfsdk:"no_on_demand_probe"`
	WmiTimeout                types.Int64 `tfsdk:"wmi_timeout"`
}

func (block *servicesBlockUserIdentificationBlockADAccessConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

//nolint:lll
type servicesBlockUserIdentificationBlockIdentityManagement struct {
	AuthenticationEntryTimeout        types.Int64                                                            `tfsdk:"authentication_entry_timeout"`
	BatchQueryInterval                types.Int64                                                            `tfsdk:"batch_query_interval"`
	BatchQueryItemsPerBatch           types.Int64                                                            `tfsdk:"batch_query_items_per_batch"`
	FilterDomain                      []types.String                                                         `tfsdk:"filter_domain"`
	FilterExcludeIPAddressBook        types.String                                                           `tfsdk:"filter_exclude_ip_address_book"`
	FilterExcludeIPAddressSet         types.String                                                           `tfsdk:"filter_exclude_ip_address_set"`
	FilterIncludeIPAddressBook        types.String                                                           `tfsdk:"filter_include_ip_address_book"`
	FilterIncludeIPAddressSet         types.String                                                           `tfsdk:"filter_include_ip_address_set"`
	InvalidAuthenticationEntryTimeout types.Int64                                                            `tfsdk:"invalid_authentication_entry_timeout"`
	IPQueryDisable                    types.Bool                                                             `tfsdk:"ip_query_disable"`
	IPQueryDelayTime                  types.Int64                                                            `tfsdk:"ip_query_delay_time"`
	Connection                        *servicesBlockUserIdentificationBlockIdentityManagementBlockConnection `tfsdk:"connection"`
}

//nolint:lll
type servicesBlockUserIdentificationBlockIdentityManagementConfig struct {
	AuthenticationEntryTimeout        types.Int64                                                            `tfsdk:"authentication_entry_timeout"`
	BatchQueryInterval                types.Int64                                                            `tfsdk:"batch_query_interval"`
	BatchQueryItemsPerBatch           types.Int64                                                            `tfsdk:"batch_query_items_per_batch"`
	FilterDomain                      types.Set                                                              `tfsdk:"filter_domain"`
	FilterExcludeIPAddressBook        types.String                                                           `tfsdk:"filter_exclude_ip_address_book"`
	FilterExcludeIPAddressSet         types.String                                                           `tfsdk:"filter_exclude_ip_address_set"`
	FilterIncludeIPAddressBook        types.String                                                           `tfsdk:"filter_include_ip_address_book"`
	FilterIncludeIPAddressSet         types.String                                                           `tfsdk:"filter_include_ip_address_set"`
	InvalidAuthenticationEntryTimeout types.Int64                                                            `tfsdk:"invalid_authentication_entry_timeout"`
	IPQueryDisable                    types.Bool                                                             `tfsdk:"ip_query_disable"`
	IPQueryDelayTime                  types.Int64                                                            `tfsdk:"ip_query_delay_time"`
	Connection                        *servicesBlockUserIdentificationBlockIdentityManagementBlockConnection `tfsdk:"connection"`
}

func (block *servicesBlockUserIdentificationBlockIdentityManagementConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type servicesBlockUserIdentificationBlockIdentityManagementBlockConnection struct {
	PrimaryAddress         types.String `tfsdk:"primary_address"`
	PrimaryClientID        types.String `tfsdk:"primary_client_id"`
	PrimaryClientSecret    types.String `tfsdk:"primary_client_secret"`
	ConnectMethod          types.String `tfsdk:"connect_method"`
	Port                   types.Int64  `tfsdk:"port"`
	PrimaryCACertificate   types.String `tfsdk:"primary_ca_certificate"`
	QueryAPI               types.String `tfsdk:"query_api"`
	SecondaryAddress       types.String `tfsdk:"secondary_address"`
	SecondaryCACertificate types.String `tfsdk:"secondary_ca_certificate"`
	SecondaryClientID      types.String `tfsdk:"secondary_client_id"`
	SecondaryClientSecret  types.String `tfsdk:"secondary_client_secret"`
	TokenAPI               types.String `tfsdk:"token_api"`
}

func (rsc *services) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AdvancedAntiMalware != nil {
		if config.AdvancedAntiMalware.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("advanced_anti_malware").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"advanced_anti_malware block is empty",
			)
		}
		if config.AdvancedAntiMalware.Connection != nil {
			if !config.AdvancedAntiMalware.Connection.SourceAddress.IsNull() &&
				!config.AdvancedAntiMalware.Connection.SourceAddress.IsUnknown() &&
				!config.AdvancedAntiMalware.Connection.SourceInterface.IsNull() &&
				!config.AdvancedAntiMalware.Connection.SourceInterface.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("connection").AtName("source_address"),
					tfdiag.ConflictConfigErrSummary,
					"source_address and source_interface cannot be configured together"+
						" in connection block in advanced_anti_malware block",
				)
			}
		}
		if config.AdvancedAntiMalware.DefaultPolicy != nil {
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsUnknown() &&
				config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_action"),
					tfdiag.MissingConfigErrSummary,
					"http_inspection_profile must be specified with http_action"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyFile.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyFile.IsUnknown() {
				if config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() ||
					config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_file"),
						tfdiag.MissingConfigErrSummary,
						"http_action and http_inspection_profile must be specified with http_client_notify_file"+
							" in default_policy block in advanced_anti_malware block",
					)
				}
				if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyMessage.IsNull() &&
					!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyMessage.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_file"),
						tfdiag.ConflictConfigErrSummary,
						"http_client_notify_file and http_client_notify_message cannot be configured together"+
							" in default_policy block in advanced_anti_malware block",
					)
				}
				if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsNull() &&
					!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_file"),
						tfdiag.ConflictConfigErrSummary,
						"http_client_notify_file and http_client_notify_redirect_url cannot be configured together"+
							" in default_policy block in advanced_anti_malware block",
					)
				}
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyMessage.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyMessage.IsUnknown() {
				if config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() ||
					config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_message"),
						tfdiag.MissingConfigErrSummary,
						"http_action and http_inspection_profile must be specified with http_client_notify_message"+
							" in default_policy block in advanced_anti_malware block",
					)
				}
				if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsNull() &&
					!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_message"),
						tfdiag.ConflictConfigErrSummary,
						"http_client_notify_message and http_client_notify_redirect_url cannot be configured together"+
							" in default_policy block in advanced_anti_malware block",
					)
				}
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPClientNotifyRedirectURL.IsUnknown() &&
				(config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() ||
					config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_client_notify_redirect_url"),
					tfdiag.MissingConfigErrSummary,
					"http_action and http_inspection_profile must be specified with http_client_notify_redirect_url"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPFileVerdictUnknown.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPFileVerdictUnknown.IsUnknown() &&
				(config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() ||
					config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_file_verdict_unknown"),
					tfdiag.MissingConfigErrSummary,
					"http_action and http_inspection_profile must be specified with http_file_verdict_unknown"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsUnknown() &&
				config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_inspection_profile"),
					tfdiag.MissingConfigErrSummary,
					"http_action must be specified with http_inspection_profile"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.HTTPNotificationLog.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.HTTPNotificationLog.IsUnknown() &&
				(config.AdvancedAntiMalware.DefaultPolicy.HTTPAction.IsNull() ||
					config.AdvancedAntiMalware.DefaultPolicy.HTTPInspectionProfile.IsNull()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("http_notification_log"),
					tfdiag.MissingConfigErrSummary,
					"http_action and http_inspection_profile must be specified with http_notification_log"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.IMAPNotificationLog.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.IMAPNotificationLog.IsUnknown() &&
				config.AdvancedAntiMalware.DefaultPolicy.IMAPInspectionProfile.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("imap_notification_log"),
					tfdiag.MissingConfigErrSummary,
					"imap_inspection_profile must be specified with imap_notification_log"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
			if !config.AdvancedAntiMalware.DefaultPolicy.SMTPNotificationLog.IsNull() &&
				!config.AdvancedAntiMalware.DefaultPolicy.SMTPNotificationLog.IsUnknown() &&
				config.AdvancedAntiMalware.DefaultPolicy.SMTPInspectionProfile.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("advanced_anti_malware").AtName("default_policy").AtName("smtp_notification_log"),
					tfdiag.MissingConfigErrSummary,
					"smtp_inspection_profile must be specified with smtp_notification_log"+
						" in default_policy block in advanced_anti_malware block",
				)
			}
		}
	}
	if config.ApplicationIdentification != nil {
		if config.ApplicationIdentification.Download != nil &&
			config.ApplicationIdentification.Download.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("application_identification").AtName("download").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"download block is empty"+
					" in application_identification block",
			)
		}
		if !config.ApplicationIdentification.NoApplicationSystemCache.IsNull() &&
			!config.ApplicationIdentification.NoApplicationSystemCache.IsUnknown() &&
			config.ApplicationIdentification.ApplicationSystemCache != nil &&
			config.ApplicationIdentification.ApplicationSystemCache.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("application_identification").AtName("no_application_system_cache"),
				tfdiag.ConflictConfigErrSummary,
				"application_system_cache and no_application_system_cache cannot be configured together"+
					" in application_identification block",
			)
		}
	}
	if config.SecurityIntelligence != nil {
		if !config.SecurityIntelligence.DefaultPolicy.IsNull() &&
			!config.SecurityIntelligence.DefaultPolicy.IsUnknown() {
			var configDefaultPolicy []servicesBlockSecurityIntelligenceBlockDefaultPolicy
			asDiags := config.SecurityIntelligence.DefaultPolicy.ElementsAs(ctx, &configDefaultPolicy, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			defaultPolicyCategoryName := make(map[string]struct{})
			for i, block := range configDefaultPolicy {
				if block.CategoryName.IsUnknown() {
					continue
				}
				categoryName := block.CategoryName.ValueString()
				if _, ok := defaultPolicyCategoryName[categoryName]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("security_intelligence").AtName("default_policy").AtListIndex(i).AtName("category_name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple default_policy blocks with the same category_name %q"+
							" in security_intelligence block", categoryName),
					)
				}
				defaultPolicyCategoryName[categoryName] = struct{}{}
			}
		}
		if !config.SecurityIntelligence.AuthenticationTLSProfile.IsNull() &&
			!config.SecurityIntelligence.AuthenticationTLSProfile.IsUnknown() &&
			!config.SecurityIntelligence.AuthenticationToken.IsNull() &&
			!config.SecurityIntelligence.AuthenticationToken.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("security_intelligence").AtName("authentication_tls_profile"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_tls_profile and authentication_token cannot be configured together"+
					" in security_intelligence block",
			)
		}
	}
	if config.UserIdentification != nil {
		if config.UserIdentification.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("user_identification").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"user_identification block is empty",
			)
		}
		if config.UserIdentification.ADAccess != nil &&
			config.UserIdentification.ADAccess.hasKnownValue() &&
			config.UserIdentification.IdentityManagement != nil &&
			config.UserIdentification.IdentityManagement.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("user_identification").AtName("identity_management"),
				tfdiag.ConflictConfigErrSummary,
				"ad_access and identity_management cannot be configured together"+
					" in user_identification block",
			)
		}
		if config.UserIdentification.IdentityManagement != nil {
			if !config.UserIdentification.IdentityManagement.FilterExcludeIPAddressBook.IsNull() &&
				!config.UserIdentification.IdentityManagement.FilterExcludeIPAddressBook.IsUnknown() &&
				config.UserIdentification.IdentityManagement.FilterExcludeIPAddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("user_identification").AtName("identity_management").AtName("filter_exclude_ip_address_book"),
					tfdiag.MissingConfigErrSummary,
					"filter_exclude_ip_address_set must be specified with filter_exclude_ip_address_book"+
						" in identity_management block in user_identification block",
				)
			}
			if !config.UserIdentification.IdentityManagement.FilterExcludeIPAddressSet.IsNull() &&
				!config.UserIdentification.IdentityManagement.FilterExcludeIPAddressSet.IsUnknown() &&
				config.UserIdentification.IdentityManagement.FilterExcludeIPAddressBook.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("user_identification").AtName("identity_management").AtName("filter_exclude_ip_address_set"),
					tfdiag.MissingConfigErrSummary,
					"filter_exclude_ip_address_book must be specified with filter_exclude_ip_address_set"+
						" in identity_management block in user_identification block",
				)
			}
			if !config.UserIdentification.IdentityManagement.FilterIncludeIPAddressBook.IsNull() &&
				!config.UserIdentification.IdentityManagement.FilterIncludeIPAddressBook.IsUnknown() &&
				config.UserIdentification.IdentityManagement.FilterIncludeIPAddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("user_identification").AtName("identity_management").AtName("filter_include_ip_address_book"),
					tfdiag.MissingConfigErrSummary,
					"filter_include_ip_address_set must be specified with filter_include_ip_address_book"+
						" in identity_management block in user_identification block",
				)
			}
			if !config.UserIdentification.IdentityManagement.FilterIncludeIPAddressSet.IsNull() &&
				!config.UserIdentification.IdentityManagement.FilterIncludeIPAddressSet.IsUnknown() &&
				config.UserIdentification.IdentityManagement.FilterIncludeIPAddressBook.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("user_identification").AtName("identity_management").AtName("filter_include_ip_address_set"),
					tfdiag.MissingConfigErrSummary,
					"filter_include_ip_address_book must be specified with filter_include_ip_address_set"+
						" in identity_management block in user_identification block",
				)
			}
			if config.UserIdentification.IdentityManagement.Connection != nil {
				if config.UserIdentification.IdentityManagement.Connection.PrimaryAddress.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("user_identification").AtName("identity_management").
							AtName("connection").AtName("primary_address"),
						tfdiag.MissingConfigErrSummary,
						"primary_address must be specified"+
							" in connection block in identity_management block in user_identification block",
					)
				}
				if config.UserIdentification.IdentityManagement.Connection.PrimaryClientID.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("user_identification").AtName("identity_management").
							AtName("connection").AtName("primary_client_id"),
						tfdiag.MissingConfigErrSummary,
						"primary_client_id must be specified"+
							" in connection block in identity_management block in user_identification block",
					)
				}
				if config.UserIdentification.IdentityManagement.Connection.PrimaryClientSecret.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("user_identification").AtName("identity_management").
							AtName("connection").AtName("primary_client_secret"),
						tfdiag.MissingConfigErrSummary,
						"primary_client_secret must be specified"+
							" in connection block in identity_management block in user_identification block",
					)
				}
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("user_identification").AtName("identity_management").AtName("connection"),
					tfdiag.MissingConfigErrSummary,
					"connection block must be specified"+
						" in identity_management block in user_identification block",
				)
			}
		}
	}
}

func (rsc *services) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadComputed = &plan
	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *services) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesData
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
		func() {
			data.CleanOnDestroy = state.CleanOnDestroy
		},
		resp,
	)
}

func (rsc *services) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan servicesData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.junosClient().FakeUpdateAlso() {
		junSess := rsc.junosClient().NewSessionWithoutNetconf(ctx)

		if err := plan.delOpts(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		if err := plan.readComputed(ctx, junSess); err != nil {
			resp.Diagnostics.AddWarning(tfdiag.ConfigReadErrSummary, err.Error())
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
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

	if err := plan.delOpts(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	if err := plan.readComputed(ctx, junSess); err != nil {
		resp.Diagnostics.AddWarning(tfdiag.ConfigReadErrSummary, err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rsc *services) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.CleanOnDestroy.ValueBool() {
		defaultResourceDelete(
			ctx,
			rsc,
			&state,
			resp,
		)
	}
}

func (rsc *services) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesData

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

func (rscData *servicesData) fillID() {
	rscData.ID = types.StringValue("services")
}

func (rscData *servicesData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)

	if rscData.AdvancedAntiMalware != nil {
		if rscData.AdvancedAntiMalware.isEmpty() {
			return path.Root("advanced_anti_malware").AtName("*"),
				errors.New("advanced_anti_malware block is empty")
		}

		configSet = append(configSet, rscData.AdvancedAntiMalware.configSet()...)
	}
	if rscData.ApplicationIdentification != nil {
		blockSet, pathErr, err := rscData.ApplicationIdentification.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.SecurityIntelligence != nil {
		if rscData.SecurityIntelligence.isEmpty() {
			return path.Root("security_intelligence").AtName("*"),
				errors.New("security_intelligence block is empty")
		}

		blockSet, pathErr, err := rscData.SecurityIntelligence.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.UserIdentification != nil {
		if rscData.UserIdentification.isEmpty() {
			return path.Root("user_identification").AtName("*"),
				errors.New("user_identification block is empty")
		}

		configSet = append(configSet, rscData.UserIdentification.configSet()...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *servicesBlockAdvancedAntiMalware) configSet() []string {
	configSet := make([]string, 0, 100)

	if block.Connection != nil {
		configSet = append(configSet, block.Connection.configSet()...)
	}
	if block.DefaultPolicy != nil {
		configSet = append(configSet, block.DefaultPolicy.configSet()...)
	}

	return configSet
}

func (block *servicesBlockAdvancedAntiMalwareBlockConnection) configSet() []string {
	setPrefix := "set services advanced-anti-malware connection "
	delPrefix := junos.DeleteW + strings.TrimPrefix(setPrefix, junos.SetW)

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if v := block.AuthTLSProfile.ValueString(); v != "" {
		// delete old value only if new value must be set
		configSet = append(configSet, delPrefix+"authentication tls-profile")
		configSet = append(configSet, setPrefix+"authentication tls-profile \""+v+"\"")
	}
	if v := block.ProxyProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if v := block.SourceInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-interface "+v)
	}
	if v := block.URL.ValueString(); v != "" {
		// delete old value only if new value must be set
		configSet = append(configSet, delPrefix+"url")
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}

	return configSet
}

func (block *servicesBlockAdvancedAntiMalwareBlockDefaultPolicy) configSet() []string {
	setPrefix := "set services advanced-anti-malware default-policy "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if block.BlacklistNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"blacklist-notification log")
	}
	if block.DefaultNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"default-notification log")
	}
	if v := block.FallbackOptionsAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"fallback-options action "+v)
	}
	if block.FallbackOptionsNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"fallback-options notification log")
	}
	if v := block.HTTPAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http action "+v)
	}
	if v := block.HTTPClientNotifyFile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http client-notify file \""+v+"\"")
	}
	if v := block.HTTPClientNotifyMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http client-notify message \""+v+"\"")
	}
	if v := block.HTTPClientNotifyRedirectURL.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http client-notify redirect-url \""+v+"\"")
	}
	if v := block.HTTPFileVerdictUnknown.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http file-verdict-unknown "+v)
	}
	if v := block.HTTPInspectionProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http inspection-profile \""+v+"\"")
	}
	if block.HTTPNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"http notification log")
	}
	if v := block.IMAPInspectionProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"imap inspection-profile \""+v+"\"")
	}
	if block.IMAPNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"imap notification log")
	}
	if v := block.SMTPInspectionProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"smtp inspection-profile \""+v+"\"")
	}
	if block.SMTPNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"smtp notification log")
	}
	if v := block.VerdictThreshold.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"verdict-threshold "+v)
	}
	if block.WhitelistNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"whitelist-notification log")
	}

	return configSet
}

func (block *servicesBlockApplicationIdentification) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix := "set services application-identification "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !block.ApplicationSystemCacheTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"application-system-cache-timeout "+
			utils.ConvI64toa(block.ApplicationSystemCacheTimeout.ValueInt64()))
	}
	if !block.GlobalOffloadByteLimit.IsNull() {
		configSet = append(configSet, setPrefix+"global-offload-byte-limit "+
			utils.ConvI64toa(block.GlobalOffloadByteLimit.ValueInt64()))
	}
	if !block.IMAPCacheSize.IsNull() {
		configSet = append(configSet, setPrefix+"imap-cache-size "+
			utils.ConvI64toa(block.IMAPCacheSize.ValueInt64()))
	}
	if !block.IMAPCacheTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"imap-cache-timeout "+
			utils.ConvI64toa(block.IMAPCacheTimeout.ValueInt64()))
	}
	if !block.MaxMemory.IsNull() {
		configSet = append(configSet, setPrefix+"max-memory "+
			utils.ConvI64toa(block.MaxMemory.ValueInt64()))
	}
	if !block.MaxTransactions.IsNull() {
		configSet = append(configSet, setPrefix+"max-transactions "+
			utils.ConvI64toa(block.MaxTransactions.ValueInt64()))
	}
	if block.MicroApps.ValueBool() {
		configSet = append(configSet, setPrefix+"micro-apps")
	}
	if block.NoApplicationSystemCache.ValueBool() {
		configSet = append(configSet, setPrefix+"no-application-system-cache")
	}
	if !block.StatisticsInterval.IsNull() {
		configSet = append(configSet, setPrefix+"statistics interval "+
			utils.ConvI64toa(block.StatisticsInterval.ValueInt64()))
	}

	if block.ApplicationSystemCache != nil {
		configSet = append(configSet, block.ApplicationSystemCache.configSet()...)
	}
	if block.Download != nil {
		if block.Download.isEmpty() {
			return configSet,
				path.Root("application_identification").AtName("download").AtName("*"),
				errors.New("download block is empty" +
					" in application_identification block")
		}

		configSet = append(configSet, block.Download.configSet()...)
	}
	if block.EnablePerformanceMode != nil {
		configSet = append(configSet, setPrefix+"enable-performance-mode")

		if !block.EnablePerformanceMode.MaxPacketThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"enable-performance-mode max-packet-threshold "+
				utils.ConvI64toa(block.EnablePerformanceMode.MaxPacketThreshold.ValueInt64()))
		}
	}
	if block.InspectionLimitTCP != nil {
		configSet = append(configSet, block.InspectionLimitTCP.configSet(setPrefix+"inspection-limit tcp ")...)
	}
	if block.InspectionLimitUDP != nil {
		configSet = append(configSet, block.InspectionLimitUDP.configSet(setPrefix+"inspection-limit udp ")...)
	}

	return configSet, path.Empty(), nil
}

func (block *servicesBlockApplicationIdentificationBlockApplicationSystemCache) configSet() []string {
	setPrefix := "set services application-identification application-system-cache "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if block.NoMiscellaneousServices.ValueBool() {
		configSet = append(configSet, setPrefix+"no-miscellaneous-services")
	}
	if block.SecurityServices.ValueBool() {
		configSet = append(configSet, setPrefix+"security-services")
	}

	return configSet
}

func (block *servicesBlockApplicationIdentificationBlockDownload) configSet() []string {
	configSet := make([]string, 0, 100)
	setPrefix := "set services application-identification download "

	if !block.AutomaticInterval.IsNull() {
		configSet = append(configSet, setPrefix+"automatic interval "+
			utils.ConvI64toa(block.AutomaticInterval.ValueInt64()))
	}
	if v := block.AutomaticStartTime.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"automatic start-time "+v)
	}
	if block.IgnoreServerValidation.ValueBool() {
		configSet = append(configSet, setPrefix+"ignore-server-validation")
	}
	if v := block.ProxyProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
	}
	if v := block.URL.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}

	return configSet
}

func (block *servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP) configSet(setPrefix string) []string {
	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !block.ByteLimit.IsNull() {
		configSet = append(configSet, setPrefix+"byte-limit "+
			utils.ConvI64toa(block.ByteLimit.ValueInt64()))
	}
	if !block.PacketLimit.IsNull() {
		configSet = append(configSet, setPrefix+"packet-limit "+
			utils.ConvI64toa(block.PacketLimit.ValueInt64()))
	}

	return configSet
}

func (block *servicesBlockSecurityIntelligence) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix := "set services security-intelligence "
	delPrefix := junos.DeleteW + strings.TrimPrefix(setPrefix, junos.SetW)
	configSet := make([]string, 0, 100)

	if v := block.AuthenticationTLSProfile.ValueString(); v != "" {
		// delete old value (tls-profile or auth-token) only if new value must be set
		configSet = append(configSet, delPrefix+"authentication")
		configSet = append(configSet, setPrefix+"authentication tls-profile \""+v+"\"")
	}
	if v := block.AuthenticationToken.ValueString(); v != "" {
		// delete old value (tls-profile or auth-token) only if new value must be set
		configSet = append(configSet, delPrefix+"authentication")
		configSet = append(configSet, setPrefix+"authentication auth-token "+v)
	}

	for _, v := range block.CategoryDisable {
		if vv := v.ValueString(); vv == "all" {
			configSet = append(configSet, setPrefix+"category all disable")
		} else {
			configSet = append(configSet, setPrefix+"category category-name "+vv+" disable")
		}
	}
	if v := block.ProxyProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
	}
	if v := block.URL.ValueString(); v != "" {
		// delete old value only if new value must be set
		configSet = append(configSet, delPrefix+"url")
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}
	if v := block.URLParameter.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"url-parameter \""+v+"\"")
	}

	defaultPolicyCategoryName := make(map[string]struct{})
	for i, subBlock := range block.DefaultPolicy {
		categoryName := subBlock.CategoryName.ValueString()
		if _, ok := defaultPolicyCategoryName[categoryName]; ok {
			return configSet,
				path.Root("security_intelligence").AtName("default_policy").AtListIndex(i).AtName("category_name"),
				fmt.Errorf("multiple default_policy blocks with the same category_name %q"+
					" in security_intelligence block", categoryName)
		}
		defaultPolicyCategoryName[categoryName] = struct{}{}

		configSet = append(configSet, setPrefix+"default-policy "+categoryName+" "+subBlock.ProfileName.ValueString())
	}

	return configSet, path.Empty(), nil
}

func (block *servicesBlockUserIdentification) configSet() []string {
	setPrefix := "set services user-identification "
	configSet := make([]string, 0, 100)

	if v := block.DeviceInfoAuthSource.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"device-information authentication-source "+v)
	}

	if block.ADAccess != nil {
		configSet = append(configSet, block.ADAccess.configSet()...)
	}
	if block.IdentityManagement != nil {
		configSet = append(configSet, block.IdentityManagement.configSet()...)
	}

	return configSet
}

func (block *servicesBlockUserIdentificationBlockADAccess) configSet() []string {
	setPrefix := "set services user-identification active-directory-access "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !block.AuthEntryTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"authentication-entry-timeout "+
			utils.ConvI64toa(block.AuthEntryTimeout.ValueInt64()))
	}
	for _, v := range block.FilterExclude {
		configSet = append(configSet, setPrefix+"filter exclude "+v.ValueString())
	}
	for _, v := range block.FilterInclude {
		configSet = append(configSet, setPrefix+"filter include "+v.ValueString())
	}
	if !block.FirewallAuthForcedTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"firewall-authentication-forced-timeout "+
			utils.ConvI64toa(block.FirewallAuthForcedTimeout.ValueInt64()))
	}
	if !block.InvalidAuthEntryTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"invalid-authentication-entry-timeout "+
			utils.ConvI64toa(block.InvalidAuthEntryTimeout.ValueInt64()))
	}
	if block.NoOnDemandProbe.ValueBool() {
		configSet = append(configSet, setPrefix+"no-on-demand-probe")
	}
	if !block.WmiTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"wmi-timeout "+
			utils.ConvI64toa(block.WmiTimeout.ValueInt64()))
	}

	return configSet
}

func (block *servicesBlockUserIdentificationBlockIdentityManagement) configSet() []string {
	configSet := make([]string, 0, 100)
	setPrefix := "set services user-identification identity-management "

	if !block.AuthenticationEntryTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"authentication-entry-timeout "+
			utils.ConvI64toa(block.AuthenticationEntryTimeout.ValueInt64()))
	}
	if !block.BatchQueryInterval.IsNull() {
		configSet = append(configSet, setPrefix+"batch-query query-interval "+
			utils.ConvI64toa(block.BatchQueryInterval.ValueInt64()))
	}
	if !block.BatchQueryItemsPerBatch.IsNull() {
		configSet = append(configSet, setPrefix+"batch-query items-per-batch "+
			utils.ConvI64toa(block.BatchQueryItemsPerBatch.ValueInt64()))
	}
	for _, v := range block.FilterDomain {
		configSet = append(configSet, setPrefix+"filter domain "+v.ValueString())
	}
	if v := block.FilterExcludeIPAddressBook.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"filter exclude-ip address-book \""+v+"\"")
	}
	if v := block.FilterExcludeIPAddressSet.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"filter exclude-ip address-set \""+v+"\"")
	}
	if v := block.FilterIncludeIPAddressBook.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"filter include-ip address-book \""+v+"\"")
	}
	if v := block.FilterIncludeIPAddressSet.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"filter include-ip address-set \""+v+"\"")
	}
	if !block.InvalidAuthenticationEntryTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"invalid-authentication-entry-timeout "+
			utils.ConvI64toa(block.InvalidAuthenticationEntryTimeout.ValueInt64()))
	}
	if block.IPQueryDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"ip-query no-ip-query")
	}
	if !block.IPQueryDelayTime.IsNull() {
		configSet = append(configSet, setPrefix+"ip-query query-delay-time "+
			utils.ConvI64toa(block.IPQueryDelayTime.ValueInt64()))
	}

	if block.Connection != nil {
		configSet = append(configSet, block.Connection.configSet()...)
	}

	return configSet
}

func (block *servicesBlockUserIdentificationBlockIdentityManagementBlockConnection) configSet() []string {
	setPrefix := "set services user-identification identity-management connection "

	configSet := make([]string, 3, 100)
	configSet[0] = setPrefix + "primary address " + block.PrimaryAddress.ValueString()
	configSet[1] = setPrefix + "primary client-id \"" + block.PrimaryClientID.ValueString() + "\""
	configSet[2] = setPrefix + "primary client-secret \"" + block.PrimaryClientSecret.ValueString() + "\""

	if v := block.ConnectMethod.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"connect-method "+v)
	}
	if !block.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(block.Port.ValueInt64()))
	}
	if v := block.PrimaryCACertificate.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"primary ca-certificate \""+v+"\"")
	}
	if v := block.QueryAPI.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"query-api \""+v+"\"")
	}
	if v := block.SecondaryAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secondary address "+v)
	}
	if v := block.SecondaryCACertificate.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secondary ca-certificate \""+v+"\"")
	}
	if v := block.SecondaryClientID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secondary client-id \""+v+"\"")
	}
	if v := block.SecondaryClientSecret.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secondary client-secret \""+v+"\"")
	}
	if v := block.TokenAPI.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"token-api \""+v+"\"")
	}

	return configSet
}

func (rscData *servicesData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services" + junos.PipeDisplaySetRelative)
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
			case strings.HasPrefix(itemTrim, "advanced-anti-malware connection"),
				bchk.StringHasOneOfPrefixes(itemTrim, append(
					servicesBlockAdvancedAntiMalware{}.junosLines(),
					servicesBlockAdvancedAntiMalware{}.junosOptionalLines()...,
				)):
				if rscData.AdvancedAntiMalware == nil {
					rscData.AdvancedAntiMalware = &servicesBlockAdvancedAntiMalware{}
				}

				rscData.AdvancedAntiMalware.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "application-identification"):
				if rscData.ApplicationIdentification == nil {
					rscData.ApplicationIdentification = &servicesBlockApplicationIdentification{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.ApplicationIdentification.read(itemTrim); err != nil {
						return err
					}
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, append(
				servicesBlockSecurityIntelligence{}.junosLines(),
				servicesBlockSecurityIntelligence{}.junosOptionalLines()...,
			)):
				if rscData.SecurityIntelligence == nil {
					rscData.SecurityIntelligence = &servicesBlockSecurityIntelligence{}
				}

				if err := rscData.SecurityIntelligence.read(itemTrim, junSess); err != nil {
					return err
				}
			case strings.HasPrefix(itemTrim, "user-identification active-directory-access"),
				bchk.StringHasOneOfPrefixes(itemTrim, servicesBlockUserIdentification{}.junosLines()):
				if rscData.UserIdentification == nil {
					rscData.UserIdentification = &servicesBlockUserIdentification{}
				}

				if err := rscData.UserIdentification.read(itemTrim, junSess); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (servicesBlockAdvancedAntiMalware) junosLines() []string {
	return []string{
		"advanced-anti-malware connection proxy-profile",
		"advanced-anti-malware connection source-address",
		"advanced-anti-malware connection source-interface",
		"advanced-anti-malware default-policy",
	}
}

func (servicesBlockAdvancedAntiMalware) junosOptionalLines() []string {
	return servicesBlockAdvancedAntiMalwareBlockConnection{}.junosOptionalLines()
}

func (servicesBlockAdvancedAntiMalwareBlockConnection) junosOptionalLines() []string {
	return []string{
		"advanced-anti-malware connection authentication tls-profile",
		"advanced-anti-malware connection url",
	}
}

func (block *servicesBlockAdvancedAntiMalware) read(itemTrim string) {
	balt.CutPrefixInString(&itemTrim, "advanced-anti-malware ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "connection"):
		if block.Connection == nil {
			block.Connection = &servicesBlockAdvancedAntiMalwareBlockConnection{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.Connection.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "default-policy"):
		if block.DefaultPolicy == nil {
			block.DefaultPolicy = &servicesBlockAdvancedAntiMalwareBlockDefaultPolicy{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.DefaultPolicy.read(itemTrim)
		}
	}
}

func (block *servicesBlockAdvancedAntiMalwareBlockConnection) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication tls-profile "):
		block.AuthTLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "proxy-profile "):
		block.ProxyProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-interface "):
		block.SourceInterface = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "url "):
		block.URL = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (block *servicesBlockAdvancedAntiMalwareBlockDefaultPolicy) read(itemTrim string) {
	switch {
	case itemTrim == "blacklist-notification log":
		block.BlacklistNotificationLog = types.BoolValue(true)
	case itemTrim == "default-notification log":
		block.DefaultNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "fallback-options action "):
		block.FallbackOptionsAction = types.StringValue(itemTrim)
	case itemTrim == "fallback-options notification log":
		block.FallbackOptionsNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "http action "):
		block.HTTPAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "http client-notify file "):
		block.HTTPClientNotifyFile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "http client-notify message "):
		block.HTTPClientNotifyMessage = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "http client-notify redirect-url "):
		block.HTTPClientNotifyRedirectURL = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "http file-verdict-unknown "):
		block.HTTPFileVerdictUnknown = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "http inspection-profile "):
		block.HTTPInspectionProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "http notification log":
		block.HTTPNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "imap inspection-profile "):
		block.IMAPInspectionProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "imap notification log":
		block.IMAPNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "smtp inspection-profile "):
		block.SMTPInspectionProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "smtp notification log":
		block.SMTPNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "verdict-threshold "):
		block.VerdictThreshold = types.StringValue(itemTrim)
	case itemTrim == "whitelist-notification log":
		block.WhitelistNotificationLog = types.BoolValue(true)
	}
}

func (block *servicesBlockApplicationIdentification) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "application-system-cache-timeout "):
		block.ApplicationSystemCacheTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "application-system-cache"):
		if block.ApplicationSystemCache == nil {
			block.ApplicationSystemCache = &servicesBlockApplicationIdentificationBlockApplicationSystemCache{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.ApplicationSystemCache.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "download "):
		if block.Download == nil {
			block.Download = &servicesBlockApplicationIdentificationBlockDownload{}
		}

		err = block.Download.read(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "enable-performance-mode"):
		if block.EnablePerformanceMode == nil {
			block.EnablePerformanceMode = &servicesBlockApplicationIdentificationBlockEnablePerformanceMode{}
		}

		if balt.CutPrefixInString(&itemTrim, " max-packet-threshold ") {
			block.EnablePerformanceMode.MaxPacketThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "global-offload-byte-limit "):
		block.GlobalOffloadByteLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "imap-cache-size "):
		block.IMAPCacheSize, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "imap-cache-timeout "):
		block.IMAPCacheTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "inspection-limit tcp"):
		if block.InspectionLimitTCP == nil {
			block.InspectionLimitTCP = &servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			err = block.InspectionLimitTCP.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "inspection-limit udp"):
		if block.InspectionLimitUDP == nil {
			block.InspectionLimitUDP = &servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			err = block.InspectionLimitUDP.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "max-memory "):
		block.MaxMemory, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "max-transactions "):
		block.MaxTransactions, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "micro-apps":
		block.MicroApps = types.BoolValue(true)
	case itemTrim == "no-application-system-cache":
		block.NoApplicationSystemCache = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "statistics interval "):
		block.StatisticsInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *servicesBlockApplicationIdentificationBlockApplicationSystemCache) read(itemTrim string) {
	switch {
	case itemTrim == "no-miscellaneous-services":
		block.NoMiscellaneousServices = types.BoolValue(true)
	case itemTrim == "security-services":
		block.SecurityServices = types.BoolValue(true)
	}
}

func (block *servicesBlockApplicationIdentificationBlockDownload) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "automatic interval "):
		block.AutomaticInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "automatic start-time "):
		block.AutomaticStartTime = types.StringValue(itemTrim)
	case itemTrim == "ignore-server-validation":
		block.IgnoreServerValidation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "proxy-profile "):
		block.ProxyProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "url "):
		block.URL = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return err
}

func (block *servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "byte-limit "):
		block.ByteLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "packet-limit "):
		block.PacketLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (servicesBlockSecurityIntelligence) junosLines() []string {
	return []string{
		"security-intelligence category",
		"security-intelligence default-policy",
		"security-intelligence proxy-profile",
		"security-intelligence url-parameter",
	}
}

func (servicesBlockSecurityIntelligence) junosOptionalLines() []string {
	return []string{
		"security-intelligence authentication auth-token",
		"security-intelligence authentication tls-profile",
		"security-intelligence url",
	}
}

func (block *servicesBlockSecurityIntelligence) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	balt.CutPrefixInString(&itemTrim, "security-intelligence ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication auth-token "):
		block.AuthenticationToken = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "authentication tls-profile "):
		block.AuthenticationTLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "category all disable":
		block.CategoryDisable = append(block.CategoryDisable, types.StringValue("all"))
	case balt.CutPrefixInString(&itemTrim, "category category-name ") &&
		balt.CutSuffixInString(&itemTrim, " disable"):
		block.CategoryDisable = append(block.CategoryDisable, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "default-policy "):
		if itemTrimFields := strings.Split(itemTrim, " "); len(itemTrimFields) == 2 { // <category_name> <profile_name>
			block.DefaultPolicy = append(block.DefaultPolicy, servicesBlockSecurityIntelligenceBlockDefaultPolicy{
				CategoryName: types.StringValue(itemTrimFields[0]),
				ProfileName:  types.StringValue(itemTrimFields[1]),
			})
		}
	case balt.CutPrefixInString(&itemTrim, "proxy-profile "):
		block.ProxyProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "url "):
		block.URL = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "url-parameter "):
		block.URLParameter, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "url-parameter")
	}

	return err
}

func (servicesBlockUserIdentification) junosLines() []string {
	r := []string{
		"user-identification device-information authentication-source",
		"user-identification identity-management",
	}
	r = append(r, servicesBlockUserIdentificationBlockADAccess{}.junosLines()...)

	return r
}

func (servicesBlockUserIdentificationBlockADAccess) junosLines() []string {
	return []string{
		"user-identification active-directory-access authentication-entry-timeout",
		"user-identification active-directory-access filter",
		"user-identification active-directory-access firewall-authentication-forced-timeout",
		"user-identification active-directory-access invalid-authentication-entry-timeout",
		"user-identification active-directory-access no-on-demand-probe",
		"user-identification active-directory-access wmi-timeout",
	}
}

func (block *servicesBlockUserIdentification) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	balt.CutPrefixInString(&itemTrim, "user-identification ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "active-directory-access"):
		if block.ADAccess == nil {
			block.ADAccess = &servicesBlockUserIdentificationBlockADAccess{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			err = block.ADAccess.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "device-information authentication-source "):
		block.DeviceInfoAuthSource = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "identity-management "):
		if block.IdentityManagement == nil {
			block.IdentityManagement = &servicesBlockUserIdentificationBlockIdentityManagement{}
		}

		err = block.IdentityManagement.read(itemTrim, junSess)
	}

	return err
}

func (block *servicesBlockUserIdentificationBlockADAccess) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication-entry-timeout "):
		block.AuthEntryTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "filter exclude "):
		block.FilterExclude = append(block.FilterExclude, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "filter include "):
		block.FilterInclude = append(block.FilterInclude, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "firewall-authentication-forced-timeout "):
		block.FirewallAuthForcedTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "invalid-authentication-entry-timeout "):
		block.InvalidAuthEntryTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "no-on-demand-probe":
		block.NoOnDemandProbe = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "wmi-timeout "):
		block.WmiTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *servicesBlockUserIdentificationBlockIdentityManagement) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication-entry-timeout "):
		block.AuthenticationEntryTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "batch-query items-per-batch "):
		block.BatchQueryItemsPerBatch, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "batch-query query-interval "):
		block.BatchQueryInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "connection "):
		if block.Connection == nil {
			block.Connection = &servicesBlockUserIdentificationBlockIdentityManagementBlockConnection{}
		}

		err = block.Connection.read(itemTrim, junSess)
	case balt.CutPrefixInString(&itemTrim, "filter domain "):
		block.FilterDomain = append(block.FilterDomain, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "filter exclude-ip address-book "):
		block.FilterExcludeIPAddressBook = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "filter exclude-ip address-set "):
		block.FilterExcludeIPAddressSet = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "filter include-ip address-book "):
		block.FilterIncludeIPAddressBook = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "filter include-ip address-set "):
		block.FilterIncludeIPAddressSet = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "invalid-authentication-entry-timeout "):
		block.InvalidAuthenticationEntryTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "ip-query no-ip-query":
		block.IPQueryDisable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ip-query query-delay-time "):
		block.IPQueryDelayTime, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *servicesBlockUserIdentificationBlockIdentityManagementBlockConnection) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "primary address "):
		block.PrimaryAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "primary client-id "):
		block.PrimaryClientID = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "primary client-secret "):
		block.PrimaryClientSecret, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "primary client-secret")
	case balt.CutPrefixInString(&itemTrim, "connect-method "):
		block.ConnectMethod = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "primary ca-certificate "):
		block.PrimaryCACertificate = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "query-api "):
		block.QueryAPI = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "secondary address "):
		block.SecondaryAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "secondary ca-certificate "):
		block.SecondaryCACertificate = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "secondary client-id "):
		block.SecondaryClientID = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "secondary client-secret "):
		block.SecondaryClientSecret, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "secondary client-secret")
	case balt.CutPrefixInString(&itemTrim, "token-api "):
		block.TokenAPI = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return err
}

func (rscData *servicesData) readComputed(
	_ context.Context, junSess *junos.Session,
) error {
	defer func() {
		// set unknown to null if still unknown after reading config
		if rscData.AdvancedAntiMalware != nil &&
			rscData.AdvancedAntiMalware.Connection != nil {
			if rscData.AdvancedAntiMalware.Connection.AuthTLSProfile.IsUnknown() {
				rscData.AdvancedAntiMalware.Connection.AuthTLSProfile = types.StringNull()
			}
			if rscData.AdvancedAntiMalware.Connection.URL.IsUnknown() {
				rscData.AdvancedAntiMalware.Connection.URL = types.StringNull()
			}
		}
		if rscData.SecurityIntelligence != nil {
			if rscData.SecurityIntelligence.AuthenticationTLSProfile.IsUnknown() {
				rscData.SecurityIntelligence.AuthenticationTLSProfile = types.StringNull()
			}
			if rscData.SecurityIntelligence.AuthenticationToken.IsUnknown() {
				rscData.SecurityIntelligence.AuthenticationToken = types.StringNull()
			}
			if rscData.SecurityIntelligence.URL.IsUnknown() {
				rscData.SecurityIntelligence.URL = types.StringNull()
			}
		}
	}()

	if !junSess.HasNetconf() {
		return nil
	}

	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "advanced-anti-malware connection authentication tls-profile ") {
				if rscData.AdvancedAntiMalware != nil &&
					rscData.AdvancedAntiMalware.Connection != nil &&
					rscData.AdvancedAntiMalware.Connection.AuthTLSProfile.IsUnknown() {
					rscData.AdvancedAntiMalware.Connection.AuthTLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			}
			if balt.CutPrefixInString(&itemTrim, "advanced-anti-malware connection url ") {
				if rscData.AdvancedAntiMalware != nil &&
					rscData.AdvancedAntiMalware.Connection != nil &&
					rscData.AdvancedAntiMalware.Connection.URL.IsUnknown() {
					rscData.AdvancedAntiMalware.Connection.URL = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			}
			if balt.CutPrefixInString(&itemTrim, "security-intelligence authentication tls-profile ") {
				if rscData.SecurityIntelligence != nil &&
					rscData.SecurityIntelligence.AuthenticationTLSProfile.IsUnknown() {
					rscData.SecurityIntelligence.AuthenticationTLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			}
			if balt.CutPrefixInString(&itemTrim, "security-intelligence authentication auth-token ") {
				if rscData.SecurityIntelligence != nil &&
					rscData.SecurityIntelligence.AuthenticationToken.IsUnknown() {
					rscData.SecurityIntelligence.AuthenticationToken = types.StringValue(itemTrim)
				}
			}
			if balt.CutPrefixInString(&itemTrim, "security-intelligence url ") {
				if rscData.SecurityIntelligence != nil &&
					rscData.SecurityIntelligence.URL.IsUnknown() {
					rscData.SecurityIntelligence.URL = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			}
		}
	}

	return nil
}

// Use plan instead of state to determine if need to delete extra config.
func (rscData *servicesData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	listLinesToDelete := make([]string, 0, 100)
	listLinesToDelete = append(listLinesToDelete, servicesBlockAdvancedAntiMalware{}.junosLines()...)
	if rscData.AdvancedAntiMalware == nil {
		listLinesToDelete = append(listLinesToDelete, servicesBlockAdvancedAntiMalware{}.junosOptionalLines()...)
		listLinesToDelete = append(listLinesToDelete, "advanced-anti-malware connection")
	} else if rscData.AdvancedAntiMalware.Connection == nil {
		listLinesToDelete = append(listLinesToDelete, "advanced-anti-malware connection")
	}
	listLinesToDelete = append(listLinesToDelete, "application-identification")
	listLinesToDelete = append(listLinesToDelete, servicesBlockSecurityIntelligence{}.junosLines()...)
	if rscData.SecurityIntelligence == nil {
		listLinesToDelete = append(listLinesToDelete, servicesBlockSecurityIntelligence{}.junosOptionalLines()...)
	}
	listLinesToDelete = append(listLinesToDelete, servicesBlockUserIdentification{}.junosLines()...)
	if rscData.UserIdentification == nil || rscData.UserIdentification.ADAccess == nil {
		listLinesToDelete = append(listLinesToDelete, "user-identification active-directory-access")
	}

	configSet := make([]string, len(listLinesToDelete))
	delPrefix := "delete services "
	for i, line := range listLinesToDelete {
		configSet[i] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *servicesData) del(
	_ context.Context, junSess *junos.Session,
) error {
	listLinesToDelete := make([]string, 0, 100)
	listLinesToDelete = append(listLinesToDelete, servicesBlockAdvancedAntiMalware{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, servicesBlockAdvancedAntiMalware{}.junosOptionalLines()...)
	listLinesToDelete = append(listLinesToDelete, "advanced-anti-malware connection")
	listLinesToDelete = append(listLinesToDelete, "application-identification")
	listLinesToDelete = append(listLinesToDelete, servicesBlockSecurityIntelligence{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, servicesBlockSecurityIntelligence{}.junosOptionalLines()...)
	listLinesToDelete = append(listLinesToDelete, servicesBlockUserIdentification{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, "user-identification active-directory-access")

	configSet := make([]string, len(listLinesToDelete))
	delPrefix := "delete services "
	for i, line := range listLinesToDelete {
		configSet[i] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}
