package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ipsecProposalOptions struct {
	lifetimeSeconds        int
	lifetimeKilobytes      int
	name                   string
	authenticatioAlgorithm string
	encryptionAlgorithm    string
	protocol               string
}

func resourceIpsecProposal() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecProposalCreate,
		ReadContext:   resourceIpsecProposalRead,
		UpdateContext: resourceIpsecProposalUpdate,
		DeleteContext: resourceIpsecProposalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpsecProposalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(180, 86400),
			},
			"lifetime_kilobytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(64, 4294967294),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"esp", "ah"}, false),
			},
		},
	}
}

func resourceIpsecProposalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ipsec proposal not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	ipsecProposalExists, err := checkIpsecProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ipsecProposalExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security ipsec proposal %v already exists", d.Get("name").(string)))
	}
	if err := setIpsecProposal(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_ipsec_proposal", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecProposalExists, err = checkIpsecProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec proposal %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecProposalReadWJnprSess(d, m, jnprSess)...)
}
func resourceIpsecProposalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIpsecProposalReadWJnprSess(d, m, jnprSess)
}
func resourceIpsecProposalReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ipsecProposalOptions, err := readIpsecProposal(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecProposalOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecProposalData(d, ipsecProposalOptions)
	}

	return nil
}
func resourceIpsecProposalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIpsecProposal(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setIpsecProposal(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_ipsec_proposal", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecProposalReadWJnprSess(d, m, jnprSess)...)
}
func resourceIpsecProposalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIpsecProposal(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_ipsec_proposal", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceIpsecProposalImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ipsecProposalExists, err := checkIpsecProposalExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ipsecProposalExists {
		return nil, fmt.Errorf("don't find security ipsec proposal with id '%v' (id must be <name>)", d.Id())
	}
	ipsecProposalOptions, err := readIpsecProposal(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIpsecProposalData(d, ipsecProposalOptions)
	result[0] = d

	return result, nil
}

func checkIpsecProposalExists(ipsecProposal string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ipsecProposalConfig, err := sess.command("show configuration"+
		" security ipsec proposal "+ipsecProposal+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ipsecProposalConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setIpsecProposal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security ipsec proposal " + d.Get("name").(string)
	if d.Get("authentication_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-algorithm "+d.Get("authentication_algorithm").(string))
	}
	if d.Get("encryption_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" encryption-algorithm "+d.Get("encryption_algorithm").(string))
	}
	if d.Get("lifetime_seconds").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-seconds "+strconv.Itoa(d.Get("lifetime_seconds").(int)))
	}
	if d.Get("lifetime_kilobytes").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-kilobytes "+strconv.Itoa(d.Get("lifetime_kilobytes").(int)))
	}
	if d.Get("protocol").(string) != "" {
		configSet = append(configSet, setPrefix+" protocol "+d.Get("protocol").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readIpsecProposal(ipsecProposal string, m interface{}, jnprSess *NetconfObject) (ipsecProposalOptions, error) {
	sess := m.(*Session)
	var confRead ipsecProposalOptions

	ipsecProposalConfig, err := sess.command("show configuration"+
		" security ipsec proposal "+ipsecProposal+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ipsecProposalConfig != emptyWord {
		confRead.name = ipsecProposal
		for _, item := range strings.Split(ipsecProposalConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "authentication-algorithm "):
				confRead.authenticatioAlgorithm = strings.TrimPrefix(itemTrim, "authentication-algorithm ")
			case strings.HasPrefix(itemTrim, "encryption-algorithm "):
				confRead.encryptionAlgorithm = strings.TrimPrefix(itemTrim, "encryption-algorithm ")
			case strings.HasPrefix(itemTrim, "lifetime-kilobytes "):
				confRead.lifetimeKilobytes, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-kilobytes "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "lifetime-seconds "):
				confRead.lifetimeSeconds, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-seconds "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			}
		}
	}

	return confRead, nil
}
func delIpsecProposal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec proposal "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillIpsecProposalData(d *schema.ResourceData, ipsecProposalOptions ipsecProposalOptions) {
	if tfErr := d.Set("name", ipsecProposalOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_algorithm", ipsecProposalOptions.authenticatioAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("encryption_algorithm", ipsecProposalOptions.encryptionAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lifetime_kilobytes", ipsecProposalOptions.lifetimeKilobytes); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lifetime_seconds", ipsecProposalOptions.lifetimeSeconds); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol", ipsecProposalOptions.protocol); tfErr != nil {
		panic(tfErr)
	}
}
