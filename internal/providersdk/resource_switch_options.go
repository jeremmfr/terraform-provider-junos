package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setSwitchOptions(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("switch_options")

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setSwitchOptions(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("switch_options")

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, clt, junSess)...)
}

func resourceSwitchOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceSwitchOptionsReadWJunSess(d, clt, junSess)
}

func resourceSwitchOptionsReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delSwitchOptions(clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSwitchOptions(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSwitchOptions(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSwitchOptions(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_switch_options", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, clt, junSess)...)
}

func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		clt := m.(*junos.Client)
		if clt.FakeDeleteAlso() {
			if err := delSwitchOptions(clt, nil); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		junSess, err := clt.StartNewSession(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		defer clt.CloseSession(junSess)
		if err := clt.ConfigLock(ctx, junSess); err != nil {
			return diag.FromErr(err)
		}
		var diagWarns diag.Diagnostics
		if err := delSwitchOptions(clt, junSess); err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := clt.CommitConf("delete resource junos_switch_options", junSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceSwitchOptionsImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
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

func setSwitchOptions(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	setPrefix := "set switch-options "
	configSet := make([]string, 0)

	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return clt.ConfigSet(configSet, junSess)
}

func delSwitchOptions(clt *junos.Client, junSess *junos.Session) error {
	listLinesToDelete := []string{"vtep-source-interface"}

	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return clt.ConfigSet(configSet, junSess)
}

func readSwitchOptions(clt *junos.Client, junSess *junos.Session) (confRead switchOptionsOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"switch-options"+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "vtep-source-interface ") {
				confRead.vtepSourceIf = itemTrim
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
