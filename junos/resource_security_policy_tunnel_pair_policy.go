package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type policyPairPolicyOptions struct {
	zoneA      string
	zoneB      string
	policyAtoB string
	policyBtoA string
}

func resourceSecurityPolicyTunnelPairPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyTunnelPairPolicyCreate,
		ReadContext:   resourceSecurityPolicyTunnelPairPolicyRead,
		DeleteContext: resourceSecurityPolicyTunnelPairPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityPolicyTunnelPairPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"zone_a": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"zone_b": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"policy_a_to_b": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"policy_b_to_a": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
		},
	}
}

func resourceSecurityPolicyTunnelPairPolicyCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security policy tunnel pair policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityPolicyExists, err := checkSecurityPolicyExists(d.Get("zone_a").(string), d.Get("zone_b").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if !securityPolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security policy from %v to %v not exists",
			d.Get("zone_a").(string), d.Get("zone_b").(string)))
	}
	securityPolicyExists, err = checkSecurityPolicyExists(d.Get("zone_b").(string), d.Get("zone_a").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if !securityPolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security policy from %v to %v not exists",
			d.Get("zone_b").(string), d.Get("zone_a").(string)))
	}
	pairPolicyExists, err := checkSecurityPolicyPairExists(d.Get("zone_a").(string), d.Get("policy_a_to_b").(string),
		d.Get("zone_b").(string), d.Get("policy_b_to_a").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if pairPolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security policy pair policy %v(%v) / %v(%v) already exists",
			d.Get("zone_a").(string), d.Get("policy_a_to_b").(string),
			d.Get("zone_b").(string), d.Get("policy_b_to_a").(string)))
	}
	if err := setSecurityPolicyTunnelPairPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_policy_tunnel_pair_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	pairPolicyExists, err = checkSecurityPolicyPairExists(d.Get("zone_a").(string), d.Get("policy_a_to_b").(string),
		d.Get("zone_b").(string), d.Get("policy_b_to_a").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if pairPolicyExists {
		d.SetId(d.Get("zone_a").(string) + idSeparator + d.Get("policy_a_to_b").(string) +
			idSeparator + d.Get("zone_b").(string) + idSeparator + d.Get("policy_b_to_a").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security policy pair policy not exists after commit "+
			"=> check your config"))...)
	}

	return append(diagWarns, resourceSecurityPolicyTunnelPairPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityPolicyTunnelPairPolicyRead(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityPolicyTunnelPairPolicyReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityPolicyTunnelPairPolicyReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	policyPairPolicyOptions, err := readSecurityPolicyTunnelPairPolicy(d.Get("zone_a").(string)+idSeparator+
		d.Get("policy_a_to_b").(string)+idSeparator+
		d.Get("zone_b").(string)+idSeparator+
		d.Get("policy_b_to_a").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if policyPairPolicyOptions.policyAtoB == "" && policyPairPolicyOptions.policyBtoA == "" {
		d.SetId("")
	} else {
		fillSecurityPolicyPairData(d, policyPairPolicyOptions)
	}

	return nil
}
func resourceSecurityPolicyTunnelPairPolicyDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityPolicyTunnelPairPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_policy_tunnel_pair_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityPolicyTunnelPairPolicyImport(d *schema.ResourceData,
	m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if len(idList) < 4 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	poliyPairPolicyExists, err := checkSecurityPolicyPairExists(idList[0], idList[1], idList[2], idList[3], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !poliyPairPolicyExists {
		return nil, fmt.Errorf("don't find policy pair policy with id %v "+
			"(id must be <zone_a>"+idSeparator+"<policy_a_to_b>"+idSeparator+
			"<zone_b>"+idSeparator+"<policy_b_to_a>)", d.Id())
	}

	policyPairPolicyOptions, err := readSecurityPolicyTunnelPairPolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityPolicyPairData(d, policyPairPolicyOptions)

	result[0] = d

	return result, nil
}

func checkSecurityPolicyPairExists(zoneA, policyAtoB, zoneB, policyBtoA string,
	m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)

	pairAtoBConfig, err := sess.command("show configuration"+
		" security policies from-zone "+zoneA+" to-zone "+zoneB+" policy "+policyAtoB+
		" then permit tunnel pair-policy | display set", jnprSess)
	if err != nil {
		return false, err
	}
	pairBtoAConfig, err := sess.command("show configuration"+
		" security policies from-zone "+zoneB+" to-zone "+zoneA+" policy "+policyBtoA+
		" then permit tunnel pair-policy | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if pairAtoBConfig == emptyWord && pairBtoAConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityPolicyTunnelPairPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 2)

	configSet = append(configSet, "set security policies from-zone "+
		d.Get("zone_a").(string)+" to-zone "+d.Get("zone_b").(string)+
		" policy "+d.Get("policy_a_to_b").(string)+
		" then permit tunnel pair-policy "+d.Get("policy_b_to_a").(string))
	configSet = append(configSet, "set security policies from-zone "+
		d.Get("zone_b").(string)+" to-zone "+d.Get("zone_a").(string)+
		" policy "+d.Get("policy_b_to_a").(string)+
		" then permit tunnel pair-policy "+d.Get("policy_a_to_b").(string))

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSecurityPolicyTunnelPairPolicy(idRessource string,
	m interface{}, jnprSess *NetconfObject) (policyPairPolicyOptions, error) {
	zone := strings.Split(idRessource, idSeparator)
	zoneA := zone[0]
	policyAtoB := zone[1]
	zoneB := zone[2]
	policyBtoA := zone[3]

	sess := m.(*Session)
	var confRead policyPairPolicyOptions

	policyPairConfig, err := sess.command("show configuration"+
		" security policies from-zone "+zoneA+" to-zone "+zoneB+" policy "+policyAtoB+
		" then permit tunnel pair-policy | display set", jnprSess)
	if err != nil {
		return confRead, err
	}
	if policyPairConfig != emptyWord {
		confRead.zoneA = zoneA
		confRead.zoneB = zoneB
		for _, item := range strings.Split(policyPairConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			if strings.Contains(item, " tunnel pair-policy ") {
				confRead.policyBtoA = strings.TrimPrefix(item,
					"set security policies from-zone "+zoneA+" to-zone "+zoneB+
						" policy "+policyAtoB+" then permit tunnel pair-policy ")
			}
		}
	}
	policyPairConfig, err = sess.command("show configuration"+
		" security policies from-zone "+zoneB+" to-zone "+zoneA+" policy "+policyBtoA+
		" then permit tunnel pair-policy | display set", jnprSess)
	if err != nil {
		return confRead, err
	}
	if policyPairConfig != emptyWord {
		confRead.zoneA = zoneA
		confRead.zoneB = zoneB
		for _, item := range strings.Split(policyPairConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			if strings.Contains(item, " tunnel pair-policy ") {
				confRead.policyAtoB = strings.TrimPrefix(item,
					"set security policies from-zone "+zoneB+" to-zone "+zoneA+
						" policy "+policyBtoA+" then permit tunnel pair-policy ")
			}
		}
	}

	return confRead, nil
}
func delSecurityPolicyTunnelPairPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 2)
	configSet = append(configSet, "delete security policies"+
		" from-zone "+d.Get("zone_a").(string)+" to-zone "+d.Get("zone_b").(string)+
		" policy "+d.Get("policy_a_to_b").(string)+
		" then permit tunnel pair-policy "+d.Get("policy_b_to_a").(string))
	configSet = append(configSet, "delete security policies"+
		" from-zone "+d.Get("zone_b").(string)+" to-zone "+d.Get("zone_a").(string)+
		" policy "+d.Get("policy_b_to_a").(string)+
		" then permit tunnel pair-policy "+d.Get("policy_a_to_b").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillSecurityPolicyPairData(d *schema.ResourceData, policyPairPolicyOptions policyPairPolicyOptions) {
	if tfErr := d.Set("zone_a", policyPairPolicyOptions.zoneA); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("zone_b", policyPairPolicyOptions.zoneB); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy_a_to_b", policyPairPolicyOptions.policyAtoB); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy_b_to_a", policyPairPolicyOptions.policyBtoA); tfErr != nil {
		panic(tfErr)
	}
}
