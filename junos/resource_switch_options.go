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
		CreateWithoutTimeout: resourceSwitchOptionsCreate,
		ReadWithoutTimeout:   resourceSwitchOptionsRead,
		UpdateWithoutTimeout: resourceSwitchOptionsUpdate,
		DeleteWithoutTimeout: resourceSwitchOptionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSwitchOptionsImport,
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
		if err := setSwitchOptions(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("switch_options")

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
	if err := setSwitchOptions(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("switch_options")

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, sess, junSess)...)
}

func resourceSwitchOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(junSess)

	return resourceSwitchOptionsReadWJunSess(d, sess, junSess)
}

func resourceSwitchOptionsReadWJunSess(d *schema.ResourceData, sess *Session, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	switchOptionsOptions, err := readSwitchOptions(sess, junSess)
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
		if err := delSwitchOptions(sess, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSwitchOptions(d, sess, nil); err != nil {
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
	if err := delSwitchOptions(sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSwitchOptions(d, sess, junSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, sess, junSess)...)
}

func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		sess := m.(*Session)
		if sess.junosFakeDeleteAlso {
			if err := delSwitchOptions(sess, nil); err != nil {
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
		if err := delSwitchOptions(sess, junSess); err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := sess.commitConf("delete resource junos_switch_options", junSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceSwitchOptionsImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	junSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	switchOptionsOptions, err := readSwitchOptions(sess, junSess)
	if err != nil {
		return nil, err
	}
	fillSwitchOptions(d, switchOptionsOptions)
	d.SetId("switch_options")
	result[0] = d

	return result, nil
}

func setSwitchOptions(d *schema.ResourceData, sess *Session, junSess *junosSession) error {
	setPrefix := "set switch-options "
	configSet := make([]string, 0)

	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return sess.configSet(configSet, junSess)
}

func delSwitchOptions(sess *Session, junSess *junosSession) error {
	listLinesToDelete := []string{"vtep-source-interface"}

	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, junSess)
}

func readSwitchOptions(sess *Session, junSess *junosSession) (switchOptionsOptions, error) {
	var confRead switchOptionsOptions

	showConfig, err := sess.command(cmdShowConfig+"switch-options"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
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
