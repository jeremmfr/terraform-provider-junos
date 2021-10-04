package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
)

type syslogFileOptions struct {
	allowDuplicates             bool
	explicitPriority            bool
	filename                    string
	match                       string
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
	archive                     []map[string]interface{}
	structuredData              []map[string]interface{}
}

func resourceSystemSyslogFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemSyslogFileCreate,
		ReadContext:   resourceSystemSyslogFileRead,
		UpdateContext: resourceSystemSyslogFileUpdate,
		DeleteContext: resourceSystemSyslogFileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemSyslogFileImport,
		},
		Schema: map[string]*schema.Schema{
			"filename": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"allow_duplicates": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"explicit_priority": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"archive": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sites": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
									"routing_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"binary_data": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"archive.0.no_binary_data"},
						},
						"no_binary_data": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"archive.0.binary_data"},
						},
						"world_readable": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"archive.0.no_world_readable"},
						},
						"no_world_readable": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"archive.0.world_readable"},
						},
						"files": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(65536, 1073741824),
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"transfer_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5, 2880),
						},
					},
				},
			},
		},
	}
}

func resourceSystemSyslogFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSystemSyslogFile(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("filename").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	syslogFileExists, err := checkSystemSyslogFileExists(d.Get("filename").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogFileExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system syslog file %v already exists", d.Get("filename").(string)))...)
	}

	if err := setSystemSyslogFile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_system_syslog_file", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	syslogFileExists, err = checkSystemSyslogFileExists(d.Get("filename").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogFileExists {
		d.SetId(d.Get("filename").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system syslog file %v not exists after commit "+
			"=> check your config", d.Get("filename").(string)))...)
	}

	return append(diagWarns, resourceSystemSyslogFileReadWJnprSess(d, m, jnprSess)...)
}

func resourceSystemSyslogFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemSyslogFileReadWJnprSess(d, m, jnprSess)
}

func resourceSystemSyslogFileReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	syslogFileOptions, err := readSystemSyslogFile(d.Get("filename").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if syslogFileOptions.filename == "" {
		d.SetId("")
	} else {
		fillSystemSyslogFileData(d, syslogFileOptions)
	}

	return nil
}

func resourceSystemSyslogFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSystemSyslogFile(d.Get("filename").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemSyslogFile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_system_syslog_file", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemSyslogFileReadWJnprSess(d, m, jnprSess)...)
}

func resourceSystemSyslogFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSystemSyslogFile(d.Get("filename").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_system_syslog_file", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemSyslogFileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	syslogFileExists, err := checkSystemSyslogFileExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !syslogFileExists {
		return nil, fmt.Errorf("don't find system syslog file with id '%v' (id must be <filename>)", d.Id())
	}
	syslogFileOptions, err := readSystemSyslogFile(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemSyslogFileData(d, syslogFileOptions)

	result[0] = d

	return result, nil
}

