package junos

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
)

type servicesOptions struct {
	securityIntelligence []map[string]interface{}
}

func resourceServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesCreate,
		ReadContext:   resourceServicesRead,
		UpdateContext: resourceServicesUpdate,
		DeleteContext: resourceServicesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesImport,
		},
		Schema: map[string]*schema.Schema{
			"security_intelligence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_token": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[a-zA-Z0-9]{32}$`),
								"Auth token must be consisted of 32 alphanumeric characters"),
							ConflictsWith: []string{"security_intelligence.0.authentication_tls_profile"},
						},
						"authentication_tls_profile": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"security_intelligence.0.authentication_token"},
						},
						"category_disable": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"default_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"profile_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"proxy_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url_parameter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceServicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServices(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("services")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setServices(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_services", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("services")

	return append(diagWarns, resourceServicesReadWJnprSess(d, m, jnprSess)...)
}
func resourceServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesReadWJnprSess(d, m, jnprSess)
}
func resourceServicesReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	servicesOptions, err := readServices(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillServices(d, servicesOptions)

	return nil
}
func resourceServicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delServices(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setServices(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_services", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesReadWJnprSess(d, m, jnprSess)...)
}
func resourceServicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceServicesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	servicesOptions, err := readServices(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServices(d, servicesOptions)
	d.SetId("services")
	result[0] = d

	return result, nil
}

func setServices(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	// setPrefix := "set services "
	configSet := make([]string, 0)

	for _, v := range d.Get("security_intelligence").([]interface{}) {
		configSetSecurityIntel, err := setServicesSecurityIntell(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetSecurityIntel...)
	}

	return sess.configSet(configSet, jnprSess)
}

func setServicesSecurityIntell(secuIntel interface{}) ([]string, error) {
	setPrefix := "set services security-intelligence "
	configSet := make([]string, 0)
	if secuIntel != nil {
		secuIntelM := secuIntel.(map[string]interface{})
		if v := secuIntelM["authentication_token"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication auth-token "+v)
		}
		if v := secuIntelM["authentication_tls_profile"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication tls-profile \""+v+"\"")
		}
		for _, v := range secuIntelM["category_disable"].(*schema.Set).List() {
			if v.(string) == "all" {
				configSet = append(configSet, setPrefix+"category all disable")
			} else {
				configSet = append(configSet, setPrefix+"category category-name "+v.(string)+" disable")
			}
		}
		for _, v := range secuIntelM["default_policy"].([]interface{}) {
			defPolicy := v.(map[string]interface{})
			configSet = append(configSet, setPrefix+"default-policy "+
				defPolicy["category_name"].(string)+" "+defPolicy["profile_name"].(string))
		}
		if v := secuIntelM["proxy_profile"].(string); v != "" {
			configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
		}
		if v := secuIntelM["url"].(string); v != "" {
			configSet = append(configSet, setPrefix+"url \""+v+"\"")
		}
		if v := secuIntelM["url_parameter"].(string); v != "" {
			configSet = append(configSet, setPrefix+"url-parameter \""+v+"\"")
		}
	} else {
		return configSet, fmt.Errorf("security_intelligence block is empty")
	}

	return configSet, nil
}

func listLinesServicesSecurityIntel() []string {
	return []string{
		"security-intelligence authentication",
		"security-intelligence category",
		"security-intelligence default-policy",
		"security-intelligence proxy-profile",
		"security-intelligence url",
		"security-intelligence url-parameter",
	}
}

func delServices(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesSecurityIntel()...)
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete services "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}
func readServices(m interface{}, jnprSess *NetconfObject) (servicesOptions, error) {
	sess := m.(*Session)
	var confRead servicesOptions

	servicesConfig, err := sess.command("show configuration services"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if servicesConfig != emptyWord {
		for _, item := range strings.Split(servicesConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if checkStringHasPrefixInList(itemTrim, listLinesServicesSecurityIntel()) {
				if err := readServicesSecurityIntel(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readServicesSecurityIntel(confRead *servicesOptions, itemTrimSecurityIntel string) error {
	itemTrim := strings.TrimPrefix(itemTrimSecurityIntel, "security-intelligence ")
	if len(confRead.securityIntelligence) == 0 {
		confRead.securityIntelligence = append(confRead.securityIntelligence, map[string]interface{}{
			"authentication_token":       "",
			"authentication_tls_profile": "",
			"category_disable":           make([]string, 0),
			"default_policy":             make([]map[string]interface{}, 0),
			"proxy_profile":              "",
			"url":                        "",
			"url_parameter":              "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "authentication auth-token "):
		confRead.securityIntelligence[0]["authentication_token"] = strings.TrimPrefix(itemTrim, "authentication auth-token ")
	case strings.HasPrefix(itemTrim, "authentication tls-profile "):
		confRead.securityIntelligence[0]["authentication_tls_profile"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"authentication tls-profile "), "\"")
	case strings.HasPrefix(itemTrim, "category "):
		if itemTrim == "category all disable" {
			confRead.securityIntelligence[0]["category_disable"] = append(
				confRead.securityIntelligence[0]["category_disable"].([]string), "all")
		} else {
			confRead.securityIntelligence[0]["category_disable"] = append(
				confRead.securityIntelligence[0]["category_disable"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(itemTrim, "category category-name "), " disable"))
		}
	case strings.HasPrefix(itemTrim, "default-policy "):
		if itemTrimSplit := strings.Split(itemTrim, " "); len(itemTrimSplit) == 3 {
			confRead.securityIntelligence[0]["default_policy"] = append(
				confRead.securityIntelligence[0]["default_policy"].([]map[string]interface{}), map[string]interface{}{
					"category_name": itemTrimSplit[1],
					"profile_name":  itemTrimSplit[2],
				})
		}
	case strings.HasPrefix(itemTrim, "proxy-profile "):
		confRead.securityIntelligence[0]["proxy_profile"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"proxy-profile "), "\"")
	case strings.HasPrefix(itemTrim, "url "):
		confRead.securityIntelligence[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"url "), "\"")
	case strings.HasPrefix(itemTrim, "url-parameter "):
		var err error
		confRead.securityIntelligence[0]["url_parameter"], err =
			jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim, "url-parameter "), "\""))
		if err != nil {
			return fmt.Errorf("failed to decode url-parameter : %w", err)
		}
	}

	return nil
}

func fillServices(d *schema.ResourceData, servicesOptions servicesOptions) {
	if tfErr := d.Set("security_intelligence", servicesOptions.securityIntelligence); tfErr != nil {
		panic(tfErr)
	}
}
