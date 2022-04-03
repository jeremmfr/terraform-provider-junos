package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type evpnOptions struct {
	routingInstanceEvpn bool
	defaultGateway      string
	encapsulation       string
	multicastMode       string
	routingInstance     string
	switchOrRIOptions   []map[string]interface{}
}

func resourceEvpn() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceEvpnCreate,
		ReadWithoutTimeout:   resourceEvpnRead,
		UpdateWithoutTimeout: resourceEvpnUpdate,
		DeleteWithoutTimeout: resourceEvpnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEvpnImport,
		},
		Schema: map[string]*schema.Schema{
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance_evpn": {
				Type:         schema.TypeBool,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"routing_instance", "switch_or_ri_options"},
			},
			"encapsulation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"mpls", "vxlan"}, false),
			},
			"default_gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"advertise", "do-not-advertise", "no-gateway-community"}, false),
			},
			"multicast_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress-replication"}, false),
			},
			"switch_or_ri_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_distinguisher": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^(\d|\.)+L?:\d+$`), "must have valid route distinguisher. Use format 'x:y'"),
						},
						"vrf_export": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vrf_import": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vrf_target": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
						},
						"vrf_target_auto": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"vrf_target_export": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
						},
						"vrf_target_import": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`), "must have valid target. Use format 'target:x:y'"),
						},
					},
				},
			},
		},
	}
}

func resourceEvpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setEvpn(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("routing_instance").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	if err := setEvpn(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_evpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("routing_instance").(string))

	return append(diagWarns, resourceEvpnReadWJnprSess(d, m, jnprSess)...)
}

func resourceEvpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceEvpnReadWJnprSess(d, m, jnprSess)
}

func resourceEvpnReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			mutex.Unlock()

			return diag.FromErr(err)
		}
		if !instanceExists {
			mutex.Unlock()

			d.SetId("")

			return nil
		}
	}
	evpnOptions, err := readEvpn(d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if evpnOptions.routingInstance == "" {
		d.SetId("")
	} else {
		fillEvpnData(d, evpnOptions)
	}

	return nil
}

func resourceEvpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delEvpn(false, d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setEvpn(d, m, nil); err != nil {
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
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delEvpn(false, d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEvpn(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_evpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEvpnReadWJnprSess(d, m, jnprSess)...)
}

func resourceEvpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delEvpn(true, d, m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delEvpn(true, d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_evpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEvpnImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if idList[0] != defaultW {
		instanceExists, err := checkRoutingInstanceExists(idList[0], m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !instanceExists {
			return nil, fmt.Errorf("routing instance %v doesn't exist", idList[0])
		}
	}
	evpnOptions, err := readEvpn(idList[0], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if evpnOptions.routingInstance == "" {
		return nil, fmt.Errorf("don't find protocols evpn with id '%v' "+
			"(id must be <routing_instance>)", d.Id())
	}
	fillEvpnData(d, evpnOptions)
	if len(idList) > 1 || idList[0] == defaultW {
		if tfErr := d.Set("switch_or_ri_options", evpnOptions.switchOrRIOptions); tfErr != nil {
			panic(tfErr)
		}
	}
	result[0] = d

	return result, nil
}

func setEvpn(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := setLS
	setPrefixSwitchRIVRF := setLS
	if d.Get("routing_instance").(string) == defaultW {
		if len(d.Get("switch_or_ri_options").([]interface{})) == 0 {
			return fmt.Errorf("`switch_or_ri_options` required if `routing_instance` = %s", defaultW)
		}
		if d.Get("routing_instance_evpn").(bool) {
			return fmt.Errorf("`routing_instance_evpn` incompatible if `routing_instance` = %s", defaultW)
		}
		if v := d.Get("default_gateway").(string); v != "" {
			return fmt.Errorf("`default_gateway` incompatible if `routing_instance` = %s", defaultW)
		}
		setPrefix += "protocols evpn "
		setPrefixSwitchRIVRF += "switch-options "
	} else {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " protocols evpn "
		setPrefixSwitchRIVRF = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}

	if d.Get("routing_instance_evpn").(bool) {
		if len(d.Get("switch_or_ri_options").([]interface{})) == 0 {
			return fmt.Errorf("`switch_or_ri_options` required if routing_instance_evpn = true")
		}
		configSet = append(configSet, setPrefixSwitchRIVRF+"instance-type evpn")
	}
	if v := d.Get("encapsulation").(string); v != "" {
		configSet = append(configSet, setPrefix+"encapsulation "+v)
	}
	if v := d.Get("default_gateway").(string); v != "" {
		configSet = append(configSet, setPrefix+"default-gateway "+v)
	}
	if v := d.Get("multicast_mode").(string); v != "" {
		configSet = append(configSet, setPrefix+"multicast-mode "+v)
	}
	for _, v := range d.Get("switch_or_ri_options").([]interface{}) {
		swOpts := v.(map[string]interface{})
		configSet = append(configSet, setPrefixSwitchRIVRF+"route-distinguisher "+swOpts["route_distinguisher"].(string))
		for _, v2 := range swOpts["vrf_export"].([]interface{}) {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-export \""+v2.(string)+"\"")
		}
		for _, v2 := range swOpts["vrf_import"].([]interface{}) {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-import \""+v2.(string)+"\"")
		}
		if v2 := swOpts["vrf_target"].(string); v2 != "" {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-target "+v2)
		}
		if swOpts["vrf_target_auto"].(bool) {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-target auto")
		}
		if v2 := swOpts["vrf_target_export"].(string); v2 != "" {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-target export "+v2)
		}
		if v2 := swOpts["vrf_target_import"].(string); v2 != "" {
			configSet = append(configSet, setPrefixSwitchRIVRF+"vrf-target import "+v2)
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readEvpn(routingInstance string, m interface{}, jnprSess *NetconfObject,
) (evpnOptions, error) {
	sess := m.(*Session)
	var confRead evpnOptions
	var showConfig string
	var showConfigSwitchRI string

	if routingInstance == defaultW {
		var err error
		showConfig, err = sess.command(cmdShowConfig+"protocols evpn"+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
		showConfigSwitchRI, err = sess.command(cmdShowConfig+"switch-options"+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		var err error
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols evpn"+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
		showConfigSwitchRI, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+
			pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	}

	if showConfig != emptyW {
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "default-gateway "):
				confRead.defaultGateway = strings.TrimPrefix(itemTrim, "default-gateway ")
			case strings.HasPrefix(itemTrim, "encapsulation "):
				confRead.encapsulation = strings.TrimPrefix(itemTrim, "encapsulation ")
			case strings.HasPrefix(itemTrim, "multicast-mode "):
				confRead.multicastMode = strings.TrimPrefix(itemTrim, "multicast-mode ")
			}
		}
	}
	if showConfigSwitchRI != emptyW {
		for _, item := range strings.Split(showConfigSwitchRI, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == "instance-type evpn":
				confRead.routingInstanceEvpn = true
			case strings.HasPrefix(itemTrim, "route-distinguisher ") ||
				strings.HasPrefix(itemTrim, "vrf-export") ||
				strings.HasPrefix(itemTrim, "vrf-import") ||
				strings.HasPrefix(itemTrim, "vrf-target"):
				if len(confRead.switchOrRIOptions) == 0 {
					confRead.switchOrRIOptions = append(confRead.switchOrRIOptions, map[string]interface{}{
						"route_distinguisher": "",
						"vrf_export":          make([]string, 0),
						"vrf_import":          make([]string, 0),
						"vrf_target":          "",
						"vrf_target_auto":     false,
						"vrf_target_export":   "",
						"vrf_target_import":   "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "route-distinguisher "):
					confRead.switchOrRIOptions[0]["route_distinguisher"] = strings.TrimPrefix(itemTrim, "route-distinguisher ")
				case strings.HasPrefix(itemTrim, "vrf-export "):
					confRead.switchOrRIOptions[0]["vrf_export"] = append(confRead.switchOrRIOptions[0]["vrf_export"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "vrf-export "), "\""))
				case strings.HasPrefix(itemTrim, "vrf-import "):
					confRead.switchOrRIOptions[0]["vrf_import"] = append(confRead.switchOrRIOptions[0]["vrf_import"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "vrf-import "), "\""))
				case itemTrim == "vrf-target auto":
					confRead.switchOrRIOptions[0]["vrf_target_auto"] = true
				case strings.HasPrefix(itemTrim, "vrf-target export "):
					confRead.switchOrRIOptions[0]["vrf_target_export"] = strings.TrimPrefix(itemTrim, "vrf-target export ")
				case strings.HasPrefix(itemTrim, "vrf-target import "):
					confRead.switchOrRIOptions[0]["vrf_target_import"] = strings.TrimPrefix(itemTrim, "vrf-target import ")
				case strings.HasPrefix(itemTrim, "vrf-target "):
					confRead.switchOrRIOptions[0]["vrf_target"] = strings.TrimPrefix(itemTrim, "vrf-target ")
				}
			}
		}
	}

	return confRead, nil
}

func delEvpn(destroy bool, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)

	delPrefix := "delete protocols evpn "
	delPrefixSwitchRIVRF := "delete switch-options "
	if d.Get("routing_instance").(string) != defaultW {
		delPrefix = delRoutingInstances + d.Get("routing_instance").(string) + " protocols evpn "
		delPrefixSwitchRIVRF = delRoutingInstances + d.Get("routing_instance").(string) + " "
	}

	listLinesToDelete := []string{
		"default-gateway",
		"encapsulation",
		"multicast-mode",
	}
	// to remove line "set protocols evpn" without options when destroy resource
	if destroy {
		listLinesToDelete = append(listLinesToDelete, "")
	}
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}

	listLinesToDelete = listLinesToDelete[:0]
	if d.Get("routing_instance_evpn").(bool) {
		listLinesToDelete = append(listLinesToDelete, "instance-type")
	}
	if len(d.Get("switch_or_ri_options").([]interface{})) != 0 {
		listLinesToDelete = append(listLinesToDelete,
			"route-distinguisher",
			"vrf-export",
			"vrf-import",
			"vrf-target",
		)
	}
	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefixSwitchRIVRF+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func fillEvpnData(d *schema.ResourceData, evpnOptions evpnOptions) {
	if tfErr := d.Set("routing_instance", evpnOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance_evpn", evpnOptions.routingInstanceEvpn); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("encapsulation", evpnOptions.encapsulation); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_gateway", evpnOptions.defaultGateway); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("multicast_mode", evpnOptions.multicastMode); tfErr != nil {
		panic(tfErr)
	}
	if _, s := d.GetOk("switch_or_ri_options"); s {
		if tfErr := d.Set("switch_or_ri_options", evpnOptions.switchOrRIOptions); tfErr != nil {
			panic(tfErr)
		}
	}
}