func checkSystemSyslogFileExists(filename string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration system syslog file "+filename+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSystemSyslogFile(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system syslog file " + d.Get("filename").(string)
	configSet := make([]string, 0)

	if d.Get("allow_duplicates").(bool) {
		configSet = append(configSet, setPrefix+" allow-duplicates")
	}
	if d.Get("explicit_priority").(bool) {
		configSet = append(configSet, setPrefix+" explicit-priority")
	}
	if d.Get("match").(string) != "" {
		configSet = append(configSet, setPrefix+" match \""+d.Get("match").(string)+"\"")
	}
	for _, v := range d.Get("match_strings").([]interface{}) {
		configSet = append(configSet, setPrefix+" match-strings \""+v.(string)+"\"")
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
	for _, v := range d.Get("archive").([]interface{}) {
		setPrefixArchive := setPrefix + " archive"
		configSet = append(configSet, setPrefixArchive)
		if v != nil {
			archive := v.(map[string]interface{})
			sitesURLList := make([]string, 0)
			for _, v2 := range archive["sites"].([]interface{}) {
				sites := v2.(map[string]interface{})
				if stringInSlice(sites["url"].(string), sitesURLList) {
					return fmt.Errorf("multiple sites blocks with the same url")
				}
				sitesURLList = append(sitesURLList, sites["url"].(string))
				setPrefixArchiveSite := setPrefixArchive + " archive-sites " + sites["url"].(string)
				configSet = append(configSet, setPrefixArchiveSite)
				if sites["password"].(string) != "" {
					configSet = append(configSet, setPrefixArchiveSite+" password \""+
						sites["password"].(string)+"\"")
				}
				if sites["routing_instance"].(string) != "" {
					configSet = append(configSet, setPrefixArchiveSite+" routing-instance "+
						sites["routing_instance"].(string))
				}
			}
			if archive["binary_data"].(bool) {
				configSet = append(configSet, setPrefixArchive+" binary-data")
			}
			if archive["no_binary_data"].(bool) {
				configSet = append(configSet, setPrefixArchive+" no-binary-data")
			}
			if archive["world_readable"].(bool) {
				configSet = append(configSet, setPrefixArchive+" world-readable")
			}
			if archive["no_world_readable"].(bool) {
				configSet = append(configSet, setPrefixArchive+" no-world-readable")
			}
			if archive["files"].(int) != 0 {
				configSet = append(configSet, setPrefixArchive+" files "+
					strconv.Itoa(archive["files"].(int)))
			}
			if archive["size"].(int) != 0 {
				configSet = append(configSet, setPrefixArchive+" size "+
					strconv.Itoa(archive["size"].(int)))
			}
			if archive["start_time"].(string) != "" {
				configSet = append(configSet, setPrefixArchive+" start-time "+
					archive["start_time"].(string))
			}
			if archive["transfer_interval"].(int) != 0 {
				configSet = append(configSet, setPrefixArchive+" transfer-interval "+
					strconv.Itoa(archive["transfer_interval"].(int)))
			}
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readSystemSyslogFile(filename string, m interface{}, jnprSess *NetconfObject) (syslogFileOptions, error) {
	sess := m.(*Session)
	var confRead syslogFileOptions

	showConfig, err := sess.command("show configuration system syslog file "+filename+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.filename = filename
		for _, item := range strings.Split(showConfig, "\n") {
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
			case itemTrim == "explicit-priority":
				confRead.explicitPriority = true
			case strings.HasPrefix(itemTrim, "match "):
				confRead.match = strings.Trim(strings.TrimPrefix(itemTrim, "match "), "\"")
			case strings.HasPrefix(itemTrim, "match-strings "):
				confRead.matchStrings = append(confRead.matchStrings,
					strings.Trim(strings.TrimPrefix(itemTrim, "match-strings "), "\""))
			case strings.HasPrefix(itemTrim, "structured-data"):
				if len(confRead.structuredData) == 0 {
					confRead.structuredData = append(confRead.structuredData, map[string]interface{}{
						"brief": false,
					})
				}
				if itemTrim == "structured-data brief" {
					confRead.structuredData[0]["brief"] = true
				}
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
			case strings.HasPrefix(itemTrim, "archive"):
				if len(confRead.archive) == 0 {
					confRead.archive = append(confRead.archive, map[string]interface{}{
						"sites":             make([]map[string]interface{}, 0),
						"binary_data":       false,
						"no_binary_data":    false,
						"world_readable":    false,
						"no_world_readable": false,
						"files":             0,
						"size":              0,
						"transfer_interval": 0,
						"start_time":        "",
					})
				}
				switch {
				case itemTrim == "archive binary-data":
					confRead.archive[0]["binary_data"] = true
				case itemTrim == "archive no-binary-data":
					confRead.archive[0]["no_binary_data"] = true
				case itemTrim == "archive world-readable":
					confRead.archive[0]["world_readable"] = true
				case itemTrim == "archive no-world-readable":
					confRead.archive[0]["no_world_readable"] = true
				case strings.HasPrefix(itemTrim, "archive files "):
					var err error
					confRead.archive[0]["files"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "archive files "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "archive size "):
					var err error
					confRead.archive[0]["size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "archive size "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "archive transfer-interval "):
					var err error
					confRead.archive[0]["transfer_interval"], err =
						strconv.Atoi(strings.TrimPrefix(itemTrim, "archive transfer-interval "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "archive start-time "):
					confRead.archive[0]["start_time"] = strings.TrimPrefix(itemTrim, "archive start-time ")
				case strings.HasPrefix(itemTrim, "archive archive-sites "):
					itemTrimArchSitesSplit := strings.Split(
						strings.TrimPrefix(itemTrim, "archive archive-sites "), " ")
					itemTrimArchSites := strings.TrimPrefix(itemTrim, "archive archive-sites "+itemTrimArchSitesSplit[0]+" ")
					sitesOptions := map[string]interface{}{
						"url":              itemTrimArchSitesSplit[0],
						"password":         "",
						"routing_instance": "",
					}
					confRead.archive[0]["sites"] = copyAndRemoveItemMapList("url", sitesOptions,
						confRead.archive[0]["sites"].([]map[string]interface{}))
					switch {
					case strings.HasPrefix(itemTrimArchSites, "password "):
						var err error
						sitesOptions["password"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
							itemTrimArchSites, "password "), "\""))
						if err != nil {
							return confRead, fmt.Errorf("failed to decode password : %w", err)
						}
					case strings.HasPrefix(itemTrimArchSites, "routing-instance "):
						sitesOptions["routing_instance"] = strings.TrimPrefix(itemTrimArchSites, "routing-instance ")
					}
					confRead.archive[0]["sites"] = append(confRead.archive[0]["sites"].([]map[string]interface{}), sitesOptions)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemSyslogFile(filename string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system syslog file "+filename)

	return sess.configSet(configSet, jnprSess)
}

func fillSystemSyslogFileData(d *schema.ResourceData, syslogFileOptions syslogFileOptions) {
	if tfErr := d.Set("filename", syslogFileOptions.filename); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_duplicates", syslogFileOptions.allowDuplicates); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("explicit_priority", syslogFileOptions.explicitPriority); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("match", syslogFileOptions.match); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("match_strings", syslogFileOptions.matchStrings); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("structured_data", syslogFileOptions.structuredData); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("any_severity", syslogFileOptions.anySeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authorization_severity", syslogFileOptions.authorizationSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("changelog_severity", syslogFileOptions.changelogSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("conflictlog_severity", syslogFileOptions.conflictlogSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("daemon_severity", syslogFileOptions.daemonSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dfc_severity", syslogFileOptions.dfcSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("external_severity", syslogFileOptions.externalSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("firewall_severity", syslogFileOptions.firewallSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ftp_severity", syslogFileOptions.ftpSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interactivecommands_severity", syslogFileOptions.interactivecommandsSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("kernel_severity", syslogFileOptions.kernelSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ntp_severity", syslogFileOptions.ntpSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pfe_severity", syslogFileOptions.pfeSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_severity", syslogFileOptions.securitySeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_severity", syslogFileOptions.userSeverity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("archive", syslogFileOptions.archive); tfErr != nil {
		panic(tfErr)
	}
}
