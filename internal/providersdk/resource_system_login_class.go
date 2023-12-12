package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type systemLoginClassOptions struct {
	allowHiddenCommands       bool
	configurationBreadcrumbs  bool
	loginAlarms               bool
	loginTip                  bool
	idleTimeout               int
	name                      string
	accessEnd                 string
	accessStart               string
	allowCommands             string
	allowConfiguration        string
	cliPrompt                 string
	denyCommands              string
	denyConfiguration         string
	logicalSystem             string
	loginScript               string
	securityRole              string
	tenant                    string
	confirmCommands           []string
	allowCommandsRegexps      []string
	allowConfigurationRegexps []string
	allowedDays               []string
	denyCommandsRegexps       []string
	denyConfigurationRegexps  []string
	noHiddenCommandsExcept    []string
	permissions               []string
}

func resourceSystemLoginClass() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemLoginClassCreate,
		ReadWithoutTimeout:   resourceSystemLoginClassRead,
		UpdateWithoutTimeout: resourceSystemLoginClassUpdate,
		DeleteWithoutTimeout: resourceSystemLoginClassDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemLoginClassImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"access_end": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_start"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^([0-1]\d|2[0-3]):([0-5]\d):([0-5]\d)$`), "must have HH:MM:SS format"),
			},
			"access_start": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_end"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(
					`^([0-1]\d|2[0-3]):([0-5]\d):([0-5]\d)$`), "must have HH:MM:SS format"),
			},
			"allow_commands": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"allow_commands_regexps"},
			},
			"allow_commands_regexps": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"allow_commands"},
			},
			"allow_configuration": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"allow_configuration_regexps"},
			},
			"allow_configuration_regexps": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"allow_configuration"},
			},
			"allow_hidden_commands": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_hidden_commands_except"},
			},
			"allowed_days": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cli_prompt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_breadcrumbs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"confirm_commands": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"deny_commands": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"deny_commands_regexps"},
			},
			"deny_commands_regexps": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"deny_commands"},
			},
			"deny_configuration": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"deny_configuration_regexps"},
			},
			"deny_configuration_regexps": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"deny_configuration"},
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4294967295),
			},
			"logical_system": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_alarms": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"login_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_tip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_hidden_commands_except": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"allow_hidden_commands"},
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_role": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"audit-administrator", "crypto-administrator", "ids-administrator", "security-administrator",
				}, false),
			},
			"tenant": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSystemLoginClassCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSystemLoginClass(d, junSess); err != nil {
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
	systemLoginClassExists, err := checkSystemLoginClassExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginClassExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("system login class %v already exists", d.Get("name").(string)))...)
	}

	if err := setSystemLoginClass(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_system_login_class")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	systemLoginClassExists, err = checkSystemLoginClassExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginClassExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system login class %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSystemLoginClassReadWJunSess(d, junSess)...)
}

func resourceSystemLoginClassRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSystemLoginClassReadWJunSess(d, junSess)
}

func resourceSystemLoginClassReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	systemLoginClassOptions, err := readSystemLoginClass(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if systemLoginClassOptions.name == "" {
		d.SetId("")
	} else {
		fillSystemLoginClassData(d, systemLoginClassOptions)
	}

	return nil
}

func resourceSystemLoginClassUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemLoginClass(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemLoginClass(d, junSess); err != nil {
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
	if err := delSystemLoginClass(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemLoginClass(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_system_login_class")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemLoginClassReadWJunSess(d, junSess)...)
}

func resourceSystemLoginClassDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSystemLoginClass(d.Get("name").(string), junSess); err != nil {
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
	if err := delSystemLoginClass(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_system_login_class")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemLoginClassImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	systemLoginClassExists, err := checkSystemLoginClassExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !systemLoginClassExists {
		return nil, fmt.Errorf("don't find system login class with id '%v' (id must be <name>)", d.Id())
	}
	systemLoginClassOptions, err := readSystemLoginClass(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSystemLoginClassData(d, systemLoginClassOptions)

	result[0] = d

	return result, nil
}

func checkSystemLoginClassExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system login class " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSystemLoginClass(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := "set system login class " + d.Get("name").(string) + " "

	if d.Get("access_end").(string) != "" {
		configSet = append(configSet, setPrefix+"access-end \""+d.Get("access_end").(string)+"\"")
	}
	if d.Get("access_start").(string) != "" {
		configSet = append(configSet, setPrefix+"access-start \""+d.Get("access_start").(string)+"\"")
	}
	if d.Get("allow_commands").(string) != "" {
		configSet = append(configSet, setPrefix+"allow-commands \""+d.Get("allow_commands").(string)+"\"")
	}
	for _, v := range d.Get("allow_commands_regexps").([]interface{}) {
		configSet = append(configSet, setPrefix+"allow-commands-regexps \""+v.(string)+"\"")
	}
	if d.Get("allow_configuration").(string) != "" {
		configSet = append(configSet, setPrefix+"allow-configuration \""+d.Get("allow_configuration").(string)+"\"")
	}
	for _, v := range d.Get("allow_configuration_regexps").([]interface{}) {
		configSet = append(configSet, setPrefix+"allow-configuration-regexps \""+v.(string)+"\"")
	}
	if d.Get("allow_hidden_commands").(bool) {
		configSet = append(configSet, setPrefix+"allow-hidden-commands")
	}
	for _, v := range d.Get("allowed_days").([]interface{}) {
		configSet = append(configSet, setPrefix+"allowed-days "+v.(string))
	}
	if d.Get("configuration_breadcrumbs").(bool) {
		configSet = append(configSet, setPrefix+"configuration-breadcrumbs")
	}
	if d.Get("cli_prompt").(string) != "" {
		configSet = append(configSet, setPrefix+"cli prompt \""+d.Get("cli_prompt").(string)+"\"")
	}
	for _, v := range d.Get("confirm_commands").([]interface{}) {
		configSet = append(configSet, setPrefix+"confirm-commands \""+v.(string)+"\"")
	}
	if d.Get("deny_commands").(string) != "" {
		configSet = append(configSet, setPrefix+"deny-commands \""+d.Get("deny_commands").(string)+"\"")
	}
	for _, v := range d.Get("deny_commands_regexps").([]interface{}) {
		configSet = append(configSet, setPrefix+"deny-commands-regexps \""+v.(string)+"\"")
	}
	if d.Get("deny_configuration").(string) != "" {
		configSet = append(configSet, setPrefix+"deny-configuration \""+d.Get("deny_configuration").(string)+"\"")
	}
	for _, v := range d.Get("deny_configuration_regexps").([]interface{}) {
		configSet = append(configSet, setPrefix+"deny-configuration-regexps \""+v.(string)+"\"")
	}
	if d.Get("idle_timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"idle-timeout "+strconv.Itoa(d.Get("idle_timeout").(int)))
	}
	if d.Get("logical_system").(string) != "" {
		configSet = append(configSet, setPrefix+"logical-system \""+d.Get("logical_system").(string)+"\"")
	}
	if d.Get("login_alarms").(bool) {
		configSet = append(configSet, setPrefix+"login-alarms")
	}
	if d.Get("login_script").(string) != "" {
		configSet = append(configSet, setPrefix+"login-script "+d.Get("login_script").(string))
	}
	if d.Get("login_tip").(bool) {
		configSet = append(configSet, setPrefix+"login-tip")
	}
	for _, v := range d.Get("no_hidden_commands_except").([]interface{}) {
		configSet = append(configSet, setPrefix+"no-hidden-commands except \""+v.(string)+"\"")
	}
	for _, v := range sortSetOfString(d.Get("permissions").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"permissions "+v)
	}
	if d.Get("security_role").(string) != "" {
		configSet = append(configSet, setPrefix+"security-role "+d.Get("security_role").(string))
	}
	if d.Get("tenant").(string) != "" {
		configSet = append(configSet, setPrefix+"tenant \""+d.Get("tenant").(string)+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readSystemLoginClass(name string, junSess *junos.Session,
) (confRead systemLoginClassOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "system login class " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "access-end "):
				confRead.accessEnd = strings.Split(strings.Trim(itemTrim, "\""), " ")[0]
			case balt.CutPrefixInString(&itemTrim, "access-start "):
				confRead.accessStart = strings.Split(strings.Trim(itemTrim, "\""), " ")[0]
			case balt.CutPrefixInString(&itemTrim, "allow-commands "):
				confRead.allowCommands = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "allow-commands-regexps "):
				confRead.allowCommandsRegexps = append(confRead.allowCommandsRegexps, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "allow-configuration "):
				confRead.allowConfiguration = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "allow-configuration-regexps "):
				confRead.allowConfigurationRegexps = append(confRead.allowConfigurationRegexps, strings.Trim(itemTrim, "\""))
			case itemTrim == "allow-hidden-commands":
				confRead.allowHiddenCommands = true
			case balt.CutPrefixInString(&itemTrim, "allowed-days "):
				confRead.allowedDays = append(confRead.allowedDays, itemTrim)
			case itemTrim == "configuration-breadcrumbs":
				confRead.configurationBreadcrumbs = true
			case balt.CutPrefixInString(&itemTrim, "cli prompt "):
				confRead.cliPrompt = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "confirm-commands "):
				confRead.confirmCommands = append(confRead.confirmCommands, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "deny-commands "):
				confRead.denyCommands = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "deny-commands-regexps "):
				confRead.denyCommandsRegexps = append(confRead.denyCommandsRegexps, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "deny-configuration "):
				confRead.denyConfiguration = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "deny-configuration-regexps "):
				confRead.denyConfigurationRegexps = append(confRead.denyConfigurationRegexps, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "idle-timeout "):
				confRead.idleTimeout, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "logical-system "):
				confRead.logicalSystem = strings.Trim(itemTrim, "\"")
			case itemTrim == "login-alarms":
				confRead.loginAlarms = true
			case balt.CutPrefixInString(&itemTrim, "login-script "):
				confRead.loginScript = itemTrim
			case itemTrim == "login-tip":
				confRead.loginTip = true
			case balt.CutPrefixInString(&itemTrim, "no-hidden-commands except "):
				confRead.noHiddenCommandsExcept = append(confRead.noHiddenCommandsExcept, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "permissions "):
				confRead.permissions = append(confRead.permissions, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "security-role "):
				confRead.securityRole = itemTrim
			case balt.CutPrefixInString(&itemTrim, "tenant "):
				confRead.tenant = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delSystemLoginClass(systemLoginClass string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system login class "+systemLoginClass)

	return junSess.ConfigSet(configSet)
}

