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
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type bgpAttrData struct {
	AcceptRemoteNexthop          types.Bool                    `tfsdk:"accept_remote_nexthop"`
	AdvertiseExternal            types.Bool                    `tfsdk:"advertise_external"`
	AdvertiseExternalConditional types.Bool                    `tfsdk:"advertise_external_conditional"`
	AdvertiseInactive            types.Bool                    `tfsdk:"advertise_inactive"`
	AdvertisePeerAS              types.Bool                    `tfsdk:"advertise_peer_as"`
	NoAdvertisePeerAS            types.Bool                    `tfsdk:"no_advertise_peer_as"`
	ASOverride                   types.Bool                    `tfsdk:"as_override"`
	AuthenticationAlgorithm      types.String                  `tfsdk:"authentication_algorithm"`
	AuthenticationKey            types.String                  `tfsdk:"authentication_key"`
	AuthenticationKeyChain       types.String                  `tfsdk:"authentication_key_chain"`
	Cluster                      types.String                  `tfsdk:"cluster"`
	Damping                      types.Bool                    `tfsdk:"damping"`
	Description                  types.String                  `tfsdk:"description"`
	Export                       []types.String                `tfsdk:"export"`
	HoldTime                     types.Int64                   `tfsdk:"hold_time"`
	Import                       []types.String                `tfsdk:"import"`
	KeepAll                      types.Bool                    `tfsdk:"keep_all"`
	KeepNone                     types.Bool                    `tfsdk:"keep_none"`
	LocalAddress                 types.String                  `tfsdk:"local_address"`
	LocalAS                      types.String                  `tfsdk:"local_as"`
	LocalASAlias                 types.Bool                    `tfsdk:"local_as_alias"`
	LocalASLoops                 types.Int64                   `tfsdk:"local_as_loops"`
	LocalASNoPrependGlobalAS     types.Bool                    `tfsdk:"local_as_no_prepend_global_as"`
	LocalASPrivate               types.Bool                    `tfsdk:"local_as_private"`
	LocalInterface               types.String                  `tfsdk:"local_interface"`
	LocalPreference              types.Int64                   `tfsdk:"local_preference"`
	LogUpdown                    types.Bool                    `tfsdk:"log_updown"`
	MetricOut                    types.Int64                   `tfsdk:"metric_out"`
	MetricOutIgp                 types.Bool                    `tfsdk:"metric_out_igp"`
	MetricOutIgpDelayMedUpdate   types.Bool                    `tfsdk:"metric_out_igp_delay_med_update"`
	MetricOutIgpOffset           types.Int64                   `tfsdk:"metric_out_igp_offset"`
	MetricOutMinimumIgp          types.Bool                    `tfsdk:"metric_out_minimum_igp"`
	MetricOutMinimumIgpOffset    types.Int64                   `tfsdk:"metric_out_minimum_igp_offset"`
	MtuDiscovery                 types.Bool                    `tfsdk:"mtu_discovery"`
	Multihop                     types.Bool                    `tfsdk:"multihop"`
	NoClientReflect              types.Bool                    `tfsdk:"no_client_reflect"`
	OutDelay                     types.Int64                   `tfsdk:"out_delay"`
	Passive                      types.Bool                    `tfsdk:"passive"`
	PeerAS                       types.String                  `tfsdk:"peer_as"`
	Preference                   types.Int64                   `tfsdk:"preference"`
	RemovePrivate                types.Bool                    `tfsdk:"remove_private"`
	TCPAggressiveTransmission    types.Bool                    `tfsdk:"tcp_aggressive_transmission"`
	BfdLivenessDetection         *bgpBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
	BgpErrorTolerance            *bgpBlockBgpErrorTolerance    `tfsdk:"bgp_error_tolerance"`
	BgpMultipath                 *bgpBlockBgpMultipath         `tfsdk:"bgp_multipath"`
	FamilyEvpn                   []bgpBlockFamily              `tfsdk:"family_evpn"`
	FamilyInet                   []bgpBlockFamily              `tfsdk:"family_inet"`
	FamilyInet6                  []bgpBlockFamily              `tfsdk:"family_inet6"`
	GracefulRestart              *bgpBlockGracefulRestart      `tfsdk:"graceful_restart"`
}

