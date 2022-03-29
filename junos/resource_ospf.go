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

type ospfOptions struct {
	disable                      bool
	forwardingAddressToBroadcast bool
	noNssaAbr                    bool
	noRfc1583                    bool
	shamLink                     bool
	externalPreference           int
	labeledPreference            int
	lsaRefreshInterval           int
	preference                   int
	prefixExportLimit            int
	domainID                     string
	referenceBandwidth           string
	ribGroup                     string
	routingInstance              string
	shamLinkLocal                string
	version                      string
	export                       []string
	importL                      []string
	databaseProtection           []map[string]interface{}
	gracefulRestart              []map[string]interface{}
	spfOptions                   []map[string]interface{}
	overload                     []map[string]interface{}
}

func resourceOspf() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOspfCreate,
		ReadContext:   resourceOspfRead,
		UpdateContext: resourceOspfUpdate,
		DeleteContext: resourceOspfDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOspfImport,
		},
		Schema: map[string]*schema.Schema{
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "v2",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v2", "v3"}, false),
			},
			"database_protection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maximum_lsa": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1000000),
						},
						"ignore_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 32),
						},
						"ignore_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 3600),
						},
						"reset_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 86400),
						},
						"warning_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"warning_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 100),
						},
					},
				},
			},
			"disable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain_id": { // only if routing_instance != default
				Type:     schema.TypeString,
				Optional: true,
			},
			"export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"external_preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"forwarding_address_to_broadcast": {
				Type:     schema.TypeBool,
				Optional: true,
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
						"helper_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"helper_disable_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"both", "restart-signaling", "standard"}, false),
							RequiredWith: []string{"graceful_restart.0.helper_disable"},
						},
						"no_strict_lsa_checking": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"graceful_restart.0.helper_disable"},
						},
						"notify_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3600),
						},
						"restart_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3600),
						},
					},
				},
			},
			"import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"labeled_preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"lsa_refresh_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(25, 50),
			},
			"no_nssa_abr": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_rfc1583": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"overload": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_route_leaking": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"as_external": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"stub_network": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 1800),
						},
					},
				},
			},
			"preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"prefix_export_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"reference_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(\d)+(m|k|g)?$`),
					`must be a bandwidth ^(\d)+(m|k|g)?$`),
			},
			"rib_group": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"sham_link": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sham_link_local": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				RequiredWith: []string{"sham_link"},
			},
			"spf_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(50, 8000),
						},
						"holddown": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2000, 20000),
						},
						"no_ignore_our_externals": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rapid_runs": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},
		},
	}
}

func resourceOspfCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setOspf(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("version").(string) + idSeparator + d.Get("routing_instance").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
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
	if err := setOspf(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_ospf", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(d.Get("version").(string) + idSeparator + d.Get("routing_instance").(string))

	return append(diagWarns, resourceOspfReadWJnprSess(d, m, jnprSess)...)
}

func resourceOspfRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceOspfReadWJnprSess(d, m, jnprSess)
}

func resourceOspfReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
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
	ospfOptions, err := readOspf(d.Get("version").(string), d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillOspfData(d, ospfOptions)

	return nil
}

func resourceOspfUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delOspf(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setOspf(d, m, nil); err != nil {
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
	if err := delOspf(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setOspf(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_ospf", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceOspfReadWJnprSess(d, m, jnprSess)...)
}

func resourceOspfDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delOspf(d, m, nil); err != nil {
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
	if err := delOspf(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_ospf", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceOspfImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if idSplit[0] != "v2" && idSplit[0] != "v3" {
		return nil, fmt.Errorf("%s is not a valid version", idSplit[0])
	}
	if idSplit[1] != defaultW {
		instanceExists, err := checkRoutingInstanceExists(idSplit[1], m, jnprSess)
		if err != nil {
			return nil, err
		}
		if !instanceExists {
			return nil, fmt.Errorf("routing instance %v doesn't exist", idSplit[1])
		}
	}
	ospfOptions, err := readOspf(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillOspfData(d, ospfOptions)
	result[0] = d

	return result, nil
}

func setOspf(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	ospfVersion := ospfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = ospfV3
	}
	setPrefix += "protocols " + ospfVersion + " "

	for _, dbPro := range d.Get("database_protection").([]interface{}) {
		dbProM := dbPro.(map[string]interface{})
		configSet = append(configSet, setPrefix+"database-protection maximum-lsa "+strconv.Itoa(dbProM["maximum_lsa"].(int)))
		if v := dbProM["ignore_count"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"database-protection ignore-count "+strconv.Itoa(v))
		}
		if v := dbProM["ignore_time"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"database-protection ignore-time "+strconv.Itoa(v))
		}
		if v := dbProM["reset_time"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"database-protection reset-time "+strconv.Itoa(v))
		}
		if dbProM["warning_only"].(bool) {
			configSet = append(configSet, setPrefix+"database-protection warning-only")
		}
		if v := dbProM["warning_threshold"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"database-protection warning-threshold "+strconv.Itoa(v))
		}
	}
	if d.Get("disable").(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if v := d.Get("domain_id").(string); v != "" {
		if d.Get("routing_instance").(string) == defaultW {
			return fmt.Errorf("domain_id not compatible with routing_instance=default")
		}
		configSet = append(configSet, setPrefix+"domain-id \""+v+"\"")
	}
	for _, v := range d.Get("export").([]interface{}) {
		configSet = append(configSet, setPrefix+"export \""+v.(string)+"\"")
	}
	if v := d.Get("external_preference").(int); v != -1 {
		configSet = append(configSet, setPrefix+"external-preference "+strconv.Itoa(v))
	}
	if d.Get("forwarding_address_to_broadcast").(bool) {
		configSet = append(configSet, setPrefix+"forwarding-address-to-broadcast")
	}
	for _, grR := range d.Get("graceful_restart").([]interface{}) {
		if grR == nil {
			return fmt.Errorf("graceful_restart block is empty")
		}
		grRM := grR.(map[string]interface{})
		if grRM["disable"].(bool) {
			configSet = append(configSet, setPrefix+"graceful-restart disable")
		}
		if grRM["helper_disable"].(bool) {
			configSet = append(configSet, setPrefix+"graceful-restart helper-disable")
			if v := grRM["helper_disable_type"].(string); v != "" {
				configSet = append(configSet, setPrefix+"graceful-restart helper-disable "+v)
			}
		} else if grRM["helper_disable_type"].(string) != "" {
			return fmt.Errorf("helper_disable need to be true with helper_disable_type")
		}
		if grRM["no_strict_lsa_checking"].(bool) {
			configSet = append(configSet, setPrefix+"graceful-restart no-strict-lsa-checking")
		}
		if v := grRM["notify_duration"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"graceful-restart notify-duration "+strconv.Itoa(v))
		}
		if v := grRM["restart_duration"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"graceful-restart restart-duration "+strconv.Itoa(v))
		}
	}
	for _, v := range d.Get("import").([]interface{}) {
		configSet = append(configSet, setPrefix+"import \""+v.(string)+"\"")
	}
	if v := d.Get("labeled_preference").(int); v != -1 {
		configSet = append(configSet, setPrefix+"labeled-preference "+strconv.Itoa(v))
	}
	if v := d.Get("lsa_refresh_interval").(int); v != 0 {
		configSet = append(configSet, setPrefix+"lsa-refresh-interval "+strconv.Itoa(v))
	}
	if d.Get("no_nssa_abr").(bool) {
		configSet = append(configSet, setPrefix+"no-nssa-abr")
	}
	if d.Get("no_rfc1583").(bool) {
		configSet = append(configSet, setPrefix+"no-rfc-1583")
	}
	for _, ovL := range d.Get("overload").([]interface{}) {
		configSet = append(configSet, setPrefix+"overload")
		if ovL != nil {
			ovLM := ovL.(map[string]interface{})
			if ovLM["allow_route_leaking"].(bool) {
				configSet = append(configSet, setPrefix+"overload allow-route-leaking")
			}
			if ovLM["as_external"].(bool) {
				configSet = append(configSet, setPrefix+"overload as-external")
			}
			if ovLM["stub_network"].(bool) {
				configSet = append(configSet, setPrefix+"overload stub-network")
			}
			if v := ovLM["timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"overload timeout "+strconv.Itoa(v))
			}
		}
	}
	if v := d.Get("preference").(int); v != -1 {
		configSet = append(configSet, setPrefix+"preference "+strconv.Itoa(v))
	}
	if v := d.Get("prefix_export_limit").(int); v != -1 {
		configSet = append(configSet, setPrefix+"prefix-export-limit "+strconv.Itoa(v))
	}
	if v := d.Get("reference_bandwidth").(string); v != "" {
		configSet = append(configSet, setPrefix+"reference-bandwidth "+v)
	}
	if v := d.Get("rib_group").(string); v != "" {
		configSet = append(configSet, setPrefix+"rib-group "+v)
	}
	if d.Get("sham_link").(bool) {
		configSet = append(configSet, setPrefix+"sham-link")
		if v := d.Get("sham_link_local").(string); v != "" {
			configSet = append(configSet, setPrefix+"sham-link local "+v)
		}
	} else if d.Get("sham_link_local").(string) != "" {
		return fmt.Errorf("sham_link need to be true with sham_link_local")
	}
	for _, spfO := range d.Get("spf_options").([]interface{}) {
		if spfO == nil {
			return fmt.Errorf("spf_options block is empty")
		}
		sfpOM := spfO.(map[string]interface{})
		if v := sfpOM["delay"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"spf-options delay "+strconv.Itoa(v))
		}
		if v := sfpOM["holddown"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"spf-options holddown "+strconv.Itoa(v))
		}
		if sfpOM["no_ignore_our_externals"].(bool) {
			configSet = append(configSet, setPrefix+"spf-options no-ignore-our-externals")
		}
		if v := sfpOM["rapid_runs"].(int); v != 0 {
			configSet = append(configSet, setPrefix+"spf-options rapid-runs "+strconv.Itoa(v))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readOspf(version, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (ospfOptions, error) {
	sess := m.(*Session)
	var confRead ospfOptions
	confRead.externalPreference = -1
	confRead.labeledPreference = -1
	confRead.preference = -1
	confRead.prefixExportLimit = -1

	var showConfig string
	ospfVersion := ospfV2
	if version == "v3" {
		ospfVersion = ospfV3
	}
	if routingInstance == defaultW {
		var err error
		showConfig, err = sess.command(cmdShowConfig+
			"protocols "+ospfVersion+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		var err error
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols "+ospfVersion+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	}

	confRead.version = version
	confRead.routingInstance = routingInstance
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "database-protection"):
				if len(confRead.databaseProtection) == 0 {
					confRead.databaseProtection = append(confRead.databaseProtection, map[string]interface{}{
						"ignore_count":      0,
						"ignore_time":       0,
						"maximum_lsa":       0,
						"reset_time":        0,
						"warning_only":      false,
						"warning_threshold": 0,
					})
				}
				dbPro := confRead.databaseProtection[0]
				itemTrimDP := strings.TrimPrefix(itemTrim, "database-protection ")
				switch {
				case strings.HasPrefix(itemTrimDP, "ignore-count "):
					var err error
					dbPro["ignore_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimDP, "ignore-count "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimDP, "ignore-time "):
					var err error
					dbPro["ignore_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrimDP, "ignore-time "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimDP, "maximum-lsa "):
					var err error
					dbPro["maximum_lsa"], err = strconv.Atoi(strings.TrimPrefix(itemTrimDP, "maximum-lsa "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrimDP, "reset-time "):
					var err error
					dbPro["reset_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrimDP, "reset-time "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrimDP == "warning-only":
					dbPro["warning_only"] = true
				case strings.HasPrefix(itemTrimDP, "warning-threshold "):
					var err error
					dbPro["warning_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrimDP, "warning-threshold "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case itemTrim == "disable":
				confRead.disable = true
			case strings.HasPrefix(itemTrim, "domain-id "):
				confRead.domainID = strings.Trim(strings.TrimPrefix(itemTrim, "domain-id "), "\"")
			case strings.HasPrefix(itemTrim, "export "):
				confRead.export = append(confRead.export, strings.Trim(strings.TrimPrefix(itemTrim, "export "), "\""))
			case strings.HasPrefix(itemTrim, "external-preference "):
				var err error
				confRead.externalPreference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "external-preference "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "forwarding-address-to-broadcast":
				confRead.forwardingAddressToBroadcast = true
			case strings.HasPrefix(itemTrim, "graceful-restart "):
				if len(confRead.gracefulRestart) == 0 {
					confRead.gracefulRestart = append(confRead.gracefulRestart, map[string]interface{}{
						"disable":                false,
						"helper_disable":         false,
						"helper_disable_type":    "",
						"no_strict_lsa_checking": false,
						"notify_duration":        0,
						"restart_duration":       0,
					})
				}
				grR := confRead.gracefulRestart[0]
				switch {
				case itemTrim == "graceful-restart disable":
					grR["disable"] = true
				case strings.HasPrefix(itemTrim, "graceful-restart helper-disable"):
					grR["helper_disable"] = true
					if strings.HasPrefix(itemTrim, "graceful-restart helper-disable ") {
						grR["helper_disable_type"] = strings.TrimPrefix(itemTrim, "graceful-restart helper-disable ")
					}
				case itemTrim == "graceful-restart no-strict-lsa-checking":
					grR["no_strict_lsa_checking"] = true
				case strings.HasPrefix(itemTrim, "graceful-restart notify-duration "):
					var err error
					grR["notify_duration"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "graceful-restart notify-duration "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "graceful-restart restart-duration "):
					var err error
					grR["restart_duration"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "graceful-restart restart-duration "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "import "):
				confRead.importL = append(confRead.importL, strings.Trim(strings.TrimPrefix(itemTrim, "import "), "\""))
			case strings.HasPrefix(itemTrim, "labeled-preference "):
				var err error
				confRead.labeledPreference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "labeled-preference "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "lsa-refresh-interval "):
				var err error
				confRead.lsaRefreshInterval, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lsa-refresh-interval "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "no-nssa-abr":
				confRead.noNssaAbr = true
			case itemTrim == "no-rfc-1583":
				confRead.noRfc1583 = true
			case strings.HasPrefix(itemTrim, "overload"):
				if len(confRead.overload) == 0 {
					confRead.overload = append(confRead.overload, map[string]interface{}{
						"allow_route_leaking": false,
						"as_external":         false,
						"stub_network":        false,
						"timeout":             0,
					})
				}
				switch {
				case itemTrim == "overload allow-route-leaking":
					confRead.overload[0]["allow_route_leaking"] = true
				case itemTrim == "overload as-external":
					confRead.overload[0]["as_external"] = true
				case itemTrim == "overload stub-network":
					confRead.overload[0]["stub_network"] = true
				case strings.HasPrefix(itemTrim, "overload timeout "):
					var err error
					confRead.overload[0]["timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "overload timeout "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "preference "):
				var err error
				confRead.preference, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preference "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "prefix-export-limit "):
				var err error
				confRead.prefixExportLimit, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "prefix-export-limit "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "reference-bandwidth "):
				confRead.referenceBandwidth = strings.TrimPrefix(itemTrim, "reference-bandwidth ")
			case strings.HasPrefix(itemTrim, "rib-group "):
				confRead.ribGroup = strings.TrimPrefix(itemTrim, "rib-group ")
			case strings.HasPrefix(itemTrim, "sham-link"):
				confRead.shamLink = true
				if strings.HasPrefix(itemTrim, "sham-link local ") {
					confRead.shamLinkLocal = strings.TrimPrefix(itemTrim, "sham-link local ")
				}
			case strings.HasPrefix(itemTrim, "spf-options "):
				if len(confRead.spfOptions) == 0 {
					confRead.spfOptions = append(confRead.spfOptions, map[string]interface{}{
						"delay":                   0,
						"holddown":                0,
						"no_ignore_our_externals": false,
						"rapid_runs":              0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "spf-options delay "):
					var err error
					confRead.spfOptions[0]["delay"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "spf-options delay "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "spf-options holddown "):
					var err error
					confRead.spfOptions[0]["holddown"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "spf-options holddown "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "spf-options no-ignore-our-externals":
					confRead.spfOptions[0]["no_ignore_our_externals"] = true
				case strings.HasPrefix(itemTrim, "spf-options rapid-runs "):
					var err error
					confRead.spfOptions[0]["rapid_runs"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "spf-options rapid-runs "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			}
		}
	}

	return confRead, nil
}

func delOspf(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)

	ospfVersion := ospfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = ospfV3
	}
	delPrefix := deleteLS
	if d.Get("routing_instance").(string) != defaultW {
		delPrefix = delRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	delPrefix += "protocols " + ospfVersion + " "

	listLinesToDelete := []string{
		"database-protection",
		"disable",
		"domain-id",
		"export",
		"external-preference",
		"forwarding-address-to-broadcast",
		"graceful-restart",
		"import",
		"labeled-preference",
		"lsa-refresh-interval",
		"no-nssa-abr",
		"no-rfc-1583",
		"overload",
		"preference",
		"prefix-export-limit",
		"reference-bandwidth",
		"rib-group",
		"sham-link",
		"spf-options",
	}

	for _, line := range listLinesToDelete {
		configSet = append(configSet, delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func fillOspfData(d *schema.ResourceData, ospfOptions ospfOptions) {
	if tfErr := d.Set("routing_instance", ospfOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ospfOptions.version); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("database_protection", ospfOptions.databaseProtection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("disable", ospfOptions.disable); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain_id", ospfOptions.domainID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export", ospfOptions.export); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("external_preference", ospfOptions.externalPreference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forwarding_address_to_broadcast", ospfOptions.forwardingAddressToBroadcast); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_restart", ospfOptions.gracefulRestart); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import", ospfOptions.importL); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("labeled_preference", ospfOptions.labeledPreference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lsa_refresh_interval", ospfOptions.lsaRefreshInterval); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_nssa_abr", ospfOptions.noNssaAbr); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_rfc1583", ospfOptions.noRfc1583); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("overload", ospfOptions.overload); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", ospfOptions.preference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefix_export_limit", ospfOptions.prefixExportLimit); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reference_bandwidth", ospfOptions.referenceBandwidth); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rib_group", ospfOptions.ribGroup); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("sham_link", ospfOptions.shamLink); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("sham_link_local", ospfOptions.shamLinkLocal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("spf_options", ospfOptions.spfOptions); tfErr != nil {
		panic(tfErr)
	}
}
