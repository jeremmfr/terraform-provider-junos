package providersdk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type ribGroupOptions struct {
	name         string
	exportRib    string
	importPolicy []string
	importRib    []string
}

func resourceRibGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRibGroupCreate,
		ReadWithoutTimeout:   resourceRibGroupRead,
		UpdateWithoutTimeout: resourceRibGroupUpdate,
		DeleteWithoutTimeout: resourceRibGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRibGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"import_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setRibGroup(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	ribGroupExists, err := checkRibGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ribGroupExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("rib-group %v already exists", d.Get("name").(string)))...)
	}
	if err := setRibGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_rib_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ribGroupExists, err = checkRibGroupExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ribGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("rib-group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceRibGroupReadWJunSess(d, junSess)...)
}

func resourceRibGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceRibGroupReadWJunSess(d, junSess)
}

func resourceRibGroupReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	junos.MutexLock()
	ribGroupOptions, err := readRibGroup(d.Get("name").(string), junSess)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if d.HasChange("import_policy") {
			if err := delRibGroupElement("import-policy", d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		if d.HasChange("import_rib") {
			if err := delRibGroupElement("import-rib", d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		if d.HasChange("export_rib") {
			if err := delRibGroupElement("export-rib", d.Get("name").(string), junSess); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := setRibGroup(d, junSess); err != nil {
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
	if d.HasChange("import_policy") {
		if err := delRibGroupElement("import-policy", d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if d.HasChange("import_rib") {
		if err := delRibGroupElement("import-rib", d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if d.HasChange("export_rib") {
		if err := delRibGroupElement("export-rib", d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if err := setRibGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_rib_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceRibGroupReadWJunSess(d, junSess)...)
}

func resourceRibGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delRibGroup(d, junSess); err != nil {
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
	if err := delRibGroup(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_rib_group")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceRibGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ribGroupExists, err := checkRibGroupExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ribGroupExists {
		return nil, fmt.Errorf("don't find rib group with id '%v' (id must be <name>)", d.Id())
	}
	rigGroupOptions, err := readRibGroup(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillRibGroupData(d, rigGroupOptions)
	result[0] = d

	return result, nil
}

func checkRibGroupExists(group string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingOptionsWS + "rib-groups " + group + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setRibGroup(d *schema.ResourceData, junSess *junos.Session) error {
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

	return junSess.ConfigSet(configSet)
}

func readRibGroup(group string, junSess *junos.Session,
) (confRead ribGroupOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingOptionsWS + "rib-groups " + group + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = group
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "import-policy "):
				confRead.importPolicy = append(confRead.importPolicy, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "import-rib "):
				confRead.importRib = append(confRead.importRib, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "export-rib "):
				confRead.exportRib = itemTrim
			}
		}
	}

	return confRead, nil
}

func delRibGroupElement(element, group string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-options rib-groups "+group+" "+element)

	return junSess.ConfigSet(configSet)
}

func delRibGroup(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-options rib-groups "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func validateRibGroup(d *schema.ResourceData) error {
	var err string
	for _, v := range d.Get("import_rib").([]interface{}) {
		if !strings.HasSuffix(v.(string), ".inet.0") && !strings.HasSuffix(v.(string), ".inet6.0") {
			err = err + "rib-group " + v.(string) + " invalid name (missing .inet.0 or .inet6.0),"
		}
	}
	if d.Get("export_rib").(string) != "" {
		if !strings.HasSuffix(d.Get("export_rib").(string), ".inet.0") &&
			!strings.HasSuffix(d.Get("export_rib").(string), ".inet6.0") {
			err = err + "rib-group " + d.Get("export_rib").(string) + " invalid name (missing .inet.0 or .inet6.0),"
		}
	}
	if err != "" {
		return errors.New(err)
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
