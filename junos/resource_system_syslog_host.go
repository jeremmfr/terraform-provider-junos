package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type syslogHostOptions struct {
	allowDuplicates             bool
	excludeHostname             bool
	explicitPriority            bool
	port                        int
	host                        string
	facilityOverride            string
	logPrefix                   string
	match                       string
	sourceAddress               string
	anySeverity                 string
	authorizationSeverity       string
	changelogSeverity           string
	conflictlogSeverity         string
	daemonSeverity              string
	dfcSeverity                 string
	externalSeverity            string
	firewallSeverity            string
	ftpSeverity                 string
	interactivecommandsSeverity string
	kernelSeverity              string
	ntpSeverity                 string
	pfeSeverity                 string
	securitySeverity            string
	userSeverity                string
	matchStrings                []string
	structuredData              []map[string]interface{}
}

func resourceSystemSyslogHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemSyslogHostCreate,
		ReadContext:   resourceSystemSyslogHostRead,
		UpdateContext: resourceSystemSyslogHostUpdate,
		DeleteContext: resourceSystemSyslogHostDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemSyslogHostImport,
		},
		Schema: map[string]*schema.Schema{
			"host": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateAddress(),
			},
			"allow_duplicates": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"exclude_hostname": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"explicit_priority": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"facility_override": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogFacility(), false),
			},
			"log_prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"match": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_strings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"source_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"structured_data": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"brief": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"any_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"authorization_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"changelog_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"conflictlog_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"daemon_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"dfc_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"external_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"firewall_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"ftp_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"interactivecommands_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"kernel_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"ntp_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"pfe_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"security_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
			"user_severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
			},
		},
	}
}

func resourceSystemSyslogHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	syslogHostExists, err := checkSystemSyslogHostExists(d.Get("host").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if syslogHostExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("system syslog host %v already exists", d.Get("host").(string)))
	}

	if err := setSystemSyslogHost(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_system_syslog_host", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	syslogHostExists, err = checkSystemSyslogHostExists(d.Get("host").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogHostExists {
		d.SetId(d.Get("host").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system syslog host %v not exists after commit "+
			"=> check your config", d.Get("host").(string)))...)
	}

	return append(diagWarns, resourceSystemSyslogHostReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemSyslogHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemSyslogHostReadWJnprSess(d, m, jnprSess)
}
func resourceSystemSyslogHostReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	syslogHostOptions, err := readSystemSyslogHost(d.Get("host").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if syslogHostOptions.host == "" {
		d.SetId("")
	} else {
		fillSystemSyslogHostData(d, syslogHostOptions)
	}

	return nil
}
func resourceSystemSyslogHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemSyslogHost(d.Get("host").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystemSyslogHost(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_system_syslog_host", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemSyslogHostReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemSyslogHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemSyslogHost(d.Get("host").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_system_syslog_host", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSystemSyslogHostImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	syslogHostExists, err := checkSystemSyslogHostExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !syslogHostExists {
		return nil, fmt.Errorf("don't find system syslog host with id '%v' (id must be <host>)", d.Id())
	}
	syslogHostOptions, err := readSystemSyslogHost(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemSyslogHostData(d, syslogHostOptions)

	result[0] = d

	return result, nil
}

func checkSystemSyslogHostExists(host string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	syslogHostConfig, err := sess.command("show configuration"+
		" system syslog host "+host+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if syslogHostConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSystemSyslogHost(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system syslog host " + d.Get("host").(string)
	configSet := make([]string, 0)

	if d.Get("allow_duplicates").(bool) {
		configSet = append(configSet, setPrefix+" allow-duplicates")
	}
	if d.Get("exclude_hostname").(bool) {
		configSet = append(configSet, setPrefix+" exclude-hostname")
	}
	if d.Get("explicit_priority").(bool) {
		configSet = append(configSet, setPrefix+" explicit-priority")
	}
	if d.Get("facility_override").(string) != "" {
		configSet = append(configSet, setPrefix+" facility-override "+d.Get("facility_override").(string))
	}
	if d.Get("log_prefix").(string) != "" {
		configSet = append(configSet, setPrefix+" log-prefix "+d.Get("log_prefix").(string))
	}
	if d.Get("match").(string) != "" {
		configSet = append(configSet, setPrefix+" match \""+d.Get("match").(string)+"\"")
	}
	for _, v := range d.Get("match_strings").([]interface{}) {
		configSet = append(configSet, setPrefix+" match-strings \""+v.(string)+"\"")
	}
	if d.Get("port").(int) != 0 {
		configSet = append(configSet, setPrefix+" port "+strconv.Itoa(d.Get("port").(int)))
	}
	if d.Get("source_address").(string) != "" {
		configSet = append(configSet, setPrefix+" source-address "+d.Get("source_address").(string))
	}
	for _, v := range d.Get("structured_data").([]interface{}) {
		configSet = append(configSet, setPrefix+" structured-data")
		if v != nil {
			m := v.(map[string]interface{})
			if m["brief"].(bool) {
				configSet = append(configSet, setPrefix+" structured-data brief")
			}
		}
	}
	if d.Get("any_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" any "+d.Get("any_severity").(string))
	}
	if d.Get("authorization_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" authorization "+d.Get("authorization_severity").(string))
	}
	if d.Get("changelog_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" change-log "+d.Get("changelog_severity").(string))
	}
	if d.Get("conflictlog_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" conflict-log "+d.Get("conflictlog_severity").(string))
	}
	if d.Get("daemon_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" daemon "+d.Get("daemon_severity").(string))
	}
	if d.Get("dfc_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" dfc "+d.Get("dfc_severity").(string))
	}
	if d.Get("external_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" external "+d.Get("external_severity").(string))
	}
	if d.Get("firewall_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" firewall "+d.Get("firewall_severity").(string))
	}
	if d.Get("ftp_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" ftp "+d.Get("ftp_severity").(string))
	}
	if d.Get("interactivecommands_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" interactive-commands "+d.Get("interactivecommands_severity").(string))
	}
	if d.Get("kernel_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" kernel "+d.Get("kernel_severity").(string))
	}
	if d.Get("ntp_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" ntp "+d.Get("ntp_severity").(string))
	}
	if d.Get("pfe_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" pfe "+d.Get("pfe_severity").(string))
	}
	if d.Get("security_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" security "+d.Get("security_severity").(string))
	}
	if d.Get("user_severity").(string) != "" {
		configSet = append(configSet, setPrefix+" user "+d.Get("user_severity").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystemSyslogHost(host string, m interface{}, jnprSess *NetconfObject) (syslogHostOptions, error) {
	sess := m.(*Session)
	var confRead syslogHostOptions

	syslogHostConfig, err := sess.command("show configuration"+
		" system syslog host "+host+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if syslogHostConfig != emptyWord {
		confRead.host = host
		for _, item := range strings.Split(syslogHostConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == "allow-duplicates":
				confRead.allowDuplicates = true
			case itemTrim == "exclude-hostname":
				confRead.excludeHostname = true
			case itemTrim == "explicit-priority":
				confRead.explicitPriority = true
			case strings.HasPrefix(itemTrim, "facility-override "):
				confRead.facilityOverride = strings.TrimPrefix(itemTrim, "facility-override ")
			case strings.HasPrefix(itemTrim, "log-prefix "):
				confRead.logPrefix = strings.TrimPrefix(itemTrim, "log-prefix ")
			case strings.HasPrefix(itemTrim, "match "):
				confRead.match = strings.Trim(strings.TrimPrefix(itemTrim, "match "), "\"")
			case strings.HasPrefix(itemTrim, "match-strings "):
				confRead.matchStrings = append(confRead.matchStrings,
					strings.Trim(strings.TrimPrefix(itemTrim, "match-strings "), "\""))
			case strings.HasPrefix(itemTrim, "port "):
				var err error
				confRead.port, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "port "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "source-address "):
				confRead.sourceAddress = strings.TrimPrefix(itemTrim, "source-address ")
			case strings.HasPrefix(itemTrim, "structured-data"):
				structuredData := map[string]interface{}{
					"brief": false,
				}
				if itemTrim == "structured-data brief" {
					structuredData["brief"] = true
				}
				// override (maxItem = 1)
				confRead.structuredData = []map[string]interface{}{structuredData}
			case strings.HasPrefix(itemTrim, "any "):
				confRead.anySeverity = strings.TrimPrefix(itemTrim, "any ")
			case strings.HasPrefix(itemTrim, "authorization "):
				confRead.authorizationSeverity = strings.TrimPrefix(itemTrim, "authorization ")
			case strings.HasPrefix(itemTrim, "change-log "):
				confRead.changelogSeverity = strings.TrimPrefix(itemTrim, "change-log ")
			case strings.HasPrefix(itemTrim, "conflict-log "):
				confRead.conflictlogSeverity = strings.TrimPrefix(itemTrim, "conflict-log ")
			case strings.HasPrefix(itemTrim, "daemon "):
				confRead.daemonSeverity = strings.TrimPrefix(itemTrim, "daemon ")
			case strings.HasPrefix(itemTrim, "dfc "):
				confRead.dfcSeverity = strings.TrimPrefix(itemTrim, "dfc ")
			case strings.HasPrefix(itemTrim, "external "):
				confRead.externalSeverity = strings.TrimPrefix(itemTrim, "external ")
			case strings.HasPrefix(itemTrim, "firewall "):
				confRead.firewallSeverity = strings.TrimPrefix(itemTrim, "firewall ")
			case strings.HasPrefix(itemTrim, "ftp "):
				confRead.ftpSeverity = strings.TrimPrefix(itemTrim, "ftp ")
			case strings.HasPrefix(itemTrim, "interactive-commands "):
				confRead.interactivecommandsSeverity = strings.TrimPrefix(itemTrim, "interactive-commands ")
			case strings.HasPrefix(itemTrim, "kernel "):
				confRead.kernelSeverity = strings.TrimPrefix(itemTrim, "kernel ")
			case strings.HasPrefix(itemTrim, "ntp "):
				confRead.ntpSeverity = strings.TrimPrefix(itemTrim, "ntp ")
			case strings.HasPrefix(itemTrim, "pfe "):
				confRead.pfeSeverity = strings.TrimPrefix(itemTrim, "pfe ")
			case strings.HasPrefix(itemTrim, "security "):
				confRead.securitySeverity = strings.TrimPrefix(itemTrim, "security ")
			case strings.HasPrefix(itemTrim, "user "):
				confRead.userSeverity = strings.TrimPrefix(itemTrim, "user ")
			}
		}
	}

	return confRead, nil
}

func delSystemSyslogHost(host string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system syslog host "+host)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillSystemSyslogHostData(d *schema.ResourceData, syslogHostOptions syslogHostOptions) {
	if tfErr := d.Set("host", syslogHostOptions.host); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_duplicates", syslogHostOptions.allowDuplicates); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("exclude_hostname", syslogHostOptions.excludeHostname); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("explicit_priority", syslogHostOptions.explicitPriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("facility_override", syslogHostOptions.facilityOverride); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("log_prefix", syslogHostOptions.logPrefix); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("match", syslogHostOptions.match); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("match_strings", syslogHostOptions.matchStrings); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port", syslogHostOptions.port); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_address", syslogHostOptions.sourceAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("structured_data", syslogHostOptions.structuredData); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("any_severity", syslogHostOptions.anySeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authorization_severity", syslogHostOptions.authorizationSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("changelog_severity", syslogHostOptions.changelogSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("conflictlog_severity", syslogHostOptions.conflictlogSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("daemon_severity", syslogHostOptions.daemonSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dfc_severity", syslogHostOptions.dfcSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("external_severity", syslogHostOptions.externalSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("firewall_severity", syslogHostOptions.firewallSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ftp_severity", syslogHostOptions.ftpSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interactivecommands_severity", syslogHostOptions.interactivecommandsSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("kernel_severity", syslogHostOptions.kernelSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ntp_severity", syslogHostOptions.ntpSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pfe_severity", syslogHostOptions.pfeSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_severity", syslogHostOptions.securitySeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_severity", syslogHostOptions.userSeverity); tfErr != nil {
		panic(tfErr)
	}
}
