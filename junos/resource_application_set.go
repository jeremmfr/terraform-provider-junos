package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type applicationSetOptions struct {
	name         string
	applications []string
}

func resourceApplicationSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceApplicationSetCreate,
		ReadWithoutTimeout:   resourceApplicationSetRead,
		UpdateWithoutTimeout: resourceApplicationSetUpdate,
		DeleteWithoutTimeout: resourceApplicationSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationSetImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"applications": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApplicationSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setApplicationSet(d, m, nil); err != nil {
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
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	appSetExists, err := checkApplicationSetExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if appSetExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("application-set %v already exists", d.Get("name").(string)))...)
	}
	if err := setApplicationSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_application_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	appSetExists, err = checkApplicationSetExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if appSetExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("application-set %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceApplicationSetReadWJnprSess(d, m, jnprSess)...)
}

func resourceApplicationSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceApplicationSetReadWJnprSess(d, m, jnprSess)
}

func resourceApplicationSetReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject,
) diag.Diagnostics {
	mutex.Lock()
	applicationSetOptions, err := readApplicationSet(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if applicationSetOptions.name == "" {
		d.SetId("")
	} else {
		fillApplicationSetData(d, applicationSetOptions)
	}

	return nil
}

func resourceApplicationSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delApplicationSet(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setApplicationSet(d, m, nil); err != nil {
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
	if err := delApplicationSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setApplicationSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_application_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceApplicationSetReadWJnprSess(d, m, jnprSess)...)
}

func resourceApplicationSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delApplicationSet(d, m, nil); err != nil {
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
	if err := delApplicationSet(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_application_set", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceApplicationSetImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	appSetExists, err := checkApplicationSetExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !appSetExists {
		return nil, fmt.Errorf("don't find application-set with id '%v' (id must be <name>)", d.Id())
	}
	applicationSetOptions, err := readApplicationSet(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillApplicationSetData(d, applicationSetOptions)
	result[0] = d

	return result, nil
}

func checkApplicationSetExists(applicationSet string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command(cmdShowConfig+
		"applications application-set "+applicationSet+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setApplicationSet(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set applications application-set " + d.Get("name").(string)
	for _, v := range d.Get("applications").([]interface{}) {
		configSet = append(configSet, setPrefix+" application "+v.(string))
	}

	return sess.configSet(configSet, jnprSess)
}

func readApplicationSet(applicationSet string, m interface{}, jnprSess *NetconfObject) (applicationSetOptions, error) {
	sess := m.(*Session)
	var confRead applicationSetOptions

	showConfig, err := sess.command(cmdShowConfig+
		"applications application-set "+applicationSet+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = applicationSet
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.HasPrefix(itemTrim, "application ") {
				confRead.applications = append(confRead.applications, strings.TrimPrefix(itemTrim, "application "))
			}
		}
	}

	return confRead, nil
}

func delApplicationSet(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application-set "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
}

func fillApplicationSetData(d *schema.ResourceData, applicationSetOptions applicationSetOptions) {
	if tfErr := d.Set("name", applicationSetOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("applications", applicationSetOptions.applications); tfErr != nil {
		panic(tfErr)
	}
}
