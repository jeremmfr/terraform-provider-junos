package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type communityOptions struct {
	invertMatch bool
	name        string
	members     []string
}

func resourcePolicyoptionsCommunity() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyoptionsCommunityCreate,
		Read:   resourcePolicyoptionsCommunityRead,
		Update: resourcePolicyoptionsCommunityUpdate,
		Delete: resourcePolicyoptionsCommunityDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsCommunityImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"members": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"invert_match": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsCommunityCreate(d *schema.ResourceData, m interface{}) error {
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
	policyoptsCommunityExists, err := checkPolicyoptionsCommunityExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if policyoptsCommunityExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("policy-options community %v already exists", d.Get("name").(string))
	}

	err = setPolicyoptionsCommunity(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_policyoptions_community", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	policyoptsCommunityExists, err = checkPolicyoptionsCommunityExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if policyoptsCommunityExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("policy-options community %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourcePolicyoptionsCommunityRead(d, m)
}
func resourcePolicyoptionsCommunityRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	communityOptions, err := readPolicyoptionsCommunity(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if communityOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsCommunityData(d, communityOptions)
	}
	return nil
}
func resourcePolicyoptionsCommunityUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsCommunity(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setPolicyoptionsCommunity(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_policyoptions_community", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourcePolicyoptionsCommunityRead(d, m)
}
func resourcePolicyoptionsCommunityDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsCommunity(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_policyoptions_community", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourcePolicyoptionsCommunityImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsCommunityExists, err := checkPolicyoptionsCommunityExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsCommunityExists {
		return nil, fmt.Errorf("don't find policy-options community with id '%v' (id must be <name>)", d.Id())
	}
	communityOptions, err := readPolicyoptionsCommunity(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsCommunityData(d, communityOptions)

	result[0] = d
	return result, nil
}

func checkPolicyoptionsCommunityExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	communityConfig, err := sess.command("show configuration policy-options community "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if communityConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setPolicyoptionsCommunity(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set policy-options community " + d.Get("name").(string) + " "
	for _, v := range d.Get("members").([]interface{}) {
		configSet = append(configSet, setPrefix+"members "+v.(string)+"\n")
	}
	if d.Get("invert_match").(bool) {
		configSet = append(configSet, setPrefix+"invert-match\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readPolicyoptionsCommunity(community string, m interface{}, jnprSess *NetconfObject) (communityOptions, error) {
	sess := m.(*Session)
	var confRead communityOptions

	communityConfig, err := sess.command("show configuration"+
		" policy-options community "+community+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if communityConfig != emptyWord {
		confRead.name = community
		for _, item := range strings.Split(communityConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "members "):
				confRead.members = append(confRead.members, strings.TrimPrefix(itemTrim, "members "))
			case strings.HasPrefix(itemTrim, "invert-match"):
				confRead.invertMatch = true
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delPolicyoptionsCommunity(community string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options community "+community+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillPolicyoptionsCommunityData(d *schema.ResourceData, communityOptions communityOptions) {
	tfErr := d.Set("name", communityOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("members", communityOptions.members)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("invert_match", communityOptions.invertMatch)
	if tfErr != nil {
		panic(tfErr)
	}
}
