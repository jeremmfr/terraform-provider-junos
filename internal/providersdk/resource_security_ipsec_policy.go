package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ipsecPolicyOptions struct {
	name        string
	pfsKeys     string
	proposalSet string
	proposals   []string
}

func resourceIpsecPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIpsecPolicyCreate,
		ReadWithoutTimeout:   resourceIpsecPolicyRead,
		UpdateWithoutTimeout: resourceIpsecPolicyUpdate,
		DeleteWithoutTimeout: resourceIpsecPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"pfs_keys": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposals": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
			"proposal_set": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listProposalSet(), false),
				ExactlyOneOf: []string{"proposals", "proposal_set"},
			},
		},
	}
}

func resourceIpsecPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setIpsecPolicy(d, clt, nil); err != nil {
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
	if !junos.CheckCompatibilitySecurity(junSess) {
		return diag.FromErr(fmt.Errorf("security ipsec policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecPolicyExists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security ipsec policy %v already exists", d.Get("name").(string)))...)
	}
	if err := setIpsecPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_security_ipsec_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecPolicyExists, err = checkIpsecPolicyExists(d.Get("name").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceIpsecPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceIpsecPolicyReadWJunSess(d, clt, junSess)
}

func resourceIpsecPolicyReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	ipsecPolicyOptions, err := readIpsecPolicy(d.Get("name").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecPolicyData(d, ipsecPolicyOptions)
	}

	return nil
}

func resourceIpsecPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delIpsecPolicy(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setIpsecPolicy(d, clt, nil); err != nil {
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
	if err := delIpsecPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIpsecPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_security_ipsec_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecPolicyReadWJunSess(d, clt, junSess)...)
}

func resourceIpsecPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delIpsecPolicy(d, clt, nil); err != nil {
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
	if err := delIpsecPolicy(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_security_ipsec_policy", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	ipsecPolicyExists, err := checkIpsecPolicyExists(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	if !ipsecPolicyExists {
		return nil, fmt.Errorf("don't find security ipsec policy with id '%v' (id must be <name>)", d.Id())
	}
	ipsecPolicyOptions, err := readIpsecPolicy(d.Id(), clt, junSess)
	if err != nil {
		return nil, err
	}
	fillIpsecPolicyData(d, ipsecPolicyOptions)
	result[0] = d

	return result, nil
}

func checkIpsecPolicyExists(ipsecPolicy string, clt *junos.Client, junSess *junos.Session) (bool, error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security ipsec policy "+ipsecPolicy+junos.PipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecPolicy(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security ipsec policy " + d.Get("name").(string)
	if d.Get("pfs_keys").(string) != "" {
		configSet = append(configSet, setPrefix+" perfect-forward-secrecy keys "+d.Get("pfs_keys").(string))
	}
	for _, v := range d.Get("proposals").([]interface{}) {
		configSet = append(configSet, setPrefix+" proposals "+v.(string))
	}
	if d.Get("proposal_set").(string) != "" {
		configSet = append(configSet, setPrefix+" proposal-set "+d.Get("proposal_set").(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func readIpsecPolicy(ipsecPolicy string, clt *junos.Client, junSess *junos.Session,
) (confRead ipsecPolicyOptions, err error) {
	showConfig, err := clt.Command(junos.CmdShowConfig+
		"security ipsec policy "+ipsecPolicy+junos.PipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ipsecPolicy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "perfect-forward-secrecy keys "):
				confRead.pfsKeys = itemTrim
			case balt.CutPrefixInString(&itemTrim, "proposals "):
				confRead.proposals = append(confRead.proposals, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proposal-set "):
				confRead.proposalSet = itemTrim
			}
		}
	}

	return confRead, nil
}

func delIpsecPolicy(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec policy "+d.Get("name").(string))

	return clt.ConfigSet(configSet, junSess)
}

func fillIpsecPolicyData(d *schema.ResourceData, ipsecPolicyOptions ipsecPolicyOptions) {
	if tfErr := d.Set("name", ipsecPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("pfs_keys", ipsecPolicyOptions.pfsKeys); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposals", ipsecPolicyOptions.proposals); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("proposal_set", ipsecPolicyOptions.proposalSet); tfErr != nil {
		panic(tfErr)
	}
}
