package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type prefixListOptions struct {
	dynamicDB bool
	name      string
	applyPath string
	prefix    []string
}

func resourcePolicyoptionsPrefixList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyoptionsPrefixListCreate,
		ReadContext:   resourcePolicyoptionsPrefixListRead,
		UpdateContext: resourcePolicyoptionsPrefixListUpdate,
		DeleteContext: resourcePolicyoptionsPrefixListDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsPrefixListImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"apply_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateCIDRNetworkFunc(),
				},
			},
		},
	}
}

func resourcePolicyoptionsPrefixListCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setPolicyoptionsPrefixList(d, m, nil); err != nil {
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
	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsPrefixListExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options prefix-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsPrefixList(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_policyoptions_prefix_list", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsPrefixListExists, err = checkPolicyoptionsPrefixListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsPrefixListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options prefix-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsPrefixListReadWJnprSess(d, m, jnprSess)...)
}

func resourcePolicyoptionsPrefixListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourcePolicyoptionsPrefixListReadWJnprSess(d, m, jnprSess)
}

func resourcePolicyoptionsPrefixListReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if prefixListOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsPrefixListData(d, prefixListOptions)
	}

	return nil
}

func resourcePolicyoptionsPrefixListUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delPolicyoptionsPrefixList(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsPrefixList(d, m, nil); err != nil {
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
	if err := delPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsPrefixList(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_policyoptions_prefix_list", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsPrefixListReadWJnprSess(d, m, jnprSess)...)
}

func resourcePolicyoptionsPrefixListDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delPolicyoptionsPrefixList(d.Get("name").(string), m, nil); err != nil {
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
	if err := delPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_policyoptions_prefix_list", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsPrefixListImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsPrefixListExists {
		return nil, fmt.Errorf("don't find policy-options prefix-list with id '%v' (id must be <name>)", d.Id())
	}
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsPrefixListData(d, prefixListOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsPrefixListExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+"policy-options prefix-list "+name+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsPrefixList(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set policy-options prefix-list " + d.Get("name").(string)
	configSet = append(configSet, setPrefix)
	if d.Get("apply_path").(string) != "" {
		replaceSign := strings.ReplaceAll(d.Get("apply_path").(string), "<", "&lt;")
		replaceSign = strings.ReplaceAll(replaceSign, ">", "&gt;")
		configSet = append(configSet, setPrefix+" apply-path \""+replaceSign+"\"")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, setPrefix+" dynamic-db")
	}
	for _, v := range sortSetOfString(d.Get("prefix").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+" "+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func readPolicyoptionsPrefixList(name string, m interface{}, jnprSess *NetconfObject) (prefixListOptions, error) {
	sess := m.(*Session)
	var confRead prefixListOptions

	showConfig, err := sess.command(cmdShowConfig+
		"policy-options prefix-list "+name+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			switch {
			case strings.HasPrefix(itemTrim, "apply-path "):
				replaceSign := strings.ReplaceAll(strings.Trim(strings.TrimPrefix(itemTrim, "apply-path "), "\""), "&lt;", "<")
				replaceSign = strings.ReplaceAll(replaceSign, "&gt;", ">")
				confRead.applyPath = replaceSign
			case itemTrim == "dynamic-db":
				confRead.dynamicDB = true
			case strings.Contains(itemTrim, "/"):
				confRead.prefix = append(confRead.prefix, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsPrefixList(prefixList string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options prefix-list "+prefixList)

	return sess.configSet(configSet, jnprSess)
}

func fillPolicyoptionsPrefixListData(d *schema.ResourceData, prefixListOptions prefixListOptions) {
	if tfErr := d.Set("name", prefixListOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("apply_path", prefixListOptions.applyPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_db", prefixListOptions.dynamicDB); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefix", prefixListOptions.prefix); tfErr != nil {
		panic(tfErr)
	}
}
