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
		CreateWithoutTimeout: resourceServicesSecurityIntellPolicyCreate,
		ReadWithoutTimeout:   resourceServicesSecurityIntellPolicyRead,
		UpdateWithoutTimeout: resourceServicesSecurityIntellPolicyUpdate,
		DeleteWithoutTimeout: resourceServicesSecurityIntellPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesSecurityIntellPolicyImport,
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

func resourceServicesSecurityIntellPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setServicesSecurityIntellPolicy(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityIntellPolicyExists, err := checkServicesSecurityIntellPolicyExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellPolicyExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services security-intelligence policy %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesSecurityIntellPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_services_security_intelligence_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityIntellPolicyExists, err = checkServicesSecurityIntellPolicyExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services security-intelligence policy %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesSecurityIntellPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceServicesSecurityIntellPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceServicesSecurityIntellPolicyReadWJunSess(d, clt, junSess)
}

func resourceServicesSecurityIntellPolicyReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	securityIntellPolicyOptions, err := readServicesSecurityIntellPolicy(d.Get("name").(string), clt, junSess)
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

func resourceServicesSecurityIntellPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delServicesSecurityIntellPolicy(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesSecurityIntellPolicy(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delServicesSecurityIntellPolicy(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesSecurityIntellPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_services_security_intelligence_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesSecurityIntellPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceServicesSecurityIntellPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delServicesSecurityIntellPolicy(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delServicesSecurityIntellPolicy(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_services_security_intelligence_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesSecurityIntellPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	securityIntellPolicyExists, err := checkServicesSecurityIntellPolicyExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !securityIntellPolicyExists {
		return nil, fmt.Errorf("don't find services security-intelligence policy with id '%v' (id must be <name>)", d.Id())
	}
	securityIntellPolicyOptions, err := readServicesSecurityIntellPolicy(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillServicesSecurityIntellPolicyData(d, securityIntellPolicyOptions)

	result[0] = d

	return result, nil
}

func checkServicesSecurityIntellPolicyExists(policy string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"services security-intelligence policy \""+policy+"\""+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setServicesSecurityIntellPolicy(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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

	return clt.configSet(configSet, junSess)
}

func readServicesSecurityIntellPolicy(policy string, clt *Client, junSess *junosSession,
) (securityIntellPolicyOptions, error) {
	var confRead securityIntellPolicyOptions

	showConfig, err := clt.command(cmdShowConfig+
		"services security-intelligence policy \""+policy+"\""+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = policy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
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

func delServicesSecurityIntellPolicy(policy string, clt *Client, junSess *junosSession) error {
	configSet := []string{"delete services security-intelligence policy \"" + policy + "\""}

	return clt.configSet(configSet, junSess)
}

func fillServicesSecurityIntellPolicyData(
	d *schema.ResourceData, securityIntellPolicyOptions securityIntellPolicyOptions,
) {
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
