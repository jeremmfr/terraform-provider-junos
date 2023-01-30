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
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceSystemSyslogFileCreate,
		ReadWithoutTimeout:   resourceSystemSyslogFileRead,
		UpdateWithoutTimeout: resourceSystemSyslogFileUpdate,
		DeleteWithoutTimeout: resourceSystemSyslogFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemSyslogFileImport,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSystemSyslogFile(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("filename").(string))

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
	syslogFileExists, err := checkSystemSyslogFileExists(d.Get("filename").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogFileExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system syslog file %v already exists", d.Get("filename").(string)))...)
	}

	if err := setSystemSyslogFile(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_system_syslog_file")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	syslogFileExists, err = checkSystemSyslogFileExists(d.Get("filename").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if syslogFileExists {
		d.SetId(d.Get("filename").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system syslog file %v not exists after commit "+
			"=> check your config", d.Get("filename").(string)))...)
	}

	return append(diagWarns, resourceSystemSyslogFileReadWJunSess(d, junSess)...)
}

func resourceSystemSyslogFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemSyslogFileReadWJunSess(d, junSess)
}

func resourceSystemSyslogFileReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	syslogFileOptions, err := readSystemSyslogFile(d.Get("filename").(string), junSess)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemSyslogFile(d.Get("filename").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemSyslogFile(d, junSess); err != nil {
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
	if err := delSystemSyslogFile(d.Get("filename").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemSyslogFile(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_system_syslog_file")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemSyslogFileReadWJunSess(d, junSess)...)
}

func resourceSystemSyslogFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemSyslogFile(d.Get("filename").(string), junSess); err != nil {
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
	if err := delSystemSyslogFile(d.Get("filename").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_system_syslog_file")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemSyslogFileImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	syslogFileExists, err := checkSystemSyslogFileExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !syslogFileExists {
		return nil, fmt.Errorf("don't find system syslog file with id '%v' (id must be <filename>)", d.Id())
	}
	syslogFileOptions, err := readSystemSyslogFile(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSystemSyslogFileData(d, syslogFileOptions)

	result[0] = d

	return result, nil
}

func checkSystemSyslogFileExists(filename string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system syslog file " + filename + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemSyslogFile(d *schema.ResourceData, junSess *junos.Session) error {
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
				if bchk.InSlice(sites["url"].(string), sitesURLList) {
					return fmt.Errorf("multiple blocks sites with the same url %s", sites["url"].(string))
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

	return junSess.ConfigSet(configSet)
}

func readSystemSyslogFile(filename string, junSess *junos.Session,
) (confRead syslogFileOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog file " + filename + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.filename = filename
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
			case itemTrim == "explicit-priority":
				confRead.explicitPriority = true
			case balt.CutPrefixInString(&itemTrim, "match "):
				confRead.match = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "match-strings "):
				confRead.matchStrings = append(confRead.matchStrings, strings.Trim(itemTrim, "\""))
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
			case balt.CutPrefixInString(&itemTrim, "archive"):
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
				case itemTrim == " binary-data":
					confRead.archive[0]["binary_data"] = true
				case itemTrim == " no-binary-data":
					confRead.archive[0]["no_binary_data"] = true
				case itemTrim == " world-readable":
					confRead.archive[0]["world_readable"] = true
				case itemTrim == " no-world-readable":
					confRead.archive[0]["no_world_readable"] = true
				case balt.CutPrefixInString(&itemTrim, " files "):
					confRead.archive[0]["files"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " size "):
					confRead.archive[0]["size"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " transfer-interval "):
					confRead.archive[0]["transfer_interval"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, " start-time "):
					confRead.archive[0]["start_time"] = strings.Split(strings.Trim(itemTrim, "\""), " ")[0]
				case balt.CutPrefixInString(&itemTrim, " archive-sites "):
					itemTrimFields := strings.Split(itemTrim, " ")
					sitesOptions := map[string]interface{}{
						"url":              itemTrimFields[0],
						"password":         "",
						"routing_instance": "",
					}
					confRead.archive[0]["sites"] = copyAndRemoveItemMapList("url", sitesOptions,
						confRead.archive[0]["sites"].([]map[string]interface{}))
					balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
					switch {
					case balt.CutPrefixInString(&itemTrim, "password "):
						sitesOptions["password"], err = jdecode.Decode(strings.Trim(itemTrim, "\""))
						if err != nil {
							return confRead, fmt.Errorf("decoding password: %w", err)
						}
					case balt.CutPrefixInString(&itemTrim, "routing-instance "):
						sitesOptions["routing_instance"] = itemTrim
					}
					confRead.archive[0]["sites"] = append(confRead.archive[0]["sites"].([]map[string]interface{}), sitesOptions)
				}
			}
		}
	}

	return confRead, nil
}

func delSystemSyslogFile(filename string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system syslog file "+filename)

	return junSess.ConfigSet(configSet)
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
