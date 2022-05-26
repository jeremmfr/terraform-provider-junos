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
		CreateWithoutTimeout: resourceIkePolicyCreate,
		ReadWithoutTimeout:   resourceIkePolicyRead,
		UpdateWithoutTimeout: resourceIkePolicyUpdate,
		DeleteWithoutTimeout: resourceIkePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIkePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
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
	if sess.junosFakeCreateSetFile != "" {
		if err := setIkePolicy(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ike policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ikePolicyExists, err := checkIkePolicyExists(d.Get("name").(string), sess, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikePolicyExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike policy %v already exists", d.Get("name").(string)))...)
	}
	if err := setIkePolicy(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikePolicyExists, err = checkIkePolicyExists(d.Get("name").(string), sess, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikePolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkePolicyReadWJnprSess(d, sess, jnprSess)...)
}

func resourceIkePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIkePolicyReadWJnprSess(d, sess, jnprSess)
}

func resourceIkePolicyReadWJnprSess(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ikePolicyOptions, err := readIkePolicy(d.Get("name").(string), sess, jnprSess)
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
	if sess.junosFakeUpdateAlso {
		if err := delIkePolicy(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setIkePolicy(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIkePolicy(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIkePolicy(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	d.Partial(false)

	return append(diagWarns, resourceIkePolicyReadWJnprSess(d, sess, jnprSess)...)
}

func resourceIkePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delIkePolicy(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIkePolicy(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_ike_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIkePolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ikePolicyExists, err := checkIkePolicyExists(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ikePolicyExists {
		return nil, fmt.Errorf("don't find security ike policy with id '%v' (id must be <name>)", d.Id())
	}
	ikePolicyOptions, err := readIkePolicy(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIkePolicyData(d, ikePolicyOptions)
	result[0] = d

	return result, nil
}

func checkIkePolicyExists(ikePolicy string, sess *Session, jnprSess *NetconfObject) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+"security ike policy "+ikePolicy+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setIkePolicy(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
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

	return sess.configSet(configSet, jnprSess)
}

func readIkePolicy(ikePolicy string, sess *Session, jnprSess *NetconfObject) (ikePolicyOptions, error) {
	var confRead ikePolicyOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security ike policy "+ikePolicy+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = ikePolicy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
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
					return confRead, fmt.Errorf("failed to decode pre-shared-key hexadecimal: %w", err)
				}
			case strings.HasPrefix(itemTrim, "pre-shared-key ascii-text "):
				confRead.preSharedKeyText, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"pre-shared-key ascii-text "), "\""))
				if err != nil {
					return confRead, fmt.Errorf("failed to decode pre-shared-key ascii-text: %w", err)
				}
			}
		}
	}

	return confRead, nil
}

func delIkePolicy(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike policy "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
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
