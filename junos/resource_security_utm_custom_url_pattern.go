package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type utmCustomURLPatternOptions struct {
	name  string
	value []string
}

func resourceSecurityUtmCustomURLPattern() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityUtmCustomURLPatternCreate,
		ReadWithoutTimeout:   resourceSecurityUtmCustomURLPatternRead,
		UpdateWithoutTimeout: resourceSecurityUtmCustomURLPatternUpdate,
		DeleteWithoutTimeout: resourceSecurityUtmCustomURLPatternDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityUtmCustomURLPatternImport,
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

func resourceSecurityUtmCustomURLPatternCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setUtmCustomURLPattern(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if !checkCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLPatternExists {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern %v already exists", d.Get("name").(string)))...)
	}
	if err := setUtmCustomURLPattern(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_utm_custom_url_pattern", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	utmCustomURLPatternExists, err = checkUtmCustomURLPatternsExists(d.Get("name").(string), sess, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if utmCustomURLPatternExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security utm custom-objects url-pattern %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityUtmCustomURLPatternReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityUtmCustomURLPatternRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceSecurityUtmCustomURLPatternReadWJunSess(d, sess, junSess)
}

func resourceSecurityUtmCustomURLPatternReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Get("name").(string), sess, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if utmCustomURLPatternOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)
	}

	return nil
}

func resourceSecurityUtmCustomURLPatternUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delUtmCustomURLPattern(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setUtmCustomURLPattern(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delUtmCustomURLPattern(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setUtmCustomURLPattern(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_utm_custom_url_pattern", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityUtmCustomURLPatternReadWJunSess(d, sess, junSess)...)
}

func resourceSecurityUtmCustomURLPatternDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delUtmCustomURLPattern(d.Get("name").(string), sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)
	if err := sess.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delUtmCustomURLPattern(d.Get("name").(string), sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_utm_custom_url_pattern", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityUtmCustomURLPatternImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	utmCustomURLPatternExists, err := checkUtmCustomURLPatternsExists(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	if !utmCustomURLPatternExists {
		return nil, fmt.Errorf("don't find security utm custom-objects url-pattern with id '%v' (id must be <name>)", d.Id())
	}
	utmCustomURLPatternOptions, err := readUtmCustomURLPattern(d.Id(), sess, junSess)
	if err != nil {
		return nil, err
	}
	fillUtmCustomURLPatternData(d, utmCustomURLPatternOptions)

	result[0] = d

	return result, nil
}

func checkUtmCustomURLPatternsExists(urlPattern string, sess *Session, junSess *junosSession) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+
		"security utm custom-objects url-pattern "+urlPattern+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setUtmCustomURLPattern(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set security utm custom-objects url-pattern " + d.Get("name").(string) + " "
	for _, v := range d.Get("value").([]interface{}) {
		configSet = append(configSet, setPrefix+"value "+v.(string))
	}

	return sess.configSet(configSet, junSess)
}

func readUtmCustomURLPattern(urlPattern string, sess *Session, junSess *junosSession,
) (utmCustomURLPatternOptions, error) {
	var confRead utmCustomURLPatternOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security utm custom-objects url-pattern "+urlPattern+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = urlPattern
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.HasPrefix(itemTrim, "value ") {
				confRead.value = append(confRead.value, strings.Trim(strings.TrimPrefix(itemTrim, "value "), "\""))
			}
		}
	}

	return confRead, nil
}

func delUtmCustomURLPattern(urlPattern string, sess *Session, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm custom-objects url-pattern "+urlPattern)

	return sess.configSet(configSet, junSess)
}

func fillUtmCustomURLPatternData(d *schema.ResourceData, utmCustomURLPatternOptions utmCustomURLPatternOptions) {
	if tfErr := d.Set("name", utmCustomURLPatternOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("value", utmCustomURLPatternOptions.value); tfErr != nil {
		panic(tfErr)
	}
}
