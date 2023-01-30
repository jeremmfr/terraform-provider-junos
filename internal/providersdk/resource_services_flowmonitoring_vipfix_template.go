package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceServicesFlowMonitoringVIPFixTemplateCreate,
		ReadWithoutTimeout:   resourceServicesFlowMonitoringVIPFixTemplateRead,
		UpdateWithoutTimeout: resourceServicesFlowMonitoringVIPFixTemplateUpdate,
		DeleteWithoutTimeout: resourceServicesFlowMonitoringVIPFixTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesFlowMonitoringVIPFixTemplateImport,
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

func resourceServicesFlowMonitoringVIPFixTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setServicesFlowMonitoringVIPFixTemplate(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	flowMonitoringVIPFixTemplateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if flowMonitoringVIPFixTemplateExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services flow-monitoring version-ipfix template "+
				" %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesFlowMonitoringVIPFixTemplate(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_services_flowmonitoring_vipfix_template")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	flowMonitoringVIPFixTemplateExists, err = checkServicesFlowMonitoringVIPFixTemplateExists(
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if flowMonitoringVIPFixTemplateExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services flow-monitoring version-ipfix template %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesFlowMonitoringVIPFixTemplateReadWJunSess(d, junSess)...)
}

func resourceServicesFlowMonitoringVIPFixTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceServicesFlowMonitoringVIPFixTemplateReadWJunSess(d, junSess)
}

func resourceServicesFlowMonitoringVIPFixTemplateReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	flowMonitoringVIPFixTemplateOptions, err := readServicesFlowMonitoringVIPFixTemplate(
		d.Get("name").(string),
		junSess,
	)
	junos.MutexUnlock()
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

func resourceServicesFlowMonitoringVIPFixTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesFlowMonitoringVIPFixTemplate(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

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
	if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesFlowMonitoringVIPFixTemplate(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_services_flowmonitoring_vipfix_template")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesFlowMonitoringVIPFixTemplateReadWJunSess(d, junSess)...)
}

func resourceServicesFlowMonitoringVIPFixTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), junSess); err != nil {
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
	if err := delServicesFlowMonitoringVIPFixTemplate(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_services_flowmonitoring_vipfix_template")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesFlowMonitoringVIPFixTemplateImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	flowMonitoringVIPFixTemplateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !flowMonitoringVIPFixTemplateExists {
		return nil, fmt.Errorf("don't find services flow-monitoring version-ipfix template with "+
			"id '%v' (id must be <name>)", d.Id())
	}
	flowMonitoringVIPFixTemplateOptions, err := readServicesFlowMonitoringVIPFixTemplate(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillServicesFlowMonitoringVIPFixTemplateData(d, flowMonitoringVIPFixTemplateOptions)

	result[0] = d

	return result, nil
}

func checkServicesFlowMonitoringVIPFixTemplateExists(template string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services flow-monitoring version-ipfix template \"" + template + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesFlowMonitoringVIPFixTemplate(d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func readServicesFlowMonitoringVIPFixTemplate(template string, junSess *junos.Session,
) (confRead flowMonitoringVIPFixTemplateOptions, err error) {
	// default -1
	confRead.observationDomainID = -1
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services flow-monitoring version-ipfix template \"" + template + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = template
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case bchk.InSlice(itemTrim, []string{"ipv4-template", "ipv6-template", "mpls-template"}):
				confRead.typeTemplate = itemTrim
			case balt.CutPrefixInString(&itemTrim, "ipv6-template export-extension "):
				confRead.typeTemplate = "ipv6-template"
				confRead.ipTemplateExportExt = append(confRead.ipTemplateExportExt, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ipv4-template export-extension "):
				confRead.typeTemplate = "ipv4-template"
				confRead.ipTemplateExportExt = append(confRead.ipTemplateExportExt, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
				confRead.flowActiveTimeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
				confRead.flowInactiveTimeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "flow-key flow-direction":
				confRead.flowKeyFlowDirection = true
			case itemTrim == "flow-key vlan-id":
				confRead.flowKeyVlanID = true
			case itemTrim == "nexthop-learning enable":
				confRead.nexthopLearningEnable = true
			case itemTrim == "nexthop-learning disable":
				confRead.nexthopLearningDisable = true
			case balt.CutPrefixInString(&itemTrim, "observation-domain-id "):
				confRead.observationDomainID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "option-refresh-rate"):
				if len(confRead.optionRefreshRate) == 0 {
					confRead.optionRefreshRate = append(confRead.optionRefreshRate, map[string]interface{}{
						"packets": 0,
						"seconds": 0,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " packets "):
					confRead.optionRefreshRate[0]["packets"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " seconds "):
					confRead.optionRefreshRate[0]["seconds"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "option-template-id "):
				confRead.optionTemplateID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "template-id "):
				confRead.templateID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
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

func delServicesFlowMonitoringVIPFixTemplate(template string, junSess *junos.Session) error {
	configSet := []string{"delete services flow-monitoring version-ipfix template \"" + template + "\""}

	return junSess.ConfigSet(configSet)
}

func fillServicesFlowMonitoringVIPFixTemplateData(
	d *schema.ResourceData, flowMonitoringVIPFixTemplateOptions flowMonitoringVIPFixTemplateOptions,
) {
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
