package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setPolicyoptionsAsPath(d, junSess); err != nil {
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
	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options as-path %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsAsPath(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_policyoptions_as_path")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsAsPathExists, err = checkPolicyoptionsAsPathExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options as-path %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsAsPathReadWJunSess(d, junSess)...)
}

func resourcePolicyoptionsAsPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourcePolicyoptionsAsPathReadWJunSess(d, junSess)
}

func resourcePolicyoptionsAsPathReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	asPathOptions, err := readPolicyoptionsAsPath(d.Get("name").(string), junSess)
	junos.MutexUnlock()
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
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delPolicyoptionsAsPath(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsAsPath(d, junSess); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsAsPath(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_policyoptions_as_path")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathReadWJunSess(d, junSess)...)
}

func resourcePolicyoptionsAsPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delPolicyoptionsAsPath(d.Get("name").(string), junSess); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_policyoptions_as_path")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

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
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathExists {
		return nil, fmt.Errorf("don't find policy-options as-path with id '%v' (id must be <name>)", d.Id())
	}
	asPathOptions, err := readPolicyoptionsAsPath(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathData(d, asPathOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsAsPathExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "policy-options as-path " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsAsPath(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" dynamic-db")
	}
	if d.Get("path").(string) != "" {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" \""+d.Get("path").(string)+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readPolicyoptionsAsPath(name string, junSess *junos.Session,
) (confRead asPathOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options as-path " + name + junos.PipeDisplaySetRelative)
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

func delPolicyoptionsAsPath(asPath string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path "+asPath)

	return junSess.ConfigSet(configSet)
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
