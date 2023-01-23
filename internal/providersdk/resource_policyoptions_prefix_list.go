package providersdk

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type prefixListOptions struct {
	dynamicDB bool
	name      string
	applyPath string
	prefix    []string
}

func resourcePolicyoptionsPrefixList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePolicyoptionsPrefixListCreate,
		ReadWithoutTimeout:   resourcePolicyoptionsPrefixListRead,
		UpdateWithoutTimeout: resourcePolicyoptionsPrefixListUpdate,
		DeleteWithoutTimeout: resourcePolicyoptionsPrefixListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyoptionsPrefixListImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"apply_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateCIDRNetworkFunc(),
				},
			},
		},
	}
}

func resourcePolicyoptionsPrefixListCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setPolicyoptionsPrefixList(d, clt, nil); err != nil {
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
	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsPrefixListExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options prefix-list %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsPrefixList(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_policyoptions_prefix_list", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	policyoptsPrefixListExists, err = checkPolicyoptionsPrefixListExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsPrefixListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("policy-options prefix-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourcePolicyoptionsPrefixListReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsPrefixListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourcePolicyoptionsPrefixListReadWJunSess(d, clt, junSess)
}

func resourcePolicyoptionsPrefixListReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if prefixListOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsPrefixListData(d, prefixListOptions)
	}

	return nil
}

func resourcePolicyoptionsPrefixListUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delPolicyoptionsPrefixList(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsPrefixList(d, clt, nil); err != nil {
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
	if err := delPolicyoptionsPrefixList(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsPrefixList(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_policyoptions_prefix_list", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsPrefixListReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsPrefixListDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delPolicyoptionsPrefixList(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delPolicyoptionsPrefixList(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_policyoptions_prefix_list", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsPrefixListImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsPrefixListExists {
		return nil, fmt.Errorf("don't find policy-options prefix-list with id '%v' (id must be <name>)", d.Id())
	}
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsPrefixListData(d, prefixListOptions)

	result[0] = d

	return result, nil
}

func checkPolicyoptionsPrefixListExists(name string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+"policy-options prefix-list "+name+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsPrefixList(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set policy-options prefix-list " + d.Get("name").(string)
	configSet = append(configSet, setPrefix)
	if d.Get("apply_path").(string) != "" {
		replaceSign := strings.ReplaceAll(d.Get("apply_path").(string), "<", "&lt;")
		replaceSign = strings.ReplaceAll(replaceSign, ">", "&gt;")
		configSet = append(configSet, setPrefix+" apply-path \""+replaceSign+"\"")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, setPrefix+" dynamic-db")
	}
	for _, v := range sortSetOfString(d.Get("prefix").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+" "+v)
	}

	return clt.ConfigSet(configSet, junSess)
}

func readPolicyoptionsPrefixList(name string, clt *junos.Client, junSess *junos.Session,
) (confRead prefixListOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"policy-options prefix-list "+name+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			switch {
			case balt.CutPrefixInString(&itemTrim, "apply-path "):
				confRead.applyPath = html.UnescapeString(strings.Trim(itemTrim, "\""))
			case itemTrim == "dynamic-db":
				confRead.dynamicDB = true
			case strings.Contains(itemTrim, "/"):
				confRead.prefix = append(confRead.prefix, itemTrim)
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsPrefixList(prefixList string, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options prefix-list "+prefixList)

	return clt.ConfigSet(configSet, junSess)
}

func fillPolicyoptionsPrefixListData(d *schema.ResourceData, prefixListOptions prefixListOptions) {
	if tfErr := d.Set("name", prefixListOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("apply_path", prefixListOptions.applyPath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_db", prefixListOptions.dynamicDB); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("prefix", prefixListOptions.prefix); tfErr != nil {
		panic(tfErr)
	}
}
