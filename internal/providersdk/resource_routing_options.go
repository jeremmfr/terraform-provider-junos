package providersdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type routingOptionsOptions struct {
	routerID         string
	instanceExport   []string
	instanceImport   []string
	autonomousSystem []map[string]interface{}
	forwardingTable  []map[string]interface{}
	gracefulRestart  []map[string]interface{}
}

func resourceRoutingOptions() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRoutingOptionsCreate,
		ReadWithoutTimeout:   resourceRoutingOptionsRead,
		UpdateWithoutTimeout: resourceRoutingOptionsUpdate,
		DeleteWithoutTimeout: resourceRoutingOptionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoutingOptionsImport,
		},
		Schema: map[string]*schema.Schema{
			"clean_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"autonomous_system": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:     schema.TypeString,
							Required: true,
						},
						"asdot_notation": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"loops": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},
			"forwarding_table": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chain_composite_max_label_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 8),
						},
						"chained_composite_next_hop_ingress": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"chained_composite_next_hop_transit": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dynamic_list_next_hop": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ecmp_fast_reroute": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.no_ecmp_fast_reroute"},
						},
						"no_ecmp_fast_reroute": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.ecmp_fast_reroute"},
						},
						"export": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
							},
						},
						"indirect_next_hop": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.no_indirect_next_hop"},
						},
						"no_indirect_next_hop": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.indirect_next_hop"},
						},
						"indirect_next_hop_change_acknowledgements": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.no_indirect_next_hop_change_acknowledgements"},
						},
						"no_indirect_next_hop_change_acknowledgements": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"forwarding_table.0.indirect_next_hop_change_acknowledgements"},
						},
						"krt_nexthop_ack_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 400),
						},
						"remnant_holdtime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 10000),
						},
						"unicast_reverse_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"active-paths", "feasible-paths"}, false),
						},
					},
				},
			},
			"forwarding_table_export_configure_singly": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"forwarding_table.0.export"},
			},
			"graceful_restart": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"restart_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(120, 10000),
						},
					},
				},
			},
			"instance_export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"instance_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"router_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceRoutingOptionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setRoutingOptions(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("routing_options")

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setRoutingOptions(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_routing_options")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("routing_options")

	return append(diagWarns, resourceRoutingOptionsReadWJunSess(d, junSess)...)
}

func resourceRoutingOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceRoutingOptionsReadWJunSess(d, junSess)
}

func resourceRoutingOptionsReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	routingOptionsOptions, err := readRoutingOptions(junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillRoutingOptions(d, routingOptionsOptions)

	return nil
}

func resourceRoutingOptionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	var diagWarns diag.Diagnostics
	fwTableExportConfigSingly := d.Get("forwarding_table_export_configure_singly").(bool)
	if d.HasChange("forwarding_table_export_configure_singly") {
		if o, _ := d.GetChange("forwarding_table_export_configure_singly"); o.(bool) {
			fwTableExportConfigSingly = o.(bool)
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Disable forwarding_table_export_configure_singly on resource already created doesn't " +
					"delete export list already configured.",
				Detail:        "So refresh resource after apply to detect export list entries that need to be deleted",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "forwarding_table_export_configure_singly"}},
			})
		} else {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Enable forwarding_table_export_configure_singly on resource already created doesn't " +
					"delete export list already configured.",
				Detail: "So add `add_it_to_forwarding_table_export` argument on each `junos_policyoptions_policy_statement` " +
					"resource to be able to manage each element of the export list",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "forwarding_table_export_configure_singly"}},
			})
		}
	}
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRoutingOptions(fwTableExportConfigSingly, junSess); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if err := setRoutingOptions(d, junSess); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		d.Partial(false)

		return diagWarns
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	if err := delRoutingOptions(fwTableExportConfigSingly, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRoutingOptions(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_routing_options")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRoutingOptionsReadWJunSess(d, junSess)...)
}

func resourceRoutingOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		clt := m.(*junos.Client)
		if clt.FakeDeleteAlso() {
			junSess := clt.NewSessionWithoutNetconf(ctx)
			if err := delRoutingOptions(d.Get("forwarding_table_export_configure_singly").(bool), junSess); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		junSess, err := clt.StartNewSession(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		defer junSess.Close()
		if err := junSess.ConfigLock(ctx); err != nil {
			return diag.FromErr(err)
		}
		var diagWarns diag.Diagnostics
		if err := delRoutingOptions(d.Get("forwarding_table_export_configure_singly").(bool), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := junSess.CommitConf(ctx, "delete resource junos_routing_options")
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceRoutingOptionsImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	routingOptionsOptions, err := readRoutingOptions(junSess)
	if err != nil {
		return nil, err
	}
	fillRoutingOptions(d, routingOptionsOptions)
	d.SetId("routing_options")
	result[0] = d

	return result, nil
}

func setRoutingOptions(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set routing-options "
	configSet := make([]string, 0)

	for _, as := range d.Get("autonomous_system").([]interface{}) {
		asM := as.(map[string]interface{})
		configSet = append(configSet, setPrefix+"autonomous-system "+asM["number"].(string))
		if asM["asdot_notation"].(bool) {
			configSet = append(configSet, setPrefix+"autonomous-system asdot-notation")
		}
		if asM["loops"].(int) > 0 {
			configSet = append(configSet, setPrefix+"autonomous-system loops "+strconv.Itoa(asM["loops"].(int)))
		}
	}
	for _, vFwTable := range d.Get("forwarding_table").([]interface{}) {
		fwTable := vFwTable.(map[string]interface{})
		if v := fwTable["chain_composite_max_label_count"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"forwarding-table chain-composite-max-label-count "+strconv.Itoa(v))
		}
		for _, v := range sortSetOfString(fwTable["chained_composite_next_hop_ingress"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"forwarding-table chained-composite-next-hop ingress "+v)
		}
		for _, v := range sortSetOfString(fwTable["chained_composite_next_hop_transit"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"forwarding-table chained-composite-next-hop transit "+v)
		}
		if fwTable["dynamic_list_next_hop"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table dynamic-list-next-hop")
		}
		if fwTable["ecmp_fast_reroute"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table ecmp-fast-reroute")
		}
		if fwTable["no_ecmp_fast_reroute"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table no-ecmp-fast-reroute")
		}
		for _, v := range fwTable["export"].([]interface{}) {
			configSet = append(configSet, setPrefix+"forwarding-table export \""+v.(string)+"\"")
		}
		if fwTable["indirect_next_hop"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table indirect-next-hop")
		}
		if fwTable["no_indirect_next_hop"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table no-indirect-next-hop")
		}
		if fwTable["indirect_next_hop_change_acknowledgements"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table indirect-next-hop-change-acknowledgements")
		}
		if fwTable["no_indirect_next_hop_change_acknowledgements"].(bool) {
			configSet = append(configSet, setPrefix+"forwarding-table no-indirect-next-hop-change-acknowledgements")
		}
		if v := fwTable["krt_nexthop_ack_timeout"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"forwarding-table krt-nexthop-ack-timeout "+strconv.Itoa(v))
		}
		if v := fwTable["remnant_holdtime"].(int); v != -1 {
			configSet = append(configSet, setPrefix+"forwarding-table remnant-holdtime "+strconv.Itoa(v))
		}
		if v := fwTable["unicast_reverse_path"].(string); v != "" {
			configSet = append(configSet, setPrefix+"forwarding-table unicast-reverse-path "+v)
		}

		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"forwarding-table ") {
			return errors.New("forwarding_table block is empty")
		}
	}
	for _, grR := range d.Get("graceful_restart").([]interface{}) {
		configSet = append(configSet, setPrefix+"graceful-restart")
		if grR != nil {
			grRM := grR.(map[string]interface{})
			if grRM["disable"].(bool) {
				configSet = append(configSet, setPrefix+"graceful-restart disable")
			}
			if grRM["restart_duration"].(int) > 0 {
				configSet = append(configSet, setPrefix+"graceful-restart restart-duration "+
					strconv.Itoa(grRM["restart_duration"].(int)))
			}
		}
	}
	for _, v := range d.Get("instance_export").([]interface{}) {
		configSet = append(configSet, setPrefix+"instance-export "+v.(string))
	}
	for _, v := range d.Get("instance_import").([]interface{}) {
		configSet = append(configSet, setPrefix+"instance-import "+v.(string))
	}
	if v := d.Get("router_id").(string); v != "" {
		configSet = append(configSet, setPrefix+"router-id "+v)
	}

	return junSess.ConfigSet(configSet)
}

func delRoutingOptions(fwTableExportConfigSingly bool, junSess *junos.Session) error {
	listLinesToDelete := []string{
		"autonomous-system",
		"graceful-restart",
		"instance-export",
		"instance-import",
		"router-id",
	}
	if fwTableExportConfigSingly {
		listLinesToDeleteFwTable := []string{
			"forwarding-table chain-composite-max-label-count",
			"forwarding-table chained-composite-next-hop",
			"forwarding-table dynamic-list-next-hop",
			"forwarding-table ecmp-fast-reroute",
			"forwarding-table no-ecmp-fast-reroute",
			"forwarding-table indirect-next-hop",
			"forwarding-table no-indirect-next-hop",
			"forwarding-table indirect-next-hop-change-acknowledgements",
			"forwarding-table no-indirect-next-hop-change-acknowledgements",
			"forwarding-table krt-nexthop-ack-timeout",
			"forwarding-table remnant-holdtime",
			"forwarding-table unicast-reverse-path",
		}
		listLinesToDelete = append(listLinesToDelete, listLinesToDeleteFwTable...)
	} else {
		listLinesToDelete = append(listLinesToDelete, "forwarding-table")
	}
	configSet := make([]string, 0)
	delPrefix := "delete routing-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return junSess.ConfigSet(configSet)
}