func (bgpAttrData) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"accept_remote_nexthop": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow import policy to specify a non-directly connected next-hop.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"advertise_external": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Advertise best external routes.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"advertise_external_conditional": schema.BoolAttribute{
			Optional:    true,
			Description: "Route matches active route upto med-comparison rule.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"advertise_inactive": schema.BoolAttribute{
			Optional:    true,
			Description: "Advertise inactive routes.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"advertise_peer_as": schema.BoolAttribute{
			Optional:    true,
			Description: "Advertise routes received from the same autonomous system.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_advertise_peer_as": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't advertise routes received from the same autonomous system.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"as_override": schema.BoolAttribute{
			Optional:    true,
			Description: "Replace neighbor AS number with our AS number.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"authentication_algorithm": schema.StringAttribute{
			Optional:    true,
			Description: "Authentication algorithm name.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"authentication_key": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "MD5 authentication key.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 126),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"authentication_key_chain": schema.StringAttribute{
			Optional:    true,
			Description: "Key chain name.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 128),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"cluster": schema.StringAttribute{
			Optional:    true,
			Description: "Cluster identifier.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress(),
			},
		},
		"damping": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable route flap damping.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Description: "Text description.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 900),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"export": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Export policy list.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.NoNullValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"hold_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Hold time used when negotiating with a peer.",
			Validators: []validator.Int64{
				int64validator.Between(3, 65535),
			},
		},
		"import": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Import policy list.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.NoNullValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"keep_all": schema.BoolAttribute{
			Optional:    true,
			Description: "Retain all routes.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"keep_none": schema.BoolAttribute{
			Optional:    true,
			Description: "Retain no routes.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"local_address": schema.StringAttribute{
			Optional:    true,
			Description: "Address of local end of BGP session.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress(),
			},
		},
		"local_as": schema.StringAttribute{
			Optional:    true,
			Description: "Local autonomous system number.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^\d+(\.\d+)?$`),
					"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
			},
		},
		"local_as_alias": schema.BoolAttribute{
			Optional:    true,
			Description: "Treat this AS as an alias to the system AS.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"local_as_loops": schema.Int64Attribute{
			Optional:    true,
			Description: "Maximum number of times this AS can be in an AS path (1..10).",
			Validators: []validator.Int64{
				int64validator.Between(1, 10),
			},
		},
		"local_as_no_prepend_global_as": schema.BoolAttribute{
			Optional:    true,
			Description: "Do not prepend global autonomous-system number in advertised paths.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"local_as_private": schema.BoolAttribute{
			Optional:    true,
			Description: "Hide this local AS in paths learned from this peering.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"local_interface": schema.StringAttribute{
			Optional:    true,
			Description: "Local interface for IPv6 link local EBGP peering.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
				tfvalidator.String1DotCount(),
			},
		},
		"local_preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Value of LOCAL_PREF path attribute.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"log_updown": schema.BoolAttribute{
			Optional:    true,
			Description: "Log a message for peer state transitions.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"metric_out": schema.Int64Attribute{
			Optional:    true,
			Description: "Route metric sent in MED.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"metric_out_igp": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Track the IGP metric.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"metric_out_igp_delay_med_update": schema.BoolAttribute{
			Optional:    true,
			Description: "Delay updating MED when IGP metric increases.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"metric_out_igp_offset": schema.Int64Attribute{
			Optional:    true,
			Description: "Metric offset for MED.",
			Validators: []validator.Int64{
				int64validator.Between(-2147483648, 2147483647),
			},
		},
		"metric_out_minimum_igp": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Track the minimum IGP metric.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"metric_out_minimum_igp_offset": schema.Int64Attribute{
			Optional:    true,
			Description: "Metric offset for MED.",
			Validators: []validator.Int64{
				int64validator.Between(-2147483648, 2147483647),
			},
		},
		"mtu_discovery": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable TCP path MTU discovery.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"multihop": schema.BoolAttribute{
			Optional:    true,
			Description: "Configure an EBGP multihop session.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_client_reflect": schema.BoolAttribute{
			Optional:    true,
			Description: "Disable intracluster route redistribution.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"out_delay": schema.Int64Attribute{
			Optional:    true,
			Description: "How long before exporting routes from routing table.",
			Validators: []validator.Int64{
				int64validator.Between(1, 65535),
			},
		},
		"passive": schema.BoolAttribute{
			Optional:    true,
			Description: "Do not send open messages to a peer.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"peer_as": schema.StringAttribute{
			Optional:    true,
			Description: "Autonomous system number.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^\d+(\.\d+)?$`),
					"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
			},
		},
		"preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Preference value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"remove_private": schema.BoolAttribute{
			Optional:    true,
			Description: "Remove well-known private AS numbers.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"tcp_aggressive_transmission": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable aggressive transmission of pure TCP ACKs and retransmissions.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (bgpAttrData) blocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"bfd_liveness_detection": bgpBlockBfdLivenessDetection{}.schema(),
		"bgp_error_tolerance":    bgpBlockBgpErrorTolerance{}.schema(),
		"bgp_multipath":          bgpBlockBgpMultipath{}.schema(),
		"family_evpn":            bgpBlockFamily{}.schema("EVPN"),
		"family_inet":            bgpBlockFamily{}.schema("IPv4"),
		"family_inet6":           bgpBlockFamily{}.schema("IPv6"),
		"graceful_restart":       bgpBlockGracefulRestart{}.schema(),
	}
}

