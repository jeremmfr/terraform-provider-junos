package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type asPathGroupOptions struct {
	dynamicDB bool
	name      string
	asPath    []map[string]interface{}
}

func resourcePolicyoptionsAsPathGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePolicyoptionsAsPathGroupCreate,
		ReadWithoutTimeout:   resourcePolicyoptionsAsPathGroupRead,
		UpdateWithoutTimeout: resourcePolicyoptionsAsPathGroupUpdate,
		DeleteWithoutTimeout: resourcePolicyoptionsAsPathGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyoptionsAsPathGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"as_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsAsPathGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setPolicyoptionsAsPathGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathGroupExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options as-path-group %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsAsPathGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsAsPathGroupExists, err = checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options as-path-group %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsAsPathGroupReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourcePolicyoptionsAsPathGroupReadWJunSess(d, clt, junSess)
}

func resourcePolicyoptionsAsPathGroupReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if asPathGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)
	}

	return nil
}

func resourcePolicyoptionsAsPathGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsAsPathGroup(d, clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsAsPathGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathGroupReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsAsPathGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathGroupExists {
		return nil, fmt.Errorf("don't find policy-options as-path-group with id '%v' (id must be <name>)", d.Id())
	}
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsAsPathGroupExists(name string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"policy-options as-path-group "+name+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsAsPathGroup(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set policy-options as-path-group " + d.Get("name").(string)
	asPathNameList := make([]string, 0)
	for _, v := range d.Get("as_path").([]interface{}) {
		asPath := v.(map[string]interface{})
		if bchk.InSlice(asPath["name"].(string), asPathNameList) {
			return fmt.Errorf("multiple blocks as_path with the same name %s", asPath["name"].(string))
		}
		asPathNameList = append(asPathNameList, asPath["name"].(string))
		configSet = append(configSet, setPrefix+
			" as-path "+asPath["name"].(string)+
			" \""+asPath["path"].(string)+"\"")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, setPrefix+" dynamic-db")
	}

	return clt.ConfigSet(configSet, junSess)
}

func readPolicyoptionsAsPathGroup(name string, clt *junos.Client, junSess *junos.Session,
) (confRead asPathGroupOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"policy-options as-path-group "+name+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "dynamic-db":
				confRead.dynamicDB = true
			case balt.CutPrefixInString(&itemTrim, "as-path "):
				itemTrimFields := strings.Split(itemTrim, " ")
				asPath := map[string]interface{}{
					"name": itemTrimFields[0],
					"path": strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "), "\""),
				}
				confRead.asPath = append(confRead.asPath, asPath)
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsAsPathGroup(asPathGroup string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path-group "+asPathGroup)

	return clt.ConfigSet(configSet, junSess)
}

func fillPolicyoptionsAsPathGroupData(d *schema.ResourceData, asPathGroupOptions asPathGroupOptions) {
	if tfErr := d.Set("name", asPathGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_path", asPathGroupOptions.asPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_db", asPathGroupOptions.dynamicDB); tfErr != nil {
		panic(tfErr)
	}
}
