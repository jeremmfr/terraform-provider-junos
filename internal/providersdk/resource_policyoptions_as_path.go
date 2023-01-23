package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type asPathOptions struct {
	dynamicDB bool
	name      string
	path      string
}

func resourcePolicyoptionsAsPath() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePolicyoptionsAsPathCreate,
		ReadWithoutTimeout:   resourcePolicyoptionsAsPathRead,
		UpdateWithoutTimeout: resourcePolicyoptionsAsPathUpdate,
		DeleteWithoutTimeout: resourcePolicyoptionsAsPathDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyoptionsAsPathImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsAsPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setPolicyoptionsAsPath(d, clt, nil); err != nil {
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
	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options as-path %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsAsPath(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsAsPathExists, err = checkPolicyoptionsAsPathExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options as-path %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsAsPathReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourcePolicyoptionsAsPathReadWJunSess(d, clt, junSess)
}

func resourcePolicyoptionsAsPathReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	asPathOptions, err := readPolicyoptionsAsPath(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if asPathOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsAsPathData(d, asPathOptions)
	}

	return nil
}

func resourcePolicyoptionsAsPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsAsPath(d, clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsAsPath(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsAsPathImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathExists {
		return nil, fmt.Errorf("don't find policy-options as-path with id '%v' (id must be <name>)", d.Id())
	}
	asPathOptions, err := readPolicyoptionsAsPath(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathData(d, asPathOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsAsPathExists(name string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"policy-options as-path "+name+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsAsPath(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" dynamic-db")
	}
	if d.Get("path").(string) != "" {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" \""+d.Get("path").(string)+"\"")
	}

	return clt.ConfigSet(configSet, junSess)
}

func readPolicyoptionsAsPath(name string, clt *junos.Client, junSess *junos.Session,
) (confRead asPathOptions, err error) {
	showConfig, err := clt.Command(
		junos.CmdShowConfig+"policy-options as-path "+name+junos.PipeDisplaySetRelative,
		junSess,
	)
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
			default:
				confRead.path = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsAsPath(asPath string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path "+asPath)

	return clt.ConfigSet(configSet, junSess)
}

func fillPolicyoptionsAsPathData(d *schema.ResourceData, asPathOptions asPathOptions) {
	if tfErr := d.Set("name", asPathOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_db", asPathOptions.dynamicDB); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("path", asPathOptions.path); tfErr != nil {
		panic(tfErr)
	}
}
