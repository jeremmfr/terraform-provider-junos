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
		CreateWithoutTimeout: resourceSecurityUtmCustomURLCategoryCreate,
		ReadWithoutTimeout:   resourceSecurityUtmCustomURLCategoryRead,
		UpdateWithoutTimeout: resourceSecurityUtmCustomURLCategoryUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmCustomURLCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmCustomURLCategoryImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
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

func resourceSecurityUtmCustomURLCategoryCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setUtmCustomURLCategory(d, m, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security utm custom-objects custom-url-category "+
			"not compatible with Junos device %s", jnprSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmCustomURLCategoryExists, err := checkUtmCustomURLCategorysExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLCategoryExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security utm custom-objects custom-url-category %v already exists", d.Get("name").(string)))...)
	}

	if err := setUtmCustomURLCategory(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

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

func resourceSecurityUtmCustomURLCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityUtmCustomURLCategoryReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityUtmCustomURLCategoryReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
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

func resourceSecurityUtmCustomURLCategoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delUtmCustomURLCategory(d.Get("name").(string), m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmCustomURLCategory(d, m, nil); err != nil {
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
	if err := delUtmCustomURLCategory(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmCustomURLCategory(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLCategoryReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityUtmCustomURLCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delUtmCustomURLCategory(d.Get("name").(string), m, nil); err != nil {
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
	if err := delUtmCustomURLCategory(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_utm_custom_url_category", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmCustomURLCategoryImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
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
	showConfig, err := sess.command(cmdShowConfig+
		"security utm custom-objects custom-url-category "+urlCategory+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
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

	return sess.configSet(configSet, jnprSess)
}

func readUtmCustomURLCategory(urlCategory string, m interface{}, jnprSess *NetconfObject,
) (utmCustomURLCategoryOptions, error) {
	sess := m.(*Session)
	var confRead utmCustomURLCategoryOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security utm custom-objects custom-url-category "+urlCategory+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = urlCategory
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
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

	return sess.configSet(configSet, jnprSess)
}

func fillUtmCustomURLCategoryData(d *schema.ResourceData, utmCustomURLCategoryOptions utmCustomURLCategoryOptions) {
	if tfErr := d.Set("name", utmCustomURLCategoryOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLCategoryOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
