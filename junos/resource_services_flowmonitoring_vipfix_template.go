package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type flowMonitoringVIPFixTemplateOptions struct {
	flowKeyFlowDirection   bool
	flowKeyVlanID          bool
	nexthopLearningEnable  bool
	nexthopLearningDisable bool
	tunnelObservationIPv4  bool
	tunnelObservationIPv6  bool
	flowActiveTimeout      int
	flowInactiveTimeout    int
	observationDomainID    int
	optionTemplateID       int
	templateID             int
	name                   string
	typeTemplate           string
	ipTemplateExportExt    []string
	optionRefreshRate      []map[string]interface{}
}

func resourceServicesFlowMonitoringVIPFixTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesFlowMonitoringVIPFixTemplateCreate,
		ReadContext:   resourceServicesFlowMonitoringVIPFixTemplateRead,
		UpdateContext: resourceServicesFlowMonitoringVIPFixTemplateUpdate,
		DeleteContext: resourceServicesFlowMonitoringVIPFixTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesFlowMonitoringVIPFixTemplateImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4-template", "ipv6-template", "mpls-template"}, false),
			},
			"flow_active_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 600),
			},
			"flow_inactive_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 600),
			},
			"flow_key_flow_direction": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"flow_key_vlan_id": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip_template_export_extension": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nexthop_learning_enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"nexthop_learning_disable"},
			},
			"nexthop_learning_disable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"nexthop_learning_enable"},
			},
			"observation_domain_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"option_refresh_rate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"packets": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 480000),
						},
						"seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 600),
						},
					},
				},
			},
			"option_template_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1024, 65535),
			},
			"template_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1024, 65535),
			},
			"tunnel_observation_ipv4": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tunnel_observation_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceServicesFlowMonitoringVIPFixTemplateCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesFlowMonitoringVIPFixTemplate(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	flowMonitoringVIPFixTemplateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(
		d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if flowMonitoringVIPFixTemplateExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services flow-monitoring version-ipfix template "+
				" %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesFlowMonitoringVIPFixTemplate(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services_flowmonitoring_vipfix_template", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	flowMonitoringVIPFixTemplateExists, err = checkServicesFlowMonitoringVIPFixTemplateExists(
		d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if flowMonitoringVIPFixTemplateExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services flow-monitoring version-ipfix template %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesFlowMonitoringVIPFixTemplateReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesFlowMonitoringVIPFixTemplateRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesFlowMonitoringVIPFixTemplateReadWJnprSess(d, m, jnprSess)
}

func resourceServicesFlowMonitoringVIPFixTemplateReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	flowMonitoringVIPFixTemplateOptions, err := readServicesFlowMonitoringVIPFixTemplate(
		d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if flowMonitoringVIPFixTemplateOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesFlowMonitoringVIPFixTemplateData(d, flowMonitoringVIPFixTemplateOptions)
	}

	return nil
}

func resourceServicesFlowMonitoringVIPFixTemplateUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesFlowMonitoringVIPFixTemplate(d, m, nil); err != nil {
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
	if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesFlowMonitoringVIPFixTemplate(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services_flowmonitoring_vipfix_template", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesFlowMonitoringVIPFixTemplateReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesFlowMonitoringVIPFixTemplateDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), m, nil); err != nil {
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
	if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_services_flowmonitoring_vipfix_template", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesFlowMonitoringVIPFixTemplateImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	flowMonitoringVIPFixTemplateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !flowMonitoringVIPFixTemplateExists {
		return nil, fmt.Errorf("don't find services flow-monitoring version-ipfix template with "+
			"id '%v' (id must be <name>)", d.Id())
	}
	flowMonitoringVIPFixTemplateOptions, err := readServicesFlowMonitoringVIPFixTemplate(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesFlowMonitoringVIPFixTemplateData(d, flowMonitoringVIPFixTemplateOptions)

	result[0] = d

	return result, nil
}

func checkServicesFlowMonitoringVIPFixTemplateExists(
	template string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" services flow-monitoring version-ipfix template \""+template+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesFlowMonitoringVIPFixTemplate(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services flow-monitoring version-ipfix template \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+d.Get("type").(string))
	for _, v := range sortSetOfString(d.Get("ip_template_export_extension").(*schema.Set).List()) {
		if v2 := d.Get("type").(string); v2 != "ipv4-template" && v2 != "ipv6-template" {
			return fmt.Errorf("ip_template_export_extension not compatible with type %s", v2)
		}
		configSet = append(configSet, setPrefix+d.Get("type").(string)+" export-extension "+v)
	}
	if v := d.Get("flow_active_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+strconv.Itoa(v))
	}
	if v := d.Get("flow_inactive_timeout").(int); v != 0 {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+strconv.Itoa(v))
	}
	if d.Get("flow_key_flow_direction").(bool) {
		configSet = append(configSet, setPrefix+"flow-key flow-direction")
	}
	if d.Get("flow_key_vlan_id").(bool) {
		configSet = append(configSet, setPrefix+"flow-key vlan-id")
	}
	if d.Get("nexthop_learning_enable").(bool) {
		configSet = append(configSet, setPrefix+"nexthop-learning enable")
	}
	if d.Get("nexthop_learning_disable").(bool) {
		configSet = append(configSet, setPrefix+"nexthop-learning disable")
	}
	if v := d.Get("observation_domain_id").(int); v != -1 {
		configSet = append(configSet, setPrefix+"observation-domain-id "+strconv.Itoa(v))
	}
	for _, v := range d.Get("option_refresh_rate").([]interface{}) {
		configSet = append(configSet, setPrefix+"option-refresh-rate")
		if v != nil {
			optRefRate := v.(map[string]interface{})
			if v2 := optRefRate["packets"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"option-refresh-rate packets "+strconv.Itoa(v2))
			}
			if v2 := optRefRate["seconds"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"option-refresh-rate seconds "+strconv.Itoa(v2))
			}
		}
	}
	if v := d.Get("option_template_id").(int); v != 0 {
		configSet = append(configSet, setPrefix+"option-template-id "+strconv.Itoa(v))
	}
	if v := d.Get("template_id").(int); v != 0 {
		configSet = append(configSet, setPrefix+"template-id "+strconv.Itoa(v))
	}
	if d.Get("tunnel_observation_ipv4").(bool) {
		configSet = append(configSet, setPrefix+"tunnel-observation ipv4")
	}
	if d.Get("tunnel_observation_ipv6").(bool) {
		configSet = append(configSet, setPrefix+"tunnel-observation ipv6")
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesFlowMonitoringVIPFixTemplate(template string, m interface{}, jnprSess *NetconfObject) (
	flowMonitoringVIPFixTemplateOptions, error) {
	sess := m.(*Session)
	var confRead flowMonitoringVIPFixTemplateOptions
	// setup default value
	confRead.observationDomainID = -1

	showConfig, err := sess.command("show configuration"+
		" services flow-monitoring version-ipfix template \""+template+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = template
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case bchk.StringInSlice(itemTrim, []string{"ipv4-template", "ipv6-template", "mpls-template"}):
				confRead.typeTemplate = itemTrim
			case strings.HasPrefix(itemTrim, "ipv6-template export-extension ") ||
				strings.HasPrefix(itemTrim, "ipv4-template export-extension "):
				itemTrimSplit := strings.Split(itemTrim, " ")
				confRead.typeTemplate = itemTrimSplit[0]
				confRead.ipTemplateExportExt = append(confRead.ipTemplateExportExt,
					strings.TrimPrefix(itemTrim, itemTrimSplit[0]+" export-extension "))
			case strings.HasPrefix(itemTrim, "flow-active-timeout "):
				var err error
				confRead.flowActiveTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "flow-active-timeout "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "flow-inactive-timeout "):
				var err error
				confRead.flowInactiveTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "flow-inactive-timeout "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == "flow-key flow-direction":
				confRead.flowKeyFlowDirection = true
			case itemTrim == "flow-key vlan-id":
				confRead.flowKeyVlanID = true
			case itemTrim == "nexthop-learning enable":
				confRead.nexthopLearningEnable = true
			case itemTrim == "nexthop-learning disable":
				confRead.nexthopLearningDisable = true
			case strings.HasPrefix(itemTrim, "observation-domain-id "):
				var err error
				confRead.observationDomainID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "observation-domain-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "option-refresh-rate"):
				if len(confRead.optionRefreshRate) == 0 {
					confRead.optionRefreshRate = append(confRead.optionRefreshRate, map[string]interface{}{
						"packets": 0,
						"seconds": 0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "option-refresh-rate packets "):
					var err error
					confRead.optionRefreshRate[0]["packets"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "option-refresh-rate packets "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "option-refresh-rate seconds "):
					var err error
					confRead.optionRefreshRate[0]["seconds"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "option-refresh-rate seconds "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "option-template-id "):
				var err error
				confRead.optionTemplateID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "option-template-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "template-id "):
				var err error
				confRead.templateID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "template-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case itemTrim == "tunnel-observation ipv4":
				confRead.tunnelObservationIPv4 = true
			case itemTrim == "tunnel-observation ipv6":
				confRead.tunnelObservationIPv6 = true
			}
		}
	}

	return confRead, nil
}