func (rscData *bgpAttrData) configSet(setPrefix string) ([]string, path.Path, error) {
	configSet := make([]string, 0, 100)

	if rscData.AcceptRemoteNexthop.ValueBool() {
		configSet = append(configSet, setPrefix+"accept-remote-nexthop")
	}
	if rscData.AdvertiseExternal.ValueBool() {
		configSet = append(configSet, setPrefix+"advertise-external")
	}
	if rscData.AdvertiseExternalConditional.ValueBool() {
		configSet = append(configSet, setPrefix+"advertise-external conditional")
	}
	if rscData.AdvertiseInactive.ValueBool() {
		configSet = append(configSet, setPrefix+"advertise-inactive")
	}
	if rscData.AdvertisePeerAS.ValueBool() {
		configSet = append(configSet, setPrefix+"advertise-peer-as")
	}
	if rscData.NoAdvertisePeerAS.ValueBool() {
		configSet = append(configSet, setPrefix+"no-advertise-peer-as")
	}
	if rscData.ASOverride.ValueBool() {
		configSet = append(configSet, setPrefix+"as-override")
	}
	if v := rscData.AuthenticationAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-algorithm "+v)
	}
	if v := rscData.AuthenticationKey.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key \""+v+"\"")
	}
	if v := rscData.AuthenticationKeyChain.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key-chain \""+v+"\"")
	}
	if v := rscData.Cluster.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"cluster "+v)
	}
	if rscData.Damping.ValueBool() {
		configSet = append(configSet, setPrefix+"damping")
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range rscData.Export {
		configSet = append(configSet, setPrefix+"export "+v.ValueString())
	}
	if !rscData.HoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"hold-time "+
			utils.ConvI64toa(rscData.HoldTime.ValueInt64()))
	}
	for _, v := range rscData.Import {
		configSet = append(configSet, setPrefix+"import "+v.ValueString())
	}
	if rscData.KeepAll.ValueBool() {
		configSet = append(configSet, setPrefix+"keep all")
	}
	if rscData.KeepNone.ValueBool() {
		configSet = append(configSet, setPrefix+"keep none")
	}
	if v := rscData.LocalAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"local-address "+v)
	}
	if v := rscData.LocalAS.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"local-as "+v)
	}
	if rscData.LocalASAlias.ValueBool() {
		configSet = append(configSet, setPrefix+"local-as alias")
	}
	if !rscData.LocalASLoops.IsNull() {
		configSet = append(configSet, setPrefix+"local-as loops "+
			utils.ConvI64toa(rscData.LocalASLoops.ValueInt64()))
	}
	if rscData.LocalASNoPrependGlobalAS.ValueBool() {
		configSet = append(configSet, setPrefix+"local-as no-prepend-global-as")
	}
	if rscData.LocalASPrivate.ValueBool() {
		configSet = append(configSet, setPrefix+"local-as private")
	}
	if v := rscData.LocalInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"local-interface "+v)
	}
	if !rscData.LocalPreference.IsNull() {
		configSet = append(configSet, setPrefix+"local-preference "+
			utils.ConvI64toa(rscData.LocalPreference.ValueInt64()))
	}
	if rscData.LogUpdown.ValueBool() {
		configSet = append(configSet, setPrefix+"log-updown")
	}
	if !rscData.MetricOut.IsNull() {
		configSet = append(configSet, setPrefix+"metric-out "+
			utils.ConvI64toa(rscData.MetricOut.ValueInt64()))
	}
	if rscData.MetricOutIgp.ValueBool() {
		configSet = append(configSet, setPrefix+"metric-out igp")
	}
	if rscData.MetricOutIgpDelayMedUpdate.ValueBool() {
		configSet = append(configSet, setPrefix+"metric-out igp delay-med-update")
	}
	if !rscData.MetricOutIgpOffset.IsNull() {
		configSet = append(configSet, setPrefix+"metric-out igp "+
			utils.ConvI64toa(rscData.MetricOutIgpOffset.ValueInt64()))
	}
	if rscData.MetricOutMinimumIgp.ValueBool() {
		configSet = append(configSet, setPrefix+"metric-out minimum-igp")
	}
	if !rscData.MetricOutMinimumIgpOffset.IsNull() {
		configSet = append(configSet, setPrefix+"metric-out minimum-igp "+
			utils.ConvI64toa(rscData.MetricOutMinimumIgpOffset.ValueInt64()))
	}
	if rscData.MtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"mtu-discovery")
	}
	if rscData.Multihop.ValueBool() {
		configSet = append(configSet, setPrefix+"multihop")
	}
	if rscData.NoClientReflect.ValueBool() {
		configSet = append(configSet, setPrefix+"no-client-reflect")
	}
	if !rscData.OutDelay.IsNull() {
		configSet = append(configSet, setPrefix+"out-delay "+
			utils.ConvI64toa(rscData.OutDelay.ValueInt64()))
	}
	if rscData.Passive.ValueBool() {
		configSet = append(configSet, setPrefix+"passive")
	}
	if v := rscData.PeerAS.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"peer-as "+v)
	}
	if !rscData.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(rscData.Preference.ValueInt64()))
	}
	if rscData.RemovePrivate.ValueBool() {
		configSet = append(configSet, setPrefix+"remove-private")
	}
	if rscData.TCPAggressiveTransmission.ValueBool() {
		configSet = append(configSet, setPrefix+"tcp-aggressive-transmission")
	}
	if rscData.BfdLivenessDetection != nil {
		if rscData.BfdLivenessDetection.isEmpty() {
			return configSet,
				path.Root("bfd_liveness_detection").AtName("*"),
				errors.New("bfd_liveness_detection block is empty")
		}

		configSet = append(configSet, rscData.BfdLivenessDetection.configSet(setPrefix)...)
	}
	if rscData.BgpErrorTolerance != nil {
		configSet = append(configSet, rscData.BgpErrorTolerance.configSet(setPrefix)...)
	}
	if rscData.BgpMultipath != nil {
		configSet = append(configSet, rscData.BgpMultipath.configSet(setPrefix)...)
	}
	familyEvpnNlriType := make(map[string]struct{})
	for i, block := range rscData.FamilyEvpn {
		nlriType := block.NlriType.ValueString()
		if _, ok := familyEvpnNlriType[nlriType]; ok {
			return configSet, path.Root("family_evpn").AtListIndex(i).AtName("nlri_type"),
				fmt.Errorf("multiple family_evpn blocks with the same nlri_type %q", nlriType)
		}
		familyEvpnNlriType[nlriType] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix+"family evpn ", path.Root("family_evpn").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	familyInetNlriType := make(map[string]struct{})
	for i, block := range rscData.FamilyInet {
		nlriType := block.NlriType.ValueString()
		if _, ok := familyInetNlriType[nlriType]; ok {
			return configSet,
				path.Root("family_inet").AtListIndex(i).AtName("nlri_type"),
				fmt.Errorf("multiple family_inet blocks with the same nlri_type %q", nlriType)
		}
		familyInetNlriType[nlriType] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix+"family inet ", path.Root("family_inet").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	familyInet6NlriType := make(map[string]struct{})
	for i, block := range rscData.FamilyInet6 {
		nlriType := block.NlriType.ValueString()
		if _, ok := familyInet6NlriType[nlriType]; ok {
			return configSet,
				path.Root("family_inet6").AtListIndex(i).AtName("nlri_type"),
				fmt.Errorf("multiple family_inet6 blocks with the same nlri_type %q", nlriType)
		}
		familyInet6NlriType[nlriType] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix+"family inet6 ", path.Root("family_inet6").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.GracefulRestart != nil {
		configSet = append(configSet, rscData.GracefulRestart.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (rscData *bgpAttrData) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	switch {
	case itemTrim == "accept-remote-nexthop":
		rscData.AcceptRemoteNexthop = types.BoolValue(true)
	case itemTrim == "advertise-external":
		rscData.AdvertiseExternal = types.BoolValue(true)
	case itemTrim == "advertise-external conditional":
		rscData.AdvertiseExternal = types.BoolValue(true)
		rscData.AdvertiseExternalConditional = types.BoolValue(true)
	case itemTrim == "advertise-inactive":
		rscData.AdvertiseInactive = types.BoolValue(true)
	case itemTrim == "advertise-peer-as":
		rscData.AdvertisePeerAS = types.BoolValue(true)
	case itemTrim == "no-advertise-peer-as":
		rscData.NoAdvertisePeerAS = types.BoolValue(true)
	case itemTrim == "as-override":
		rscData.ASOverride = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "authentication-algorithm "):
		rscData.AuthenticationAlgorithm = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "authentication-key "):
		rscData.AuthenticationKey, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "authentication-key-chain "):
		rscData.AuthenticationKeyChain = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "cluster "):
		rscData.Cluster = types.StringValue(itemTrim)
	case itemTrim == "damping":
		rscData.Damping = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "description "):
		rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "export "):
		rscData.Export = append(rscData.Export, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "hold-time "):
		rscData.HoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "import "):
		rscData.Import = append(rscData.Import, types.StringValue(itemTrim))
	case itemTrim == "keep all":
		rscData.KeepAll = types.BoolValue(true)
	case itemTrim == "keep none":
		rscData.KeepNone = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "local-address "):
		rscData.LocalAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-as "):
		switch {
		case itemTrim == "private":
			rscData.LocalASPrivate = types.BoolValue(true)
		case itemTrim == "alias":
			rscData.LocalASAlias = types.BoolValue(true)
		case itemTrim == "no-prepend-global-as":
			rscData.LocalASNoPrependGlobalAS = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "loops "):
			rscData.LocalASLoops, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		default:
			rscData.LocalAS = types.StringValue(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "local-interface "):
		rscData.LocalInterface = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		rscData.LocalPreference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "log-updown":
		rscData.LogUpdown = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "metric-out "):
		switch {
		case balt.CutPrefixInString(&itemTrim, "igp"):
			rscData.MetricOutIgp = types.BoolValue(true)
			switch {
			case itemTrim == " delay-med-update":
				rscData.MetricOutIgpDelayMedUpdate = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, " "):
				rscData.MetricOutIgpOffset, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-igp"):
			rscData.MetricOutMinimumIgp = types.BoolValue(true)
			if balt.CutPrefixInString(&itemTrim, " ") {
				rscData.MetricOutMinimumIgpOffset, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		default:
			rscData.MetricOut, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case itemTrim == "mtu-discovery":
		rscData.MtuDiscovery = types.BoolValue(true)
	case itemTrim == "multihop":
		rscData.Multihop = types.BoolValue(true)
	case itemTrim == "no-client-reflect":
		rscData.NoClientReflect = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "out-delay "):
		rscData.OutDelay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "passive":
		rscData.Passive = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "peer-as "):
		rscData.PeerAS = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "preference "):
		rscData.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "remove-private":
		rscData.RemovePrivate = types.BoolValue(true)
	case itemTrim == "tcp-aggressive-transmission":
		rscData.TCPAggressiveTransmission = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
		if rscData.BfdLivenessDetection == nil {
			rscData.BfdLivenessDetection = &bgpBlockBfdLivenessDetection{}
		}
		if err := rscData.BfdLivenessDetection.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "bgp-error-tolerance"):
		if rscData.BgpErrorTolerance == nil {
			rscData.BgpErrorTolerance = &bgpBlockBgpErrorTolerance{}
		}
		if err := rscData.BgpErrorTolerance.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "family evpn "):
		nlriType := tfdata.FirstElementOfJunosLine(itemTrim)
		rscData.FamilyEvpn = tfdata.AppendPotentialNewBlock(rscData.FamilyEvpn, types.StringValue(nlriType))
		familyEvpn := &rscData.FamilyEvpn[len(rscData.FamilyEvpn)-1]

		if balt.CutPrefixInString(&itemTrim, nlriType+" ") {
			if err := familyEvpn.read(itemTrim); err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "family inet "):
		nlriType := tfdata.FirstElementOfJunosLine(itemTrim)
		rscData.FamilyInet = tfdata.AppendPotentialNewBlock(rscData.FamilyInet, types.StringValue(nlriType))
		familyInet := &rscData.FamilyInet[len(rscData.FamilyInet)-1]

		if balt.CutPrefixInString(&itemTrim, nlriType+" ") {
			if err := familyInet.read(itemTrim); err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "family inet6 "):
		nlriType := tfdata.FirstElementOfJunosLine(itemTrim)
		rscData.FamilyInet6 = tfdata.AppendPotentialNewBlock(rscData.FamilyInet6, types.StringValue(nlriType))
		familyInet6 := &rscData.FamilyInet6[len(rscData.FamilyInet6)-1]

		if balt.CutPrefixInString(&itemTrim, nlriType+" ") {
			if err := familyInet6.read(itemTrim); err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "multipath"):
		if rscData.BgpMultipath == nil {
			rscData.BgpMultipath = &bgpBlockBgpMultipath{}
		}
		rscData.BgpMultipath.read(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "graceful-restart"):
		if rscData.GracefulRestart == nil {
			rscData.GracefulRestart = &bgpBlockGracefulRestart{}
		}
		if err := rscData.GracefulRestart.read(itemTrim); err != nil {
			return err
		}
	}

	return nil
}

