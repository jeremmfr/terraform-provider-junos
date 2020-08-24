package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type utmCustomURLPatternOptions struct {
	name  string
	value []string
}

func resourceSecurityUtmCustomURLPattern() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityUtmCustomURLPatternCreate,
		Read:   resourceSecurityUtmCustomURLPatternRead,
		Update: resourceSecurityUtmCustomURLPatternUpdate,
		Delete: resourceSecurityUtmCustomURLPatternDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityUtmCustomURLPatternImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"value": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityUtmCustomURLPatternCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security utm custom-objects url-pattern "+
			"not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if utmCustomURLPatternExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("security utm custom-objects url-pattern %v already exists", d.Get("name").(string))
	}

	err = setUtmCustomURLPattern(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_security_utm_custom_url_pattern", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	mutex.Lock()
	utmCustomURLPatternExists, err = checkUtmCustomURLPatternsExists(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmCustomURLPatternExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security utm custom-objects url-pattern %v "+
			"not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceSecurityUtmCustomURLPatternRead(d, m)
}
func resourceSecurityUtmCustomURLPatternRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmCustomURLPatternOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)
	}

	return nil
}
func resourceSecurityUtmCustomURLPatternUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delUtmCustomURLPattern(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setUtmCustomURLPattern(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_security_utm_custom_url_pattern", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSecurityUtmCustomURLPatternRead(d, m)
}
func resourceSecurityUtmCustomURLPatternDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delUtmCustomURLPattern(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_security_utm_custom_url_pattern", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSecurityUtmCustomURLPatternImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLPatternExists {
		return nil, fmt.Errorf("don't find security utm custom-objects url-pattern with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLPatternsExists(urlPattern string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	urlPatternConfig, err := sess.command("show configuration security utm custom-objects url-pattern "+
		urlPattern+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if urlPatternConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setUtmCustomURLPattern(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects url-pattern " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readUtmCustomURLPattern(urlPattern string, m interface{}, jnprSess *NetconfObject) (
	utmCustomURLPatternOptions, error) {
	sess := m.(*Session)
	var confRead utmCustomURLPatternOptions

	urlPatternConfig, err := sess.command("show configuration"+
		" security utm custom-objects url-pattern "+urlPattern+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if urlPatternConfig != emptyWord {
		confRead.name = urlPattern
		for _, item := range strings.Split(urlPatternConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "value ") {
				confRead.value = append(confRead.value, strings.Trim(strings.TrimPrefix(itemTrim, "value "), "\""))
			}
		}
	} else {
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}

func delUtmCustomURLPattern(urlPattern string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects url-pattern "+urlPattern+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillUtmCustomURLPatternData(d *schema.ResourceData, utmCustomURLPatternOptions utmCustomURLPatternOptions) {
	tfErr := d.Set("name", utmCustomURLPatternOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("value", utmCustomURLPatternOptions.value)
	if tfErr != nil {
		panic(tfErr)
	}
}
