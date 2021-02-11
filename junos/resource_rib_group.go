package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ribGroupOptions struct {
	name         string
	exportRib    string
	importPolicy []string
	importRib    []string
}

func resourceRibGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRibGroupCreate,
		ReadContext:   resourceRibGroupRead,
		UpdateContext: resourceRibGroupUpdate,
		DeleteContext: resourceRibGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRibGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"import_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"import_rib": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"export_rib": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRibGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := validateRibGroup(d); err != nil {
		return diag.FromErr(err)
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	ribGroupExists, err := checkRibGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ribGroupExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("rib-group %v already exists", d.Get("name").(string)))
	}
	if err := setRibGroup(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_rib_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	ribGroupExists, err = checkRibGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ribGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("rib-group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceRibGroupReadWJnprSess(d, m, jnprSess)...)
}
func resourceRibGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceRibGroupReadWJnprSess(d, m, jnprSess)
}
func resourceRibGroupReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ribGroupOptions, err := readRibGroup(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ribGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillRibGroupData(d, ribGroupOptions)
	}

	return nil
}
func resourceRibGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	if err := validateRibGroup(d); err != nil {
		return diag.FromErr(err)
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if d.HasChange("import_policy") {
		err = delRibGroupElement("import-policy", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
	}
	if d.HasChange("import_rib") {
		err = delRibGroupElement("import-rib", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
	}
	if d.HasChange("export_rib") {
		err = delRibGroupElement("export-rib", d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
	}
	if err := setRibGroup(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_rib_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRibGroupReadWJnprSess(d, m, jnprSess)...)
}
func resourceRibGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delRibGroup(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_rib_group", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceRibGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ribGroupExists, err := checkRibGroupExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ribGroupExists {
		return nil, fmt.Errorf("don't find rib group with id '%v' (id must be <name>)", d.Id())
	}
	rigGroupOptions, err := readRibGroup(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRibGroupData(d, rigGroupOptions)
	result[0] = d

	return result, nil
}

func checkRibGroupExists(group string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	rigGroupConfig, err := sess.command("show configuration routing-options rib-groups "+group+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if rigGroupConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setRibGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set routing-options rib-groups " + d.Get("name").(string) + " "
	for _, v := range d.Get("import_policy").([]interface{}) {
		configSet = append(configSet, setPrefix+"import-policy "+v.(string))
	}
	for _, v := range d.Get("import_rib").([]interface{}) {
		configSet = append(configSet, setPrefix+"import-rib "+v.(string))
	}
	if d.Get("export_rib").(string) != "" {
		configSet = append(configSet, setPrefix+"export-rib "+d.Get("export_rib").(string))
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readRibGroup(group string, m interface{}, jnprSess *NetconfObject) (ribGroupOptions, error) {
	sess := m.(*Session)
	var confRead ribGroupOptions

	ribGroupConfig, err := sess.command("show configuration"+
		" routing-options rib-groups "+group+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ribGroupConfig != emptyWord {
		confRead.name = group
		for _, item := range strings.Split(ribGroupConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "import-policy "):
				confRead.importPolicy = append(confRead.importPolicy, strings.TrimPrefix(itemTrim, "import-policy "))
			case strings.HasPrefix(itemTrim, "import-rib "):
				confRead.importRib = append(confRead.importRib, strings.TrimPrefix(itemTrim, "import-rib "))
			case strings.HasPrefix(itemTrim, "export-rib "):
				confRead.exportRib = strings.TrimPrefix(itemTrim, "export-rib ")
			}
		}
	}

	return confRead, nil
}
func delRibGroupElement(element string, group string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-options rib-groups "+group+" "+element)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func delRibGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-options rib-groups "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func validateRibGroup(d *schema.ResourceData) error {
	var errors string
	for _, v := range d.Get("import_rib").([]interface{}) {
		if !strings.HasSuffix(v.(string), ".inet.0") && !strings.HasSuffix(v.(string), ".inet6.0") {
			errors = errors + "rib-group " + v.(string) + " invalid name (missing .inet.0 or .inet6.0),"
		}
	}
	if d.Get("export_rib").(string) != "" {
		if !strings.HasSuffix(d.Get("export_rib").(string), ".inet.0") &&
			!strings.HasSuffix(d.Get("export_rib").(string), ".inet6.0") {
			errors = errors + "rib-group " + d.Get("export_rib").(string) + " invalid name (missing .inet.0 or .inet6.0),"
		}
	}
	if errors != "" {
		return fmt.Errorf(errors)
	}

	return nil
}
func fillRibGroupData(d *schema.ResourceData, ribGroupOptions ribGroupOptions) {
	if tfErr := d.Set("name", ribGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import_policy", ribGroupOptions.importPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import_rib", ribGroupOptions.importRib); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export_rib", ribGroupOptions.exportRib); tfErr != nil {
		panic(tfErr)
	}
}
