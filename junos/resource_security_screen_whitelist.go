package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type screenWhiteListOptions struct {
	name    string
	address []string
}

func resourceSecurityScreenWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityScreenWhiteListCreate,
		ReadContext:   resourceSecurityScreenWhiteListRead,
		UpdateContext: resourceSecurityScreenWhiteListUpdate,
		DeleteContext: resourceSecurityScreenWhiteListDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityScreenWhiteListImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityScreenWhiteListCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security screen white-list not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityScreenWhiteListExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security screen white-list %v already exists", d.Get("name").(string)))
	}

	if err := setSecurityScreenWhiteList(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_screen_whitelist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityScreenWhiteListExists, err = checkSecurityScreenWhiteListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenWhiteListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen white-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityScreenWhiteListRead(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityScreenWhiteListReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityScreenWhiteListReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	whiteListOptions, err := readSecurityScreenWhiteList(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if whiteListOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityScreenWhiteListData(d, whiteListOptions)
	}

	return nil
}
func resourceSecurityScreenWhiteListUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delSecurityScreenWhiteList(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setSecurityScreenWhiteList(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_screen_whitelist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityScreenWhiteListReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityScreenWhiteListDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityScreenWhiteList(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_screen_whitelist", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityScreenWhiteListImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityScreenWhiteListExists, err := checkSecurityScreenWhiteListExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityScreenWhiteListExists {
		return nil, fmt.Errorf("don't find screen white-list with id '%v' (id must be <name>)", d.Id())
	}
	whiteListOptions, err := readSecurityScreenWhiteList(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityScreenWhiteListData(d, whiteListOptions)

	result[0] = d

	return result, nil
}

func checkSecurityScreenWhiteListExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	whiteListConfig, err := sess.command("show configuration"+
		" security screen white-list "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if whiteListConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityScreenWhiteList(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security screen white-list " + d.Get("name").(string) + " "

	for _, v := range d.Get("address").([]interface{}) {
		if err := validateCIDRNetwork(v.(string)); err != nil {
			return err
		}
		configSet = append(configSet, setPrefix+"address "+v.(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func readSecurityScreenWhiteList(name string, m interface{}, jnprSess *NetconfObject) (screenWhiteListOptions, error) {
	sess := m.(*Session)
	var confRead screenWhiteListOptions

	whiteListConfig, err := sess.command("show configuration security screen white-list "+
		name+" | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	if whiteListConfig != emptyWord {
		confRead.name = name
		for _, item := range strings.Split(whiteListConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "address ") {
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			}
		}
	}

	return confRead, nil
}

func delSecurityScreenWhiteList(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security screen white-list "+name)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillSecurityScreenWhiteListData(d *schema.ResourceData, whiteListOptions screenWhiteListOptions) {
	if tfErr := d.Set("name", whiteListOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", whiteListOptions.address); tfErr != nil {
		panic(tfErr)
	}
}
