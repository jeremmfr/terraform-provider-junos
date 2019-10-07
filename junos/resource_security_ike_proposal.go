package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type ikeProposalOptions struct {
	lifetimeSeconds         int
	name                    string
	authenticationAlgorithm string
	authenticationMethod    string
	dhGroup                 string
	encryptionAlgorithm     string
}

func resourceIkeProposal() *schema.Resource {
	return &schema.Resource{
		Create: resourceIkeProposalCreate,
		Read:   resourceIkeProposalRead,
		Update: resourceIkeProposalUpdate,
		Delete: resourceIkeProposalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIkeProposalImport,
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
			"authentication_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dh_group": {
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
		},
	}
}

func resourceIkeProposalCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security ike proposal not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	ikeProposalExists, err := checkIkeProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if ikeProposalExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("ike proposal %v already exists", d.Get("name").(string))
	}
	err = setIkeProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	ikeProposalExists, err = checkIkeProposalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if ikeProposalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("ike proposal %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceIkeProposalRead(d, m)
}
func resourceIkeProposalRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	ikeProposalOptions, err := readIkeProposal(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if ikeProposalOptions.name == "" {
		d.SetId("")
	} else {
		fillIkeProposalData(d, ikeProposalOptions)
	}
	return nil
}
func resourceIkeProposalUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delIkeProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setIkeProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceIkeProposalRead(d, m)
}
func resourceIkeProposalDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delIkeProposal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceIkeProposalImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ikeProposalExists, err := checkIkeProposalExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ikeProposalExists {
		return nil, fmt.Errorf("don't find ike proposal with id '%v' (id must be <name>)", d.Id())
	}
	ikeProposalOptions, err := readIkeProposal(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIkeProposalData(d, ikeProposalOptions)
	result[0] = d
	return result, nil
}

func checkIkeProposalExists(ikeProposal string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ikeProposalConfig, err := sess.command("show configuration"+
		" security ike proposal "+ikeProposal+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ikeProposalConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setIkeProposal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security ike proposal " + d.Get("name").(string)
	if d.Get("authentication_method").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-method "+d.Get("authentication_method").(string)+"\n")
	}
	if d.Get("authentication_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" authentication-algorithm "+d.Get("authentication_algorithm").(string)+"\n")
	}
	if d.Get("dh_group").(string) != "" {
		configSet = append(configSet, setPrefix+" dh-group "+d.Get("dh_group").(string)+"\n")
	}
	if d.Get("encryption_algorithm").(string) != "" {
		configSet = append(configSet, setPrefix+" encryption-algorithm "+d.Get("encryption_algorithm").(string)+"\n")
	}
	if d.Get("lifetime_seconds").(int) != 0 {
		configSet = append(configSet, setPrefix+" lifetime-seconds "+strconv.Itoa(d.Get("lifetime_seconds").(int))+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readIkeProposal(ikeProposal string, m interface{}, jnprSess *NetconfObject) (ikeProposalOptions, error) {
	sess := m.(*Session)
	var confRead ikeProposalOptions

	ikeProposalConfig, err := sess.command("show configuration"+
		" security ike proposal "+ikeProposal+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ikeProposalConfig != emptyWord {
		confRead.name = ikeProposal
		for _, item := range strings.Split(ikeProposalConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "authentication-algorithm "):
				confRead.authenticationAlgorithm = strings.TrimPrefix(itemTrim, "authentication-algorithm ")
			case strings.HasPrefix(itemTrim, "authentication-method "):
				confRead.authenticationMethod = strings.TrimPrefix(itemTrim, "authentication-method ")
			case strings.HasPrefix(itemTrim, "dh-group "):
				confRead.dhGroup = strings.TrimPrefix(itemTrim, "dh-group ")
			case strings.HasPrefix(itemTrim, "encryption-algorithm"):
				confRead.encryptionAlgorithm = strings.TrimPrefix(itemTrim, "encryption-algorithm ")
			case strings.HasPrefix(itemTrim, "lifetime-seconds"):
				confRead.lifetimeSeconds, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "lifetime-seconds "))
				if err != nil {
					return confRead, err
				}
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delIkeProposal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike proposal "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillIkeProposalData(d *schema.ResourceData, ikeProposalOptions ikeProposalOptions) {
	tfErr := d.Set("name", ikeProposalOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_algorithm", ikeProposalOptions.authenticationAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_method", ikeProposalOptions.authenticationMethod)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("dh_group", ikeProposalOptions.dhGroup)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("encryption_algorithm", ikeProposalOptions.encryptionAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("lifetime_seconds", ikeProposalOptions.lifetimeSeconds)
	if tfErr != nil {
		panic(tfErr)
	}
}
