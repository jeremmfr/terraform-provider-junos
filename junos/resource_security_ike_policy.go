package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
)

type ikePolicyOptions struct {
	name             string
	mode             string
	preSharedKeyText string
	preSharedKeyHexa string
	proposalSet      string
	proposals        []string
}

func listProposalSet() []string {
	return []string{
		"basic",
		"compatible",
		"prime-128",
		"prime-256",
		"standard",
		"suiteb-gcm-128",
		"suiteb-gcm-256",
	}
}

func resourceIkePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIkePolicyCreate,
		ReadContext:   resourceIkePolicyRead,
		UpdateContext: resourceIkePolicyUpdate,
		DeleteContext: resourceIkePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIkePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
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
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "main",
				ValidateFunc: validation.StringInSlice([]string{"main", "aggressive"}, false),
			},
			"pre_shared_key_text": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pre_shared_key_hexa"},
				Sensitive:     true,
			},
			"pre_shared_key_hexa": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pre_shared_key_text"},
				Sensitive:     true,
			},
		},
	}
}

func resourceIkePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ike policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	ikePolicyExists, err := checkIkePolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ikePolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security ike policy %v already exists", d.Get("name").(string)))
	}
	if err := setIkePolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikePolicyExists, err = checkIkePolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikePolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkePolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceIkePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIkePolicyReadWJnprSess(d, m, jnprSess)
}
func resourceIkePolicyReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ikePolicyOptions, err := readIkePolicy(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ikePolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIkePolicyData(d, ikePolicyOptions)
	}

	return nil
}
func resourceIkePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIkePolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setIkePolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceIkePolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceIkePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIkePolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceIkePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ikePolicyExists, err := checkIkePolicyExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ikePolicyExists {
		return nil, fmt.Errorf("don't find security ike policy with id '%v' (id must be <name>)", d.Id())
	}
	ikePolicyOptions, err := readIkePolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIkePolicyData(d, ikePolicyOptions)
	result[0] = d

	return result, nil
}

func checkIkePolicyExists(ikePolicy string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ikePolicyConfig, err := sess.command("show configuration security ike policy "+ikePolicy+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ikePolicyConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setIkePolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security ike policy " + d.Get("name").(string)
	if d.Get("mode").(string) != "" {
		if d.Get("mode").(string) != "main" && d.Get("mode").(string) != "aggressive" {
			return fmt.Errorf("unknown ike mode %v", d.Get("mode").(string))
		}
		configSet = append(configSet, setPrefix+" mode "+d.Get("mode").(string))
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string))
	}
	if d.Get("proposal_set").(string) != "" {
		configSet = append(configSet, setPrefix+" proposal-set "+d.Get("proposal_set").(string))
	}
	if d.Get("pre_shared_key_text").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key ascii-text "+d.Get("pre_shared_key_text").(string))
	}
	if d.Get("pre_shared_key_hexa").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key hexadecimal "+d.Get("pre_shared_key_hexa").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readIkePolicy(ikePolicy string, m interface{}, jnprSess *NetconfObject) (ikePolicyOptions, error) {
	sess := m.(*Session)
	var confRead ikePolicyOptions

	ikePolicyConfig, err := sess.command("show configuration"+
		" security ike policy "+ikePolicy+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ikePolicyConfig != emptyWord {
		confRead.name = ikePolicy
		for _, item := range strings.Split(ikePolicyConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "mode "):
				confRead.mode = strings.TrimPrefix(itemTrim, "mode ")
			case strings.HasPrefix(itemTrim, "proposals "):
				confRead.proposals = append(confRead.proposals, strings.TrimPrefix(itemTrim, "proposals "))
			case strings.HasPrefix(itemTrim, "proposal-set "):
				confRead.proposalSet = strings.TrimPrefix(itemTrim, "proposal-set ")
			case strings.HasPrefix(itemTrim, "pre-shared-key hexadecimal "):
				confRead.preSharedKeyHexa, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"pre-shared-key hexadecimal "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode pre-shared-key hexadecimal : %w", err)
				}
			case strings.HasPrefix(itemTrim, "pre-shared-key ascii-text "):
				confRead.preSharedKeyText, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"pre-shared-key ascii-text "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode pre-shared-key ascii-text : %w", err)
				}
			}
		}
	}

	return confRead, nil
}
func delIkePolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike policy "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillIkePolicyData(d *schema.ResourceData, ikePolicyOptions ikePolicyOptions) {
	if tfErr := d.Set("name", ikePolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mode", ikePolicyOptions.mode); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pre_shared_key_text", ikePolicyOptions.preSharedKeyText); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pre_shared_key_hexa", ikePolicyOptions.preSharedKeyHexa); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposals", ikePolicyOptions.proposals); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposal_set", ikePolicyOptions.proposalSet); tfErr != nil {
		panic(tfErr)
	}
}
