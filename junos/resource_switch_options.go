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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSwitchOptions(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("switch_options")

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setSwitchOptions(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("switch_options")

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, clt, junSess)...)
}

func resourceSwitchOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSwitchOptionsReadWJunSess(d, clt, junSess)
}

func resourceSwitchOptionsReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	switchOptionsOptions, err := readSwitchOptions(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSwitchOptions(d, switchOptionsOptions)

	return nil
}

func resourceSwitchOptionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSwitchOptions(clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSwitchOptions(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSwitchOptions(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSwitchOptions(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, clt, junSess)...)
}

func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		clt := m.(*Client)
		if clt.fakeDeleteAlso {
			if err := delSwitchOptions(clt, nil); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		junSess, err := clt.startNewSession(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		defer clt.closeSession(junSess)
		if err := clt.configLock(ctx, junSess); err != nil {
			return diag.FromErr(err)
		}
		var diagWarns diag.Diagnostics
		if err := delSwitchOptions(clt, junSess); err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := clt.commitConf("delete resource junos_switch_options", junSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceSwitchOptionsImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	switchOptionsOptions, err := readSwitchOptions(clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSwitchOptions(d, switchOptionsOptions)
	d.SetId("switch_options")
	result[0] = d

	return result, nil
}

func setSwitchOptions(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set switch-options "
	configSet := make([]string, 0)

	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return clt.configSet(configSet, junSess)
}

func delSwitchOptions(clt *Client, junSess *junosSession) error {
	listLinesToDelete := []string{"vtep-source-interface"}

	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return clt.configSet(configSet, junSess)
}

func readSwitchOptions(clt *Client, junSess *junosSession) (switchOptionsOptions, error) {
	var confRead switchOptionsOptions

	showConfig, err := clt.command(cmdShowConfig+"switch-options"+pipeDisplaySetRelative, junSess)
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
