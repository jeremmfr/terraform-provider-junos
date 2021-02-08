package junos

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type switchOptionsOptions struct {
	routeDistinguisher string
	vrfExport []string
	vrfImport []string
	vrfTarget string
	vtepSourceInterface string
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
			"route_distinguisher": {
				Type:             schema.TypeString,
				Optional:         true,
			},
			"vrf_import": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrf_export": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrf_target": {
				Type:             schema.TypeString,
				Optional:         true,
			},
			"vtep_source_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSwitchOptionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setSwitchOptions(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("create resource junos_switch_options", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	d.SetId("switch_options")

	return resourceSwitchOptionsReadWJnprSess(d, m, jnprSess)
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
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSwitchOptions(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSwitchOptions(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("update resource junos_switch_options", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceSwitchOptionsReadWJnprSess(d, m, jnprSess)
}
func resourceSwitchOptionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.SetId("routing_options")
	result[0] = d

	return result, nil
}

func setSwitchOptions(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set switch-options "
	if d.Get("route_distinguisher").(string) != "" {
		configSet = append(configSet, setPrefix+"route-distinguisher "+d.Get("route_distinguisher").(string))
	}
	for _, v := range d.Get("vrf_import").([]interface{}) {
		configSet = append(configSet, setPrefix+" vrf-import "+v.(string))
	}
	for _, v := range d.Get("vrf_export").([]interface{}) {
			configSet = append(configSet, setPrefix+" vrf-export "+v.(string))
	}
	if d.Get("vrf_target").(string) != "" {
		configSet = append(configSet, setPrefix+"vrf-target "+d.Get("vrf_target").(string))
	}
	if d.Get("vtep_source_interface").(string) != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+d.Get("vtep_source_interface").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func delSwitchOptions(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := []string{
		"route-distinguisher",
		"vrf-import",
		"vrf-export",
		"vrf-target",
		"vtep-source-interface",
	}
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete switch-options "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSwitchOptions(m interface{}, jnprSess *NetconfObject) (switchOptionsOptions, error) {
	sess := m.(*Session)
	var confRead switchOptionsOptions

	switchOptionsConfig, err := sess.command("show configuration switch-options"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if switchOptionsConfig != emptyWord {
		for _, item := range strings.Split(switchOptionsConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
				case strings.HasPrefix(itemTrim, "route-distinguisher "):
					confRead.routeDistinguisher = strings.TrimPrefix(itemTrim, "route-distinguisher ")
				case strings.HasPrefix(itemTrim, "vtep-source-interface "):
					confRead.vtepSourceInterface = strings.TrimPrefix(itemTrim, "vtep-source-interface ")
				case strings.HasPrefix(itemTrim, "vrf-target "):
					confRead.vrfTarget = strings.TrimPrefix(itemTrim, "vrf-target ")
				case strings.HasPrefix(itemTrim, "vrf-import "):
					confRead.vrfImport = append(confRead.vrfImport, strings.TrimPrefix(itemTrim, "vrf-import "))
				case strings.HasPrefix(itemTrim, "vrf-export "):
					confRead.vrfExport = append(confRead.vrfExport, strings.TrimPrefix(itemTrim, "vrf-export "))
				}
		}
	}
	return confRead, nil
}

func fillSwitchOptions(d *schema.ResourceData, switchOptionsOptions switchOptionsOptions) {
	if tfErr := d.Set("route_distinguisher", switchOptionsOptions.routeDistinguisher); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_import", switchOptionsOptions.vrfImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_export", switchOptionsOptions.vrfExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vrf_target", switchOptionsOptions.vrfTarget); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vtep_source_interface", switchOptionsOptions.vtepSourceInterface); tfErr != nil {
		panic(tfErr)
	}
}
