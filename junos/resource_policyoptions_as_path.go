package junos

import (
	"context"
	"fmt"
	"strings"

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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setPolicyoptionsAsPath(d, clt, nil); err != nil {
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
	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if policyoptsAsPathExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("policy-options as-path %v already exists", d.Get("name").(string)))...)
	}

	if err := setPolicyoptionsAsPath(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

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
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourcePolicyoptionsAsPathReadWJunSess(d, clt, junSess)
}

func resourcePolicyoptionsAsPathReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession,
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setPolicyoptionsAsPath(d, clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setPolicyoptionsAsPath(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourcePolicyoptionsAsPathReadWJunSess(d, clt, junSess)...)
}

func resourcePolicyoptionsAsPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, nil); err != nil {
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
	if err := delPolicyoptionsAsPath(d.Get("name").(string), clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_policyoptions_as_path", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourcePolicyoptionsAsPathImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
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

func checkPolicyoptionsAsPathExists(name string, clt *Client, junSess *junosSession) (bool, error) {
	showConfig, err := clt.command(cmdShowConfig+"policy-options as-path "+name+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setPolicyoptionsAsPath(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" dynamic-db")
	}
	if d.Get("path").(string) != "" {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" \""+d.Get("path").(string)+"\"")
	}

	return clt.configSet(configSet, junSess)
}

func readPolicyoptionsAsPath(name string, clt *Client, junSess *junosSession) (asPathOptions, error) {
	var confRead asPathOptions

	showConfig, err := clt.command(cmdShowConfig+"policy-options as-path "+name+pipeDisplaySetRelative, junSess)
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
			default:
				confRead.path = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delPolicyoptionsAsPath(asPath string, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path "+asPath)

	return clt.configSet(configSet, junSess)
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
