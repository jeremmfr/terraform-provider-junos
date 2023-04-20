package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSwitchOptions(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("switch_options")

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setSwitchOptions(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_switch_options")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("switch_options")

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, junSess)...)
}

func resourceSwitchOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSwitchOptionsReadWJunSess(d, junSess)
}

func resourceSwitchOptionsReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	switchOptionsOptions, err := readSwitchOptions(junSess)
	junos.MutexUnlock()
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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSwitchOptions(junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSwitchOptions(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSwitchOptions(junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSwitchOptions(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_switch_options")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSwitchOptionsReadWJunSess(d, junSess)...)
}

func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		clt := m.(*junos.Client)
		if clt.FakeDeleteAlso() {
			junSess := clt.NewSessionWithoutNetconf(ctx)
			if err := delSwitchOptions(junSess); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		junSess, err := clt.StartNewSession(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		defer junSess.Close()
		if err := junSess.ConfigLock(ctx); err != nil {
			return diag.FromErr(err)
		}
		var diagWarns diag.Diagnostics
		if err := delSwitchOptions(junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := junSess.CommitConf("delete resource junos_switch_options")
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

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
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	switchOptionsOptions, err := readSwitchOptions(junSess)
	if err != nil {
		return nil, err
	}
	fillSwitchOptions(d, switchOptionsOptions)
	d.SetId("switch_options")
	result[0] = d

	return result, nil
}

func setSwitchOptions(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := "set switch-options "
	configSet := make([]string, 0)

	if v := d.Get("vtep_source_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return junSess.ConfigSet(configSet)
}

func delSwitchOptions(junSess *junos.Session) error {
	listLinesToDelete := []string{"vtep-source-interface"}

	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return junSess.ConfigSet(configSet)
}

func readSwitchOptions(junSess *junos.Session,
) (confRead switchOptionsOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "switch-options" + junos.PipeDisplaySetRelative)
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