func (rscData bgpAttrData) configOptsToDel(delPrefix string) []string {
	return []string{
		delPrefix + "accept-remote-nexthop",
		delPrefix + "advertise-external",
		delPrefix + "advertise-inactive",
		delPrefix + "advertise-peer-as",
		delPrefix + "no-advertise-peer-as",
		delPrefix + "as-override",
		delPrefix + "authentication-algorithm",
		delPrefix + "authentication-key",
		delPrefix + "authentication-key-chain",
		delPrefix + "cluster",
		delPrefix + "damping",
		delPrefix + "description",
		delPrefix + "export",
		delPrefix + "hold-time",
		delPrefix + "import",
		delPrefix + "keep",
		delPrefix + "local-address",
		delPrefix + "local-as",
		delPrefix + "local-interface",
		delPrefix + "local-preference",
		delPrefix + "log-updown",
		delPrefix + "metric-out",
		delPrefix + "mtu-discovery",
		delPrefix + "multihop",
		delPrefix + "multipath",
		delPrefix + "no-client-reflect",
		delPrefix + "out-delay",
		delPrefix + "passive",
		delPrefix + "peer-as",
		delPrefix + "preference",
		delPrefix + "remove-private",
		delPrefix + "tcp-aggressive-transmission",
		delPrefix + "bfd-liveness-detection",
		delPrefix + "bgp-error-tolerance",
		delPrefix + "family evpn",
		delPrefix + "family inet",
		delPrefix + "family inet6",
		delPrefix + "graceful-restart",
	}
}

