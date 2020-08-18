package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	jdecode "github.com/jeremmfr/junosdecode"
)

type ikePolicyOptions struct {
	name             string
	mode             string
	preSharedKeyText string
	preSharedKeyHexa string
	proposals        []string
}

func resourceIkePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIkePolicyCreate,
		Read:   resourceIkePolicyRead,
		Update: resourceIkePolicyUpdate,
		Delete: resourceIkePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIkePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"proposals": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "main",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"main", "aggressive"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'main' or 'aggressive'", value, k))
					}

					return
				},
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

func resourceIkePolicyCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security ike policy not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	ikePolicyExists, err := checkIkePolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if ikePolicyExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("security ike policy %v already exists", d.Get("name").(string))
	}
	err = setIkePolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_security_ike_policy", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	ikePolicyExists, err = checkIkePolicyExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if ikePolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security ike policy %v not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceIkePolicyRead(d, m)
}
func resourceIkePolicyRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	ikePolicyOptions, err := readIkePolicy(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if ikePolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIkePolicyData(d, ikePolicyOptions)
	}

	return nil
}
func resourceIkePolicyUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delIkePolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setIkePolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_security_ike_policy", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceIkePolicyRead(d, m)
}
func resourceIkePolicyDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delIkePolicy(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_security_ike_policy", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
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
		configSet = append(configSet, setPrefix+" mode "+d.Get("mode").(string)+"\n")
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string)+"\n")
	}
	if d.Get("pre_shared_key_text").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key ascii-text "+d.Get("pre_shared_key_text").(string)+"\n")
	}
	if d.Get("pre_shared_key_hexa").(string) != "" {
		configSet = append(configSet, setPrefix+" pre-shared-key hexadecimal "+d.Get("pre_shared_key_hexa").(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
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
			case strings.HasPrefix(itemTrim, "pre-shared-key hexadecimal "):
				confRead.preSharedKeyHexa, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"pre-shared-key hexadecimal "), "\""))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "pre-shared-key ascii-text "):
				confRead.preSharedKeyText, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"pre-shared-key ascii-text "), "\""))
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
func delIkePolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike policy "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillIkePolicyData(d *schema.ResourceData, ikePolicyOptions ikePolicyOptions) {
	tfErr := d.Set("name", ikePolicyOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("mode", ikePolicyOptions.mode)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("pre_shared_key_text", ikePolicyOptions.preSharedKeyText)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("pre_shared_key_hexa", ikePolicyOptions.preSharedKeyHexa)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("proposals", ikePolicyOptions.proposals)
	if tfErr != nil {
		panic(tfErr)
	}
}
