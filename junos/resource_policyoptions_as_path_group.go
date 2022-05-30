package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setPolicyoptionsAsPathGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathGroupExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options as-path-group %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsAsPathGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

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
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourcePolicyoptionsAsPathGroupReadWJunSess(d, clt, junSess)
}

func resourcePolicyoptionsAsPathGroupReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsAsPathGroup(d, clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsAsPathGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathGroupReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPathGroup(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_policyoptions_as_path_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsAsPathGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
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

func checkPolicyoptionsAsPathGroupExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+
		"policy-options as-path-group "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsAsPathGroup(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := "set policy-options as-path-group " + d.Get("name").(string)
	asPathNameList := make([]string, 0)
	for _, v := range d.Get("as_path").([]interface{}) {
		asPath := v.(map[string]interface{})
		if bchk.StringInSlice(asPath["name"].(string), asPathNameList) {
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

	return clt.configSet(configSet, junSess)
}

func readPolicyoptionsAsPathGroup(name string, clt *Client, junSess *junosSession) (asPathGroupOptions, error) {
	var confRead asPathGroupOptions

	showConfig, err := clt.command(cmdShowConfig+
		"policy-options as-path-group "+name+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case itemTrim == "dynamic-db":
				confRead.dynamicDB = true
			case strings.HasPrefix(itemTrim, "as-path "):
				asPath := map[string]interface{}{
					"name": "",
					"path": "",
				}
				itemSplit := strings.Split(strings.TrimPrefix(itemTrim, "as-path "), " ")
				asPath["name"] = itemSplit[0]
				asPath["path"] = strings.Trim(strings.TrimPrefix(itemTrim,
					"as-path "+asPath["name"].(string)+" "), "\"")
				confRead.asPath = append(confRead.asPath, asPath)
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsAsPathGroup(asPathGroup string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path-group "+asPathGroup)

	return clt.configSet(configSet, junSess)
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