type bgpAttrConfig struct {
	AcceptRemoteDesktop          types.Bool                    `tfsdk:"accept_remote_nexthop"`
	AdvertiseExternal            types.Bool                    `tfsdk:"advertise_external"`
	AdvertiseExternalConditional types.Bool                    `tfsdk:"advertise_external_conditional"`
	AdvertiseInactive            types.Bool                    `tfsdk:"advertise_inactive"`
	AdvertisePeerAS              types.Bool                    `tfsdk:"advertise_peer_as"`
	NoAdvertisePeerAS            types.Bool                    `tfsdk:"no_advertise_peer_as"`
	ASOverride                   types.Bool                    `tfsdk:"as_override"`
	AuthenticationAlgorithm      types.String                  `tfsdk:"authentication_algorithm"`
	AuthenticationKey            types.String                  `tfsdk:"authentication_key"`
	AuthenticationKeyChain       types.String                  `tfsdk:"authentication_key_chain"`
	Cluster                      types.String                  `tfsdk:"cluster"`
	Damping                      types.Bool                    `tfsdk:"damping"`
	Description                  types.String                  `tfsdk:"description"`
	Export                       types.List                    `tfsdk:"export"`
	HoldTime                     types.Int64                   `tfsdk:"hold_time"`
	Import                       types.List                    `tfsdk:"import"`
	KeepAll                      types.Bool                    `tfsdk:"keep_all"`
	KeepNone                     types.Bool                    `tfsdk:"keep_none"`
	LocalAddress                 types.String                  `tfsdk:"local_address"`
	LocalAS                      types.String                  `tfsdk:"local_as"`
	LocalASAlias                 types.Bool                    `tfsdk:"local_as_alias"`
	LocalASLoops                 types.Int64                   `tfsdk:"local_as_loops"`
	LocalASNoPrependGlobalAS     types.Bool                    `tfsdk:"local_as_no_prepend_global_as"`
	LocalASPrivate               types.Bool                    `tfsdk:"local_as_private"`
	LocalInterface               types.String                  `tfsdk:"local_interface"`
	LocalPreference              types.Int64                   `tfsdk:"local_preference"`
	LogUpdown                    types.Bool                    `tfsdk:"log_updown"`
	MetricOut                    types.Int64                   `tfsdk:"metric_out"`
	MetricOutIgp                 types.Bool                    `tfsdk:"metric_out_igp"`
	MetricOutIgpDelayMedUpdate   types.Bool                    `tfsdk:"metric_out_igp_delay_med_update"`
	MetricOutIgpOffset           types.Int64                   `tfsdk:"metric_out_igp_offset"`
	MetricOutMinimumIgp          types.Bool                    `tfsdk:"metric_out_minimum_igp"`
	MetricOutMinimumIgpOffset    types.Int64                   `tfsdk:"metric_out_minimum_igp_offset"`
	MtuDiscovery                 types.Bool                    `tfsdk:"mtu_discovery"`
	Multihop                     types.Bool                    `tfsdk:"multihop"`
	NoClientReflect              types.Bool                    `tfsdk:"no_client_reflect"`
	OutDelay                     types.Int64                   `tfsdk:"out_delay"`
	Passive                      types.Bool                    `tfsdk:"passive"`
	PeerAS                       types.String                  `tfsdk:"peer_as"`
	Preference                   types.Int64                   `tfsdk:"preference"`
	RemotePrivate                types.Bool                    `tfsdk:"remove_private"`
	TCPAggressiveTransmission    types.Bool                    `tfsdk:"tcp_aggressive_transmission"`
	BfdLivenessDetection         *bgpBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
	BgpErrorTolerance            *bgpBlockBgpErrorTolerance    `tfsdk:"bgp_error_tolerance"`
	BgpMultipah                  *bgpBlockBgpMultipath         `tfsdk:"bgp_multipath"`
	FamilyEvpn                   types.List                    `tfsdk:"family_evpn"`
	FamilyInet                   types.List                    `tfsdk:"family_inet"`
	FamilyInet6                  types.List                    `tfsdk:"family_inet6"`
	GracefulRestart              *bgpBlockGracefulRestart      `tfsdk:"graceful_restart"`
}

