package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type svcSSLInitiationProfileOptions struct {
	enableFlowTracing  bool
	enableSessionCache bool
	clientCertificate  string
	name               string
	preferredCiphers   string
	protocolVersion    string
	customCiphers      []string
	trustedCA          []string
	actions            []map[string]interface{}
}

func resourceServicesSSLInitiationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesSSLInitiationProfileCreate,
		ReadContext:   resourceServicesSSLInitiationProfileRead,
		UpdateContext: resourceServicesSSLInitiationProfileUpdate,
		DeleteContext: resourceServicesSSLInitiationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesSSLInitiationProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crl_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"crl_if_not_present": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"allow", "drop"}, false),
						},
						"crl_ignore_hold_instruction_code": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ignore_server_auth_failure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_ciphers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_flow_tracing": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_session_cache": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"preferred_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"custom", "medium", "strong", "weak",
				}, false),
			},
			"protocol_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trusted_ca": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServicesSSLInitiationProfileCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesSSLInitiationProfile(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	svcSSLInitiationProfileExists, err := checkServicesSSLInitiationProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcSSLInitiationProfileExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf(
				"services ssl initiation profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesSSLInitiationProfile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services_ssl_initiation_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	svcSSLInitiationProfileExists, err = checkServicesSSLInitiationProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if svcSSLInitiationProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"services ssl initiation profile %v "+
				"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesSSLInitiationProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSSLInitiationProfileRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesSSLInitiationProfileReadWJnprSess(d, m, jnprSess)
}

func resourceServicesSSLInitiationProfileReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	svcSSLInitiationProfileOptions, err := readServicesSSLInitiationProfile(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if svcSSLInitiationProfileOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesSSLInitiationProfileData(d, svcSSLInitiationProfileOptions)
	}

	return nil
}

func resourceServicesSSLInitiationProfileUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delServicesSSLInitiationProfile(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesSSLInitiationProfile(d, m, nil); err != nil {
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
	if err := delServicesSSLInitiationProfile(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesSSLInitiationProfile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services_ssl_initiation_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesSSLInitiationProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSSLInitiationProfileDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delServicesSSLInitiationProfile(d.Get("name").(string), m, nil); err != nil {
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
	if err := delServicesSSLInitiationProfile(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_services_ssl_initiation_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesSSLInitiationProfileImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	svcSSLInitiationProfileExists, err := checkServicesSSLInitiationProfileExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !svcSSLInitiationProfileExists {
		return nil, fmt.Errorf("don't find services ssl initiation profile with id '%v' (id must be <name>)", d.Id())
	}
	svcSSLInitiationProfileOptions, err := readServicesSSLInitiationProfile(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesSSLInitiationProfileData(d, svcSSLInitiationProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesSSLInitiationProfileExists(
	profile string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" services ssl initiation profile \""+profile+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesSSLInitiationProfile(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services ssl initiation profile \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	for _, v := range d.Get("actions").([]interface{}) {
		if v == nil {
			return fmt.Errorf("actions block is empty")
		}
		actions := v.(map[string]interface{})
		if actions["crl_disable"].(bool) {
			configSet = append(configSet, setPrefix+"actions crl disable")
		}
		if v2 := actions["crl_if_not_present"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"actions crl if-not-present "+v2)
		}
		if actions["crl_ignore_hold_instruction_code"].(bool) {
			configSet = append(configSet, setPrefix+"actions crl ignore-hold-instruction-code")
		}
		if actions["ignore_server_auth_failure"].(bool) {
			configSet = append(configSet, setPrefix+"actions ignore-server-auth-failure")
		}
	}
	if v := d.Get("client_certificate").(string); v != "" {
		configSet = append(configSet, setPrefix+"client-certificate \""+v+"\"")
	}
	for _, v := range sortSetOfString(d.Get("custom_ciphers").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"custom-ciphers "+v)
	}
	if d.Get("enable_flow_tracing").(bool) {
		configSet = append(configSet, setPrefix+"enable-flow-tracing")
	}
	if d.Get("enable_session_cache").(bool) {
		configSet = append(configSet, setPrefix+"enable-session-cache")
	}
	if v := d.Get("preferred_ciphers").(string); v != "" {
		configSet = append(configSet, setPrefix+"preferred-ciphers "+v)
	}
	if v := d.Get("protocol_version").(string); v != "" {
		configSet = append(configSet, setPrefix+"protocol-version "+v)
	}
	for _, v := range sortSetOfString(d.Get("trusted_ca").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"trusted-ca \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesSSLInitiationProfile(profile string, m interface{}, jnprSess *NetconfObject) (
	svcSSLInitiationProfileOptions, error) {
	sess := m.(*Session)
	var confRead svcSSLInitiationProfileOptions

	showConfig, err := sess.command("show configuration"+
		" services ssl initiation profile \""+profile+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "actions "):
				if len(confRead.actions) == 0 {
					confRead.actions = append(confRead.actions, map[string]interface{}{
						"crl_disable":                      false,
						"crl_if_not_present":               "",
						"crl_ignore_hold_instruction_code": false,
						"ignore_server_auth_failure":       false,
					})
				}
				switch {
				case itemTrim == "actions crl disable":
					confRead.actions[0]["crl_disable"] = true
				case strings.HasPrefix(itemTrim, "actions crl if-not-present "):
					confRead.actions[0]["crl_if_not_present"] = strings.TrimPrefix(itemTrim, "actions crl if-not-present ")
				case itemTrim == "actions crl ignore-hold-instruction-code":
					confRead.actions[0]["crl_ignore_hold_instruction_code"] = true
				case itemTrim == "actions ignore-server-auth-failure":
					confRead.actions[0]["ignore_server_auth_failure"] = true
				}
			case strings.HasPrefix(itemTrim, "client-certificate "):
				confRead.clientCertificate = strings.Trim(strings.TrimPrefix(itemTrim, "client-certificate "), "\"")
			case strings.HasPrefix(itemTrim, "custom-ciphers "):
				confRead.customCiphers = append(confRead.customCiphers, strings.TrimPrefix(itemTrim, "custom-ciphers "))
			case itemTrim == "enable-flow-tracing":
				confRead.enableFlowTracing = true
			case itemTrim == "enable-session-cache":
				confRead.enableSessionCache = true
			case strings.HasPrefix(itemTrim, "preferred-ciphers "):
				confRead.preferredCiphers = strings.TrimPrefix(itemTrim, "preferred-ciphers ")
			case strings.HasPrefix(itemTrim, "protocol-version "):
				confRead.protocolVersion = strings.TrimPrefix(itemTrim, "protocol-version ")
			case strings.HasPrefix(itemTrim, "trusted-ca "):
				confRead.trustedCA = append(confRead.trustedCA,
					strings.Trim(strings.TrimPrefix(itemTrim, "trusted-ca "), "\""))
			}
		}
	}

	return confRead, nil
}

func delServicesSSLInitiationProfile(profile string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{
		"delete services ssl initiation profile \"" + profile + "\"",
	}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesSSLInitiationProfileData(
	d *schema.ResourceData, svcSSLInitiationProfileOptions svcSSLInitiationProfileOptions) {
	if tfErr := d.Set("name", svcSSLInitiationProfileOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("actions", svcSSLInitiationProfileOptions.actions); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("client_certificate", svcSSLInitiationProfileOptions.clientCertificate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("custom_ciphers", svcSSLInitiationProfileOptions.customCiphers); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("enable_flow_tracing", svcSSLInitiationProfileOptions.enableFlowTracing); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("enable_session_cache", svcSSLInitiationProfileOptions.enableSessionCache); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preferred_ciphers", svcSSLInitiationProfileOptions.preferredCiphers); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol_version", svcSSLInitiationProfileOptions.protocolVersion); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("trusted_ca", svcSSLInitiationProfileOptions.trustedCA); tfErr != nil {
		panic(tfErr)
	}
}