func delServicesFlowMonitoringVIPFixTemplate(template string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete services flow-monitoring version-ipfix template \"" + template + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesFlowMonitoringVIPFixTemplateData(
	d *schema.ResourceData, flowMonitoringVIPFixTemplateOptions flowMonitoringVIPFixTemplateOptions) {
	if tfErr := d.Set("name", flowMonitoringVIPFixTemplateOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("type", flowMonitoringVIPFixTemplateOptions.typeTemplate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ip_template_export_extension",
		flowMonitoringVIPFixTemplateOptions.ipTemplateExportExt); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("flow_active_timeout", flowMonitoringVIPFixTemplateOptions.flowActiveTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("flow_inactive_timeout", flowMonitoringVIPFixTemplateOptions.flowInactiveTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("flow_key_flow_direction", flowMonitoringVIPFixTemplateOptions.flowKeyFlowDirection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("flow_key_vlan_id", flowMonitoringVIPFixTemplateOptions.flowKeyVlanID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("nexthop_learning_enable", flowMonitoringVIPFixTemplateOptions.nexthopLearningEnable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("nexthop_learning_disable",
		flowMonitoringVIPFixTemplateOptions.nexthopLearningDisable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("observation_domain_id", flowMonitoringVIPFixTemplateOptions.observationDomainID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("option_refresh_rate", flowMonitoringVIPFixTemplateOptions.optionRefreshRate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("option_template_id", flowMonitoringVIPFixTemplateOptions.optionTemplateID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("template_id", flowMonitoringVIPFixTemplateOptions.templateID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tunnel_observation_ipv4", flowMonitoringVIPFixTemplateOptions.tunnelObservationIPv4); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tunnel_observation_ipv6", flowMonitoringVIPFixTemplateOptions.tunnelObservationIPv6); tfErr != nil {
		panic(tfErr)
	}
}