func fillSystemLoginClassData(d *schema.ResourceData, systemLoginClassOptions systemLoginClassOptions) {
	if tfErr := d.Set("name", systemLoginClassOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("access_end", systemLoginClassOptions.accessEnd); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("access_start", systemLoginClassOptions.accessStart); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_commands", systemLoginClassOptions.allowCommands); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_commands_regexps", systemLoginClassOptions.allowCommandsRegexps); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_configuration", systemLoginClassOptions.allowConfiguration); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_configuration_regexps", systemLoginClassOptions.allowConfigurationRegexps); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allow_hidden_commands", systemLoginClassOptions.allowHiddenCommands); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("allowed_days", systemLoginClassOptions.allowedDays); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cli_prompt", systemLoginClassOptions.cliPrompt); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("configuration_breadcrumbs", systemLoginClassOptions.configurationBreadcrumbs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("confirm_commands", systemLoginClassOptions.confirmCommands); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("deny_commands", systemLoginClassOptions.denyCommands); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("deny_commands_regexps", systemLoginClassOptions.denyCommandsRegexps); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("deny_configuration", systemLoginClassOptions.denyConfiguration); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("deny_configuration_regexps", systemLoginClassOptions.denyConfigurationRegexps); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("idle_timeout", systemLoginClassOptions.idleTimeout); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("logical_system", systemLoginClassOptions.logicalSystem); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("login_alarms", systemLoginClassOptions.loginAlarms); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("login_script", systemLoginClassOptions.loginScript); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("login_tip", systemLoginClassOptions.loginTip); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_hidden_commands_except", systemLoginClassOptions.noHiddenCommandsExcept); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("permissions", systemLoginClassOptions.permissions); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_role", systemLoginClassOptions.securityRole); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tenant", systemLoginClassOptions.tenant); tfErr != nil {
		panic(tfErr)
	}
}