func readRoutingOptions(junSess *junos.Session,
) (confRead routingOptionsOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"routing-options" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
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
			switch {
			case balt.CutPrefixInString(&itemTrim, "autonomous-system "):
				if len(confRead.autonomousSystem) == 0 {
					confRead.autonomousSystem = append(confRead.autonomousSystem, map[string]interface{}{
						"number":         "",
						"asdot_notation": false,
						"loops":          0,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "loops "):
					confRead.autonomousSystem[0]["loops"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "asdot-notation":
					confRead.autonomousSystem[0]["asdot_notation"] = true
				default:
					confRead.autonomousSystem[0]["number"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "forwarding-table "):
				if len(confRead.forwardingTable) == 0 {
					confRead.forwardingTable = append(confRead.forwardingTable, map[string]interface{}{
						"chain_composite_max_label_count":              0,
						"chained_composite_next_hop_ingress":           make([]string, 0),
						"chained_composite_next_hop_transit":           make([]string, 0),
						"dynamic_list_next_hop":                        false,
						"ecmp_fast_reroute":                            false,
						"no_ecmp_fast_reroute":                         false,
						"export":                                       make([]string, 0),
						"indirect_next_hop":                            false,
						"no_indirect_next_hop":                         false,
						"indirect_next_hop_change_acknowledgements":    false,
						"no_indirect_next_hop_change_acknowledgements": false,
						"krt_nexthop_ack_timeout":                      0,
						"remnant_holdtime":                             -1,
						"unicast_reverse_path":                         "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "chain-composite-max-label-count "):
					confRead.forwardingTable[0]["chain_composite_max_label_count"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "chained-composite-next-hop ingress "):
					confRead.forwardingTable[0]["chained_composite_next_hop_ingress"] = append(
						confRead.forwardingTable[0]["chained_composite_next_hop_ingress"].([]string),
						itemTrim,
					)
				case balt.CutPrefixInString(&itemTrim, "chained-composite-next-hop transit "):
					confRead.forwardingTable[0]["chained_composite_next_hop_transit"] = append(
						confRead.forwardingTable[0]["chained_composite_next_hop_transit"].([]string),
						itemTrim,
					)
				case itemTrim == "dynamic-list-next-hop":
					confRead.forwardingTable[0]["dynamic_list_next_hop"] = true
				case itemTrim == "ecmp-fast-reroute":
					confRead.forwardingTable[0]["ecmp_fast_reroute"] = true
				case itemTrim == "no-ecmp-fast-reroute":
					confRead.forwardingTable[0]["no_ecmp_fast_reroute"] = true
				case balt.CutPrefixInString(&itemTrim, "export "):
					confRead.forwardingTable[0]["export"] = append(
						confRead.forwardingTable[0]["export"].([]string),
						strings.Trim(itemTrim, "\""),
					)
				case itemTrim == "indirect-next-hop":
					confRead.forwardingTable[0]["indirect_next_hop"] = true
				case itemTrim == "no-indirect-next-hop":
					confRead.forwardingTable[0]["no_indirect_next_hop"] = true
				case itemTrim == "indirect-next-hop-change-acknowledgements":
					confRead.forwardingTable[0]["indirect_next_hop_change_acknowledgements"] = true
				case itemTrim == "no-indirect-next-hop-change-acknowledgements":
					confRead.forwardingTable[0]["no_indirect_next_hop_change_acknowledgements"] = true
				case balt.CutPrefixInString(&itemTrim, "krt-nexthop-ack-timeout "):
					confRead.forwardingTable[0]["krt_nexthop_ack_timeout"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "remnant-holdtime "):
					confRead.forwardingTable[0]["remnant_holdtime"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "unicast-reverse-path "):
					confRead.forwardingTable[0]["unicast_reverse_path"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "graceful-restart"):
				if len(confRead.gracefulRestart) == 0 {
					confRead.gracefulRestart = append(confRead.gracefulRestart, map[string]interface{}{
						"disable":          false,
						"restart_duration": 0,
					})
				}
				switch {
				case itemTrim == " disable":
					confRead.gracefulRestart[0]["disable"] = true
				case balt.CutPrefixInString(&itemTrim, " restart-duration "):
					confRead.gracefulRestart[0]["restart_duration"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "instance-export "):
				confRead.instanceExport = append(confRead.instanceExport, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "instance-import "):
				confRead.instanceImport = append(confRead.instanceImport, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "router-id "):
				confRead.routerID = itemTrim
			}
		}
	}

	return confRead, nil
}

func fillRoutingOptions(d *schema.ResourceData, routingOptionsOptions routingOptionsOptions) {
	if tfErr := d.Set("autonomous_system", routingOptionsOptions.autonomousSystem); tfErr != nil {
		panic(tfErr)
	}
	if d.Get("forwarding_table_export_configure_singly").(bool) && len(routingOptionsOptions.forwardingTable) > 0 {
		forwardingTable := routingOptionsOptions.forwardingTable
		forwardingTable[0]["export"] = make([]string, 0)
		if tfErr := d.Set("forwarding_table", forwardingTable); tfErr != nil {
			panic(tfErr)
		}
	} else if tfErr := d.Set("forwarding_table", routingOptionsOptions.forwardingTable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_restart", routingOptionsOptions.gracefulRestart); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("instance_export", routingOptionsOptions.instanceExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("instance_import", routingOptionsOptions.instanceImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("router_id", routingOptionsOptions.routerID); tfErr != nil {
		panic(tfErr)
	}
}