//nolint:gocognit,gocyclo
func (config *bgpAttrConfig) validateConfig(
	ctx context.Context, resp *resource.ValidateConfigResponse,
) {
	if !config.AdvertisePeerAS.IsNull() && !config.AdvertisePeerAS.IsUnknown() &&
		!config.NoAdvertisePeerAS.IsNull() && !config.NoAdvertisePeerAS.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("advertise_peer_as"),
			tfdiag.ConflictConfigErrSummary,
			"advertise_peer_as and no_advertise_peer_as can't be true in same time",
		)
	}
	if !config.KeepAll.IsNull() && !config.KeepAll.IsUnknown() &&
		!config.KeepNone.IsNull() && !config.KeepNone.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("keep_all"),
			tfdiag.ConflictConfigErrSummary,
			"keep_all and keep_none can't be true in same time",
		)
	}
	if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() {
		if !config.AuthenticationAlgorithm.IsNull() && !config.AuthenticationAlgorithm.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_algorithm"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_algorithm and authentication_key cannot be configured together",
			)
		}
		if !config.AuthenticationKeyChain.IsNull() && !config.AuthenticationKeyChain.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key_chain"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key_chain and authentication_key cannot be configured together",
			)
		}
	}
	if !config.LocalASAlias.IsNull() && !config.LocalASAlias.IsUnknown() {
		if !config.LocalASPrivate.IsNull() && !config.LocalASPrivate.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("local_as_alias"),
				tfdiag.ConflictConfigErrSummary,
				"local_as_alias and local_as_private cannot be configured together",
			)
		}
		if !config.LocalASNoPrependGlobalAS.IsNull() && !config.LocalASNoPrependGlobalAS.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("local_as_alias"),
				tfdiag.ConflictConfigErrSummary,
				"local_as_alias and local_as_no_prepend_global_as cannot be configured together",
			)
		}
	}
	if !config.LocalASPrivate.IsNull() && !config.LocalASPrivate.IsUnknown() {
		if !config.LocalASNoPrependGlobalAS.IsNull() && !config.LocalASNoPrependGlobalAS.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("local_as_private"),
				tfdiag.ConflictConfigErrSummary,
				"local_as_private and local_as_no_prepend_global_as cannot be configured together",
			)
		}
	}
	if !config.MetricOut.IsNull() && !config.MetricOut.IsUnknown() {
		if !config.MetricOutIgp.IsNull() && !config.MetricOutIgp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_igp"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out and metric_out_igp cannot be configured together",
			)
		}
		if !config.MetricOutIgpDelayMedUpdate.IsNull() && !config.MetricOutIgpDelayMedUpdate.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_igp_delay_med_update"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out and metric_out_igp_delay_med_update cannot be configured together",
			)
		}
		if !config.MetricOutIgpOffset.IsNull() && !config.MetricOutIgpOffset.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_igp_offset"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out and metric_out_igp_offset cannot be configured together",
			)
		}
		if !config.MetricOutMinimumIgp.IsNull() && !config.MetricOutMinimumIgp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out and metric_out_minimum_igp cannot be configured together",
			)
		}
		if !config.MetricOutMinimumIgpOffset.IsNull() && !config.MetricOutMinimumIgpOffset.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp_offset"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out and metric_out_minimum_igp_offset cannot be configured together",
			)
		}
	}
	if !config.MetricOutIgp.IsNull() && !config.MetricOutIgp.IsUnknown() {
		if !config.MetricOutMinimumIgp.IsNull() && !config.MetricOutMinimumIgp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp and metric_out_minimum_igp cannot be configured together",
			)
		}
		if !config.MetricOutMinimumIgpOffset.IsNull() && !config.MetricOutMinimumIgpOffset.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp_offset"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp and metric_out_minimum_igp_offset cannot be configured together",
			)
		}
	}
	if !config.MetricOutIgpDelayMedUpdate.IsNull() && !config.MetricOutIgpDelayMedUpdate.IsUnknown() {
		if !config.MetricOutMinimumIgp.IsNull() && !config.MetricOutMinimumIgp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp_delay_med_update and metric_out_minimum_igp cannot be configured together",
			)
		}
		if !config.MetricOutMinimumIgpOffset.IsNull() && !config.MetricOutMinimumIgpOffset.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp_offset"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp_delay_med_update and metric_out_minimum_igp_offset cannot be configured together",
			)
		}
	}
	if !config.MetricOutIgpOffset.IsNull() && !config.MetricOutIgpOffset.IsUnknown() {
		if !config.MetricOutMinimumIgp.IsNull() && !config.MetricOutMinimumIgp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp_offset and metric_out_minimum_igp cannot be configured together",
			)
		}
		if !config.MetricOutMinimumIgpOffset.IsNull() && !config.MetricOutMinimumIgpOffset.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("metric_out_minimum_igp_offset"),
				tfdiag.ConflictConfigErrSummary,
				"metric_out_igp_offset and metric_out_minimum_igp_offset cannot be configured together",
			)
		}
	}
	if config.BfdLivenessDetection != nil {
		if config.BfdLivenessDetection.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("bfd_liveness_detection").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"bfd_liveness_detection block is empty",
			)
		}
	}
	if config.BgpErrorTolerance != nil {
		if !config.BgpErrorTolerance.MalformedRouteLimit.IsNull() &&
			!config.BgpErrorTolerance.MalformedRouteLimit.IsUnknown() &&
			!config.BgpErrorTolerance.NoMalformedRouteLimit.IsNull() &&
			!config.BgpErrorTolerance.NoMalformedRouteLimit.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("bgp_error_tolerance").AtName("no_malformed_route_limit"),
				tfdiag.ConflictConfigErrSummary,
				"malformed_route_limit and no_malformed_route_limit cannot be configured together"+
					" in bgp_error_tolerance block",
			)
		}
	}
	if !config.FamilyEvpn.IsNull() && !config.FamilyEvpn.IsUnknown() {
		var configFamilyEvpn []bgpBlockFamily
		asDiags := config.FamilyEvpn.ElementsAs(ctx, &configFamilyEvpn, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		familyEvpnNlriType := make(map[string]struct{})
		for i, block := range configFamilyEvpn {
			if !block.NlriType.IsUnknown() {
				nlriType := block.NlriType.ValueString()
				if _, ok := familyEvpnNlriType[nlriType]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_evpn").AtListIndex(i).AtName("nlri_type"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple family_evpn blocks with the same nlri_type %q", nlriType),
					)
				}
				familyEvpnNlriType[nlriType] = struct{}{}
			}
			if block.AcceptedPrefixLimit != nil {
				if block.AcceptedPrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_evpn").AtListIndex(i).AtName("accepted_prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in accepted_prefix_limit block in family_evpn block",
					)
				}
				if !block.AcceptedPrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_evpn").AtListIndex(i).AtName("accepted_prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in accepted_prefix_limit block in family_evpn block",
					)
				}
			}
			if block.PrefixLimit != nil {
				if block.PrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_evpn").AtListIndex(i).AtName("prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in prefix_limit block in family_evpn block",
					)
				}
				if !block.PrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_evpn").AtListIndex(i).AtName("prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in prefix_limit block family_evpn block",
					)
				}
			}
		}
	}
	if !config.FamilyInet.IsNull() && !config.FamilyInet.IsUnknown() {
		var configFamilyInet []bgpBlockFamily
		asDiags := config.FamilyInet.ElementsAs(ctx, &configFamilyInet, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		familyInetNlriType := make(map[string]struct{})
		for i, block := range configFamilyInet {
			if !block.NlriType.IsUnknown() {
				nlriType := block.NlriType.ValueString()
				if _, ok := familyInetNlriType[nlriType]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet").AtListIndex(i).AtName("nlri_type"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple family_inet blocks with the same nlri_type %q", nlriType),
					)
				}
				familyInetNlriType[nlriType] = struct{}{}
			}
			if block.AcceptedPrefixLimit != nil {
				if block.AcceptedPrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet").AtListIndex(i).AtName("accepted_prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in accepted_prefix_limit block in family_inet block",
					)
				}
				if !block.AcceptedPrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet").AtListIndex(i).AtName("accepted_prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in accepted_prefix_limit block in family_inet block",
					)
				}
			}
			if block.PrefixLimit != nil {
				if block.PrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet").AtListIndex(i).AtName("prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in prefix_limit block in family_inet block",
					)
				}
				if !block.PrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet").AtListIndex(i).AtName("prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in prefix_limit block family_inet block",
					)
				}
			}
		}
	}
	if !config.FamilyInet6.IsNull() && !config.FamilyInet6.IsUnknown() {
		var configFamilyInet6 []bgpBlockFamily
		asDiags := config.FamilyInet6.ElementsAs(ctx, &configFamilyInet6, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		familyInet6NlriType := make(map[string]struct{})
		for i, block := range configFamilyInet6 {
			if !block.NlriType.IsUnknown() {
				nlriType := block.NlriType.ValueString()
				if _, ok := familyInet6NlriType[nlriType]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6").AtListIndex(i).AtName("nlri_type"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple family_inet6 blocks with the same nlri_type %q", nlriType),
					)
				}
				familyInet6NlriType[nlriType] = struct{}{}
			}
			if block.AcceptedPrefixLimit != nil {
				if block.AcceptedPrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6").AtListIndex(i).AtName("accepted_prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in accepted_prefix_limit block in family_inet6 block",
					)
				}
				if !block.AcceptedPrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.AcceptedPrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6").AtListIndex(i).AtName("accepted_prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in accepted_prefix_limit block in family_inet6 block",
					)
				}
			}
			if block.PrefixLimit != nil {
				if block.PrefixLimit.Maximum.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6").AtListIndex(i).AtName("prefix_limit").AtName("maximum"),
						tfdiag.MissingConfigErrSummary,
						"maximum must be specified in prefix_limit block in family_inet6 block",
					)
				}
				if !block.PrefixLimit.TeardownIdleTimeout.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeout.IsUnknown() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsNull() &&
					!block.PrefixLimit.TeardownIdleTimeoutForever.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6").AtListIndex(i).AtName("prefix_limit").AtName("teardown_idle_timeout"),
						tfdiag.ConflictConfigErrSummary,
						"teardown_idle_timeout and teardown_idle_timeout_forever cannot be configured together"+
							" in prefix_limit block family_inet6 block",
					)
				}
			}
		}
	}
	if config.GracefulRestart != nil {
		if !config.GracefulRestart.Disable.IsNull() && !config.GracefulRestart.Disable.IsUnknown() {
			if !config.GracefulRestart.RestartTime.IsNull() && !config.GracefulRestart.RestartTime.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("graceful_restart").AtName("restart_time"),
					tfdiag.ConflictConfigErrSummary,
					"restart_time and disable cannot be configured together"+
						" in graceful_restart block",
				)
			}
			if !config.GracefulRestart.StaleRouteTime.IsNull() && !config.GracefulRestart.StaleRouteTime.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("graceful_restart").AtName("stale_route_time"),
					tfdiag.ConflictConfigErrSummary,
					"stale_route_time and disable cannot be configured together"+
						" in graceful_restart block",
				)
			}
		}
	}
}

