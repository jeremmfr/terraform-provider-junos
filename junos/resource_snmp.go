package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type snmpOptions struct {
	arp                         bool
	arpHostNameResolution       bool
	filterDuplicates            bool
	filterInternalInterfaces    bool
	ifCountWithFilterInterfaces bool
	routingInstanceAccess       bool
	contact                     string
	description                 string
	engineID                    string
	location                    string
	filterInterfaces            []string
	interFace                   []string
	routingInstanceAccessList   []string
	healthMonitor               []map[string]interface{}
}

func resourceSnmp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnmpCreate,
		ReadContext:   resourceSnmpRead,
		UpdateContext: resourceSnmpUpdate,
		DeleteContext: resourceSnmpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSnmpImport,
		},
		Schema: map[string]*schema.Schema{
			"clean_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"arp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"arp_host_name_resolution": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"arp"},
			},
			"contact": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_id": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^(use-default-ip-address|use-mac-address|local .+)$`),
					"must have 'use-default-ip-address', 'use-mac-address' or 'local ...'"),
			},
			"filter_duplicates": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter_internal_interfaces": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_monitor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"falling_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"idp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"idp_falling_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"health_monitor.0.idp"},
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"idp_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"health_monitor.0.idp"},
							ValidateFunc: validation.IntBetween(1, 2147483647),
						},
						"idp_rising_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"health_monitor.0.idp"},
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2147483647),
						},
						"rising_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
			"if_count_with_filter_interfaces": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routing_instance_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routing_instance_access_list": {
				Type:         schema.TypeSet,
				Optional:     true,
				RequiredWith: []string{"routing_instance_access"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSnmpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSnmp(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("snmp")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setSnmp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_snmp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("snmp")

	return append(diagWarns, resourceSnmpReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSnmpReadWJnprSess(d, m, jnprSess)
}

func resourceSnmpReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	snmpOptions, err := readSnmp(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSnmp(d, snmpOptions)

	return nil
}

func resourceSnmpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSnmp(m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmp(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSnmp(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmp(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_snmp", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpReadWJnprSess(d, m, jnprSess)...)
}

func resourceSnmpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		sess := m.(*Session)
		if sess.junosFakeDeleteAlso {
			if err := delSnmp(m, nil); err != nil {
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
		if err := delSnmp(m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := sess.commitConf("delete resource junos_snmp", jnprSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceSnmpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	snmpOptions, err := readSnmp(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSnmp(d, snmpOptions)
	d.SetId("snmp")
	result[0] = d

	return result, nil
}

func setSnmp(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set snmp "
	configSet := make([]string, 0)

	if d.Get("arp").(bool) {
		configSet = append(configSet, setPrefix+"arp")
	}
	if d.Get("arp_host_name_resolution").(bool) {
		configSet = append(configSet, setPrefix+"arp host-name-resolution")
	}
	if v := d.Get("contact").(string); v != "" {
		configSet = append(configSet, setPrefix+"contact \""+v+"\"")
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("engine_id").(string); v != "" {
		configSet = append(configSet, setPrefix+"engine-id "+v)
	}
	if d.Get("filter_duplicates").(bool) {
		configSet = append(configSet, setPrefix+"filter-duplicates")
	}
	for _, v := range sortSetOfString(d.Get("filter_interfaces").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"filter-interfaces interfaces \""+v+"\"")
	}
	if d.Get("filter_internal_interfaces").(bool) {
		configSet = append(configSet, setPrefix+"filter-interfaces all-internal-interfaces")
	}
	for _, v := range d.Get("health_monitor").([]interface{}) {
		configSet = append(configSet, setPrefix+"health-monitor")
		if v != nil {
			hMon := v.(map[string]interface{})
			if v2 := hMon["falling_threshold"].(int); v2 != -1 {
				configSet = append(configSet, setPrefix+"health-monitor falling-threshold "+strconv.Itoa(v2))
			}
			if hMon["idp"].(bool) {
				configSet = append(configSet, setPrefix+"health-monitor idp")
			}
			if v2 := hMon["idp_falling_threshold"].(int); v2 != -1 {
				configSet = append(configSet, setPrefix+"health-monitor idp falling-threshold "+strconv.Itoa(v2))
			}
			if v2 := hMon["idp_interval"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"health-monitor idp interval "+strconv.Itoa(v2))
			}
			if v2 := hMon["idp_rising_threshold"].(int); v2 != -1 {
				configSet = append(configSet, setPrefix+"health-monitor idp rising-threshold "+strconv.Itoa(v2))
			}
			if v2 := hMon["interval"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"health-monitor interval "+strconv.Itoa(v2))
			}
			if v2 := hMon["rising_threshold"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"health-monitor rising-threshold "+strconv.Itoa(v2))
			}
		}
	}
	if d.Get("if_count_with_filter_interfaces").(bool) {
		configSet = append(configSet, setPrefix+"if-count-with-filter-interfaces")
	}
	for _, v := range sortSetOfString(d.Get("interface").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"interface "+v)
	}
	if v := d.Get("location").(string); v != "" {
		configSet = append(configSet, setPrefix+"location \""+v+"\"")
	}
	if d.Get("routing_instance_access").(bool) {
		configSet = append(configSet, setPrefix+"routing-instance-access")
	}
	for _, v := range sortSetOfString(d.Get("routing_instance_access_list").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"routing-instance-access access-list \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func delSnmp(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := []string{
		"arp",
		"contact",
		"description",
		"engine-id",
		"filter-duplicates",
		"filter-interfaces",
		"health-monitor",
		"if-count-with-filter-interfaces",
		"interface",
		"location",
		"routing-instance-access",
	}
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete snmp "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func readSnmp(m interface{}, jnprSess *NetconfObject) (snmpOptions, error) {
	sess := m.(*Session)
	var confRead snmpOptions

	showConfig, err := sess.command("show configuration snmp | display set relative", jnprSess)
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
			case itemTrim == "arp":
				confRead.arp = true
			case itemTrim == "arp host-name-resolution":
				confRead.arp = true
				confRead.arpHostNameResolution = true
			case strings.HasPrefix(itemTrim, "contact "):
				confRead.contact = strings.Trim(strings.TrimPrefix(itemTrim, "contact "), "\"")
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "engine-id "):
				confRead.engineID = strings.TrimPrefix(itemTrim, "engine-id ")
			case itemTrim == "filter-duplicates":
				confRead.filterDuplicates = true
			case strings.HasPrefix(itemTrim, "filter-interfaces interfaces "):
				confRead.filterInterfaces = append(confRead.filterInterfaces,
					strings.Trim(strings.TrimPrefix(itemTrim, "filter-interfaces interfaces "), "\""))
			case itemTrim == "filter-interfaces all-internal-interfaces":
				confRead.filterInternalInterfaces = true
			case strings.HasPrefix(itemTrim, "health-monitor"):
				if len(confRead.healthMonitor) == 0 {
					confRead.healthMonitor = append(confRead.healthMonitor, map[string]interface{}{
						"falling_threshold":     -1,
						"idp":                   false,
						"idp_falling_threshold": -1,
						"idp_interval":          0,
						"idp_rising_threshold":  -1,
						"interval":              0,
						"rising_threshold":      0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "health-monitor falling-threshold "):
					var err error
					confRead.healthMonitor[0]["falling_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "health-monitor falling-threshold "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case itemTrim == "health-monitor idp":
					confRead.healthMonitor[0]["idp"] = true
				case strings.HasPrefix(itemTrim, "health-monitor idp falling-threshold "):
					confRead.healthMonitor[0]["idp"] = true
					var err error
					confRead.healthMonitor[0]["idp_falling_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "health-monitor idp falling-threshold "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "health-monitor idp interval "):
					confRead.healthMonitor[0]["idp"] = true
					var err error
					confRead.healthMonitor[0]["idp_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "health-monitor idp interval "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "health-monitor idp rising-threshold "):
					confRead.healthMonitor[0]["idp"] = true
					var err error
					confRead.healthMonitor[0]["idp_rising_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "health-monitor idp rising-threshold "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "health-monitor interval "):
					var err error
					confRead.healthMonitor[0]["interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "health-monitor interval "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "health-monitor rising-threshold "):
					var err error
					confRead.healthMonitor[0]["rising_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "health-monitor rising-threshold "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
			case itemTrim == "if-count-with-filter-interfaces":
				confRead.ifCountWithFilterInterfaces = true
			case strings.HasPrefix(itemTrim, "interface "):
				confRead.interFace = append(confRead.interFace, strings.TrimPrefix(itemTrim, "interface "))
			case strings.HasPrefix(itemTrim, "location "):
				confRead.location = strings.Trim(strings.TrimPrefix(itemTrim, "location "), "\"")
			case itemTrim == "routing-instance-access":
				confRead.routingInstanceAccess = true
			case strings.HasPrefix(itemTrim, "routing-instance-access access-list "):
				confRead.routingInstanceAccess = true
				confRead.routingInstanceAccessList = append(confRead.routingInstanceAccessList,
					strings.Trim(strings.TrimPrefix(itemTrim, "routing-instance-access access-list "), "\""))
			}
		}
	}

	return confRead, nil
}

func fillSnmp(d *schema.ResourceData, snmpOptions snmpOptions) {
	if tfErr := d.Set("arp", snmpOptions.arp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("arp_host_name_resolution", snmpOptions.arpHostNameResolution); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("contact", snmpOptions.contact); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", snmpOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("engine_id", snmpOptions.engineID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_duplicates", snmpOptions.filterDuplicates); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_interfaces", snmpOptions.filterInterfaces); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("filter_internal_interfaces", snmpOptions.filterInternalInterfaces); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("health_monitor", snmpOptions.healthMonitor); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("if_count_with_filter_interfaces", snmpOptions.ifCountWithFilterInterfaces); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", snmpOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("location", snmpOptions.location); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance_access", snmpOptions.routingInstanceAccess); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance_access_list", snmpOptions.routingInstanceAccessList); tfErr != nil {
		panic(tfErr)
	}
}
