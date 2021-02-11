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
		CreateContext: resourceSystemLoginClassCreate,
		ReadContext:   resourceSystemLoginClassRead,
		UpdateContext: resourceSystemLoginClassUpdate,
		DeleteContext: resourceSystemLoginClassDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemLoginClassImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"access_end": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_end"},
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
					"audit-administrator", "crypto-administrator", "ids-administrator", "security-administrator"}, false),
			},
			"tenant": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSystemLoginClassCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	systemLoginClassExists, err := checkSystemLoginClassExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if systemLoginClassExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("system login class %v already exists", d.Get("name").(string)))
	}

	if err := setSystemLoginClass(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_system_login_class", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	systemLoginClassExists, err = checkSystemLoginClassExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemLoginClassExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("system login class %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSystemLoginClassReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemLoginClassRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemLoginClassReadWJnprSess(d, m, jnprSess)
}
func resourceSystemLoginClassReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	systemLoginClassOptions, err := readSystemLoginClass(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemLoginClass(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystemLoginClass(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_system_login_class", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemLoginClassReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemLoginClassDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystemLoginClass(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_system_login_class", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSystemLoginClassImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	systemLoginClassExists, err := checkSystemLoginClassExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !systemLoginClassExists {
		return nil, fmt.Errorf("don't find system login class with id '%v' (id must be <name>)", d.Id())
	}
	systemLoginClassOptions, err := readSystemLoginClass(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemLoginClassData(d, systemLoginClassOptions)

	result[0] = d

	return result, nil
}

func checkSystemLoginClassExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	systemLoginClassConfig, err := sess.command("show configuration system login class "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if systemLoginClassConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSystemLoginClass(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
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
	for _, v := range d.Get("permissions").(*schema.Set).List() {
		configSet = append(configSet, setPrefix+"permissions "+v.(string))
	}
	if d.Get("security_role").(string) != "" {
		configSet = append(configSet, setPrefix+"security-role "+d.Get("security_role").(string))
	}
	if d.Get("tenant").(string) != "" {
		configSet = append(configSet, setPrefix+"tenant \""+d.Get("tenant").(string)+"\"")
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystemLoginClass(class string, m interface{}, jnprSess *NetconfObject) (systemLoginClassOptions, error) {
	sess := m.(*Session)
	var confRead systemLoginClassOptions

	systemLoginClassConfig, err := sess.command("show configuration system login class "+
		class+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if systemLoginClassConfig != emptyWord {
		confRead.name = class
		for _, item := range strings.Split(systemLoginClassConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "access-end "):
				accessSplit := strings.Split(strings.Trim(strings.TrimPrefix(itemTrim, "access-end "), "\""), " ")
				confRead.accessEnd = accessSplit[0]
			case strings.HasPrefix(itemTrim, "access-start "):
				accessSplit := strings.Split(strings.Trim(strings.TrimPrefix(itemTrim, "access-start "), "\""), " ")
				confRead.accessStart = accessSplit[0]
			case strings.HasPrefix(itemTrim, "allow-commands "):
				confRead.allowCommands = strings.Trim(strings.TrimPrefix(itemTrim, "allow-commands "), "\"")
			case strings.HasPrefix(itemTrim, "allow-commands-regexps "):
				confRead.allowCommandsRegexps = append(confRead.allowCommandsRegexps,
					strings.Trim(strings.TrimPrefix(itemTrim, "allow-commands-regexps "), "\""))
			case strings.HasPrefix(itemTrim, "allow-configuration "):
				confRead.allowConfiguration = strings.Trim(strings.TrimPrefix(itemTrim, "allow-configuration "), "\"")
			case strings.HasPrefix(itemTrim, "allow-configuration-regexps "):
				confRead.allowConfigurationRegexps = append(confRead.allowConfigurationRegexps,
					strings.Trim(strings.TrimPrefix(itemTrim, "allow-configuration-regexps "), "\""))
			case itemTrim == "allow-hidden-commands":
				confRead.allowHiddenCommands = true
			case strings.HasPrefix(itemTrim, "allowed-days "):
				confRead.allowedDays = append(confRead.allowedDays, strings.TrimPrefix(itemTrim, "allowed-days "))
			case itemTrim == "configuration-breadcrumbs":
				confRead.configurationBreadcrumbs = true
			case strings.HasPrefix(itemTrim, "cli prompt "):
				confRead.cliPrompt = strings.Trim(strings.TrimPrefix(itemTrim, "cli prompt "), "\"")
			case strings.HasPrefix(itemTrim, "confirm-commands "):
				confRead.confirmCommands = append(confRead.confirmCommands,
					strings.Trim(strings.TrimPrefix(itemTrim, "confirm-commands "), "\""))
			case strings.HasPrefix(itemTrim, "deny-commands "):
				confRead.denyCommands = strings.Trim(strings.TrimPrefix(itemTrim, "deny-commands "), "\"")
			case strings.HasPrefix(itemTrim, "deny-commands-regexps "):
				confRead.denyCommandsRegexps = append(confRead.denyCommandsRegexps,
					strings.Trim(strings.TrimPrefix(itemTrim, "deny-commands-regexps "), "\""))
			case strings.HasPrefix(itemTrim, "deny-configuration "):
				confRead.denyConfiguration = strings.Trim(strings.TrimPrefix(itemTrim, "deny-configuration "), "\"")
			case strings.HasPrefix(itemTrim, "deny-configuration-regexps "):
				confRead.denyConfigurationRegexps = append(confRead.denyConfigurationRegexps,
					strings.Trim(strings.TrimPrefix(itemTrim, "deny-configuration-regexps "), "\""))
			case strings.HasPrefix(itemTrim, "idle-timeout "):
				var err error
				confRead.idleTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "idle-timeout "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "logical-system "):
				confRead.logicalSystem = strings.Trim(strings.TrimPrefix(itemTrim, "logical-system "), "\"")
			case itemTrim == "login-alarms":
				confRead.loginAlarms = true
			case strings.HasPrefix(itemTrim, "login-script "):
				confRead.loginScript = strings.TrimPrefix(itemTrim, "login-script ")
			case itemTrim == "login-tip":
				confRead.loginTip = true
			case strings.HasPrefix(itemTrim, "no-hidden-commands except "):
				confRead.noHiddenCommandsExcept = append(confRead.noHiddenCommandsExcept,
					strings.Trim(strings.TrimPrefix(itemTrim, "no-hidden-commands except "), "\""))
			case strings.HasPrefix(itemTrim, "permissions "):
				confRead.permissions = append(confRead.permissions,
					strings.TrimPrefix(itemTrim, "permissions "))
			case strings.HasPrefix(itemTrim, "security-role "):
				confRead.securityRole = strings.TrimPrefix(itemTrim, "security-role ")
			case strings.HasPrefix(itemTrim, "tenant "):
				confRead.tenant = strings.Trim(strings.TrimPrefix(itemTrim, "tenant "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSystemLoginClass(systemLoginClass string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system login class "+systemLoginClass)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
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