func (config *bgpAttrConfig) modifyPlan(
	_ context.Context, plan *bgpAttrConfig,
) {
	if config.AdvertiseExternal.IsNull() {
		if config.AdvertiseExternalConditional.IsNull() {
			plan.AdvertiseExternal = types.BoolNull()
		} else if !plan.AdvertiseExternalConditional.IsNull() &&
			!plan.AdvertiseExternalConditional.IsUnknown() {
			plan.AdvertiseExternal = types.BoolValue(true)
		}
	}
	if config.MetricOutIgp.IsNull() {
		if config.MetricOutIgpDelayMedUpdate.IsNull() &&
			config.MetricOutIgpOffset.IsNull() {
			plan.MetricOutIgp = types.BoolNull()
		} else {
			if !plan.MetricOutIgpDelayMedUpdate.IsNull() &&
				!plan.MetricOutIgpDelayMedUpdate.IsUnknown() {
				plan.MetricOutIgp = types.BoolValue(true)
			}
			if !plan.MetricOutIgpOffset.IsNull() &&
				!plan.MetricOutIgpOffset.IsUnknown() {
				plan.MetricOutIgp = types.BoolValue(true)
			}
		}
	}
	if config.MetricOutMinimumIgp.IsNull() {
		if config.MetricOutMinimumIgpOffset.IsNull() {
			plan.MetricOutMinimumIgp = types.BoolNull()
		} else if !plan.MetricOutMinimumIgpOffset.IsNull() &&
			!plan.MetricOutMinimumIgpOffset.IsUnknown() {
			plan.MetricOutMinimumIgp = types.BoolValue(true)
		}
	}
}
