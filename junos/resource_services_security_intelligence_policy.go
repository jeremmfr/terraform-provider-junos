package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type securityIntellPolicyOptions struct {
	name        string
	description string
	category    []map[string]interface{}
}

func resourceServicesSecurityIntellPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesSecurityIntellPolicyCreate,
		ReadContext:   resourceServicesSecurityIntellPolicyRead,
		UpdateContext: resourceServicesSecurityIntellPolicyUpdate,
		DeleteContext: resourceServicesSecurityIntellPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesSecurityIntellPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"category": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"profile_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceServicesSecurityIntellPolicyCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesSecurityIntellPolicy(d, m, nil); err != nil {
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
	securityIntellPolicyExists, err := checkServicesSecurityIntellPolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellPolicyExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services security-intelligence policy %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesSecurityIntellPolicy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services_security_intelligence_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityIntellPolicyExists, err = checkServicesSecurityIntellPolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services security-intelligence policy %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesSecurityIntellPolicyReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSecurityIntellPolicyRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesSecurityIntellPolicyReadWJnprSess(d, m, jnprSess)
}

func resourceServicesSecurityIntellPolicyReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	securityIntellPolicyOptions, err := readServicesSecurityIntellPolicy(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if securityIntellPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesSecurityIntellPolicyData(d, securityIntellPolicyOptions)
	}

	return nil
}

func resourceServicesSecurityIntellPolicyUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delServicesSecurityIntellPolicy(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesSecurityIntellPolicy(d, m, nil); err != nil {
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
	if err := delServicesSecurityIntellPolicy(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesSecurityIntellPolicy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services_security_intelligence_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesSecurityIntellPolicyReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSecurityIntellPolicyDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delServicesSecurityIntellPolicy(d.Get("name").(string), m, nil); err != nil {
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
	if err := delServicesSecurityIntellPolicy(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_services_security_intelligence_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesSecurityIntellPolicyImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityIntellPolicyExists, err := checkServicesSecurityIntellPolicyExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityIntellPolicyExists {
		return nil, fmt.Errorf("don't find services security-intelligence policy with id '%v' (id must be <name>)", d.Id())
	}
	securityIntellPolicyOptions, err := readServicesSecurityIntellPolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesSecurityIntellPolicyData(d, securityIntellPolicyOptions)

	result[0] = d

	return result, nil
}

func checkServicesSecurityIntellPolicyExists(policy string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" services security-intelligence policy \""+policy+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesSecurityIntellPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services security-intelligence policy \"" + d.Get("name").(string) + "\" "
	categoryNameList := make([]string, 0)
	for _, v := range d.Get("category").([]interface{}) {
		category := v.(map[string]interface{})
		if bchk.StringInSlice(category["name"].(string), categoryNameList) {
			return fmt.Errorf("multiple blocks category with the same name %s", category["name"].(string))
		}
		categoryNameList = append(categoryNameList, category["name"].(string))
		configSet = append(configSet,
			setPrefix+category["name"].(string)+" \""+category["profile_name"].(string)+"\"")
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesSecurityIntellPolicy(policy string, m interface{}, jnprSess *NetconfObject) (
	securityIntellPolicyOptions, error) {
	sess := m.(*Session)
	var confRead securityIntellPolicyOptions

	showConfig, err := sess.command("show configuration"+
		" services security-intelligence policy \""+policy+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = policy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case len(strings.Split(itemTrim, " ")) == 2:
				lineCut := strings.Split(itemTrim, " ")
				confRead.category = append(confRead.category, map[string]interface{}{
					"name":         lineCut[0],
					"profile_name": strings.Trim(lineCut[1], "\""),
				})
			}
		}
	}

	return confRead, nil
}

func delServicesSecurityIntellPolicy(policy string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete services security-intelligence policy \"" + policy + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesSecurityIntellPolicyData(
	d *schema.ResourceData, securityIntellPolicyOptions securityIntellPolicyOptions) {
	if tfErr := d.Set("name", securityIntellPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("category", securityIntellPolicyOptions.category); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", securityIntellPolicyOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
