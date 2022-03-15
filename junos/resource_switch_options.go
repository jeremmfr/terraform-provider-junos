package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type switchOptionsOptions struct {
	vtepSourceIf string
}

func resourceSwitchOptions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwitchOptionsCreate,
		ReadContext:   resourceSwitchOptionsRead,
		UpdateContext: resourceSwitchOptionsUpdate,
		DeleteContext: resourceSwitchOptionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSwitchOptionsImport,
		},
		Schema: map[string]*schema.Schema{
			"clean_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vtep_source_interface": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") != 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q need to have 1 dot", value, k))
					}

					return
				},
			},
		},
	}
}

func resourceSwitchOptionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSwitchOptions(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("switch_options")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setSwitchOptions(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_switch_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("switch_options")

	return append(diagWarns, resourceSwitchOptionsReadWJnprSess(d, m, jnprSess)...)
}

func resourceSwitchOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSwitchOptionsReadWJnprSess(d, m, jnprSess)
}

func resourceSwitchOptionsReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	switchOptionsOptions, err := readSwitchOptions(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSwitchOptions(d, switchOptionsOptions)

	return nil
}

func resourceSwitchOptionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delSwitchOptions(m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSwitchOptions(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSwitchOptions(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSwitchOptions(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_switch_options", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSwitchOptionsReadWJnprSess(d, m, jnprSess)...)
}

func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		sess := m.(*Session)
		if sess.junosFakeDeleteAlso {
			if err := delSwitchOptions(m, nil); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		jnprSess, err := sess.startNewSession()
		if err != nil {
			return diag.FromErr(err)
		}
		defer sess.closeSession(jnprSess)
		sess.configLock(jnprSess)
		var diagWarns diag.Diagnostics
		if err := delSwitchOptions(m, jnprSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := sess.commitConf("delete resource junos_switch_options", jnprSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceSwitchOptionsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	switchOptionsOptions, err := readSwitchOptions(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSwitchOptions(d, switchOptionsOptions)
	d.SetId("switch_options")
	result[0] = d

	return result, nil
}

func setSwitchOptions(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set switch-options "
	configSet := make([]string, 0)

	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return sess.configSet(configSet, jnprSess)
}

func delSwitchOptions(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := []string{"vtep-source-interface"}

	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func readSwitchOptions(m interface{}, jnprSess *NetconfObject) (switchOptionsOptions, error) {
	sess := m.(*Session)
	var confRead switchOptionsOptions

	showConfig, err := sess.command(cmdShowConfig+"switch-options | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "vtep-source-interface ") {
				confRead.vtepSourceIf = strings.TrimPrefix(itemTrim, "vtep-source-interface ")
			}
		}
	}

	return confRead, nil
}

func fillSwitchOptions(d *schema.ResourceData, switchOptionsOptions switchOptionsOptions) {
	if tfErr := d.Set("vtep_source_interface", switchOptionsOptions.vtepSourceIf); tfErr != nil {
		panic(tfErr)
	}
}
