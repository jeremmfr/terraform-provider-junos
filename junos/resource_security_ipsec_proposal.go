package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Create: resourceIpsecProposalCreate,
		Read:   resourceIpsecProposalRead,
		Update: resourceIpsecProposalUpdate,
		Delete: resourceIpsecProposalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpsecProposalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
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
				ValidateFunc: validateIntRange(180, 86400),
			},
			"lifetime_kilobytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(64, 4294967294),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"esp", "ah"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'esp' or 'ah'", value, k))
					}
					return
				},
			},
		},
	}
}

func resourceIpsecProposalCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security ipsec proposal not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	ipsecProposalExists, err := checkIpsecProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if ipsecProposalExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("ipsec proposal %v already exists", d.Get("name").(string))
	}
	err = setIpsecProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_ipsec_proposal", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	ipsecProposalExists, err = checkIpsecProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if ipsecProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("ipsec proposal %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceIpsecProposalRead(d, m)
}
func resourceIpsecProposalRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	ipsecProposalOptions, err := readIpsecProposal(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if ipsecProposalOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecProposalData(d, ipsecProposalOptions)
	}
	return nil
}
func resourceIpsecProposalUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delIpsecProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setIpsecProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_ipsec_policy", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceIpsecProposalRead(d, m)
}
func resourceIpsecProposalDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delIpsecProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_ipsec_policy", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
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
		return nil, fmt.Errorf("don't find ipsec proposal with id '%v' (id must be <name>)", d.Id())
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
		configSet = append(configSet, setPrefix+" authentication-algorithm "+d.Get("authentication_algorithm").(string)+"\n")
	}
	if d.Get("encryption_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" encryption-algorithm "+d.Get("encryption_algorithm").(string)+"\n")
	}
	if d.Get("lifetime_seconds").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-seconds "+strconv.Itoa(d.Get("lifetime_seconds").(int))+"\n")
	}
	if d.Get("lifetime_kilobytes").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-kilobytes "+strconv.Itoa(d.Get("lifetime_kilobytes").(int))+"\n")
	}
	if d.Get("protocol").(string) != "" {
		configSet = append(configSet, setPrefix+" protocol "+d.Get("protocol").(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
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
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "lifetime-seconds "):
				confRead.lifetimeSeconds, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-seconds "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delIpsecProposal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec proposal "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillIpsecProposalData(d *schema.ResourceData, ipsecProposalOptions ipsecProposalOptions) {
	tfErr := d.Set("name", ipsecProposalOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_algorithm", ipsecProposalOptions.authenticatioAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("encryption_algorithm", ipsecProposalOptions.encryptionAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("lifetime_kilobytes", ipsecProposalOptions.lifetimeKilobytes)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("lifetime_seconds", ipsecProposalOptions.lifetimeSeconds)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("protocol", ipsecProposalOptions.protocol)
	if tfErr != nil {
		panic(tfErr)
	}
}
