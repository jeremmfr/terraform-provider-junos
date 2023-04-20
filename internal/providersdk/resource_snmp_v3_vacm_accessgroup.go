package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type snmpV3VacmAccessGroupOptions struct {
	name                 string
	defaultContextPrefix []map[string]interface{}
	contextPrefix        []map[string]interface{}
}

func resourceSnmpV3VacmAccessGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSnmpV3VacmAccessGroupCreate,
		ReadWithoutTimeout:   resourceSnmpV3VacmAccessGroupRead,
		UpdateWithoutTimeout: resourceSnmpV3VacmAccessGroupUpdate,
		DeleteWithoutTimeout: resourceSnmpV3VacmAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnmpV3VacmAccessGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"context_prefix": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"context_prefix", "default_context_prefix"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_config": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"any", "usm", "v1", "v2c"}, false),
									},
									"level": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"authentication", "none", "privacy"}, false),
									},
									"context_match": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringInSlice([]string{"exact", "prefix"}, false),
									},
									"notify_view": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"read_view": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"write_view": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
								},
							},
						},
					},
				},
			},
			"default_context_prefix": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"context_prefix", "default_context_prefix"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"any", "usm", "v1", "v2c"}, false),
						},
						"level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"authentication", "none", "privacy"}, false),
						},
						"context_match": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validation.StringInSlice([]string{"exact", "prefix"}, false),
						},
						"notify_view": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"read_view": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"write_view": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
	}
}

func resourceSnmpV3VacmAccessGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSnmpV3VacmAccessGroup(d, junSess); err != nil {
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
	snmpV3VacmAccessGroupExists, err := checkSnmpV3VacmAccessGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3VacmAccessGroupExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"snmp v3 vacm access group %v already exists", d.Get("name").(string)))...)
	}

	if err := setSnmpV3VacmAccessGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_snmp_v3_vacm_accessgroup")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	snmpV3VacmAccessGroupExists, err = checkSnmpV3VacmAccessGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if snmpV3VacmAccessGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("snmp v3 vacm access group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSnmpV3VacmAccessGroupReadWJunSess(d, junSess)...)
}

func resourceSnmpV3VacmAccessGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSnmpV3VacmAccessGroupReadWJunSess(d, junSess)
}

func resourceSnmpV3VacmAccessGroupReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	snmpV3VacmAccessGroupOptions, err := readSnmpV3VacmAccessGroup(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if snmpV3VacmAccessGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillSnmpV3VacmAccessGroupData(d, snmpV3VacmAccessGroupOptions)
	}

	return nil
}

