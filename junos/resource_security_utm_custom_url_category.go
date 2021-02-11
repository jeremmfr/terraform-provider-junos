package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type utmCustomURLCategoryOptions struct {
	name  string
	value []string
}

func resourceSecurityUtmCustomURLCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityUtmCustomURLCategoryCreate,
		ReadContext:   resourceSecurityUtmCustomURLCategoryRead,
		UpdateContext: resourceSecurityUtmCustomURLCategoryUpdate,
		DeleteContext: resourceSecurityUtmCustomURLCategoryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityUtmCustomURLCategoryImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
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

func resourceSecurityUtmCustomURLCategoryCreate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category "+
			"not compatible with Junos device %s", jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if utmCustomURLCategoryExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf(
			"security utm custom-objects custom-url-category %v already exists", d.Get("name").(string)))
	}

	if err := setUtmCustomURLCategory(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmCustomURLCategoryExists, err = checkUtmCustomURLCategorysExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityUtmCustomURLCategoryRead(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityUtmCustomURLCategoryReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityUtmCustomURLCategoryReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmCustomURLCategoryOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmCustomURLCategoryData(d, utmCustomURLCategoryOptions)
	}

	return nil
}
func resourceSecurityUtmCustomURLCategoryUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delUtmCustomURLCategory(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setUtmCustomURLCategory(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityUtmCustomURLCategoryDelete(
	ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delUtmCustomURLCategory(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityUtmCustomURLCategoryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLCategoryExists {
		return nil, fmt.Errorf(
			"missing security utm custom-objects custom-url-category with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLCategoryOptions, err := readUtmCustomURLCategory(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLCategoryData(d, utmCustomURLCategoryOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLCategorysExists(urlCategory string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	urlCategoryConfig, err := sess.command("show configuration security utm custom-objects custom-url-category "+
		urlCategory+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if urlCategoryConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setUtmCustomURLCategory(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects custom-url-category " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readUtmCustomURLCategory(urlCategory string, m interface{}, jnprSess *NetconfObject) (
	utmCustomURLCategoryOptions, error) {
	sess := m.(*Session)
	var confRead utmCustomURLCategoryOptions

	urlCategoryConfig, err := sess.command("show configuration"+
		" security utm custom-objects custom-url-category "+urlCategory+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if urlCategoryConfig != emptyWord {
		confRead.name = urlCategory
		for _, item := range strings.Split(urlCategoryConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "value ") {
				confRead.value = append(confRead.value, strings.TrimPrefix(itemTrim, "value "))
			}
		}
	}

	return confRead, nil
}

func delUtmCustomURLCategory(urlCategory string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects custom-url-category "+urlCategory)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillUtmCustomURLCategoryData(d *schema.ResourceData, utmCustomURLCategoryOptions utmCustomURLCategoryOptions) {
	if tfErr := d.Set("name", utmCustomURLCategoryOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLCategoryOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
