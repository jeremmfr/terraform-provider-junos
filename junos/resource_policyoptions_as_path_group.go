package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type asPathGroupOptions struct {
	dynamicDB bool
	name      string
	asPath    []map[string]interface{}
}

func resourcePolicyoptionsAsPathGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyoptionsAsPathGroupCreate,
		ReadContext:   resourcePolicyoptionsAsPathGroupRead,
		UpdateContext: resourcePolicyoptionsAsPathGroupUpdate,
		DeleteContext: resourcePolicyoptionsAsPathGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsAsPathGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"as_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsAsPathGroupCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if policyoptsAsPathGroupExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("policy-options as-path-group %v already exists", d.Get("name").(string)))
	}

	if err := setPolicyoptionsAsPathGroup(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_policyoptions_as_path_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsAsPathGroupExists, err = checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options as-path-group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsAsPathGroupReadWJnprSess(d, m, jnprSess)...)
}
func resourcePolicyoptionsAsPathGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourcePolicyoptionsAsPathGroupReadWJnprSess(d, m, jnprSess)
}
func resourcePolicyoptionsAsPathGroupReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if asPathGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)
	}

	return nil
}
func resourcePolicyoptionsAsPathGroupUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setPolicyoptionsAsPathGroup(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_policyoptions_as_path_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathGroupReadWJnprSess(d, m, jnprSess)...)
}
func resourcePolicyoptionsAsPathGroupDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_policyoptions_as_path_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourcePolicyoptionsAsPathGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathGroupExists {
		return nil, fmt.Errorf("don't find policy-options as-path-group with id '%v' (id must be <name>)", d.Id())
	}
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsAsPathGroupExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	asPathGroupConfig, err := sess.command("show configuration "+
		"policy-options as-path-group "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if asPathGroupConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setPolicyoptionsAsPathGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set policy-options as-path-group " + d.Get("name").(string)
	for _, v := range d.Get("as_path").([]interface{}) {
		asPath := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+
			" as-path "+asPath["name"].(string)+
			" \""+asPath["path"].(string)+"\"")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, setPrefix+" dynamic-db")
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readPolicyoptionsAsPathGroup(asPathGroup string,
	m interface{}, jnprSess *NetconfObject) (asPathGroupOptions, error) {
	sess := m.(*Session)
	var confRead asPathGroupOptions

	asPathGroupConfig, err := sess.command("show configuration"+
		" policy-options as-path-group "+asPathGroup+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if asPathGroupConfig != emptyWord {
		confRead.name = asPathGroup
		for _, item := range strings.Split(asPathGroupConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == dynamicDB:
				confRead.dynamicDB = true
			case strings.HasPrefix(itemTrim, "as-path "):
				asPath := map[string]interface{}{
					"name": "",
					"path": "",
				}
				itemSplit := strings.Split(strings.TrimPrefix(itemTrim, "as-path "), " ")
				asPath["name"] = itemSplit[0]
				asPath["path"] = strings.Trim(strings.TrimPrefix(itemTrim,
					"as-path "+asPath["name"].(string)+" "), "\"")
				confRead.asPath = append(confRead.asPath, asPath)
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsAsPathGroup(asPathGroup string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path-group "+asPathGroup)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func fillPolicyoptionsAsPathGroupData(d *schema.ResourceData, asPathGroupOptions asPathGroupOptions) {
	if tfErr := d.Set("name", asPathGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path", asPathGroupOptions.asPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_db", asPathGroupOptions.dynamicDB); tfErr != nil {
		panic(tfErr)
	}
}