func resourceSnmpV3VacmAccessGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSnmpV3VacmAccessGroup(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSnmpV3VacmAccessGroup(d, junSess); err != nil {
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
	if err := delSnmpV3VacmAccessGroup(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSnmpV3VacmAccessGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_snmp_v3_vacm_accessgroup")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSnmpV3VacmAccessGroupReadWJunSess(d, junSess)...)
}

func resourceSnmpV3VacmAccessGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSnmpV3VacmAccessGroup(d.Get("name").(string), junSess); err != nil {
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
	if err := delSnmpV3VacmAccessGroup(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_snmp_v3_vacm_accessgroup")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSnmpV3VacmAccessGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	snmpV3VacmAccessGroupExists, err := checkSnmpV3VacmAccessGroupExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !snmpV3VacmAccessGroupExists {
		return nil, fmt.Errorf("don't find snmp v3 vacm access group with id '%v' (id must be <name>)", d.Id())
	}
	snmpV3VacmAccessGroupOptions, err := readSnmpV3VacmAccessGroup(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSnmpV3VacmAccessGroupData(d, snmpV3VacmAccessGroupOptions)

	result[0] = d

	return result, nil
}

func checkSnmpV3VacmAccessGroupExists(name string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm access group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSnmpV3VacmAccessGroup(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set snmp v3 vacm access group \"" + d.Get("name").(string) + "\" "
	configSet := make([]string, 0)

	defaultContextPrefixList := make([]string, 0)
	for _, mDefCtxPref := range d.Get("default_context_prefix").(*schema.Set).List() {
		defaultContextPrefix := mDefCtxPref.(map[string]interface{})
		if bchk.InSlice(defaultContextPrefix["model"].(string)+junos.IDSeparator+defaultContextPrefix["level"].(string),
			defaultContextPrefixList) {
			return fmt.Errorf("multiple blocks default_context_prefix with the same model '%s' and level '%s'",
				defaultContextPrefix["model"].(string), defaultContextPrefix["level"].(string))
		}
		defaultContextPrefixList = append(defaultContextPrefixList,
			defaultContextPrefix["model"].(string)+junos.IDSeparator+defaultContextPrefix["level"].(string))
		setPrefixDefCtxPref := setPrefix + " default-context-prefix security-model " +
			defaultContextPrefix["model"].(string) + " security-level " + defaultContextPrefix["level"].(string) + " "
		if v := defaultContextPrefix["context_match"].(string); v != "" {
			configSet = append(configSet, setPrefixDefCtxPref+"context-match "+v)
		}
		if v := defaultContextPrefix["notify_view"].(string); v != "" {
			configSet = append(configSet, setPrefixDefCtxPref+"notify-view \""+v+"\"")
		}
		if v := defaultContextPrefix["read_view"].(string); v != "" {
			configSet = append(configSet, setPrefixDefCtxPref+"read-view \""+v+"\"")
		}
		if v := defaultContextPrefix["write_view"].(string); v != "" {
			configSet = append(configSet, setPrefixDefCtxPref+"write-view \""+v+"\"")
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixDefCtxPref) {
			return fmt.Errorf("missing argument to default_context_prefix with model %s and level %s",
				defaultContextPrefix["model"].(string), defaultContextPrefix["level"].(string))
		}
	}
	contextPrefixList := make([]string, 0)
	for _, mCtxPref := range d.Get("context_prefix").([]interface{}) {
		contextPrefix := mCtxPref.(map[string]interface{})
		if bchk.InSlice(contextPrefix["prefix"].(string), contextPrefixList) {
			return fmt.Errorf("multiple blocks context_prefix with the same prefix '%s'", contextPrefix["prefix"].(string))
		}
		contextPrefixList = append(contextPrefixList, contextPrefix["prefix"].(string))
		accessConfigList := make([]string, 0)
		for _, mAccConf := range contextPrefix["access_config"].(*schema.Set).List() {
			accessConfig := mAccConf.(map[string]interface{})
			if bchk.InSlice(accessConfig["model"].(string)+junos.IDSeparator+accessConfig["level"].(string), accessConfigList) {
				return fmt.Errorf(
					"multiple blocks access_config with the same model '%s' and level '%s' in context_prefix with prefix '%s'",
					accessConfig["model"].(string), accessConfig["level"].(string), contextPrefix["prefix"].(string))
			}
			accessConfigList = append(accessConfigList,
				accessConfig["model"].(string)+junos.IDSeparator+accessConfig["level"].(string))
			setPrefixCtxPref := setPrefix + " context-prefix \"" + contextPrefix["prefix"].(string) +
				"\" security-model " + accessConfig["model"].(string) + " security-level " + accessConfig["level"].(string) + " "
			if v := accessConfig["context_match"].(string); v != "" {
				configSet = append(configSet, setPrefixCtxPref+"context-match "+v)
			}
			if v := accessConfig["notify_view"].(string); v != "" {
				configSet = append(configSet, setPrefixCtxPref+"notify-view \""+v+"\"")
			}
			if v := accessConfig["read_view"].(string); v != "" {
				configSet = append(configSet, setPrefixCtxPref+"read-view \""+v+"\"")
			}
			if v := accessConfig["write_view"].(string); v != "" {
				configSet = append(configSet, setPrefixCtxPref+"write-view \""+v+"\"")
			}
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixCtxPref) {
				return fmt.Errorf("missing argument to access_config with model %s and level %s in context_prefix with prefix %s",
					accessConfig["model"].(string), accessConfig["level"].(string), contextPrefix["prefix"].(string))
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func readSnmpV3VacmAccessGroup(name string, junSess *junos.Session,
) (confRead snmpV3VacmAccessGroupOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm access group \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "default-context-prefix security-model "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <model> security-level <level>
					return confRead,
						fmt.Errorf(junos.CantReadValuesNotEnoughFields, "default-context-prefix security-model", itemTrim)
				}
				defaultContextPrefix := map[string]interface{}{
					"model":         itemTrimFields[0],
					"level":         itemTrimFields[2],
					"context_match": "",
					"notify_view":   "",
					"read_view":     "",
					"write_view":    "",
				}
				confRead.defaultContextPrefix = copyAndRemoveItemMapList2("model", "level", defaultContextPrefix,
					confRead.defaultContextPrefix)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" security-level "+itemTrimFields[2]+" ")
				readSnmpV3VacmAccessGroupContextPrefixConfig(itemTrim, defaultContextPrefix)
				confRead.defaultContextPrefix = append(confRead.defaultContextPrefix, defaultContextPrefix)
			case balt.CutPrefixInString(&itemTrim, "context-prefix "):
				itemTrimFields := strings.Split(itemTrim, " ")
				contextPrefix := map[string]interface{}{
					"prefix":        strings.Trim(itemTrimFields[0], "\""),
					"access_config": make([]map[string]interface{}, 0),
				}
				confRead.contextPrefix = copyAndRemoveItemMapList("prefix", contextPrefix, confRead.contextPrefix)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				if balt.CutPrefixInString(&itemTrim, "security-model ") {
					itemTrimSecModelFields := strings.Split(itemTrim, " ")
					if len(itemTrimSecModelFields) < 3 { // <model> security-level <level>
						return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "security-model", itemTrim)
					}
					contextPrefixAccessConfig := map[string]interface{}{
						"model":         itemTrimSecModelFields[0],
						"level":         itemTrimSecModelFields[2],
						"context_match": "",
						"notify_view":   "",
						"read_view":     "",
						"write_view":    "",
					}
					contextPrefix["access_config"] = copyAndRemoveItemMapList2("model", "level", contextPrefixAccessConfig,
						contextPrefix["access_config"].([]map[string]interface{}))
					balt.CutPrefixInString(&itemTrim, itemTrimSecModelFields[0]+" security-level "+itemTrimSecModelFields[2]+" ")
					readSnmpV3VacmAccessGroupContextPrefixConfig(itemTrim, contextPrefixAccessConfig)
					contextPrefix["access_config"] = append(
						contextPrefix["access_config"].([]map[string]interface{}),
						contextPrefixAccessConfig,
					)
				}
				confRead.contextPrefix = append(confRead.contextPrefix, contextPrefix)
			}
		}
	}

	return confRead, nil
}

func readSnmpV3VacmAccessGroupContextPrefixConfig(itemTrim string, config map[string]interface{}) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "context-match "):
		config["context_match"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "notify-view "):
		config["notify_view"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "read-view "):
		config["read_view"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "write-view "):
		config["write_view"] = strings.Trim(itemTrim, "\"")
	}
}

func delSnmpV3VacmAccessGroup(name string, junSess *junos.Session) error {
	configSet := []string{"delete snmp v3 vacm access group \"" + name + "\""}

	return junSess.ConfigSet(configSet)
}

func fillSnmpV3VacmAccessGroupData(d *schema.ResourceData, snmpV3VacmAccessGroupOptions snmpV3VacmAccessGroupOptions) {
	if tfErr := d.Set("name", snmpV3VacmAccessGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_context_prefix", snmpV3VacmAccessGroupOptions.defaultContextPrefix); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("context_prefix", snmpV3VacmAccessGroupOptions.contextPrefix); tfErr != nil {
		panic(tfErr)
	}
}
