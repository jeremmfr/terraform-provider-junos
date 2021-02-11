package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ipsecPolicyOptions struct {
	name        string
	pfsKeys     string
	proposalSet string
	proposals   []string
}

func resourceIpsecPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecPolicyCreate,
		ReadContext:   resourceIpsecPolicyRead,
		UpdateContext: resourceIpsecPolicyUpdate,
		DeleteContext: resourceIpsecPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpsecPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"pfs_keys": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposals": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
			"proposal_set": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listProposalSet(), false),
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
		},
	}
}

func resourceIpsecPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ipsec policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ipsecPolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security ipsec policy %v already exists", d.Get("name").(string)))
	}
	if err := setIpsecPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_ipsec_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecPolicyExists, err = checkIpsecPolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceIpsecPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIpsecPolicyReadWJnprSess(d, m, jnprSess)
}
func resourceIpsecPolicyReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ipsecPolicyOptions, err := readIpsecPolicy(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecPolicyData(d, ipsecPolicyOptions)
	}

	return nil
}
func resourceIpsecPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIpsecPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setIpsecPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_ipsec_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceIpsecPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIpsecPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_ipsec_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceIpsecPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ipsecPolicyExists {
		return nil, fmt.Errorf("don't find security ipsec policy with id '%v' (id must be <name>)", d.Id())
	}
	ipsecPolicyOptions, err := readIpsecPolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIpsecPolicyData(d, ipsecPolicyOptions)
	result[0] = d

	return result, nil
}

func checkIpsecPolicyExists(ipsecPolicy string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ipsecPolicyConfig, err := sess.command("show configuration"+
		" security ipsec policy "+ipsecPolicy+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ipsecPolicyConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setIpsecPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security ipsec policy " + d.Get("name").(string)
	if d.Get("pfs_keys").(string) != "" {
		configSet = append(configSet, setPrefix+" perfect-forward-secrecy keys "+d.Get("pfs_keys").(string))
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string))
	}
	if d.Get("proposal_set").(string) != "" {
		configSet = append(configSet, setPrefix+" proposal-set "+d.Get("proposal_set").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readIpsecPolicy(ipsecPolicy string, m interface{}, jnprSess *NetconfObject) (ipsecPolicyOptions, error) {
	sess := m.(*Session)
	var confRead ipsecPolicyOptions

	ipsecPolicyConfig, err := sess.command("show configuration"+
		" security ipsec policy "+ipsecPolicy+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ipsecPolicyConfig != emptyWord {
		confRead.name = ipsecPolicy
		for _, item := range strings.Split(ipsecPolicyConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "perfect-forward-secrecy keys "):
				confRead.pfsKeys = strings.TrimPrefix(itemTrim, "perfect-forward-secrecy keys ")
			case strings.HasPrefix(itemTrim, "proposals "):
				confRead.proposals = append(confRead.proposals, strings.TrimPrefix(itemTrim, "proposals "))
			case strings.HasPrefix(itemTrim, "proposal-set "):
				confRead.proposalSet = strings.TrimPrefix(itemTrim, "proposal-set ")
			}
		}
	}

	return confRead, nil
}
func delIpsecPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec policy "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillIpsecPolicyData(d *schema.ResourceData, ipsecPolicyOptions ipsecPolicyOptions) {
	if tfErr := d.Set("name", ipsecPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pfs_keys", ipsecPolicyOptions.pfsKeys); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposals", ipsecPolicyOptions.proposals); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposal_set", ipsecPolicyOptions.proposalSet); tfErr != nil {
		panic(tfErr)
	}
}
