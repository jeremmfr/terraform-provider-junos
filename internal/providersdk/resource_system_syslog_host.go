package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceSystemSyslogHostCreate,
		ReadWithoutTimeout:   resourceSystemSyslogHostRead,
		UpdateWithoutTimeout: resourceSystemSyslogHostUpdate,
		DeleteWithoutTimeout: resourceSystemSyslogHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemSyslogHostImport,
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
				ValidateFunc: validation.StringInSlice(junos.SyslogFacilities(), false),
			},
			"log_prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSystemSyslogHost(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("host").(string))

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
	syslogHostExists, err := checkSystemSyslogHostExists(d.Get("host").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogHostExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("system syslog host %v already exists", d.Get("host").(string)))...)
	}

	if err := setSystemSyslogHost(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_system_syslog_host")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	syslogHostExists, err = checkSystemSyslogHostExists(d.Get("host").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogHostExists {
		d.SetId(d.Get("host").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system syslog host %v not exists after commit "+
			"=> check your config", d.Get("host").(string)))...)
	}

	return append(diagWarns, resourceSystemSyslogHostReadWJunSess(d, junSess)...)
}

func resourceSystemSyslogHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemSyslogHostReadWJunSess(d, junSess)
}

func resourceSystemSyslogHostReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	syslogHostOptions, err := readSystemSyslogHost(d.Get("host").(string), junSess)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemSyslogHost(d.Get("host").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemSyslogHost(d, junSess); err != nil {
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
	if err := delSystemSyslogHost(d.Get("host").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemSyslogHost(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_system_syslog_host")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemSyslogHostReadWJunSess(d, junSess)...)
}

func resourceSystemSyslogHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemSyslogHost(d.Get("host").(string), junSess); err != nil {
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
	if err := delSystemSyslogHost(d.Get("host").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_system_syslog_host")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemSyslogHostImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	syslogHostExists, err := checkSystemSyslogHostExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !syslogHostExists {
		return nil, fmt.Errorf("don't find system syslog host with id '%v' (id must be <host>)", d.Id())
	}
	syslogHostOptions, err := readSystemSyslogHost(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSystemSyslogHostData(d, syslogHostOptions)

	result[0] = d

	return result, nil
}

func checkSystemSyslogHostExists(host string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system syslog host " + host + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemSyslogHost(d *schema.ResourceData, junSess *junos.Session) error {
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
			ma := v.(map[string]interface{})
			if ma["brief"].(bool) {
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

	return junSess.ConfigSet(configSet)
}

func readSystemSyslogHost(host string, junSess *junos.Session,
) (confRead syslogHostOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system syslog host " + host + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.host = host
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "allow-duplicates":
				confRead.allowDuplicates = true
			case itemTrim == "exclude-hostname":
				confRead.excludeHostname = true
			case itemTrim == "explicit-priority":
				confRead.explicitPriority = true
			case balt.CutPrefixInString(&itemTrim, "facility-override "):
				confRead.facilityOverride = itemTrim
			case balt.CutPrefixInString(&itemTrim, "log-prefix "):
				confRead.logPrefix = itemTrim
			case balt.CutPrefixInString(&itemTrim, "match "):
				confRead.match = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "match-strings "):
				confRead.matchStrings = append(confRead.matchStrings, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "port "):
				confRead.port, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "source-address "):
				confRead.sourceAddress = itemTrim
			case balt.CutPrefixInString(&itemTrim, "structured-data"):
				if len(confRead.structuredData) == 0 {
					confRead.structuredData = append(confRead.structuredData, map[string]interface{}{
						"brief": false,
					})
				}
				if itemTrim == " brief" {
					confRead.structuredData[0]["brief"] = true
				}
			case balt.CutPrefixInString(&itemTrim, "any "):
				confRead.anySeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "authorization "):
				confRead.authorizationSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "change-log "):
				confRead.changelogSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "conflict-log "):
				confRead.conflictlogSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "daemon "):
				confRead.daemonSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "dfc "):
				confRead.dfcSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "external "):
				confRead.externalSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "firewall "):
				confRead.firewallSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "ftp "):
				confRead.ftpSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "interactive-commands "):
				confRead.interactivecommandsSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "kernel "):
				confRead.kernelSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "ntp "):
				confRead.ntpSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "pfe "):
				confRead.pfeSeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "security "):
				confRead.securitySeverity = itemTrim
			case balt.CutPrefixInString(&itemTrim, "user "):
				confRead.userSeverity = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSystemSyslogHost(host string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system syslog host "+host)

	return junSess.ConfigSet(configSet)
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
