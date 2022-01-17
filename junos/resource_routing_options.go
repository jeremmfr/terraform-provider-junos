package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		CreateContext: resourceRoutingOptionsCreate,
		ReadContext:   resourceRoutingOptionsRead,
		UpdateContext: resourceRoutingOptionsUpdate,
		DeleteContext: resourceRoutingOptionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoutingOptionsImport,
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
							Elem:     &schema.Schema{Type: schema.TypeString},
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setRoutingOptions(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("routing_options")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setRoutingOptions(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_routing_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("routing_options")

	return append(diagWarns, resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)...)
}

func resourceRoutingOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)
}

func resourceRoutingOptionsReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	routingOptionsOptions, err := readRoutingOptions(m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delRoutingOptions(fwTableExportConfigSingly, m, nil); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if err := setRoutingOptions(d, m, nil); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		d.Partial(false)

		return diagWarns
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delRoutingOptions(fwTableExportConfigSingly, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setRoutingOptions(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_routing_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRoutingOptionsReadWJnprSess(d, m, jnprSess)...)
}

func resourceRoutingOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		sess := m.(*Session)
		if sess.junosFakeDeleteAlso {
			if err := delRoutingOptions(d.Get("forwarding_table_export_configure_singly").(bool), m, nil); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		jnprSess, err := sess.startNewSession()
		if err != nil {
			return diag.FromErr(err)
		}
		defer sess.closeSession(jnprSess)
		sess.configLock(jnprSess)
		var diagWarns diag.Diagnostics
		if err := delRoutingOptions(d.Get("forwarding_table_export_configure_singly").(bool), m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := sess.commitConf("delete resource junos_routing_options", jnprSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceRoutingOptionsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	routingOptionsOptions, err := readRoutingOptions(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRoutingOptions(d, routingOptionsOptions)
	d.SetId("routing_options")
	result[0] = d

	return result, nil
}

func setRoutingOptions(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

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
			return fmt.Errorf("forwarding_table block is empty")
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

	return sess.configSet(configSet, jnprSess)
}

func delRoutingOptions(fwTableExportConfigSingly bool, m interface{}, jnprSess *NetconfObject) error {
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
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete routing-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func readRoutingOptions(m interface{}, jnprSess *NetconfObject) (routingOptionsOptions, error) {
	sess := m.(*Session)
	var confRead routingOptionsOptions

	showConfig, err := sess.command("show configuration"+
		" routing-options"+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "autonomous-system "):
				if len(confRead.autonomousSystem) == 0 {
					confRead.autonomousSystem = append(confRead.autonomousSystem, map[string]interface{}{
						"number":         "",
						"asdot_notation": false,
						"loops":          0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "autonomous-system loops "):
					var err error
					confRead.autonomousSystem[0]["loops"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "autonomous-system loops "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case itemTrim == "autonomous-system asdot-notation":
					confRead.autonomousSystem[0]["asdot_notation"] = true
				default:
					confRead.autonomousSystem[0]["number"] = strings.TrimPrefix(itemTrim, "autonomous-system ")
				}
			case strings.HasPrefix(itemTrim, "forwarding-table "):
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
				case strings.HasPrefix(itemTrim, "forwarding-table chain-composite-max-label-count "):
					var err error
					confRead.forwardingTable[0]["chain_composite_max_label_count"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "forwarding-table chain-composite-max-label-count "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "forwarding-table chained-composite-next-hop ingress "):
					confRead.forwardingTable[0]["chained_composite_next_hop_ingress"] = append(
						confRead.forwardingTable[0]["chained_composite_next_hop_ingress"].([]string),
						strings.TrimPrefix(itemTrim, "forwarding-table chained-composite-next-hop ingress "))
				case strings.HasPrefix(itemTrim, "forwarding-table chained-composite-next-hop transit "):
					confRead.forwardingTable[0]["chained_composite_next_hop_transit"] = append(
						confRead.forwardingTable[0]["chained_composite_next_hop_transit"].([]string),
						strings.TrimPrefix(itemTrim, "forwarding-table chained-composite-next-hop transit "))
				case itemTrim == "forwarding-table dynamic-list-next-hop":
					confRead.forwardingTable[0]["dynamic_list_next_hop"] = true
				case itemTrim == "forwarding-table ecmp-fast-reroute":
					confRead.forwardingTable[0]["ecmp_fast_reroute"] = true
				case itemTrim == "forwarding-table no-ecmp-fast-reroute":
					confRead.forwardingTable[0]["no_ecmp_fast_reroute"] = true
				case strings.HasPrefix(itemTrim, "forwarding-table export "):
					confRead.forwardingTable[0]["export"] = append(confRead.forwardingTable[0]["export"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "forwarding-table export "), "\""))
				case itemTrim == "forwarding-table indirect-next-hop":
					confRead.forwardingTable[0]["indirect_next_hop"] = true
				case itemTrim == "forwarding-table no-indirect-next-hop":
					confRead.forwardingTable[0]["no_indirect_next_hop"] = true
				case itemTrim == "forwarding-table indirect-next-hop-change-acknowledgements":
					confRead.forwardingTable[0]["indirect_next_hop_change_acknowledgements"] = true
				case itemTrim == "forwarding-table no-indirect-next-hop-change-acknowledgements":
					confRead.forwardingTable[0]["no_indirect_next_hop_change_acknowledgements"] = true
				case strings.HasPrefix(itemTrim, "forwarding-table krt-nexthop-ack-timeout "):
					var err error
					confRead.forwardingTable[0]["krt_nexthop_ack_timeout"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "forwarding-table krt-nexthop-ack-timeout "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "forwarding-table remnant-holdtime "):
					var err error
					confRead.forwardingTable[0]["remnant_holdtime"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "forwarding-table remnant-holdtime "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "forwarding-table unicast-reverse-path "):
					confRead.forwardingTable[0]["unicast_reverse_path"] =
						strings.TrimPrefix(itemTrim, "forwarding-table unicast-reverse-path ")
				}
			case strings.HasPrefix(itemTrim, "graceful-restart"):
				if len(confRead.gracefulRestart) == 0 {
					confRead.gracefulRestart = append(confRead.gracefulRestart, map[string]interface{}{
						"disable":          false,
						"restart_duration": 0,
					})
				}
				switch {
				case itemTrim == "graceful-restart disable":
					confRead.gracefulRestart[0]["disable"] = true
				case strings.HasPrefix(itemTrim, "graceful-restart restart-duration "):
					var err error
					confRead.gracefulRestart[0]["restart_duration"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "graceful-restart restart-duration "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "instance-export "):
				confRead.instanceExport = append(confRead.instanceExport, strings.TrimPrefix(itemTrim, "instance-export "))
			case strings.HasPrefix(itemTrim, "instance-import "):
				confRead.instanceImport = append(confRead.instanceImport, strings.TrimPrefix(itemTrim, "instance-import "))
			case strings.HasPrefix(itemTrim, "router-id "):
				confRead.routerID = strings.TrimPrefix(itemTrim, "router-id ")
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
